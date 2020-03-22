
package natspublisher

import (
	"github.com/sirupsen/logrus"
	nats "github.com/nats-io/nats.go"
)

type PluginConfig struct {
	ServerUrls string
	Topic string
}

type NatsPublisher struct {
	Config *PluginConfig
	Logger *logrus.Logger
	Conn *nats.Conn
}

func NewPlugin(config *PluginConfig, logger *logrus.Logger) (*NatsPublisher, error) {

	conn, err := nats.Connect(config.ServerUrls)
	if(err != nil) {
		return nil, err
	}
	return &NatsPublisher{
		Config: config,
		Logger: logger,
		Conn: conn,
	}, nil
}