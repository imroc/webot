# webot - 企业微信机器人 SDK

## 快速上手

**Install**
```go
go get github.com/imroc/webot
```

**Import**

```go
import github.com/imroc/webot
```

**Usage**

发送纯文本:

```go
bot := webot.New(webhookURL) // webhookURL 是添加机器人时自动生成的
bot.Debug(true) // 开启调试，可以看到所有请求和响应内容
resp, err := bot.SendText(&webot.TextMessage{Content: "hello world", MentionedList: []string{"@all"}})
```
