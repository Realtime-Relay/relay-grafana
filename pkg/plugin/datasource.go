package plugin

import (
	"context"
	"encoding/json"
	// "path"
	"time"
	"fmt"
	// "reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
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
func NewDatasource(ctx context.Context, s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Error("NewDatasource()")

	settings, secureSettings, err := getDatasourceSettings(s)
	if err != nil {
		return nil, err 
	}

	return &Datasource{
		Path: settings.Path,
		ApiKey: secureSettings.ApiKey,
		SecretKey: secureSettings.SecretKey,
	}, nil
}

type Options struct {
	Path string `json:"path"`
}

type DecryptedSecureJSONData struct {
	SecretKey string `json:"secretKey"`
	ApiKey string `json:"apiKey"`
}

func getDatasourceSettings(s backend.DataSourceInstanceSettings) (*Options, *DecryptedSecureJSONData, error) {
	settings := &Options{}
	secureSettings := &DecryptedSecureJSONData{}

	if err := json.Unmarshal(s.JSONData, settings); err != nil {
		logObject("CUSTOM_DEBUG", err)
		return nil, nil, err
	}

	// Convert DecryptedSecureJSONData map to JSON string
	secureJSONBytes, err := json.Marshal(s.DecryptedSecureJSONData)
	if err != nil {
		logObject("CUSTOM_DEBUG", err)
		return nil, nil, err
	}

	// Unmarshal the JSON string into secureSettings struct
	if err := json.Unmarshal(secureJSONBytes, secureSettings); err != nil {
		logObject("CUSTOM_DEBUG", err)
		return nil, nil, err
	}

	return settings, secureSettings, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	Path string
	ApiKey string
	SecretKey string
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
	log.DefaultLogger.Error("DIPOSE()")
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Connection successfully established"

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
	_, _, err := InitNewClient(d)
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

type ReqData struct {
	Topic string `json:"topic"`
	StartTime string `json:"start_time"`
	Path int64 `json:"path"`
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	// for simplicity on any error the function returns and ends the streaming
	natsClient, namespace, err := InitNewClient(d)

	if err != nil {
		logObject("RELAY_DEBUG", err)
		return err
	}
	
	js, _ := jetstream.New(natsClient)

	log.DefaultLogger.Info("Connected to Relay!")
	logObject("RELAY_DEBUG_REQ", req)

	reqData := ReqData{}

	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		logObject("RELAY_DEBUG_JSON_ERR", reqData)
		return err
	}

	logObject("RELAY_DEBUG_REQ_DATA", reqData)

	var streamName = fmt.Sprintf("%s_stream", namespace)
	var topic = fmt.Sprintf("%s_%s", streamName, reqData.Topic)

	newStream, sErr := js.CreateStream(ctx, jetstream.StreamConfig{
		Name: streamName,
		Subjects: []string{topic},
	})

	logObject("RELAY_DEBUG_JS_ERR", sErr)
	logObject("RELAY_DEBUG_JS_STREAM", newStream)

	log.DefaultLogger.Info("TIMESTAMP:", reqData.StartTime)

	parsedTime, err := time.Parse(time.RFC3339, reqData.StartTime)
	if err != nil {
		log.DefaultLogger.Info("TIMESTAMP ERR:", err)
		return err
	}

	log.DefaultLogger.Info("CONSUMER NAME", reqData.Path)

	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name: fmt.Sprintf("%s", reqData.Path), // Create unique consumer for particular consumer
		FilterSubject: topic,
		DeliverPolicy: jetstream.DeliverByStartTimePolicy,
		AckPolicy: jetstream.AckExplicitPolicy,
		ReplayPolicy: jetstream.ReplayInstantPolicy,
		OptStartTime: &parsedTime,
	})

	if err != nil {
		log.DefaultLogger.Info("CONSUMER ERR", err)
	}

	var ticker = time.NewTicker(100)

	log.DefaultLogger.Info("TICKER INIT")

	for {
		select {
			case <-ctx.Done():
				log.DefaultLogger.Info("CONTEXT DONE")
				return ctx.Err()
			case <-ticker.C:
				log.DefaultLogger.Info("TICKET START")
				iter, err := consumer.Messages(jetstream.PullMaxMessages(1))
				log.DefaultLogger.Info("MESSAGES ERR", err)

				if err != nil {
					log.DefaultLogger.Info("Error retrieving message pull")
					continue
				}

				msg, err := iter.Next()
				log.DefaultLogger.Info("ITER ERR", err)
				if err != nil {
					log.DefaultLogger.Info("Error retrieving message")
					continue
				}else{
					msg.Ack()

					log.DefaultLogger.Info(string(msg.Data()))
			
					var jsonMap map[string]interface{}
					json.Unmarshal(msg.Data(), &jsonMap)
			
					messageMap := jsonMap["message"].(map[string]interface{})
					rawMsg, _ := json.Marshal(messageMap)
			
					var message json.RawMessage = rawMsg
			
					err := sender.SendFrame(
						data.NewFrame(
							"response",
							data.NewField("value", nil, []json.RawMessage{message}),
						),
						data.IncludeAll,
					)
			
					if err != nil {
						Logger.Error("Failed send frame", "error", err)
					}
				}
		}
	}

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