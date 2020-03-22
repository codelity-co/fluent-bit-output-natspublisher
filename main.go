package main

// #include <stdlib.h>

import (
	"C"
	"encoding/json"
	"os"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/codelity-co/fluentbit-plugin-natspublisher/pkg/natspublisher"
)

var (
	plugins []*natspublisher.NatsPublisher
	logger  *logrus.Logger
)

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "natspublisher", "NATS publisher plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	pluginConfig := &natspublisher.PluginConfig{
		ServerUrls: output.FLBPluginConfigKey(ctx, "ServerUrls"),
		Topic:      output.FLBPluginConfigKey(ctx, "Topic"),
	}
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	plugin, err := natspublisher.NewPlugin(pluginConfig, logger)
	if err != nil {
		logger.Error(err)
		return output.FLB_ERROR
	}
	output.FLBPluginSetContext(ctx, plugin)
	plugins = append(plugins, plugin)
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, _ *C.char) int {
	plugin := output.FLBPluginGetContext(ctx).(*natspublisher.NatsPublisher)
	if plugin == nil {
		logger.Error("plugin not initialized")
		return output.FLB_ERROR
	}

	decoder := output.NewDecoder(data, int(length))

	for {

		ret, _, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}
		payload := map[string]interface{}{}
		for k, v := range record {
			switch v.(type) {
				case []uint8:
					payload[k.(string)] = string(v.([]uint8))
				default:
					payload[k.(string)] = v
			}
		}

		jsonString, err := json.Marshal(payload)
		if err != nil {
			logger.Error(err)
		}

		if len(plugin.Config.Topic) > 0 {
			msg := &nats.Msg{
				Subject: plugin.Config.Topic,
				Data:    jsonString,
			}
			if err := plugin.Conn.PublishMsg(msg); err != nil {
				logger.Error(err)
			}
		} 
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	for _, plugin := range plugins {
		if plugin.Conn != nil {
			if err := plugin.Conn.Drain(); err != nil {
				logger.Error(err)
			}
			plugin.Conn.Close()
		}
	}
	return output.FLB_OK
}

func main() {}
