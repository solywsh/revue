# revue

> 基于go-cqhttp，在使用前请先配置并运行[go-cqhttp](https://github.com/Mrs4s/go-cqhttp)。本项目为go-cqhttp反向端。

## 开始

先配置并运行`go-cqhttp`，其中默认开启中间件密钥以及反向post密钥，注意配置以下内容：

> 注意：这里给的示例由于go-cqhttp端和revue端都部署在同一台服务器

`[go-cqhttp]config.yaml`

```yaml
# 默认中间件锚点
default-middlewares: &default
  # 访问密钥, 强烈推荐在公网的服务器设置
  access-token: '<secret>'
  
... 

servers:
  - http: # HTTP 通信设置
      host: 127.0.0.1 # 服务端监听地址
      port: 5700      # 服务端监听端口
      timeout: 5      # 反向 HTTP 超时时间, 单位秒，<5 时将被忽略
      long-polling:   # 长轮询拓展
        enabled: false       # 是否开启
        max-queue-size: 2000 # 消息队列大小，0 表示不限制队列大小，谨慎使用
      middlewares:
        <<: *default # 引用默认中间件
      post:           # 反向HTTP POST地址列表
        - url: 'http://127.0.0.1:5000'  # 地址
          secret: '<secret>'            	# 密钥
```

接下来配置`[revue]config.yaml`

```yaml
# 监听端口
listenPort: '5000'
# api的url,也就是部署go-cqhttp的url
urlHeader: 'http://127.0.0.1:5700'
# qq机器人的qq号
selfId: 'xxxxxxxx'
# 设置管理员
adminUser:
  - 'xxxxxxxxx'
# 管理员命令头,默认"$",则命令类似于"$start"
# 用于区分是否为命令(虽然具体的命令没有实现就是了)
adminUserOrderHeader: '$'
# 监听的qq群
listenGroup:
  - 'xxxxxxxxx'
# 正向鉴权,即向go-cqhttp客户端发送消息时进行鉴权
forwardAuthentication:
  enable: true # true\false
  token: '<secret>' 
# 反向鉴权,即接收go-cqhttp客户端消息时进行鉴权
reverseAuthentication:
  enable: true # true\false
  secret: '<secret>'
# revue机器人接口相关
revue:
  enable: true # true\false
# 数据库相关
Database:
  # 数据库路径
  path: './data.db'
```

