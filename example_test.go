package mqueue_test

import (
	"context"

	"github.com/khafsmk/mqueue"
)

func ExampleClient() {
	var (
		squadName   = "squad"
		serviceName = "service"
		domain      = "domain"
	)
	client := mqueue.New(domain, squadName, serviceName, mqueue.NewEventBridgeHandler("bus-name", mqueue.FSAWSConfig))

	// ...

	input := map[string]string{"key": "value"}
	err := client.Publish(context.Background(), input)
	if err != nil {
		panic(err)
	}
}
