package main

// #include <stdlib.h>

import (
	"C"
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"

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
		ServerUrls:      output.FLBPluginConfigKey(ctx, "ServerUrls"),
		EnableStreaming: output.FLBPluginConfigKey(ctx, "EnableStreaming"),
		ClusterId:       output.FLBPluginConfigKey(ctx, "ClusterId"),
		ChannelId:       output.FLBPluginConfigKey(ctx, "ChannelId"),
		Debug:           output.FLBPluginConfigKey(ctx, "Debug"),
	}

	if pluginConfig.ServerUrls == "" {
		pluginConfig.ServerUrls = "nats://localhost:4222"
	}

	if pluginConfig.EnableStreaming == "" {
		pluginConfig.EnableStreaming = "off"
	}

	if pluginConfig.Debug == "" {
		pluginConfig.Debug = "off"
	}

	plugin, err := natspublisher.NewPlugin(pluginConfig)
	if err != nil {
		fmt.Println(err)
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

	nc := plugin.Conn
	sc := plugin.StanConn

	for {

		// Receive records from fluent-bit filter plugin
		ret, _, record := output.GetRecord(decoder)
		if ret != 0 {
			break
		}
		plugin.Logger.Debug(fmt.Sprintf("record = %v", record))

		// Determine NATS topic
		var subject string = string(record["topic"].([]uint8))
		subject = strings.ReplaceAll(subject, "/", ".")
		delete(record, "topic")
		plugin.Logger.Debug(fmt.Sprintf("subject = %v", subject))

		// Create NATS payload
		payload := make(map[string]interface{})
		for k, v := range record {
			payload[k.(string)] = v
		}
		payload["subject"] = subject
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			plugin.Logger.Error(err)
			return output.FLB_ERROR
		}
		plugin.Logger.Debug(fmt.Sprintf("jsonPayload = %v", string(jsonPayload)))

		// Publish NATS message
		if plugin.Config.EnableStreaming == "on" {
			plugin.Logger.Debug(fmt.Sprintf("channelId = %v", plugin.Config.ChannelId))
			plugin.Logger.Debug(fmt.Sprintf("jsonPayload = %v", string(jsonPayload)))
			_, err = sc.PublishAsync(plugin.Config.ChannelId, jsonPayload, func(ackedNuid string, err error) {
				if err != nil {
					plugin.Logger.Error(err)
				}
				plugin.Logger.Debug(ackedNuid)
			})
			if err != nil {
				plugin.Logger.Error(err)
				return output.FLB_ERROR
			}
		} else {
			if err = nc.Publish(subject, jsonPayload); err != nil {
				plugin.Logger.Error(err)
				return output.FLB_ERROR
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
