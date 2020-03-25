package notification_queue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (this *NotificationQueue) messageHandle(messageChannel chan amqp.Delivery) {
	for message := range messageChannel {
		var notification QueueItem
		err := json.Unmarshal(message.Body, &notification)
		if err != nil {
			message.Ack(false)
			log.Printf("Failed to parse to object from json. Payload: %s Err: %s", string(message.Body), err)
			continue
		}

		if err, canSkip := this.notifyRemote(&notification); err != nil {
			log.Println(err)
			// 判断这个错误是否可以忽略不做重试操作
			if !canSkip {
				this.reRegisterNotification(&notification)
				message.Ack(false)
			}
			continue
		}

		message.Ack(false)
	}
}

// notifyRemote 根据QueueItem请求远端, 并根据错误返回是否可以忽略这个通知
func (this *NotificationQueue) notifyRemote(message *QueueItem) (error, bool) {
	method := "GET"
	if message.NotifyMethod != "GET" {
		method = "POST"
	}

	request, err := http.NewRequest(method, message.NotifyUrl, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("NewRequest failed. %v", err)), true
	}
	request.Close = true

	switch message.NotifyMethod {
	case "GET":
		query := request.URL.Query()
		for k, v := range message.NotifyArgs {
			query.Add(k, fmt.Sprintf("%v", v))
		}
		request.URL.RawQuery = query.Encode()
	case "JSON":
		request.Header.Set("Content-Type", "application/json")
		jsonValue, err := json.Marshal(message.NotifyArgs)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to parse to json. %s", err)), true
		}
		request.Body = ioutil.NopCloser(bytes.NewReader(jsonValue))
	}

	// 发送远端请求
	response, err := this.httpClient.Do(request)
	if err != nil {
		return err, false
	}
	defer response.Body.Close()

	// 读取远端返回是否正确
	responseBytes := make([]byte, 2)
	_, err = response.Body.Read(responseBytes)
	if err != nil {
		return err, false
	}

	if strings.ToLower(string(responseBytes)) != "ok" {
		return errors.New("Unexpected response from remote."), false
	}

	return nil, true
}
