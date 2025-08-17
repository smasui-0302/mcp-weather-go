package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAlertsParams struct {
	State string `json:"state" jsonschema:"Two-Letter US state code (e.g. CA, NY)"`
}

func formatAlert(feature map[string]any) string {
	props, _ := feature["properties"].(map[string]any)
	get := func(k string) string {
		if v, ok := props[k]; ok && v != nil {
			return fmt.Sprintf("%v", v)
		}
		return "unknown"
	}
	return fmt.Sprintf("Event: %s\nArea: %s\nSeverity: %s\nDescription: %s\nInstructions: %s",
		get("event"),
		get("areaDesc"),
		get("severity"),
		get("description"),
		get("instructions"))
}

func GetAlerts(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[GetAlertsParams]) (*mcp.CallToolResultFor[any], error) {
	state := req.Arguments.State
	if state == "" {
		return &mcp.CallToolResultFor[any]{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "State parameter is required"}}}, nil
	}
	url := fmt.Sprintf("%s/alerts/active/area/%s", newAPIBase, state)
	data, err := getJSON(ctx, url)
	if err != nil {
		log.Printf("get_alerts error: %v", err)
		return &mcp.CallToolResultFor[any]{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch alerts"}}}, nil
	}

	features, _ := data["features"].([]any)
	if len(features) == 0 {
		return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: "No active alerts"}}}, nil
	}

	out := ""
	for i, f := range features {
		if m, ok := f.(map[string]any); ok {
			if i > 0 {
				out += "\n---\n"
			}
			out += formatAlert(m)
		}
	}
	return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: out}}}, nil
}

func registerAlerts(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_alerts",
		Description: "Get active weather alerts for a specific US state",
	}, GetAlerts)
}
