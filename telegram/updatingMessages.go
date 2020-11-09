package telegram

import (
	"errors"
)

// https://core.telegram.org/bots/api#updating-messages
//通过以下方法，您可以更改消息历史记录中的现有消息，而不是通过操作结果发送新消息。这对于使用带回调查询的嵌入式键盘的消息最有用，但也可以帮助减少与常规聊天机器人进行对话时的混乱情况。
//请注意，当前仅可以在不使用Reply_markup或嵌入式键盘的情况下编辑邮件。

// EditMessageTextOptional EditMessageText 可选参数
type EditMessageTextOptional struct {
	ChatID    string `json:"chat_id"`    // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID int64  `json:"message_id"` // 如果未指定inline_message_id，则为必需。要编辑的消息的标识符

	InlineMessageID string `json:"inline_message_id"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符

	ParseMode             string                `json:"parse_mode"`               // 消息文本中的实体解析模式。有关更多详细信息，请参见格式化选项。
	Entities              []MessageEntity       `json:"entities"`                 // 出现在消息文本中的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"` // 禁用此消息中链接的链接预览
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup"`             // 嵌入式键盘的JSON序列化对象。
}

// EditMessageText 使用此方法来编辑文本和游戏消息。成功后，如果编辑后的消息不是嵌入式消息，则返回编辑后的 Message，否则返回True。
// https://core.telegram.org/bots/api#editmessagetext
func (a API) EditMessageText(text string, optional EditMessageTextOptional) (*Message, error) {
	if (optional.ChatID == "" && optional.MessageID == 0) || optional.InlineMessageID == "" {
		// 让 api 自己检测？
	}
	result := &Message{}
	err := a.handleOptional("/editMessageText", map[string]interface{}{"text": text}, &optional, result)
	return result, err
}

// EditMessageCaptionOptional EditMessageCaption 可选参数
type EditMessageCaptionOptional struct {
	ChatID    string `json:"chat_id"`    // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID int64  `json:"message_id"` // 如果未指定inline_message_id，则为必需。要编辑的消息的标识符

	InlineMessageID string `json:"inline_message_id"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符

	Caption         string                `json:"caption"`          // 消息的新标题，实体解析后为0-1024个字符
	ParseMode       string                `json:"parse_mode"`       // 消息标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities []MessageEntity       `json:"caption_entities"` // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup"`     // 嵌入式键盘的JSON序列化对象。
}

// EditMessageCaption 使用此方法来编辑邮件的标题。成功后，如果编辑后的消息不是嵌入式消息，则返回编辑后的 Message，否则返回True。
// https://core.telegram.org/bots/api#editmessagecaption
func (a API) EditMessageCaption(optional EditMessageCaptionOptional) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/editMessageCaption", nil, &optional, result)
	return result, err
}

// EditMessageMediaOptional EditMessageMedia 可选参数
type EditMessageMediaOptional struct {
	ChatID    string `json:"chat_id"`    // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID int64  `json:"message_id"` // 如果未指定inline_message_id，则为必需。要编辑的消息的标识符

	InlineMessageID string `json:"inline_message_id"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"` // 用于新的嵌入式键盘的JSON序列化对象。
}

// EditMessageMedia 使用此方法可以编辑动画，音频，文档，照片或视频消息。如果消息是消息专辑的一部分，则只能将其编辑为音频专辑的音频，仅可以编辑为文档专辑的文档，否则为照片或视频。编辑嵌入式消息时，无法上载新文件。通过文件file_id使用先前上传的文件或指定URL。成功后，如果已编辑的消息是由漫游器发送的，则返回已编辑的 Message，否则返回True。
// https://core.telegram.org/bots/api#editmessagemedia
func (a API) EditMessageMedia(media InputMedia, optional EditMessageMediaOptional) (*Message, error) { // TODO
	return nil, errors.New("TODO")
}

// EditMessageReplyMarkupOptional EditMessageReplyMarkup 可选参数
type EditMessageReplyMarkupOptional struct {
	ChatID    string `json:"chat_id"`    // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID int64  `json:"message_id"` // 如果未指定inline_message_id，则为必需。要编辑的消息的标识符

	InlineMessageID string `json:"inline_message_id"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"` // 嵌入式键盘的JSON序列化对象。
}

// EditMessageReplyMarkup 使用此方法可以编辑动画，音频，文档，照片或视频消息。如果消息是消息专辑的一部分，则只能将其编辑为音频专辑的音频，仅可以编辑为文档专辑的文档，否则为照片或视频。编辑嵌入式消息时，无法上载新文件。通过文件file_id使用先前上传的文件或指定URL。成功后，如果已编辑的消息是由漫游器发送的，则返回已编辑的 Message，否则返回True。
// https://core.telegram.org/bots/api#editmessagemedia
func (a API) EditMessageReplyMarkup(optional EditMessageReplyMarkupOptional) (*Message, error) {
	result := &Message{}
	err := a.handleOptional("/editMessageReplyMarkup", nil, &optional, result)
	return result, err
}

// StopPollOptional StopPoll 可选参数
type StopPollOptional struct {
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"` // 一个用于新消息嵌入式键盘的JSON序列化对象。
}

// StopPoll 使用此方法停止由漫游器发送的轮询。成功后，将返回已停止的具有最终结果的 Poll。
// https://core.telegram.org/bots/api#stoppoll
func (a API) StopPoll(chatID string, messageID int64, optional StopPollOptional) (*Poll, error) {
	result := &Poll{}
	err := a.handleOptional("/stopPoll", map[string]interface{}{"chat_id": chatID, "message_id": messageID}, &optional, result)
	return result, err
}

// DeleteMessage 使用此方法可以删除一条消息，包括服务消息，但有以下限制：
//-仅在少于48小时前发送过的消息才能删除。
//-仅在24小时前发送的私聊中的骰子消息才能删除。
//-漫游器可以删除私人聊天，群组和超级群组中的传出消息。
//-漫游器可以删除私人聊天中的传入消息。
//-授予can_post_messages权限的漫游器可以删除通道中的传出消息。
//-如果漫游器是网上论坛的管理员，则可以在其中删除任何消息。
//-如果漫游器在超级组或频道中具有can_delete_messages权限，则可以在其中删除任何消息。
//返回True 成功。
// https://core.telegram.org/bots/api#deletemessage
func (a API) DeleteMessage(chatID string, messageID int64) (bool, error) {
	var result bool
	err := a.handleOptional("/deleteMessage", nil, nil, &result)
	return result, err
}
