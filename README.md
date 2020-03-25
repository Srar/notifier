## Notifier

基于rabbitmq回调消息队列，将队列消息转换为HTTP回调并带有定时重试通知。

当重试次数小于等于5次时，会已每30秒重试一次。

当重试次数大于等于6次时小于等于10，会已每60秒重试一次。

超过10次重试失败后存入失败队列(TTL 12小时)。


## 注册新消息

POST到`http://notifier:1234/registerNotification`并附带以下数据

### Get

```json
{
    "method":"get",
    "url":"http://www.x-speed.cc/1231111",
    "args":{
        "a":"aaaaa",
        "b":"bbbbb"
    }
}
```

消息将会通知到`GET http://www.x-speed.cc/1231111?a=aaaaa&b=bbbbb`

### Json

````json
{
    "method":"json",
    "url":"http://www.x-speed.cc/1231111",
    "args":{
        "a":"aaaaa",
        "b":"bbbbb"
    }
}
````

消息将会通知到`POST http://www.x-speed.cc/1231111`并附带JSON数据:

```json
{
    "a":"aaaaa",
    "b":"bbbbb"
}
```

