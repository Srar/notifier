package notification_queue

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"strings"
	"time"
)

func (this *NotificationQueue) RegisterNotification(notificationId, method, url string, args map[string]interface{}) error {
	method = strings.ToUpper(method)
	if method != "GET" && method != "JSON" {
		return errors.New("Unsupported method yet. \n")
	}

	if args == nil {
		args = make(map[string]interface{})
	}
	args["notificationId"] = notificationId

	this.registerNotificationChannel <- &QueueItem{
		NotificationId: notificationId,
		RetryCount:     0,
		NotifyMethod:   method,
		NotifyUrl:      url,
		NotifyArgs:     args,
		Metadata:       "",
		ResultCallback: "",
	}

	return nil
}

func (this *NotificationQueue) registerNotificationLoop() {
	for notification := range this.registerNotificationChannel {
		ch, err := this.rabbitmq.GetChannel()
		if err != nil {
			this.registerNotificationChannel <- notification
			time.Sleep(time.Second)
			continue
		}

		notificationBytes, _ := json.Marshal(notification)
		err = ch.Publish(exchange, activeQueueKey, false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        notificationBytes,
			})
		if err != nil {
			this.registerNotificationChannel <- notification
			time.Sleep(time.Second)
			continue
		}
	}
}

func (this *NotificationQueue) reRegisterNotification(item *QueueItem) {
	ch, err := this.rabbitmq.GetChannel()
	if err != nil {
		return
	}

	item.RetryCount++
	itemBytes, _ := json.Marshal(item)

	if item.RetryCount <= 5 {
		ch.Publish(exchange, padding30SecQueueKey, false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        itemBytes,
			})
		return
	}

	if item.RetryCount > 5 && item.RetryCount <= 10 {
		ch.Publish(exchange, padding60SecQueueKey, false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        itemBytes,
			})
		return
	}

	ch.Publish(exchange, failedQueueKey, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        itemBytes,
		})
}
