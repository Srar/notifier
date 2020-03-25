package notification_queue

import "github.com/streadway/amqp"

func (this *NotificationQueue) onRabbitmqConnected(ch *amqp.Channel)  {
	// 注册Exchange
	ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	// 注册推送队列
	ch.QueueDeclare(
		activeQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	// 注册推送失败队列
	ch.QueueDeclare(
		failedQueueName,
		true,
		false,
		false,
		false,
		map[string]interface{}{
			// 如果推送失败队列消息未在12小时内被消费则视为放弃
			"x-message-ttl": int32(1000 * 3600 * 12),
		},
	)

	// 注册推送失败30秒等待队列
	ch.QueueDeclare(padding30SecQueueName, true, false, false, false,
		map[string]interface{}{
			"x-message-ttl":             1000 * 30,
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": activeQueueKey,
		},
	)

	// 注册推送失败60秒等待队列
	ch.QueueDeclare(padding60SecQueueName, true, false, false, false,
		map[string]interface{}{
			"x-message-ttl":             1000 * 60,
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": activeQueueKey,
		},
	)

	// 绑定队列到Exchange
	// 推送失败策略 1~5次 等待30秒, 6~10 等待60秒
	ch.QueueBind(activeQueueName, activeQueueKey, exchange, false, nil)
	ch.QueueBind(failedQueueName, failedQueueKey, exchange, false, nil)
	ch.QueueBind(padding30SecQueueName, padding30SecQueueKey, exchange, false, nil)
	ch.QueueBind(padding60SecQueueName, padding60SecQueueKey, exchange, false, nil)
}