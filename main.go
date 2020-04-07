package main

// #include <stdlib.h>

import (
	"C"
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	nats "github.com/nats-io/nats.go"

	"github.com/codelity-co/fluentbit-plugin-natspublisher/pkg/natspublisher"
)

var (
	plugins []*natspublisher.NatsPublisher
)

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "natspublisher", "NATS publisher plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	pluginConfig := &natspublisher.PluginConfig{
		ServerUrls: output.FLBPluginConfigKey(ctx, "ServerUrls"),
		Debug: output.FLBPluginConfigKey(ctx, "Debug"),
	}

	if pluginConfig.ServerUrls == "" {
		pluginConfig.ServerUrls = "nats://localhost:4222"
	}

	if pluginConfig.Debug == "" {
		pluginConfig.Debug = "off"
	}

	plugin, err := natspublisher.NewPlugin(pluginConfig)
	if err != nil {
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
		plugin.Logger.Error("plugin not initialized")
		return output.FLB_ERROR
	}

	decoder := output.NewDecoder(data, int(length))

	for {

		ret, _, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}
		plugin.Logger.Debug(fmt.Sprintf("record = %v", record))

		var subject string = string(record["topic"].([]uint8))
		subject = strings.ReplaceAll(subject, "/", ".")
		delete(record, "topic")

		payload := make(map[string]interface{})
		for k, v := range record {
			payload[k.(string)] = v
		}
		payload["subject"] = subject
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			plugin.Logger.Error(err)
		} else {
			msg := &nats.Msg{
				Subject: subject,
				Data:    jsonPayload,
			}
			plugin.Logger.Debug(fmt.Sprintf("msg = %v", msg))
			if err = plugin.Conn.PublishMsg(msg); err != nil {
				plugin.Logger.Error(err)
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
				plugin.Logger.Error(err)
			}
			plugin.Conn.Close()
		}
	}
	return output.FLB_OK
}

func main() {}
