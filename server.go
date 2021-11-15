package plugin_sdk

import (
	"github.com/blinkops/blink-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin/config"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
	"github.com/blinkops/blink-sdk/plugin/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strconv"
)

const (
	ListenMode = "tcp"
)

func registerNetworkListener() (*net.Listener, error) {
	listenConfiguration := ":" + config.GetServerPort()
	log.Infof("Starting grpc listener on port %s\n", listenConfiguration)
	listener, err := net.Listen(ListenMode, listenConfiguration)
	if err != nil {
		log.Error("Failed to register network listener: ", err)
		return nil, err
	}

	return &listener, nil
}

func Start(pluginImplementation plugin.Implementation) error {
	description := pluginImplementation.Describe()
	log.Infof("Starting %s server.\n", description.Name)
	grpcServer := grpc.NewServer()

	log.Infof("Registering %s service!\n", description.Name)
	pluginServiceImplementation := server.NewPluginServiceImplementation(pluginImplementation)
	pb.RegisterPluginServer(grpcServer, pluginServiceImplementation)

	if val, _ := strconv.ParseBool(os.Getenv("USE_REFLECTION")); val {
		// enable reflection only if the env variable is present.
		reflection.Register(grpcServer)
		log.Warn("Server is using reflection! [USE_REFLECTION=true]")
	}

	networkListener, err := registerNetworkListener()
	if err != nil {
		return err
	}

	log.Infoln("Server is starting to serve requests...")
	return grpcServer.Serve(*networkListener)
}
