package plugin_sdk

import (
	"github.com/blinkops/plugin-sdk/plugin"
	"github.com/blinkops/plugin-sdk/plugin/config"
	"github.com/blinkops/plugin-sdk/plugin/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"

	pb "github.com/blinkops/plugin-sdk/plugin/proto"
)

const (
	ListenMode = "tcp"
)

func registerNetworkListener() (*net.Listener, error) {
	listenConfiguration := ":" + config.GetConfig().Server.Port
	listener, err := net.Listen(ListenMode, listenConfiguration)
	if err != nil {
		logrus.Error("Failed to register network listener: ", err)
		return nil, err
	}

	return &listener, nil
}

func Start(pluginImplementation plugin.Implementation) error {

	logrus.Infoln("Starting description server.")
	grpcServer := grpc.NewServer()

	logrus.Infoln("Registering description service!")
	pluginServiceImplementation := server.NewPluginServiceImplementation(pluginImplementation)
	pb.RegisterPluginServer(grpcServer, pluginServiceImplementation)

	networkListener, err := registerNetworkListener()
	if err != nil {
		return err
	}

	logrus.Infoln("Server is starting to serve requests...")
	return grpcServer.Serve(*networkListener)
}
