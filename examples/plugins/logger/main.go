package main

import (
	"fmt"
	"time"
)

type LoggerPlugin struct{}

var Plugin LoggerPlugin

func (p LoggerPlugin) OnEnterRule(data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("[%s] Starting rule processing: %v\n", time.Now().Format(time.RFC3339), data)
	data["timestamp"] = time.Now().Unix()
	return data, nil
}

func (p LoggerPlugin) OnEnterLine(data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("[%s] Processing line: %v\n", time.Now().Format(time.RFC3339), data)
	data["line_start_time"] = time.Now().Unix()
	return data, nil
}

func (p LoggerPlugin) OnExitLine(data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("[%s] Finished processing line\n", time.Now().Format(time.RFC3339))
	if mapped, ok := data["mapped"].(map[string]interface{}); ok {
		if startTime, ok := data["line_start_time"].(int64); ok {
			mapped["processing_duration"] = time.Now().Unix() - startTime
			data["mapped"] = mapped
		}
	}
	return data, nil
}

func (p LoggerPlugin) OnExitRule(data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("[%s] Finished rule processing\n", time.Now().Format(time.RFC3339))
	if payloads, ok := data["payloads"].([]map[string]interface{}); ok {
		data["total_processed"] = len(payloads)
	}
	return data, nil
}
