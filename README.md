# fpm-iot-go-middleware

使用 Golang 实现当前的 Nodejs 版本

https://github.com/team4yf/fpm-iot-cloud-middleware

主要内容:

- 实现一个HTTP(S)的服务 `/push/:device/:brand/:event` 用于接收设备平台发送的推送数据
- 实现 TCP 服务,监听 5001 端口,接受和发送 HEX 数据
- 接入MQTT服务端,用于PUB/SUB设备数据
- 接入Kafka服务器,用于PUB/SUB业务或者日志数据
- 提供一个websocket服务，用于转发消息内容


核心目标:

1. 无论数据采用什么形式接入到服务端,都经过清洗和转化,推入到MQTT,给相关的应用进行订阅消费.

2. 从应用推送的控制指令,也通过相应的方式转换到对应的数据通信方式给到设备平台.


目前主要的接入目标:

 `/push/:device/:brand/:event`

 `$ curl -H "Content-Type: application/json" -XPOST -d '{"data":1}' localhost:9000/push/light/lb/beat`

 需要使用 jsonPath 来获取数据中的设备 ID

根据设备的ID来获取设备对应的 appid，用来区分不同的应用，该信息保存在 redis 中
key: `device:type:brand: {deviceId: appid,}`
 `^push/:appid/event`
 
JSON数据格式

| 名称 | 描述 | 数据类型 |
| --- | ----- | --- |
| origin | 源消息体 | Object |
| event | 设备事件 | String |
| aid | 设备对应的应用服务平台id | String |
| pid | 设备对应的在服务中的项目id | String |
| sn | 设备的编码 | String |
| type | 设备对应的类型 | String |
| brand | 设备对应的品牌 | String |
| bind | 设备绑定的静态数据 | Object |

### 代码目录说明

- build
  - docker: 用于打包Docker镜像需要用到的文件
- conf: 用于开发时运行的docker容器的配置文件
- internal
  - model: 数据结构
  - repository: 持久层交互用到
  - service: 底层服务，用于聚合一些数据交互
- pkg: 常用的工具类
- router: 路由信息
  - middleware: 中间件
- scripts: 常用的脚本