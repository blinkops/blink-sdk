package plugin_sdk

import (
	"github.com/blinkops/plugin-sdk/plugin"
	"github.com/blinkops/plugin-sdk/plugin/config"
	"github.com/blinkops/plugin-sdk/plugin/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"

	pb "github.com/blinkops/plugin-sdk/plugin/proto"
)

const (
	ListenMode = "tcp"
)

func registerNetworkListener() (*net.Listener, error) {
	listenConfiguration := ":" + config.GetConfig().Server.Port
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

	networkListener, err := registerNetworkListener()
	if err != nil {
		return err
	}

	log.Infoln("Server is starting to serve requests...")
	return grpcServer.Serve(*networkListener)
}
