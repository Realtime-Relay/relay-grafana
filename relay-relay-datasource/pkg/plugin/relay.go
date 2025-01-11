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
	var creds = getCreds(o.ApiKey, o.SecretKey)
	log.DefaultLogger.Info(creds)

	var endpoint = fmt.Sprintf("nats://%s:4222", o.Path);
	log.DefaultLogger.Info(endpoint)

	var natsClient, err = nats.Connect(endpoint, nats.UserCredentials(creds))

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

func getCreds(jwt string, secret string) string{
	var creds = fmt.Sprintf(`
-----BEGIN NATS USER JWT-----
%s
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
%s
------END USER NKEY SEED------

*************************************************************
	`, jwt, secret)

	file, err := os.Create("user.creds")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer file.Close() // Ensure the file is closed when the function ends

	// Write some data to the file
	_, err = file.WriteString(creds)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return ""
	}

	fmt.Println("File created successfully: user.creds")

	absPath, err := filepath.Abs("user.creds")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	return absPath
}