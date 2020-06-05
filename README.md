# fpm-iot-go-middleware

使用 Golang 实现当前的 Nodejs 版本

https://github.com/team4yf/fpm-iot-cloud-middleware

主要内容:

- 实现一个HTTP(S)的服务 `/push/:device/:brand/:event` 用于接收设备平台发送的推送数据
- 实现 TCP 服务,监听 5001 端口,接受和发送 HEX 数据
- 接入MQTT服务端,用于PUB/SUB设备数据
- 接入Kafka服务器,用于PUB/SUB业务或者日志数据


核心目标:

1. 无论数据采用什么形式接入到服务端,都经过清洗和转化,推入到MQTT,给相关的应用进行订阅消费.

2. 从应用推送的控制指令,也通过相应的方式转换到对应的数据通信方式给到设备平台.


目前主要的接入目标:

 `/push/:device/:brand/:event`


 `$ curl -H "Content-Type: application/json" -XPOST -d '{"data":1}' localhost:9000/push/light/lb/beat`