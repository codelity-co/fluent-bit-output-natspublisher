
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
  StanConn stan.Conn
}

func NewPlugin(config *PluginConfig) (*NatsPublisher, error) {

	natsPublisher := NatsPublisher{}

	opts := nats.GetDefaultOptions()
	opts.Url = config.ServerUrls
	opts.AllowReconnect = true
	opts.ReconnectWait = time.Duration(1*time.Second)
	opts.MaxReconnect=-1

	nc, err := opts.Connect()
	if(err != nil) {
		return nil, err
	}
	natsPublisher.Conn = nc

	if len(config.EnableStreaming) > 0 {
		config.EnableStreaming = strings.ToLower(config.EnableStreaming)
		if !contains([]string{"on", "off", "true", "false"},config.EnableStreaming) {
			return nil, fmt.Errorf("EnableStreaming setting must be on/off")
		}
		switch config.EnableStreaming {
		case "true":
			config.EnableStreaming = "on"
		case "fase":
			config.EnableStreaming = "off"
		}
	}

	natsPublisher.Config = config
	
	if config.EnableStreaming == "on" {
		if len(config.ClusterId) == 0 {
			return nil, fmt.Errorf("ClusterId cannot be empty.")
		}
		
		if len(config.ChannelId) == 0 {
			return nil, fmt.Errorf("ClusterId cannot be empty.")
		}

		var hostname string
		hostname, _ = os.Hostname()
		natsPublisher.StanConn, err = stan.Connect(config.ClusterId, hostname, stan.NatsConn(nc))
		if err != nil {
			return nil, err
		}
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	if strings.ToLower(config.Debug) == "on" || strings.ToLower(config.Debug) == "true"  {
		logger.SetLevel(logrus.DebugLevel)
	}
	natsPublisher.Logger = logger

	return &natsPublisher, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}