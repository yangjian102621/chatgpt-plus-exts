## ChatGPT-PLUS Extensions

ChatGPT-Plus 扩展模块，包括微信机器人，Midjourney 机器人，为 ChatGPT-Plus 项目提供扩展服务。

## 部署
本项目主要用来配合 [ChatGPT-Plus](https://github.com/yangjian102621/chatgpt-plus) 项目使用，作为 ChatGPT-PLUS 项目的扩展项目。

当然，你也可以集成到其他需要这些机器人的项目中，都非常好使。

首先将项目中的 sample 配置文件 `src/config.toml.sample` 拷贝一份
```shell
cd src
cp config.toml.sample config.toml
``` 
修改其中的几处配置，如果你不知道如何获取 Discord 用户 Token 和 Bot Token 请查参考 [Midjourney｜如何集成到自己的平台](https://zhuanlan.zhihu.com/p/631079476)。
```toml
Listen = "0.0.0.0:9001"
DataDir = "./data"
Debug = true # 生产环境可以关闭调试模式
ProxyURL = "http://172.22.11.47:7777" # 魔法地址，discord API 访问需要魔法，你懂的o(*￣︶￣*)o
CallbackToken = "e8hgcjp4fsn6kq1pbe6hkgkvf89svvbi" # ChatPlus 项目回调授权 Token，这里需要配置跟 ChatPlus 项目一样的，用于双方相互通信授权

[MidJourneyConfig]
  Enabled = true # 是否启动 MidJourney 机器人
  UserToken = "YOUR_USER_TOTEN" # 用户登录 Token
  BotToken = "YOUR_BOT_TOKEN" # 机器人登录 Token，经测试有的账号这里也填 UserToken 也可以正常运行，如果启动授权报错请单独设置 BotToken
  GuildId = "YOUR_GUILD_ID" # discord 服务器 ID
  ChanelId = "YOUR_CHANEL_ID" # discord 频道 ID
  CallbackUrl = "http://localhost:5678/api/mj/notify" # MidJourney 绘画消息推送接口

[WeChatConfig]
  Enabled = true # 是否启动微信机器人
  CallbackUrl = "http://localhost:5678/api/reward/notify" # 微信转账消息推送接口

[RedisConfig] # redis 配置，用来存储简易消息队列
  Address = "localhost:6379"
  Password = ""
  Db = 0
```

配置修改完成之后就可以启动服务了。

```shell
go run main.go
```

**注意：一定要确认有启动成功的日志，** 如果你的魔法不行，机器人可能卡在那里好久都不报错，所以不要看着暂时没有报错就觉得启动成功了，要确认成功启动并且返回了，才可以确定启动成功了。

```log
2023-08-15T08:58:31.629+0800    INFO    mj/bot.go:58    Starting MidJourney Bot...
2023-08-15T08:58:31.629+0800    INFO    mj/bot.go:69    Starting consume MidJourney messages...
2023-08-15T08:58:31.629+0800    INFO    core/app_server.go:32   http://0.0.0.0:9001
2023-08-15T08:58:32.838+0800    INFO    mj/bot.go:64    Starting MidJourney Bot successfully! # 需要有这句日志次才表示 MidJourney 机器人真正启动成功了
```

如果是启动微信机器人，终端控制台应该会输出一个命令行的登录二维码，你直接用微信扫码登录即可监控当前登录账号所有的转账信息，实现动账通知的功能。

```log
9:4:49 app         | 2023-08-15T09:04:49.789+0800       INFO    weixin/bot.go:35        Starting WeChat Bot...
9:4:49 app         | 2023-08-15T09:04:49.789+0800       INFO    core/app_server.go:32   http://0.0.0.0:9001
9:4:49 app         | 2023-08-15T09:04:49.790+0800       INFO    mj/bot.go:58    Starting MidJourney Bot...
9:4:49 app         | 2023-08-15T09:04:49.789+0800       INFO    weixin/bot.go:57        Starting consume wechat messages...
9:4:49 app         | 2023-08-15T09:04:49.789+0800       INFO    mj/bot.go:69    Starting consume MidJourney messages...
9:4:49 app         | 2023-08-15T09:04:49.961+0800       INFO    weixin/bot.go:113       请使用微信扫描下面二维码登录
9:4:49 app         | 2023-08-15T09:04:49.961+0800       INFO    weixin/bot.go:115

        ██████████████  ██          ██  ██    ████  ██████████████
        ██          ██  ██████████████  ████        ██          ██
        ██  ██████  ██      ████  ██    ██  ████    ██  ██████  ██
        ██  ██████  ██  ██  ████  ██      ████  ██  ██  ██████  ██
        ██  ██████  ██        ██  ██  ████    ██    ██  ██████  ██
        ██          ██    ████  ████      ██████    ██          ██
        ██████████████  ██  ██  ██  ██  ██  ██  ██  ██████████████
                        ██    ██          ██  ██
        ██  ████  ██████    ████  ██████  ████      ██    ██  ████
          ██████  ██  ████████    ████      ████  ████████      ██
          ██████    ██          ██  ██          ██████      ████
                  ██    ██  ████████        ██  ██      ██      ██
        ██        ████  ██  ██      ██    ██  ██  ██      ████
          ██  ██        ██              ██    ██    ██      ██████
        ██████  ████████        ████      ██████████  ████  ██████
        ██    ██  ██    ████  ██    ██  ██        ████████    ██
            ████  ████  ██████  ████████    ████████  ██████  ██
          ████              ████    ██    ██    ████  ██  ██████
        ██  ████  ██████    ██      ██████            ██    ██
            ██████    ██                  ████  ██████  ██  ██
          ██    ██  ██    ██          ██  ██    ██████████████
                        ██████  ██████████████████      ██████████
        ██████████████  ██  ██          ██  ██████  ██  ████  ██
        ██          ██  ██      ██████  ██  ██  ██      ████
        ██  ██████  ██    ██          ██████    ██████████  ██████
        ██  ██████  ██  ████  ████████  ██    ██████  ██████    ██
        ██  ██████  ██  ████  ████  ██      ██████    ██    ██  ██
        ██          ██        ██  ████  ██  ██  ██      ████  ██
        ██████████████  ██  ██  ██    ██      ██████    ██    ██
```