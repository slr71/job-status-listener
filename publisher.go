package main

import (
	"github.com/cyverse-de/messaging"
)

// JobUpdatePublisher is the interface for types that need to publish a job
// update.
type JobUpdatePublisher interface {
	PublishJobUpdate(m *messaging.UpdateMessage) error
	Reconnect() error
	Close()
}

// DefaultJobUpdatePublisher provides a wrapper around messaging.Client that adds support for
// reestablishing stale connections.
type DefaultJobUpdatePublisher struct {
	uri      string
	exchange string
	client   *messaging.Client
}

func newMessagingClient(uri, exchange string, reconnect bool) (*messaging.Client, error) {
	client, err := messaging.NewClient(uri, reconnect)
	if err != nil {
		return nil, err
	}

	err = client.SetupPublishing(exchange)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewDefaultJobUpdatePublisher returns a new instance of DefaultJobUpdatePublisher.
func NewDefaultJobUpdatePublisher(uri, exchange string) (*DefaultJobUpdatePublisher, error) {
	client, err := newMessagingClient(uri, exchange, true)
	if err != nil {
		return nil, err
	}

	publisher := &DefaultJobUpdatePublisher{
		uri:      uri,
		exchange: exchange,
		client:   client,
	}
	return publisher, nil
}

// PublishJobUpdate simply forwards the function call to messaging.Client.PublishJobUpdate.
func (c *DefaultJobUpdatePublisher) PublishJobUpdate(m *messaging.UpdateMessage) error {
	return c.client.PublishJobUpdate(m)
}

// Reconnect closes the existing messaging client connection and establishes a new one.
func (c *DefaultJobUpdatePublisher) Reconnect() error {
	c.client.Close()

	client, err := newMessagingClient(c.uri, c.exchange, false)
	if err != nil {
		return err
	}

	c.client = client
	return nil
}

// Close closes the messaging client connection.
func (c *DefaultJobUpdatePublisher) Close() {
	c.client.Close()
}
