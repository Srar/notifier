package notification_queue

type QueueItem struct {
	NotificationId string `json:"notification_id"`

	RetryCount int `json:"retry_count"`

	NotifyMethod string `json:"notify_method"`
	NotifyUrl    string `json:"notify_url"`
	NotifyArgs   map[string]interface{} `json:"notify_args"`
	
	Metadata string `json:"metadata"`
	ResultCallback string `json:"result_callback"`
}
