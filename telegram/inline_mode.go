package telegram

// pack

/*Inline mode
以下方法和对象使您的机器人可以以串联模式工作。有关更多详细信息，
请参见我们的内联机器人简介。
要启用此选项，请将/setinline命令发送到@BotFather并提供用户在键入您的机器人名称后将在输入字段中看到的占位符文本。
https://core.telegram.org/bots/api#inline-mode*/

// InlineQuery 内联查询
// https://core.telegram.org/bots/api#inlinequery
type InlineQuery struct {
	ID       string    `json:"id"`       // 此查询的唯一标识符
	From     *User     `json:"from"`     // 发件人
	Location *Location `json:"location"` // 可选的。发件人位置，仅适用于请求用户位置的机器人
	Query    string    `json:"query"`    // 查询文字（最多256个字符）
	Offset   string    `json:"offset"`   // 要返回的结果的偏移量，可以由机器人控制
}

// ChosenInlineResult 用户选择并发送给其聊天伙伴的内联查询的结果
// https://core.telegram.org/bots/api#choseninlineresult
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`         // 所选结果的唯一标识符
	From            *User     `json:"from"`              // 选择结果的用户
	Location        *Location `json:"location"`          // 可选的。发件人位置，仅适用于需要用户位置的机器人
	InlineMessageID string    `json:"inline_message_id"` // 可选的。发送的内联消息的标识符。仅当消息上装有嵌入式键盘时可用。也将在回调查询中接收到，并可用于编辑消息。
	Query           string    `json:"query"`             // 用于获取结果的查询
}
