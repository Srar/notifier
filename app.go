package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sony/sonyflake"
	"log"
	"notifier/notification_queue"
	"notifier/utiles/conf"
)

var (
	uri string
	sf  *sonyflake.Sonyflake
)

func init() {
	mode := flag.String("mode", "dev", "Running mode.")
	flag.Parse()

	conf.LoadConfig(fmt.Sprintf("config_%s.conf", *mode))
	uri = fmt.Sprintf("amqp://%s:%s@%s/", *conf.ReadString("rabbitmq", "user"), *conf.ReadString("rabbitmq", "pass"), *conf.ReadString("rabbitmq", "host"))

	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	notificationQueue := notification_queue.NewNotificationQueue(uri)
	notificationQueue.Run()

	r := gin.Default()
	r.POST("/registerNotification", func(c *gin.Context) {
		notificationId, err := sf.NextID()
		if err != nil {
			c.JSON(200, gin.H{"code": -1, "message": err})
			return
		}

		var request struct {
			Method string                 `json:"method"`
			Url    string                 `json:"url"`
			Args   map[string]interface{} `json:"args"`
		}
		if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
			c.JSON(200, gin.H{"code": -1, "message": err})
			return
		}

		notificationQueue.RegisterNotification(fmt.Sprintf("%d", notificationId), request.Method, request.Url, request.Args)
		c.JSON(200, gin.H{"code": 0, "message": ""})
	})
	r.Run(*conf.ReadString("http", "port"))
}
