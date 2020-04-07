
package natspublisher

import (
	"os"
	"strings"
	
	"github.com/sirupsen/logrus"
	nats "github.com/nats-io/nats.go"
)

type PluginConfig struct {
	ServerUrls string
	Debug string
}

type NatsPublisher struct {
	Config *PluginConfig
	Logger *logrus.Logger
	Conn *nats.Conn
}

func NewPlugin(config *PluginConfig) (*NatsPublisher, error) {

	conn, err := nats.Connect(config.ServerUrls)
	if(err != nil) {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	if strings.ToLower(config.Debug) == "on" || strings.ToLower(config.Debug) == "true"  {
		logger.SetLevel(logrus.DebugLevel)
	}

	return &NatsPublisher{
		Config: config,
		Logger: logger,
		Conn: conn,
	}, nil
}