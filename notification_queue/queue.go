package notification_queue

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"notifier/notification_queue/rabbitmq_wrapper"
	"time"
)

var (
	exchange = "police"

	activeQueueName = "police_active"
	activeQueueKey  = "active"

	failedQueueName = "police_failed"
	failedQueueKey  = "failed"

	padding30SecQueueName = "police_padding_30sec"
	padding60SecQueueName = "police_padding_60sec"
	padding30SecQueueKey  = "padding_30sec"
	padding60SecQueueKey  = "padding_60sec"
)

type NotificationQueue struct {
	uri      string
	httpClient *http.Client
	rabbitmq *rabbitmq_wrapper.RabbitmqWrapper
	registerNotificationChannel chan *QueueItem
}

func NewNotificationQueue(uri string) *NotificationQueue {
	return &NotificationQueue{
		uri:        uri,
		rabbitmq:   rabbitmq_wrapper.NewRabbitmqWrapper(uri),
		httpClient: &http.Client{Timeout: 2 * time.Second},
		registerNotificationChannel: make(chan *QueueItem, 1024),
	}
}

func (this *NotificationQueue) Run() {
	// 注册Rabbitmq连接成功时间
	this.rabbitmq.RegisterOnConnectedCallback(this.onRabbitmqConnected)
	this.rabbitmq.Run()
	time.Sleep(2 * time.Second)

	go this.registerNotificationLoop()

	// 回调信息接受协程
	go func() {
		messageChannel := make(chan amqp.Delivery, 16)
		for i := uint16(0); i < 128; i++ {
			go this.messageHandle(messageChannel)
		}

		for {
			// 防止Rabbitmq无法链接时候 尝试速度过快占满CPU
			time.Sleep(time.Second)

			ch, err := this.rabbitmq.GetChannel()
			if err != nil {
				log.Printf("Failed to get channel from rabbitmq. %s", err)
				continue
			}

			msgCh, err := ch.Consume(
				activeQueueName, // queue
				"",              // consumer
				false,           // auto-ack
				false,           // exclusive
				false,           // no-local
				false,           // no-wait
				nil,             // args
			)
			if err != nil {
				log.Printf("Failed to get consume from channel. %s", err)
				continue
			}

			for message := range msgCh {
				log.Printf("Got %d from msgCh", message.DeliveryTag)
				messageChannel <- message
			}
		}
	}()
}

