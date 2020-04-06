package main

// #include <stdlib.h>

import (
	"C"
	"encoding/json"
	"os"
	"strings"
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

		var subject string = string(record["topic"].([]uint8))
		subject = strings.ReplaceAll(subject, "/", ".")

		jsonPayload, err := json.Marshal(record["payload"])
		if err != nil {
			logger.Error(err)
		} else {
			msg := &nats.Msg{
				Subject: subject,
				Data:    jsonPayload,
			}
			if err = plugin.Conn.PublishMsg(msg); err != nil {
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
