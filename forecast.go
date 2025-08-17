package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetForecastParams struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude of the location"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude of the location"`
}

func GetForecast(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[GetForecastParams]) (*mcp.CallToolResultFor[any], error) {
	lat := req.Arguments.Latitude
	lon := req.Arguments.Longitude
	pointsURL := fmt.Sprintf("%s/forecast/point/%f,%f", newAPIBase, lat, lon)

	points, err := getJSON(ctx, pointsURL)
	if err != nil {
		log.Printf("points error: %v", err)
		return &mcp.CallToolResultFor[any]{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch forecast points"}}}, nil
	}
	props, _ := points["properties"].(map[string]any)
	fcURL, _ := props["forecast"].(string)
	if fcURL == "" {
		log.Printf("forecast URL not found")
		return &mcp.CallToolResultFor[any]{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "Forecast URL not found"}}}, nil
	}
	fc, err := getJSON(ctx, fcURL)
	if err != nil {
		log.Printf("forecast error: %v", err)
		return &mcp.CallToolResultFor[any]{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "Unable to fetch forecast"}}}, nil
	}
	fcProps, _ := fc["properties"].(map[string]any)
	periods, _ := fcProps["periods"].([]any)
	if len(periods) == 0 {
		return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: "No forecast data available"}}}, nil
	}

	// 返却件数を暫定で5件に制限する
	max := 5
	if len(periods) < max {
		max = len(periods)
	}
	out := ""
	for i := 0; i < max; i++ {
		p, _ := periods[i].(map[string]any)
		name := fmt.Sprintf("%v", p["name"])
		temp := fmt.Sprintf("%v", p["temperature"])
		unit := fmt.Sprintf("%v", p["temperatureUnit"])
		wind := fmt.Sprintf("%v %v", p["windSpeed"], p["windDirection"])
		df := fmt.Sprintf("%v", p["detailedForecast"])
		if i > 0 {
			out += "\n---\n"
		}
		out += fmt.Sprintf("%s: Temperature: %s°%s Wind: %s Forecast: %s", name, temp, unit, wind, df)
	}
	return &mcp.CallToolResultFor[any]{Content: []mcp.Content{&mcp.TextContent{Text: out}}}, nil
}

func registerForecast(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get weather forecast for a location",
	}, GetForecast)
}
