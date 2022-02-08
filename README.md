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

发送纯文本消息:

```go
bot := webot.New(webhookURL) // webhookURL 是添加机器人时自动生成的
bot.Debug(true) // 开启调试，可以看到所有请求和响应内容

// 发送文本消息
resp, err := bot.SendTextContent("hello world")
//...

// 发送消息同时 @all
resp, err = bot.SendText(&webot.TextMessage{Content: "hello world", MentionedList: []string{"@all"}})
```

发送 Markdown 消息:

```go
bot := webot.New(webhookURL) // webhookURL 是添加机器人时自动生成的
bot.Debug(true) // 开启调试，可以看到所有请求和响应内容

// 发送文本消息
content := `
实时新增用户反馈<font color=\"warning\">132例</font>，请相关同事注意。
> 类型:<font color=\"comment\">用户反馈</font>
> 普通用户反馈:<font color=\"comment\">117例</font>
> VIP用户反馈:<font color=\"comment\">15例</font>`
resp, err := bot.SendMarkdownContent(content)
```