# webot - 企业微信机器人 golang SDK

webot 是根据企业微信机器人 [官方API文档](https://developer.work.weixin.qq.com/document/path/91770)，基于 [req](https://github.com/imroc/req) 封装的 golang SDK。

## 快速上手

**Install**
```go
go get -u github.com/imroc/webot
```

**Import**

```go
import github.com/imroc/webot
```

**Usage**

```go
bot := webot.New(webhookURL) // webhookURL 是添加机器人时自动生成的
bot.Debug(true) // 开启调试，可以看到所有请求和响应内容

// 发送文本消息
resp, err := bot.SendTextContent("hello world")
// ...

// 发送文本消息同时 @all
resp, err = bot.SendText(&webot.TextMessage{Content: "hello world", MentionedList: []string{"@all"}})
// ...

// 发送 markdown 格式消息
content := `
新增 <font color="warning">3个</font> 新客户:
1. 小霸王电脑公司 - 已消费 <font color="green">22000</font> 元
2. 富贵鸟皮鞋公司 - 已消费 <font color="green">18000</font> 元
3. 帝王大酒店 - 已消费 <font color="green">3000</font> 元
`
resp, err = bot.SendMarkdownContent(content)
// ...

// 发送文件消息
resp, err := bot.SendFileContent("hello.txt", []byte("hello world"))
// ...
```
