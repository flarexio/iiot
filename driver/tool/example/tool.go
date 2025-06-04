package example

import (
	"context"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"

	"github.com/flarexio/iiot/machine"
)

type Tool interface {
	Schema(ctx context.Context) ([]byte, error)
	Instruction(ctx context.Context) (string, error)
	ReadPoints(ctx context.Context, req *ReadPointsRequest) ([]any, error)
}

type Point struct {
	Name  string
	Value any
}

type ReadPointsRequest struct {
	Points []*Point
}

func NewTool() Tool {
	m := minify.New()
	m.AddFunc("application/json", json.Minify)
	return &tool{m}
}

type tool struct {
	m *minify.M
}

func (t *tool) Schema(ctx context.Context) ([]byte, error) {
	return t.m.Bytes("application/json", schema)
}

func (t *tool) Instruction(ctx context.Context) (string, error) {
	instruction := `This tool reads points from a controller named "TEMP".
	
	To use this tool, provide a list of points with their names and values.
	Example:
	{
		"points": [
			{
				"name": "temperature",
				"value": 22.5
			},
			{
				"name": "humidity",
				"value": 45
			}
		]
	}`

	return instruction, nil
}

func (t *tool) ReadPoints(ctx context.Context, req *ReadPointsRequest) ([]any, error) {
	points := make([]*machine.Point, len(req.Points))
	for i, point := range req.Points {
		points[i] = &machine.Point{
			Name: point.Name,
			Options: map[string]any{
				"value": point.Value,
			},
		}
	}

	controller := &machine.Controller{
		ControllerID: "TEMP",
		Points:       points,
	}

	svc := NewService()
	svc.AddControllers(controller)

	pointNames := make([]string, len(req.Points))
	for i, point := range req.Points {
		pointNames[i] = point.Name
	}

	return svc.ReadPoints(ctx, "TEMP", pointNames)
}

var schema = []byte(`{
	"$schema": "http://json-schema.org/2020-12/schema",
	"title": "Example Tool Schema",
	"type": "object",
	"properties": {
		"points": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"description": "The name of the point"
					},
					"value": {
						"type": ["string", "number", "boolean"],
						"description": "The value of the point, can be string, number or boolean"
					}
				},
				"required": ["name", "value"],
				"additionalProperties": false
			},
			"description": "List of points to read"
		}
	}
}`)
