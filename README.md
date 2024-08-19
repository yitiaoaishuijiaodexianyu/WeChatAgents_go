
WeChatAgents go语言客户端  
--

go语言版本使用的为1.22.0

目前已经实现

🟢点歌功能

🟢猜歌名小游戏

🟢发刺激刺激图片

🟢星火ai

🟢退群监控

其余功能需自行实现 非完善版本 遇见BUG 可能需要自行修复


为Bot爱好者，搭建免费的微信Bot Agents平台 ，万物皆Agent!

请合理合法使用，简单托管，无任何复杂配置环境，配置好SOCKS5代理登录成功后，WebSocket对接实时消息和调用功能。

尽量不要用自己的主号、或大号，追风号！S5只是绕过风控的一个环节，且S5必须要运行在用户的本地环境，云服务不可🙅

Gradio App 开源，后端核心，闭源。开源有开源的好处，比如：共同协作，Bug追踪。闭源可以躲避Ai风控追踪，代码审查，等等。不喜勿用。用户无需担心隐私问题，风控问题。

当然我们也不会，采集用户隐私,仅提供安全，稳定的Bot服务，开放必要的API能力， 提供发送✅音频✅的能力，暂不提供❌发送小视频❌的能力(预防动作片的导演)。

登录流程

![img.png](https://z.wiki/autoupload/20240714/7P3R/1542X1461/Untitled.png?type=ha)


| 平台版本 | 支持情况 | 登录类型     | 支持情况 | 消息结构 | 支持情况 | 开放能力         | 支持情况 | 开放事件       | 支持情况 |
| -------- | -------- | ------------ | -------- | -------- | -------- | ---------------- | -------- | -------------- | -------- |
| Windows  | 🟢        | 扫码登录     | 🟢        | 接收文字 | 🟢        | 发送文字         | ✅        | 收到好友请求   | 🟢        |
| macOS    | 🟢        | 推送登录     | 🟢        | 接收表情 | 🟢        | 发送表情         | ✅        | 收到拍一拍     | 🟢        |
| Linux    | 🔴        | 账号密码登录 | 🔴        | 接收图片 | 🟢        | 发送图片         | ✅        | 机器人消息     | 🟢        |
|          |          |              |          | 接收APP  | 🟢        | 发送APP          | ✅        | 实时消息事件   | 🟢        |
|          |          |              |          | 接收语音 | 🟢        | 发送语音         | ✅        | 实时朋友圈事件 | 🟢        |
|          |          |              |          | 接收视频 | 🟢        | 通过好友请求     | ❌        | 邀请进群事件   | 🟢        |
|          |          |              |          |          |          | 邀请成员入群     | ❌        |                |          |
|          |          |              |          |          |          | 发送视频         | ❌        |                |          |
|          |          |              |          |          |          | 发送CDN图片      | ✅        |                |          |
|          |          |              |          |          |          | 发送CDN文件      | ✅        |                |          |
|          |          |              |          |          |          | 下载CDN图片      | ✅        |                |          |
|          |          |              |          |          |          | 下载CDN文件      | ✅        |                |          |
|          |          |              |          |          |          | 小程序           | ❌        |                |          |
|          |          |              |          |          |          | 撤回消息         | ✅        |                |          |
|          |          |              |          |          |          | 拍一拍           | ✅        |                |          |


具体使用说明：https://aiagents-wechatagents.hf.space/

下面指令为win10系统操作的编译64位可执行文件
--

---
#x86 win 编译  
set GOOS=windows  
set GOARCH=amd64  
go build -ldflags="-s -w" -o win_x86_x64.exe main.go  
---
#arm win 编译  
set GOOS=windows  
set GOARCH=arm64  
go build -ldflags="-s -w" -o win_arm_x64.exe main.go  
---
#x86 linux 编译  
set GOOS=linux  
set GOARCH=amd64  
go build -ldflags="-s -w" -o linux_x86_x64 main.go  
---
#arm linux 编译  
set GOOS=linux  
set GOARCH=arm64  
go build -ldflags="-s -w" -o linux_arm_x64 main.go  
---
#x86 macOS 编译  
set GOOS=darwin  
set GOARCH=amd64  
go build -ldflags="-s -w" -o darwin_x86_x64 main.go  
---
#arm macOS 编译  
set GOOS=darwin  
set GOARCH=arm64  
go build -ldflags="-s -w" -o darwin_arm_x64 main.go  
---

