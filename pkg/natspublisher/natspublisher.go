
package natspublisher

import (
	"os"
	"strings"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
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

	opts := nats.GetDefaultOptions()
	opts.Url = config.ServerUrls
	opts.AllowReconnect = true
	opts.ReconnectWait = time.Duration(1*time.Second)
	opts.MaxReconnect=-1
	conn, err := opts.Connect()
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