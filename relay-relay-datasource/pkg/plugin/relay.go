package plugin

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"sync"
	// "time"

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

// type Options struct {
// 	URI           string `json:"url"`
// 	Username      string `json:"username"`
// 	Password      string `json:"password"`
// }

type TopicMap struct {
	sync.Map
}

func InitNewClient(o *Datasource) (*nats.Conn, error){
	var natsClient, err = nats.Connect("nats://host.docker.internal:4222", nats.UserInfo(o.Username, o.Password))

	if err != nil{
		return nil, err;
	}
	
	log.DefaultLogger.Info("Connected to Relay!")

	return natsClient, nil
}