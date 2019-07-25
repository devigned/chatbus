# GopherCon 2019 Azure Service Bus chat application (chatbus)

You are challenged to fill in the send and receive sides of a chat application using Azure Service Bus.
There are two functions to fill in before you are able to chat with your fellow Gophers. You can find both
functions within ./cmd/join.go.

```go
func sendMessage(ctx context.Context, topic *servicebus.Topic, name, message string) error {
	ctx, span := tab.StartSpan(ctx, "join.sendMessage")
	defer span.End()

	// TODO: Create a ChatMessage

	// TODO: marshal ChatMessage into bits

	// TODO: send message to topic
	return nil
}

func listenForAMessage(ctx context.Context, subscription *servicebus.Subscription) (*ChatMessage, error) {
	ctx, span := tab.StartSpan(ctx, "join.listenForAMessage")
	defer span.End()

	// TODO: receive one message from the subscription and unmarshal into ChatMessage

	// TODO return ChatMessage
	return nil, nil
}
```

In `sendMessage` you are required to send a message to an Azure Service Bus Topic. A Topic is a target for
broadcast messages to be listened to many subscribers.

In `listenForMessage` you are required to listen to a subscription on a topic to receive messages from the topic.

## Documentation
If you are up for the challenge, there is help to be had in https://godoc.org/github.com/Azure/azure-service-bus-go.

## Build and run
- `git clone github.com/devigned/chatbus`
- `go run  .`


