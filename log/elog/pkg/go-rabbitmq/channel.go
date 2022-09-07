package gorabbitmq

import (
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type channelManager struct {
	url                 string
	channel             *amqp.Channel
	connection          *amqp.Connection
	amqpConfig          *amqp.Config
	channelMux          *sync.RWMutex
	notifyCancelOrClose chan error
	reconnectInterval   time.Duration
}

func newChannelManager(url string, config amqp.Config, reconnectInterval time.Duration) (*channelManager, error) {
	conn, ch, err := newChannel(url, config)
	if err != nil {
		return nil, err
	}

	chManager := channelManager{
		url:                 url,
		channel:             ch,
		connection:          conn,
		amqpConfig:          &config,
		channelMux:          &sync.RWMutex{},
		notifyCancelOrClose: make(chan error),
		reconnectInterval:   reconnectInterval,
	}

	go chManager.startNotifyCancleOrClosed()

	return &chManager, nil
}

func newChannel(url string, config amqp.Config) (*amqp.Connection, *amqp.Channel, error) {
	amqpConn, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, nil, err
	}

	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return amqpConn, ch, nil
}

func (c *channelManager) startNotifyCancleOrClosed() {
	notfiyCloseChan := c.channel.NotifyClose(make(chan *amqp.Error, 1))
	notfiyCancelChan := c.channel.NotifyCancel(make(chan string, 1))

	select {
	case err := <-notfiyCloseChan:
		if err != nil && err.Server {
			fmt.Print("attempting to reconnect to amqp server after close")
			c.reconnectLoop()
			fmt.Print("succcessfully reconnected to amqp server after close")
			c.notifyCancelOrClose <- err
		} else if err != nil && err.Reason == "EOF" {
			fmt.Print("attempting to reconnect to amqp server after eof")
			c.reconnectLoop()
			fmt.Print("succcessfully reconnect to amqp server after close")
			c.notifyCancelOrClose <- err
		} else if err != nil {
			fmt.Print("not attempting to reconnect to amqp server because closure was initiated by client")
		} else if err == nil {
			fmt.Print("amqp channel closed gracefully")
		}

		if err != nil {
			fmt.Printf("not attempting to reconnect to amqp server because closure was initiated by client")
		}

		if err == nil {
			fmt.Printf("amqp channel closed gracefully")
		}
	case err := <-notfiyCancelChan:
		fmt.Printf("attempting to reconnect to amqp server after cancel")
		c.reconnectLoop()
		fmt.Printf("succcessfully reconnect to amqp server after cancel")
		c.notifyCancelOrClose <- errors.New(err)
	}
}

func (c *channelManager) reconnectLoop() {
	for {
		fmt.Print("waiting %s seconds to attempt to reconnect to amqp server", c.reconnectInterval)
		time.Sleep(c.reconnectInterval)
		if err := c.reconnect(); err != nil {
			fmt.Printf("error reconnecting to amqp server: %s", err.Error())
		} else {
			return
		}

	}
}

// reconnect safely closes the current channel and get a new one
func (c *channelManager) reconnect() error {
	c.channelMux.Lock()
	defer c.channelMux.Unlock()

	newConn, newChannel, err := newChannel(c.url, *c.amqpConfig)
	if err != nil {
		return err
	}

	// channel -> connection
	c.channel.Close()
	c.connection.Close()

	c.connection = newConn
	c.channel = newChannel

	go c.startNotifyCancleOrClosed()
	return nil
}

func (c *channelManager) close() error {
	c.channelMux.Lock()
	defer c.channelMux.Unlock()

	if err := c.channel.Close(); err != nil {
		return err
	}

	if err := c.connection.Close(); err != nil {
		return err
	}

	return nil
}
