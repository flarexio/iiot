package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/xeipuuv/gojsonschema"

	"github.com/flarexio/iiot/driver/tool/example"
	"github.com/flarexio/iiot/driver/tool/stdio"
)

type exampleTestSuite struct {
	suite.Suite
	ctx       context.Context
	cancel    context.CancelFunc
	serverIn  *bytes.Buffer
	serverOut *bytes.Buffer
}

func (suite *exampleTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.Background())
	suite.ctx = ctx
	suite.cancel = cancel

	tool := example.NewTool()

	in := new(bytes.Buffer)
	out := new(bytes.Buffer)

	suite.serverIn = in
	suite.serverOut = out

	server := stdio.NewStdioServer()
	server.AddHandler("driver.schema", SchemaHandler(tool))
	server.AddHandler("driver.readPoints", ReadPointsHandler(tool))
	server.SetIO(in, out)
	go server.Listen(ctx)
}

func (suite *exampleTestSuite) TestSchema() {
	req := json.RawMessage(`{
		"points": [
			{"name": "temperature", "value": 1200},
			{"name": "pressure", "value": 150},
			{"name": "humidity", "value": 75.5},
			{"name": "status", "value": "Running"}
		]
	}`)

	handler := suite.Handler()
	executor := stdio.NewTestableExecutor(handler)
	client := stdio.NewStdioClient(executor)

	schema, err := client.Schema(suite.ctx, "example")
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(req)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	suite.True(result.Valid())
}

func (suite *exampleTestSuite) TestReadPoints() {
	req := json.RawMessage(`{
		"points": [
			{"name": "temperature", "value": 1200},
			{"name": "pressure", "value": 150},
			{"name": "humidity", "value": 75.5},
			{"name": "status", "value": "Running"}
		]
	}`)

	handler := suite.Handler()
	executor := stdio.NewTestableExecutor(handler)
	client := stdio.NewStdioClient(executor)

	points, err := client.ReadPoints(suite.ctx, "example", req)
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	suite.Len(points, 4)
	suite.Equal(1200.0, points[0])
	suite.Equal(150.0, points[1])
	suite.Equal(75.5, points[2])
	suite.Equal("Running", points[3])
}

func (suite *exampleTestSuite) Handler() stdio.ExecuteHandler {
	return func(ctx context.Context, program string, input io.Reader, output io.Writer) error {
		if _, err := io.Copy(suite.serverIn, input); err != nil {
			return err
		}

		time.Sleep(1000 * time.Millisecond)

		if _, err := io.Copy(output, suite.serverOut); err != nil {
			return err
		}

		return nil
	}
}

func (suite *exampleTestSuite) TearDownTest() {
	suite.cancel()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(exampleTestSuite))
}
