package plugin

import (
	"context"
	"encoding/json"
	// "path"
	// "time"
	// "fmt"
	// "reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	// "github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/nats-io/nats.go/jetstream"
)

var Logger = log.DefaultLogger

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil) // Streaming data source needs to implement this
)

// NewDatasource creates a new datasource instance.
func NewDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := getDatasourceSettings(s)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		Path: settings.Path,
		Username: settings.Username,
		Password: settings.Password,
	}, nil
}

type Options struct {
	Path     string `json:"path"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func getDatasourceSettings(s backend.DataSourceInstanceSettings) (*Options, error) {
	settings := &Options{}

	logObject("CUSTOM_DEBUG", s)

	if err := json.Unmarshal(s.JSONData, settings); err != nil {
		logObject("CUSTOM_DEBUG", err)
		return nil, err
	}

	logObject("CUSTOM_DEBUG", settings)

	return settings, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	Path string
	Username string
	Password string
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Connection can be successfully established"

	if !d.canConnect() {
		status = backend.HealthStatusError
		message = "Connection not working"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (d *Datasource) canConnect() bool {
	_, err := InitNewClient(d)
	if err != nil {
		return false
	}

	return true
}

// SubscribeStream just returns an ok in this case, since we will always allow the user to successfully connect.
// Permissions verifications could be done here. Check backend.StreamHandler docs for more details.
func (d *Datasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
	}, nil
}

// PublishStream just returns permission denied in this case, since in this example we don't want the user to send stream data.
// Permissions verifications could be done here. Check backend.StreamHandler docs for more details.
func (d *Datasource) PublishStream(context.Context, *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	// for simplicity on any error the function returns and ends the streaming
	natsClient, err := InitNewClient(d)
	js, _ := jetstream.New(natsClient)

	if err != nil {
		logObject("RELAY_DEBUG", err)
		return err
	}

	logObject("RELAY_DEBUG_NC", natsClient)

	log.DefaultLogger.Info("Connected to Relay!")

	newStream, sErr := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "test-namespace-stream",
		Subjects: []string{"test-topic"},
	})

	logObject("RELAY_DEBUG_JS", newStream)
	logObject("RELAY_DEBUG_JS_ERR", sErr)

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	default:
			

	// 		msg := Message{}
	// 		rawMsg, err := ws.ReadMessage()

	// 		if err != nil {
	// 			return err
	// 		}

	// 		if err := json.Unmarshal(rawMsg, &msg); err != nil {
	// 			return err
	// 		}

	// 		err = sender.SendFrame(
	// 			data.NewFrame(
	// 				"response",
	// 				data.NewField("time", nil, []time.Time{time.UnixMilli(msg.Time)}),
	// 				data.NewField("value", nil, []float64{msg.Value})),
	// 			data.IncludeAll,
	// 		)

	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func logObject(key string, obj interface{}) {
	// Convert the object to JSON
	objJSON, err := json.Marshal(obj)
	if err != nil {
		log.DefaultLogger.Error("Failed to serialize object", "error", err)
		return
	}

	// Log the serialized object
	log.DefaultLogger.Info("Logging object", key, string(objJSON))
}