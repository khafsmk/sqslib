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
	client := mqueue.New(
		mqueue.NewEventBridgeHandler("bus-name", mqueue.FSAWSConfig),
		mqueue.WithSquadName(squadName),
		mqueue.WithServiceName(serviceName),
		mqueue.WithDomain(domain),
	)

	// ...

	input := map[string]string{"key": "value"}
	err := client.Publish(context.Background(), mqueue.EventLoanCreate, input)
	if err != nil {
		panic(err)
	}
}
