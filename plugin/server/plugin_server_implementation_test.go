package server

import (
	"context"
	"testing"
	"time"

	"github.com/blinkops/blink-sdk/plugin"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
	"github.com/stretchr/testify/suite"
)

type ServerImplementationTestSuite struct {
	suite.Suite
	service *PluginGRPCService
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServerImplementationSuite(t *testing.T) {
	suite.Run(t, &ServerImplementationTestSuite{
		service: NewPluginServiceImplementation(plugin.MockPlugin{}),
	})
}

func (suite *ServerImplementationTestSuite) TestLastUse() {
	startTime := suite.service.getLastUse()

	executeAction := &pb.ExecuteActionRequest{
		Timeout: 1,
	}

	go func() {
		if _, err := suite.service.ExecuteAction(context.Background(), executeAction); err != nil {
			suite.Require().Nil(err)
		}
	}()

	time.Sleep(time.Millisecond * 50)
	suite.Require().True(suite.service.lastUse != suite.service.getLastUse())
	suite.Require().True(suite.service.getLastUse() > startTime)

	for suite.service.activeWorkers > 0 {
		time.Sleep(time.Millisecond * 50)
	}

	time.Sleep(time.Millisecond * 50)
	suite.Require().True(suite.service.lastUse == suite.service.getLastUse())

	suite.Require().True(suite.service.lastUse >= startTime)
	suite.Require().True(suite.service.lastUse < (startTime + (time.Second * 2).Nanoseconds()))
}
