package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	servicebus "github.com/Azure/azure-service-bus-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	joinCmd.Flags().StringVarP(&joinParams.name, "name", "n", "", "the name you would like to be addressed by")
	joinCmd.Flags().StringVarP(&joinParams.topic, "topic", "t", "general", "topic of the conversation")
	rootCmd.AddCommand(joinCmd)
}

type (
	JoinParams struct {
		name  string
		topic string
	}

	ChatMessage struct {
		Message string `json:"message,omitempty"`
		Name    string `json:"name,omitempty"`
	}
)

var (
	joinParams JoinParams
	joinCmd    = &cobra.Command{
		Use:   "join",
		Short: "Join the chatbus chat",
		Args: func(cmd *cobra.Command, args []string) error {
			if joinParams.name == "" {
				return errors.New("name parameter can't be blank")
			}
			return checkAuthFlags()
		},
		Run: RunWithCtx(func(ctx context.Context, cmd *cobra.Command, args []string) {
			topic, subscription, err := buildTopicAndSubscription(ctx, joinParams.topic, joinParams.name)
			if err != nil {
				log.Error(err)
				return
			}

			fmt.Println("type then press enter to send a message, " + joinParams.name)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				for {
					select {
					case <-ctx.Done():
					default:
						// listen to the topic
						msg, err := listenForAMessage(ctx, subscription)
						if err != nil {
							log.Error(err)
							cancel()
							return
						}

						if msg.Name != joinParams.name {
						// print messages not from me
						fmt.Println(fmt.Sprintf("%s: %s", msg.Name, msg.Message))
						}
					}
				}
			}()

			go func() {
				// send to the topic
				for {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						if err := sendMessage(ctx, topic, joinParams.name, scanner.Text()); err != nil {
							log.Error(err)
							cancel()
							return
						}
					}
				}
			}()

			<-ctx.Done()
		}),
	}
)

func sendMessage(ctx context.Context, topic *servicebus.Topic, name, message string) error {
	// TODO: Fill in with send functionality
	return nil
}

func listenForAMessage(ctx context.Context, subscription *servicebus.Subscription) (*ChatMessage, error) {
	// TODO: Fill in with receive functionality
	return nil, nil
}

func buildTopicAndSubscription(ctx context.Context, topic, sub string) (*servicebus.Topic, *servicebus.Subscription, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		return nil, nil, err
	}

	tm := ns.NewTopicManager()
	if _, err := ensureTopic(ctx, tm, topic); err != nil {
		return nil, nil, err
	}

	t, err := ns.NewTopic(topic)
	if err != nil {
		return nil, nil, err
	}

	sm, err := ns.NewSubscriptionManager(topic)
	if err != nil {
		return nil, nil, err
	}

	if _, err := ensureSubscription(ctx, sm, sub); err != nil {
		return nil, nil, err
	}

	s, err := t.NewSubscription(sub)
	if err != nil {
		_ = t.Close(ctx)
		return nil, nil, err
	}

	return t, s, nil
}

func ensureTopic(ctx context.Context, tm *servicebus.TopicManager, name string, opts ...servicebus.TopicManagementOption) (*servicebus.TopicEntity, error) {
	te, err := tm.Get(ctx, name)
	if err == nil {
		return te, err
	}

	te, err = tm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return te, nil
}

func ensureSubscription(ctx context.Context, sm *servicebus.SubscriptionManager, name string, opts ...servicebus.SubscriptionManagementOption) (*servicebus.SubscriptionEntity, error) {
	subEntity, err := sm.Get(ctx, name)
	if err == nil {
		return subEntity, err
	}

	subEntity, err = sm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return subEntity, nil
}
