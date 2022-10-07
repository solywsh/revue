# revue

> 基于go-cqhttp，在使用前请先配置并运行[go-cqhttp](https://github.com/Mrs4s/go-cqhttp)。本项目为go-cqhttp反向端。

## 支持功能

- API：私信发送API
- 通用
  - 答案查询：`搜索答案{question}`
  - 音乐搜索（目前只支持163）：`音乐搜索{keywords}`
  - 程序员黄历：`程序员黄历`
  - 求运势：`求签`
  - 涩图服务：
    - `无内鬼来点涩图`：非R18图片
    - `无内鬼来点色图`：R18图片
    - `无内鬼来点{keyword1|keyword2|...}`：根据关键词搜索，exp：`无内鬼来点白毛`
  - 我在校园自动签到：`wzxy -h`查看更多
- 私聊
  - 消息发送apiToken：`getToken/resetToken/deleteToken`
- 群聊
  - 关键词回复
    - 添加：`开始添加`
    - 删除：`删除自动回复:[keyword]`
    - 触发：`keyword`
- 管理员
  - bash命令执行：`$bash <command>`
  - 我在校园签到服务token：`$wzxy -h`
  - 添加监听群组(群聊功能只对监听群组触发)：`$lg -h`

## 私信发送接口

revue提供了消息发送接口，为方便测试，这里提供一个已经部署好的示例，请按照以下步骤操作：

1. 添加revue测试QQ机器人`3056159050`为好友。（如果没有及时通过好友申请，请邮件告知我：me@kfccrazythursday.buzz）
2. 向revues私聊发送`/help`根据提示获取`token`，或直接发送`/getToken`获取。
3. 向`http://revue.magicode123.cn:5000/send_private_msg`发送对应字段。

| key         | 功能                                                         |
| ----------- | ------------------------------------------------------------ |
| token       | 获取的token                                                  |
| ~~user_id~~ | ~~qq号，token和qq为绑定状态，也就是一个token只能对一个qq号发送消息~~（出于安全考虑，目前已经不需要传入qq号了） |
| message     | 消息内容，也可以支持表情，语音，短视频等内容，发送格式为CQ码，参照[CQcode\|帮助中心 ](https://docs.go-cqhttp.org/cqcode/#cqcode) |

- 示例

```json
{
    "token":"e0c405ae-95e9-4039-9f1f-4f39f7e6bde4",
    "message":"测试"
}
```

### 使用场景

> 提醒这次ssh登录与上次ip不一致的情况，防止陌生人登录

将下面文件保存，然后在`.bashrc`或`.zshrc`最后加上`bash <文件路径>`即可，这样每次启动时都可对登陆者ip进行检查。

```shell
#!/bin/bash
revue_token="" # 填入自己的revue token，向机器人申请
LAST_LOGIN_IP=$(lastlog -u $USER | awk 'NR==2{print $3}')
THIS_LOGIN_IP=$(who | grep -P '([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))(\.([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))){3}' -o | awk 'NR==1')
if [ $LAST_LOGIN_IP != $THIS_LOGIN_IP ]; then
    THIS_LOGIN_PLACE=$(curl -s "cip.cc/$THIS_LOGIN_IP" | awk 'NR==7{print $3}')
    HOSTNAME=$(hostname)
    curl -s --location --request POST 'http://revue.magicode123.cn:5000/send_private_msg' \
    --header 'Content-Type: text/plain' \
    --data-raw '{
        "token":"'$revue_token'",
        "message":"你的服务器'$HOSTNAME'在'$THIS_LOGIN_PLACE'登录了,IP地址为'$THIS_LOGIN_IP',与上一次登录IP地址'$LAST_LOGIN_IP'不同"
    }' | grep "&*(#&$*($"
fi
```

**效果：**

![img](http://cdnimg.violetwsh.com/img/A2FD185EF8DFC04DD368F995DE323819.png)

### python-requests

```python
import requests
import json

url = "http://revue.magicode123.cn:5000/send_private_msg"

payload = json.dumps({
  "token": "<token>",
  "user_id": "<QQ号>",
  "message": "<消息内容|支持CQcode>"
})
headers = {
  'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
```

### Go-resty

```go
package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	url := "http://revue.magicode123.cn:5000/send_private_msg"
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
	}).SetBody(map[string]string{
		"token":   "<token>",
		"user_id": "<QQ号>",
		"message": "<消息内容|支持CQcode>",
	}).Post(url)
	if err != nil {
		return
	}
	fmt.Println(string(post.Body()))
}
```

### NodeJs-Axios

```js
var axios = require('axios');
var data = JSON.stringify({
  "token": "<token>",
  "user_id": "<QQ号>",
  "message": "<消息内容|支持CQcode>"
});

var config = {
  method: 'post',
  url: 'http://revue.magicode123.cn:5000/send_private_msg',
  headers: { 
    'Content-Type': 'application/json'
  },
  data : data
};

axios(config)
.then(function (response) {
  console.log(JSON.stringify(response.data));
})
.catch(function (error) {
  console.log(error);
});
```

## 自行配置运行

先配置并运行`go-cqhttp`，其中默认开启中间件密钥以及反向post密钥，注意配置以下内容：

> 注意：这里给的示例假设go-cqhttp端和revue端都部署在同一台服务器

`[go-cqhttp]config.yaml`

```yaml
# 默认中间件锚点
default-middlewares: &default
  # 访问密钥, 强烈推荐在公网的服务器设置
  access-token: '<正向鉴权密钥>'
  
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
          secret: '<反向鉴权密钥>'        # 密钥
```

接下来配置`[revue]config.yaml`，在`qqBot-revue`下创建`config.yaml`文件，配置以下内容:

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
# 正向鉴权,即向go-cqhttp客户端发送消息时进行鉴权
forwardAuthentication:
  enable: true # true\false
  token: '<正向鉴权密钥>' 
# 反向鉴权,即接收go-cqhttp客户端消息时进行鉴权
reverseAuthentication:
  enable: true # true\false
  secret: '<反向鉴权密钥>'
# revue机器人接口相关
revue:
  enable: true # true\false
# 数据库相关
Database:
  # 当sqlite和mysql都为enable时,sqlite优先级大于mysql
  sqlite:
    enable: false # true\false
    path: './data.db'
  mysql:
    enable: true # true\false
    address: '' # 地址:端口
    dbname: '' # 数据库名
    charset: 'utf8mb4'
    username: '' # 数据库用户名
    password: '' # 数据库密码
  mongo:
    # 涩图数据库,如果为false则直接调用公共接口
    # 如果为true则调用自己的接口
    hImgDB:
      enable: true # true\false
      url: ''
```

