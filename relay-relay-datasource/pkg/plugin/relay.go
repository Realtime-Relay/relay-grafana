package plugin

import (
	// "context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"os"
	"path/filepath"

	nats "github.com/nats-io/nats.go"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Client interface {
	// GetTopic(string) (*Topic, bool)
	IsConnected() bool
	// Subscribe(string) *Topic
	Unsubscribe(string)
	Dispose()
}

type TopicMap struct {
	sync.Map
}

type APIKey struct {
	APIKey string `json:"api_key"` // JSON tag to map the field name to the JSON key
}

type Namespace struct {
	Status string `json:"status"`
	Data NamespaceData `json:"data"`
}

type NamespaceData struct {
	Namespace string `json:"namespace"`
}

func InitNewClient(o *Datasource) (*nats.Conn, string, error){
	var endpoint = fmt.Sprintf("nats://%s:4222", o.Path);
	log.DefaultLogger.Info(endpoint)

	var natsClient, err = nats.Connect(endpoint, nats.UserJWTAndSeed(o.ApiKey, o.SecretKey))

	if err != nil{
		log.DefaultLogger.Info("Unable to connected to Relay: ", err)
		return nil, "", err;
	}
	
	log.DefaultLogger.Info("Connected to Relay!")
	log.DefaultLogger.Info("Getting stream...")

	apiKeyData := APIKey{
		APIKey: o.ApiKey,
	}

	jsonBytes, err := json.Marshal(apiKeyData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, "", err
	}

	msg, err := natsClient.Request("accounts.user.get_namespace", jsonBytes, 5*time.Second)
	var namespace = &Namespace{}

	if err := json.Unmarshal([]byte(string(msg.Data)), &namespace); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, "", err
	}

	log.DefaultLogger.Info(namespace.Data.Namespace)

	return natsClient, namespace.Data.Namespace, nil
}