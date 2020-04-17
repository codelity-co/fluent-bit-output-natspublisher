
package natspublisher

import (
	"fmt"
	"os"
	"strings"
	"time"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type PluginConfig struct {
	ServerUrls string
	EnableStreaming string
	ClusterId string
	ChannelId string
	Debug string
}

type NatsPublisher struct {
	Config *PluginConfig
	Logger *logrus.Logger
	Conn *nats.Conn
	StanConn *stan.Conn
}

func NewPlugin(config *PluginConfig) (*NatsPublisher, error) {

	opts := nats.GetDefaultOptions()
	opts.Url = config.ServerUrls
	opts.AllowReconnect = true
	opts.ReconnectWait = time.Duration(1*time.Second)
	opts.MaxReconnect=-1

	nc, err := opts.Connect()
	if(err != nil) {
		return nil, err
	}

	var sc stan.Conn
	if strings.ToLower(config.EnableStreaming) == "on" || strings.ToLower(config.EnableStreaming) == "true"  {
		if len(config.ClusterId) == 0 {
			return nil, fmt.Errorf("ClusterId cannot be empty.")
		}
		
		if len(config.ChannelId) == 0 {
			return nil, fmt.Errorf("ClusterId cannot be empty.")
		}

		sc, err = stan.Connect(config.ClusterId, config.ChannelId, stan.NatsConn(nc))
		if err != nil {
			return nil, fmt.Errorf("Cannot connect NATS Streaming cluster %v with channel %v", config.ChannelId, config.ChannelId)
		}
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
		Conn: nc,
		StanConn: &sc,
	}, nil
}