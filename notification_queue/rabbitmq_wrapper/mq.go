package rabbitmq_wrapper

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var ChannelNotAvailableErr = errors.New("The channel isn't available yet.")

type RabbitmqWrapper struct {
	url             string
	rabbitmq        *amqp.Connection
	rabbitmqChannel *amqp.Channel

	onConnected    func(ch *amqp.Channel)
}

func NewRabbitmqWrapper(url string) *RabbitmqWrapper {
	return &RabbitmqWrapper{
		url: url,
	}
}

func (this *RabbitmqWrapper) Run()  {
	go this.reConnect()
}

func (this *RabbitmqWrapper) reConnect() {
	for {
		newConnection, err := amqp.Dial(this.url)
		if err != nil {
			log.Println("Failed to connect to rabbitmq. ", err)
			time.Sleep(1 * time.Second)
			continue
		}

		newChannel, err := newConnection.Channel()
		if err != nil {
			newConnection.Close()
			log.Println("Failed to get channel of rabbitmq. ", err)
			time.Sleep(1 * time.Second)
			continue
		}

		this.rabbitmq = newConnection
		this.rabbitmqChannel = newChannel

		// 如不使用队列channel会导致阻塞
		channelEvent := newChannel.NotifyClose(make(chan *amqp.Error, 1))
		connectionEvent := newConnection.NotifyClose(make(chan *amqp.Error, 1))

		if this.onConnected != nil {
			this.onConnected(this.rabbitmqChannel)
		}

		select {
		case err = <-channelEvent:
			break
		case err = <-connectionEvent:
			break
		}

		newChannel.Close()
		newConnection.Close()
		log.Println("Lost the connection of rabbitmq.", err)
	}
}

func (this *RabbitmqWrapper) GetChannel() (*amqp.Channel, error) {
	if this.rabbitmq == nil || this.rabbitmqChannel == nil || this.rabbitmq.IsClosed() {
		return nil, ChannelNotAvailableErr
	}
	return this.rabbitmqChannel, nil
}


func (this *RabbitmqWrapper) RegisterOnConnectedCallback(cb func(ch *amqp.Channel))  {
	this.onConnected = cb
}

