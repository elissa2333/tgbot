package telegram

// pack

/*Available methods
Bot API中的所有方法都不区分大小写。我们支持GET和POST HTTP方法。使用URL查询字符串或application/json或application/x-www-form-urlencoded或multipart/form-data在Bot API请求中传递参数。
成功调用后，将返回一个包含结果的JSON对象。
https://core.telegram.org/bots/api#available-methods*/

import (
	"errors"
	"io"
	"reflect"

	"github.com/elissa2333/httpc"

	"github.com/elissa2333/tgbot/utils"
)

// GetMe 获取 bot 信息
// https://core.telegram.org/bots/api#making-requests
func (a API) GetMe() (*User, error) {
	res, err := a.HTTPClient.Get("/getMe")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	m := &User{}
	err = HandleResp(res, m)
	return m, err
}

// LogOut 在本地启动机器人之前，请使用此方法从云Bot API服务器注销。您必须先注销该机器人，然后才能在本地运行该机器人，否则无法保证该机器人将收到更新。成功拨打电话后，您可以立即登录本地服务器，但在10分钟内将无法重新登录到云Bot API服务器。成功返回True。不需要参数。
// https://core.telegram.org/bots/api#logout
func (a API) LogOut() (bool, error) {
	res, err := a.HTTPClient.Post("/logOut")
	if err != nil {
		return false, err
	}

	var m bool
	err = HandleResp(res, &m)
	return m, err
}

// Close 使用此方法关闭bot实例之前，将其从一台本地服务器移至另一台本地服务器。您需要在调用此方法之前删除webhook，以确保在服务器重新启动后不会再次启动该bot。在启动漫游器后的前10分钟，该方法将返回错误429。成功返回True。不需要参数。
// https://core.telegram.org/bots/api#close
func (a API) Close() (bool, error) {
	res, err := a.HTTPClient.Post("/close")
	if err != nil {
		return false, err
	}

	var m bool
	err = HandleResp(res, &m)
	return m, err
}

// SendMessageOptional SendMessage可选参数
type SendMessageOptional struct {
	ParseMode                string          `json:"parse_mode,omitempty"`               // 消息文本中的实体解析模式。有关更多详细信息，请参见格式化选项。
	Entities                 []MessageEntity `json:"entities"`                           // 出现在消息文本中的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableWebPagePreview    bool            `json:"disable_web_page_preview,omitempty"` // 禁用此消息中链接的链接预览
	DisableNotification      bool            `json:"disable_notification,omitempty"`     // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id,omitempty"`      // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`        // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`             // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendMessage 发送短信。成功后，将返回发送的消息
// https://core.telegram.org/bots/api#sendmessage
func (a API) SendMessage(chatID string /*目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）*/, text string /*待发送消息的文本，实体解析后为1-4096个字符*/, optional *SendMessageOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "text": text}
	if optional != nil {
		optMap, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range optMap {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/SendMessage")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	msg := &Message{}
	err = HandleResp(res, msg)

	return msg, err
}

// Formatting options
// Bot API支持消息的基本格式。您可以在漫游器的消息中使用粗体，斜体，下划线和删除线文本，以及内联链接和预格式化的代码。电报客户端将相应地呈现它们。您可以使用markdown样式或HTML样式格式。
// https://core.telegram.org/bots/api#formatting-options

const (
	// FormatTypeAtMarkdownV2 要使用此模式，请在parse_mode字段中传递MarkdownV2
	FormatTypeAtMarkdownV2 = "MarkdownV2" // https://core.telegram.org/bots/api#markdownv2-style
	// FormatTypeAtHTML 要使用此模式，请在parse_mode字段中传递HTML
	FormatTypeAtHTML = "HTML" // https://core.telegram.org/bots/api#html-style
	// FormatTypeAtMarkdown 这是旧版模式，保留下来是为了向后兼容。要使用此模式，请在parse_mode字段中传递Markdown
	FormatTypeAtMarkdown = "Markdown" // https://core.telegram.org/bots/api#markdown-style
)

// ForwardMessageOptional ForwardMessage 可选参数
type ForwardMessageOptional struct {
	DisableNotification bool `json:"disable_notification"` // 静默发送消息。用户将收到没有声音的通知。
}

// ForwardMessage 转发任何类型的消息。成功后，将返回发送的 Message
// https://core.telegram.org/bots/api#forwardmessage
func (a API) ForwardMessage(chatID string, fromChatID string, messageID int64, optional *ForwardMessageOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "from_chat_id": fromChatID, "message_id": messageID}
	result := &Message{}
	err := a.handleOptional("/forwardMessage", m, optional, result)
	return result, err
}

// CopyMessageOptional CopyMessage 可选参数
type CopyMessageOptional struct {
	Caption                  string          `json:"caption"`                     // 媒体的新标题，实体解析后为0-1024个字符。如果未指定，则保留原始标题
	ParseMode                string          `json:"parse_mode"`                  // 在新标题中解析实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`            // 出现在新标题中的特殊实体的列表，可以指定这些实体而不是parse_mode
	DisableNotification      bool            `json:"disable_notification"`        // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id"`         // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"` // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup"`                // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// CopyMessage 使用此方法可以复制任何类型的消息。该方法类似于forwardMessages方法，但是复制的消息没有指向原始消息的链接。成功返回已发送消息的MessageId。
// https://core.telegram.org/bots/api#copymessage
func (a API) CopyMessage(chatID string, fromChatID string, messageID int64, optional *CopyMessageOptional) (int64, error) {
	result := &struct {
		MessageID int64 `json:"message_id"`
	}{}
	err := a.handleOptional("/copyMessage", map[string]interface{}{"chat_id": chatID, "from_chat_id": fromChatID, "message_id": messageID}, optional, result)
	return result.MessageID, err
}

// handleSendMedia 处理媒体发送
func (a API) handleSendMedia(uri string, chatID string, fileKey string, file io.Reader, optional interface{}) (*Message, error) {

	var rows []httpc.FromDataRow

	rows = append(rows, httpc.FromDataRow{
		Key:   "chat_id",
		Value: chatID,
	})
	rows = append(rows, httpc.FromDataRow{
		Key:  fileKey,
		Data: file,
	})

	if !reflect.ValueOf(optional).IsNil() {
		m, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}

		for k, v := range m {
			reader, ok := v.(io.Reader)
			if ok {
				rows = append(rows, httpc.FromDataRow{
					Key:  k,
					Data: reader,
				})
			} else {
				rows = append(rows, httpc.FromDataRow{
					Key:   k,
					Value: utils.ToString(v),
				})
			}
		}
	}

	res, err := a.HTTPClient.SetFromData(rows...).Post(uri)
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	err = HandleResp(res, msg)
	return msg, err
}

// SendPhotoOptional SendPhoto 可选参数
type SendPhotoOptional struct {
	Caption                  string          `json:"caption,omitempty"`              // 图片标题（当通过file_id重新发送照片时也可以使用），在实体解析后为0-1024个字符
	ParseMode                string          `json:"parse_mode,omitempty"`           // 解析照片标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableNotification      bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象
}

// SendPhoto 发送照片
// https://core.telegram.org/bots/api#sendphoto
func (a API) SendPhoto(chatID string, photo InputFile, optional *SendPhotoOptional) (*Message, error) {
	return a.handleSendMedia("/SendPhoto", chatID, "photo", photo, optional)
}

// SendAudioOptional SendAudio可选参数
type SendAudioOptional struct {
	Caption                  string          `json:"caption,omitempty"`              // 音频字幕，实体解析后0-1024个字符
	ParseMode                string          `json:"parse_mode,omitempty"`           // 解析音频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Duration                 int             `json:"duration,omitempty"`             // 音频持续时间（以秒为单位）
	Performer                string          `json:"performer,omitempty"`            // 演员
	Title                    string          `json:"title,omitempty"`                // 曲目名称
	Thumb                    InputFile       `json:"thumb,omitempty"`                // 已发送文件的缩略图；如果服务器端支持文件的缩略图生成，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”
	DisableNotification      bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int             `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendAudio 使用此方法发送音频文件。您的音频必须为.MP3或.M4A格式。成功后，将返回发送的消息。机器人目前最多可以发送50MB的音频文件，以后可能会更改此限制
// https://core.telegram.org/bots/api#sendaudio
func (a API) SendAudio(chatID string, audio InputFile, optional *SendAudioOptional) (*Message, error) {
	return a.handleSendMedia("/sendAudio", chatID, "audio", audio, optional)
}

// SendDocumentOptional SendDocument可选参数
type SendDocumentOptional struct {
	Thumb                       InputFile       `json:"thumb,omitempty"`                // 已发送文件的缩略图；如果服务器端支持文件的缩略图生成，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“attach://<file_attach_name>”。
	Caption                     string          `json:"caption,omitempty"`              // 文档标题（在通过file_id重新发送文档时也可以使用），实体解析后为0-1024个字符
	ParseMode                   string          `json:"parse_mode,omitempty"`           // 解析文档标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities             []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableContentTypeDetection bool            `json:"disable_content_type_detection"` // 对使用 multipart/form-data 上传的文件禁用服务器端内容类型自动检测
	DisableNotification         bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID            int64           `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply    bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup                 interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendDocument 发送常规文件。成功后，将返回发送的消息。漫游器当前可以发送最大50 MB的任何类型的文件，以后可能会更改此限制。
// https://core.telegram.org/bots/api#senddocument
func (a API) SendDocument(chatID string, document InputFile, optional *SendDocumentOptional) (*Message, error) {
	return a.handleSendMedia("/sendDocument", chatID, "document", document, optional)
}

// SendVideoOptional SendVideo可选参数
type SendVideoOptional struct {
	Duration                 int64           `json:"duration,omitempty"`             // 发送视频的持续时间（以秒为单位）
	Width                    int             `json:"width,omitempty"`                // 影片宽度
	Height                   int             `json:"height,omitempty"`               // 影片高度
	Thumb                    InputFile       `json:"thumb,omitempty"`                // 已发送文件的缩略图；如果服务器端支持文件的缩略图生成，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”
	Caption                  string          `json:"caption,omitempty"`              // 视频标题（当通过file_id重新发送视频时也可以使用），实体解析后为0-1024个字符
	ParseMode                string          `json:"parse_mode,omitempty"`           // 视频字幕中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	SupportsStreaming        bool            `json:"supports_streaming,omitempty"`   // 如果上传的视频适合流式传输，则通过True
	DisableNotification      bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendVideo 发送视频文件，Telegram客户端支持mp4视频（其他格式也可以作为Document发送）。成功后，将返回发送的消息。机器人目前最多可以发送50MB的视频文件，以后可能会更改此限制
// https://core.telegram.org/bots/api#sendvideo
func (a API) SendVideo(chatID string, video InputFile, optional *SendVideoOptional) (*Message, error) {
	return a.handleSendMedia("/sendVideo", chatID, "video", video, optional)
}

// SendAnimationOptional SendAnimation可选参数
type SendAnimationOptional struct {
	Duration                 int64           `json:"duration,omitempty"`             // 发送动画的持续时间（以秒为单位）
	Width                    int             `json:"width,omitempty"`                // 动画宽度
	Height                   int             `json:"height,omitempty"`               // 动画高度
	Thumb                    InputFile       `json:"thumb,omitempty"`                // 已发送文件的缩略图；如果服务器端支持文件的缩略图生成，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”
	Caption                  string          `json:"caption,omitempty"`              // 动画标题（当通过file_id重新发送动画时也可以使用），在实体解析后为0-1024个字符
	ParseMode                string          `json:"parse_mode,omitempty"`           // 解析动画标题中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableNotification      bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendAnimation 发送动画文件（无声音的GIF或H.264/MPEG-4 AVC视频）。成功后，将返回发送的消息。机器人目前最多可以发送50MB的动画文件，以后可能会更改此限制。
// https://core.telegram.org/bots/api#sendanimation
func (a API) SendAnimation(chatID string, animation InputFile, optional *SendAnimationOptional) (*Message, error) {
	return a.handleSendMedia("/sendAnimation", chatID, "animation", animation, optional)
}

// SendVoiceOptional SendVoice可选参数
type SendVoiceOptional struct {
	Caption                  string          `json:"caption,omitempty"`              // 语音消息标题，实体解析后0-1024个字符
	ParseMode                string          `json:"parse_mode,omitempty"`           // 语音消息标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities          []MessageEntity `json:"caption_entities"`               // 标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Duration                 int64           `json:"duration,omitempty"`             // 语音留言的持续时间（以秒为单位）
	DisableNotification      bool            `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64           `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{}     `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendVoice 如果希望Telegram客户端将文件显示为可播放的语音消息，请使用此方法发送音频文件。为此，您的音频必须是使用OPUS编码的.OGG文件（其他格式可能以Audio或Document的形式发送）。成功后，将返回发送的消息。漫游器当前可以发送最大50 MB的语音消息，将来可能会更改此限制。
// https://core.telegram.org/bots/api#sendvoice
func (a API) SendVoice(chatID string, voice InputFile, optional *SendVoiceOptional) (*Message, error) {
	return a.handleSendMedia("/sendVoice", chatID, "voice", voice, optional)
}

// OptionalVideoNoteOptional SendVideoNote可选参数
type OptionalVideoNoteOptional struct {
	Duration                 int64       `json:"duration,omitempty"`             // 发送视频的持续时间（以秒为单位）
	Length                   int         `json:"length,omitempty"`               // 视频宽度和高度，即视频消息的直径
	Thumb                    InputFile   `json:"thumb,omitempty"`                // 已发送文件的缩略图；如果服务器端支持文件的缩略图生成，则可以忽略。缩略图应为JPEG格式，并且大小应小于200 kB。缩略图的宽度和高度不应超过320。如果未使用multipart / form-data上传文件，则忽略该缩略图。缩略图不能重复使用，只能作为新文件上传，因此如果缩略图是使用<file_attach_name>下的multipart / form-data上传的，则可以传递“ attach：// <file_attach_name>”
	DisableNotification      bool        `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendVideoNote 从v.4.0开始，Telegram客户端支持最长1分钟的圆形mp4方形视频。使用此方法发送视频消息。
// https://core.telegram.org/bots/api#sendvideonote
func (a API) SendVideoNote(chatID string, videoNote InputFile, optional *OptionalVideoNoteOptional) (*Message, error) {
	return a.handleSendMedia("/sendVideoNote", chatID, "video_note", videoNote, optional)
}

// SendMediaGroupOptional SendMediaGroup 可选参数
type SendMediaGroupOptional struct {
	DisableNotification      bool  `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64 `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool  `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
}

//SendMediaGroup 将一组照片或视频作为相册发送。成功后，将返回已发送消息的数组。
// https://core.telegram.org/bots/api#sendmediagroup
func (a API) SendMediaGroup(chatID string, media interface{} /*Array of InputMediaAudio, InputMediaDocument, InputMediaPhoto and InputMediaVideo*/, optional *SendMediaGroupOptional) (*Message, error) { // TODO media 未测试
	//return a.handleSendMedia("/sendMediaGroup", chatID, "media", media, optional)

	return nil, errors.New("TODO")
}

// SendLocationOptional sendLocation 可选参数
type SendLocationOptional struct {
	HorizontalAccuracy   float64     `json:"horizontal_accuracy"`            // 位置的不确定性半径，以米为单位；0-1500
	LivePeriod           int64       `json:"live_period,omitempty"`          // 位置将被更新的秒数（请参阅实时位置，应在60到86400之间）。
	Heading              int64       `json:"heading"`                        // 对于现场位置，用户移动的方向（以度为单位）。如果指定，则必须介于1到360之间。
	ProximityAlertRadius int64       `json:"proximity_alert_radius"`         // 对于实时位置，有关接近另一个聊天成员的接近警报的最大距离（以米为单位）。如果指定，则必须介于1到100000之间。
	DisableNotification  bool        `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID     int64       `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	ReplyMarkup          interface{} `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendLocation 在地图上发送点。成功后，将返回发送的消息。
// https://core.telegram.org/bots/api#sendlocation
func (a API) SendLocation(chatID string, latitude float64, longitude float64, optional *SendLocationOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "latitude": latitude, "longitude": longitude}
	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/sendlocation")
	if err != nil {
		return nil, err
	}

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

// EditMessageLiveLocationOptional EditMessageLiveLocation 可选参数
type EditMessageLiveLocationOptional struct {
	ChatID               string                `json:"chat_id,omitempty"`           // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID            int64                 `json:"message_id,omitempty"`        // 如果未指定inline_message_id，则为必需。要编辑的消息的标识符
	InlineMessageID      string                `json:"inline_message_id,omitempty"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符
	HorizontalAccuracy   float64               `json:"horizontal_accuracy"`         // 位置的不确定性半径，以米为单位；0-1500
	Heading              int64                 `json:"heading"`                     // 用户移动的方向（以度为单位）。如果指定，则必须介于1到360之间。
	ProximityAlertRadius int64                 `json:"proximity_alert_radius"`      // 有关接近另一个聊天成员的接近警报的最大距离（以米为单位）。如果指定，则必须介于1到100000之间。
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // 用于新的嵌入式键盘的JSON序列化对象。
}

// EditMessageLiveLocation 编辑实时位置消息。可以编辑位置，直到其live_period到期，或者通过调用stopMessageLiveLocation显式禁用编辑。成功后，如果已编辑的消息是由bot发送的，则返回已编辑的消息，否则返回True
// https://core.telegram.org/bots/api#editmessagelivelocation
func (a API) EditMessageLiveLocation(latitude float64, longitude float64, optional *EditMessageLiveLocationOptional) (*Message, error) {
	m := map[string]interface{}{"latitude": latitude, "longitude": longitude}
	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/sendlocation")
	if err != nil {
		return nil, err
	}

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

// StopMessageLiveLocationOptional StopMessageLiveLocation 可选参数
type StopMessageLiveLocationOptional struct {
	ChatID          string                `json:"chat_id,omitempty"`           // 如果未指定inline_message_id，则为必需。目标聊天的唯一标识符或目标频道的用户名（格式为@channelusername）
	MessageID       int64                 `json:"message_id,omitempty"`        // 如果未指定inline_message_id，则为必需。带有停止位置的消息的标识符
	InlineMessageID string                `json:"inline_message_id,omitempty"` // 如果未指定chat_id和message_id，则为必需。内联消息的标识符
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // 用于新的嵌入式键盘的 JSON序列化对象。
}

// StopMessageLiveLocation 在live_period到期之前停止更新实时位置消息。成功后，如果该消息是由漫游器发送的，则返回已发送的消息，否则返回True。
// https://core.telegram.org/bots/api#stopmessagelivelocation
func (a API) StopMessageLiveLocation(optional *StopMessageLiveLocationOptional) (*Message, error) {
	if optional == nil {
		return nil, nil
	}

	m, err := utils.StructToMap(optional)
	if err != nil {
		return nil, err
	}

	res, err := a.HTTPClient.SetBody(m).Post("/stopmessagelivelocation")
	if err != nil {
		return nil, err
	}

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

// SendVenueOptional SendVenue 可选参数
type SendVenueOptional struct {
	FoursquareID             string `json:"foursquare_id,omitempty"`        // 场地的Foursquare标识符
	FoursquareType           string `json:"foursquare_type,omitempty"`      // 场地的Foursquare类型（如果已知）。（例如，“ arts_entertainment /默认”，“ arts_entertainment /水族馆”或“食品/冰淇淋”。）
	GooglePlaceID            string `json:"google_place_id"`                // 场地的Google地方信息标识符
	GooglePlaceType          string `json:"google_place_type"`              // 场所的Google地方信息类型。（请参阅支持的类型。）
	DisableNotification      string `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64  `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool   `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              int64  `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendVenue 发送有关场地的信息。成功后，将返回发送的消息。
// https://core.telegram.org/bots/api#sendvenue
func (a API) SendVenue(chatID string, latitude float64, longitude float64, title string, address string, optional *SendVenueOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "latitude": latitude, "longitude": longitude, "title": title, "address": address}

	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/sendvenue")
	if err != nil {
		return nil, err
	}

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

// SendContactOptional SendContact 可选参数
type SendContactOptional struct {
	LastName                 string      `json:"last_name,omitempty"`            // 联系人的姓氏
	Vcard                    string      `json:"vcard,omitempty"`                // 有关vCard形式的联系人的其他数据，0-2048字节
	DisableNotification      bool        `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除键盘或强制用户回复的说明的JSON序列化对象。
}

// SendContact 发送电话联系人。成功后，将返回发送的消息。
// https://core.telegram.org/bots/api#sendcontact
func (a API) SendContact(chatID string, phoneNumber string, firstName string, optional *SendContactOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "phone_number": phoneNumber, "first_name": firstName}

	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/sendcontact")
	if err != nil {
		return nil, err
	}

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

// SendPollOptional SendPoll 可选参数
type SendPollOptional struct {
	IsAnonymous              bool        `json:"is_anonymous,omitempty"`            // 如果轮询需要匿名，则默认为True
	Type                     string      `json:"type,omitempty"`                    // 投票类型，“测验”或“常规”，默认为“常规”
	AllowsMultipleAnswers    bool        `json:"allows_multiple_answers,omitempty"` // 正确，如果轮询允许多个答案，则在测验模式下被轮询忽略，默认为False
	CorrectOptionID          int64       `json:"correct_option_id,omitempty"`       // 测验模式下的轮询所需的基于0的正确答案选项的标识符
	Explanation              string      `json:"explanation,omitempty"`             // 当用户选择错误的答案或轻按测验样式的民意测验中的灯泡图标时显示的文本，实体解析后最多0至200个字符，最多2个换行
	ExplanationParseMode     string      `json:"explanation_parse_mode,omitempty"`  // 解释中的实体解析模式。有关更多详细信息，请参见格式化选项。
	OpenPeriod               int64       `json:"open_period,omitempty"`             // 创建后轮询将在5-600秒钟内激活的时间（以秒为单位）。不能与close_date一起使用。
	CloseDate                int64       `json:"close_date,omitempty"`              // 轮询将自动关闭的时间点（Unix时间戳）。将来必须至少为5秒且不超过600秒。不能与open_period一起使用。
	DisableNotification      bool        `json:"disable_notification,omitempty"`    // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`     // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply"`       // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`            // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendPoll 发送投票。成功后，将返回发送的 Message
// https://core.telegram.org/bots/api#sendpoll
func (a API) SendPoll(chatID string, question string, options []string, optional *PollOption) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID, "question": question, "options": options}
	result := &Message{}
	err := a.handleOptional("/sendPoll", m, optional, result)
	return result, err
}

// SendDiceOptional SendDice 可选参数
type SendDiceOptional struct {
	Emoji                    string      `json:"emoji,omitempty"`                // 掷骰子动画所基于的表情符号。当前，必须是“ 🎲”，“ 🎯”，“ 🏀”，“ ⚽”或“ 🎰”之一。骰子的“ 🎲”和“ 🎯”值可以为1-6，“ ”和“ ”的值可以为1-5，“ 🏀”的⚽值可以为1-64 🎰。默认为“ 🎲”
	DisableNotification      bool        `json:"disable_notification,omitempty"` // 静默发送消息。用户将收到没有声音的通知。
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`  // 如果消息是答复，则为原始消息的ID
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply"`    // 如果未发送指定的回复消息也应发送消息，则传递True
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`         // 其他界面选项。内联键盘，自定义回复键盘，删除回复键盘或强制用户回复的说明的JSON序列化对象。
}

// SendDice 发送动画的表情符号，它将显示随机值。成功后，将返回发送的消息。
// https://core.telegram.org/bots/api#senddice
func (a API) SendDice(chatID string, optional *SendDiceOptional) (*Message, error) {
	m := map[string]interface{}{"chat_id": chatID}
	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/sendDice")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &Message{}
	err = HandleResp(res, result)
	return result, err
}

const (
	// ActionTypeAtTyping  text messages
	ActionTypeAtTyping = "typing"
	// ActionTypeAtUploadPhoto  photos
	ActionTypeAtUploadPhoto = "upload_photo"
	// ActionTypeAtUploadVideo  videos
	ActionTypeAtUploadVideo = "upload_video"
	// ActionTypeAtUploadAudio  audio files
	ActionTypeAtUploadAudio = "upload_audio"
	// ActionTypeAtUploadDocument  general files
	ActionTypeAtUploadDocument = "upload_document"
	// ActionTypeAtFindLocation  location data
	ActionTypeAtFindLocation = "find_location"
	// ActionTypeAtUploadVideoNote  video notes
	ActionTypeAtUploadVideoNote = "upload_video_note"
)

// SendChatAction 当您需要告诉用户bot一侧正在发生某些事情时，请使用此方法。状态设置为5秒或更短（当从您的bot收到消息时，Telegram客户端会清除其键入状态）。成功返回True。
// https://core.telegram.org/bots/api#sendchataction
func (a API) SendChatAction(chatID string, action string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "action": action}).Post("/sendChatAction")
	if err != nil {
		return false, err
	}

	var result bool
	err = HandleResp(res, &result)
	return result, err
}

// GetUserProfilePhotosOptional GetUserProfilePhotos 可选参数
type GetUserProfilePhotosOptional struct {
	Offset int `json:"offset,omitempty"` // 要返回的第一张照片的序号。默认情况下，将返回所有照片。
	Limit  int `json:"limit,omitempty"`  // 限制要检索的照片数量。可接受1-100之间的值。默认为100。
}

//GetUserProfilePhotos 使用此方法可获取用户的个人资料图片列表
// https://core.telegram.org/bots/api#getuserprofilephotos
func (a API) GetUserProfilePhotos(userID string, optional *GetUserProfilePhotosOptional) (*UserProfilePhotos, error) {
	m := map[string]interface{}{"user_id": userID}
	if optional != nil {
		o, err := utils.StructToMap(optional)
		if err != nil {
			return nil, err
		}
		for k, v := range o {
			m[k] = v
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post("/getUserProfilePhotos")
	if err != nil {
		return nil, err
	}

	result := &UserProfilePhotos{}
	err = HandleResp(res, result)
	return result, err
}

// GetFile 获得有关文件的基本信息，并准备将其下载。目前，漫游器可以下载最大20MB的文件。成功后，将返回一个File对象。该文件然后可以通过链接下载https://api.telegram.org/file/bot<token>/<file_path>，其中<file_path>从响应服用。可以确保链接至少有效1个小时。当链接过期时，可以通过再次调用getFile来请求一个新的链接。
// https://core.telegram.org/bots/api#getfile
func (a API) GetFile(fileID string) (*File, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"file_id": fileID}).Post("/getFile")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &File{}
	err = HandleResp(res, result)
	return result, err
}

// handleOptional 处理可选参数
func (a API) handleOptional(url string, m map[string]interface{}, optional interface{}, result interface{}) error {
	if m == nil {
		m = map[string]interface{}{}
	}

	o, err := utils.StructToMap(optional)
	if err != nil && err.Error() != "input data is nil" && err.Error() != "input data type is not a struct" {
		return err
	}

	for k, v := range o {
		m[k] = v
	}

	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}

	res, err := a.HTTPClient.SetBody(m).Post(url)
	if err != nil {
		return err
	}
	return HandleResp(res, result)
}

// KickChatMemberOptional KickChatMember 可选参数
type KickChatMemberOptional struct {
	UntilDate int64 `json:"until_date,omitempty"` // 用户将被禁止的日期，Unix时间。如果从当前时间起，用户被禁止超过366天或少于30秒，则将其视为永远被禁止
}

// KickChatMember 将用户踢出组，超组或通道。在超级组和频道的情况下，除非先取消禁止，否则用户将无法使用邀请链接等自行返回该组。该漫游器必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限
// https://core.telegram.org/bots/api#kickchatmember
func (a API) KickChatMember(chatID string, userID int64, optional *KickChatMemberOptional) (bool, error) {
	m := map[string]interface{}{"chat_id": chatID, "user_id": userID}

	var result bool
	err := a.handleOptional("/kickChatMember", m, optional, &result)
	return result, err
}

// UnbanChatMemberOptional UnbanChatMember 可选参数
type UnbanChatMemberOptional struct {
	OnlyIfBanned bool `json:"only_if_banned"` // 如果不禁止用户，则不执行任何操作
}

// UnbanChatMember 取消超级组或频道中以前被踢过的用户的权限。用户将不会自动返回到组或频道，但将能够通过链接等加入。bot必须是管理员才能起作用。
// https://core.telegram.org/bots/api#unbanchatmember
func (a API) UnbanChatMember(chatID string, userID string, optional *UnbanChatMemberOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/unbanChatMember", map[string]interface{}{"chat_id": chatID, "user_id": userID}, optional, &result)
	return result, err
}

//RestrictChatMemberOptional RestrictChatMember 可选参数
type RestrictChatMemberOptional struct {
	UntilDate int `json:"until_date,omitempty"` // 取消限制日期的日期，Unix时间。如果从当前时间起限制用户使用超过366天或少于30秒，则将其视为永远受限制
}

// RestrictChatMember 限制超组中的用户。bot必须是超级组中的管理员才能起作用，并且必须具有适当的管理员权限。为所有权限传递True
// https://core.telegram.org/bots/api#restrictchatmember
func (a API) RestrictChatMember(chatID string, userID string, permissions ChatPermissions, optional *RestrictChatMemberOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/restrictChatMember", map[string]interface{}{"chat_id": chatID, "permissions": permissions}, optional, &result)
	return result, err
}

// PromotionChatMemberOptional PromotionChatMember 可选参数
type PromotionChatMemberOptional struct {
	CanChangeInfo      bool `json:"can_change_info,omitempty"`      // 如果管理员可以更改聊天标题，照片和其他设置，则通过True
	CanPostMessages    bool `json:"can_post_messages,omitempty"`    // 如果管理员可以创建频道帖子，则通过True，仅频道
	CanEditMessages    bool `json:"can_edit_messages,omitempty"`    // 如果管理员可以编辑其他用户的消息并可以固定消息和通道，则通过True
	CanDeleteMessages  bool `json:"can_delete_messages,omitempty"`  // 如果管理员可以删除其他用户的消息，则通过True
	CanInviteUsers     bool `json:"can_invite_users,omitempty"`     // 如果管理员可以邀请新用户加入聊天，请通过True
	CanRestrictMembers bool `json:"can_restrict_members,omitempty"` // 如果管理员可以限制，禁止或取消禁止聊天成员，则通过True
	CanPinMessages     bool `json:"can_pin_messages,omitempty"`     // 如果管理员可以固定消息，则通过True，仅超级组
	CanPromoteMembers  bool `json:"can_promote_members,omitempty"`  // 如果管理员可以添加具有自己特权的子集的新管理员，或者将他已经直接或间接提升的管理员降级（由他任命的管理员提升），则通过True
}

// PromotionChatMember 提升或降级超组或渠道中的用户。该漫游器必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限。为所有布尔参数传递False降级用户
// https://core.telegram.org/bots/api#promotechatmember
func (a API) PromotionChatMember(chatID string, userID int64, optional *PromotionChatMemberOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/promoteChatMember", map[string]interface{}{"chat_id": chatID, "user_id": userID}, optional, &result)
	return result, err
}

// SetChatAdministratorCustomTitle 由bot提升的超组中的管理员设置自定义标题
// https://core.telegram.org/bots/api#setchatadministratorcustomtitle
func (a API) SetChatAdministratorCustomTitle(chatID string, userID string, customTitle string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "user_id": userID, "custom_title": customTitle}).Post("/setChatAdministratorCustomTitle")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// SetChatPermissions 为所有成员设置默认的聊天权限。bot必须是组或超级组中的管理员才能起作用，并且必须具有 can_restrict_members 管理员权限
// https://core.telegram.org/bots/api#setchatpermissions
func (a API) SetChatPermissions(chatID string, permissions ChatPermissions) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "permissions": permissions}).Post("/setChatPermissions")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// ExportChatInviteLink 为聊天生成新的邀请链接；先前生成的所有链接均被吊销。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限。成功返回新邀请链接
//https://core.telegram.org/bots/api#exportchatinvitelink
func (a API) ExportChatInviteLink(chatID string) (string, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/exportChatInviteLink")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result string
	err = HandleResp(res, &result)

	return result, err
}

// SetChatPhoto 为聊天设置新的个人资料照片。私人聊天无法更改照片。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限
// https://core.telegram.org/bots/api#setchatphoto
func (a API) SetChatPhoto(chatID string, photo InputFile) (bool, error) {
	var rows []httpc.FromDataRow
	rows = append(rows, httpc.FromDataRow{
		Key:   "chat_id",
		Value: chatID,
	})
	rows = append(rows, httpc.FromDataRow{
		Key:  "photo",
		Data: photo,
	})
	res, err := a.HTTPClient.SetFromData(rows...).Post("/setChatPhoto")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)
	return result, err
}

// DeleteChatPhoto 删除聊天照片。私人聊天无法更改照片。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限
// https://core.telegram.org/bots/api#deletechatphoto
func (a API) DeleteChatPhoto(chatID string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/deleteChatPhoto")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// SetChatTitle 使用此方法可以更改聊天标题。私人聊天的标题无法更改。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限。成功返回True。
// https://core.telegram.org/bots/api#setchattitle
func (a API) SetChatTitle(chatID string, title string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "title": title}).Post("/setChatTitle")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// SetChatDescriptionOptional SetChatDescription 可选参数
type SetChatDescriptionOptional struct {
	Description string `json:"description,omitempty"` // 新的聊天描述，0-255个字符
}

// SetChatDescription 使用此方法可以更改组，超组或通道的描述。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限
// https://core.telegram.org/bots/api#setchatdescription
func (a API) SetChatDescription(chatID string, optional *SetChatDescriptionOptional) (bool, error) { // TODO 看一下可不可以不要 可选参数
	var result bool
	err := a.handleOptional("/setChatDescription", map[string]interface{}{"chat_id": chatID}, optional, &result)
	return result, err
}

// PinChatMessageOptional PinChatMessage 可选参数
type PinChatMessageOptional struct {
	DisableNotification bool `json:"disable_notification,omitempty"` // 如果不需要向所有聊天成员发送有关新固定消息的通知，则传递True。通道中始终禁用通知。
}

// PinChatMessage 将消息固定在组，超组或通道中。该bot必须是聊天中的管理员才能正常工作，并且必须在超级组中具有“ can_pin_messages”管理员权限，或在该频道中具有“ can_edit_messages”管理员权限
// https://core.telegram.org/bots/api#pinchatmessage
func (a API) PinChatMessage(chatID string, messageID string, optional PinChatMessageOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/pinChatMessage", map[string]interface{}{"chat_id": chatID, "message_id": messageID}, optional, &result)
	return result, err
}

// UnpinChatMessageOptional UnpinChatMessage 可选参数
type UnpinChatMessageOptional struct {
	MessageID int64 `json:"message_id"` // 要取消固定的消息的标识符。如果未指定，最新的固定消息（按发送日期）将被取消固定。
}

// UnpinChatMessage 使用此方法可以取消固定组，超组或通道中的消息。该bot必须是聊天中的管理员才能正常工作，并且必须在超级组中具有“can_pin_messages”管理员权限，或在该频道中具有“can_edit_messages”管理员权限。
// https://core.telegram.org/bots/api#unpinchatmessage
func (a API) UnpinChatMessage(chatID string, optional *UnpinChatMessageOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/unpinChatMessage", map[string]interface{}{"chat_id": chatID}, optional, &result)
	return result, err
}

// UnpinAllChatMessages 使用此方法可以清除聊天中的固定消息列表。如果该聊天不是私人聊天，则该bot必须是该聊天的管理员才能正常工作，并且必须在超级组中具有“can_pin_messages”管理员权限，或者在一个频道中必须具有“can_edit_messages”管理员权限。成功返回True。
// https://core.telegram.org/bots/api#unpinallchatmessages
func (a API) UnpinAllChatMessages(chatID string) (bool, error) {
	var result bool
	err := a.handleOptional("/unpinAllChatMessages", map[string]interface{}{"chat_id": chatID}, nil, &result)
	return result, err
}

// LeaveChat 使用此方法可以离开组，超组或频道
// https://core.telegram.org/bots/api#leavechat
func (a API) LeaveChat(chatID string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/leaveChat")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// GetChat 获取有关聊天的最新信息
// https://core.telegram.org/bots/api#getchat
func (a API) GetChat(chatID string) (*Chat, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/getChat")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &Chat{}
	err = HandleResp(res, result)

	return result, err
}

// GetChatAdministrators 获取聊天中的管理员列表。成功后，返回一个ChatMember对象数组，其中包含有关除其他bot以外的所有聊天管理员的信息。如果聊天是组或超组，并且没有任命管理员，则仅返回创建者
// https://core.telegram.org/bots/api#getchatadministrators
func (a API) GetChatAdministrators(chatID string) ([]ChatMember, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/getChatAdministrators")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result []ChatMember
	err = HandleResp(res, &result)

	return result, err
}

// GetChatMembersCount 使用此方法获取聊天中的成员数
// https://core.telegram.org/bots/api#getchatmemberscount
func (a API) GetChatMembersCount(chatID string) (int64, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/getChatAdministrators")
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var result int64
	err = HandleResp(res, &result)

	return result, err
}

// GetChatMember 获取有关聊天成员的信息
// https://core.telegram.org/bots/api#getchatmember
func (a API) GetChatMember(chatID string, userID int64) (*ChatMember, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "user_id": userID}).Post("/getChatMember")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &ChatMember{}
	err = HandleResp(res, result)

	return result, err
}

// SetChatStickerSet 使用此方法可以为超组设置新的组标签集。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限。使用getChat请求中可选返回的字段can_set_sticker_set来检查bot是否可以使用此方法
// https://core.telegram.org/bots/api#setchatstickerset
func (a API) SetChatStickerSet(chatID string, stickerSetName string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID, "sticker_set_name": stickerSetName}).Post("/setChatStickerSet")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// DeleteChatStickerSet 使用此方法从超组中删除组标签集。该bot必须是聊天中的管理员才能起作用，并且必须具有适当的管理员权限。使用getChat请求中可选返回的字段can_set_sticker_set来检查bot是否可以使用此方法
// https://core.telegram.org/bots/api#deletechatstickerset
func (a API) DeleteChatStickerSet(chatID string) (bool, error) {
	res, err := a.HTTPClient.SetBody(map[string]interface{}{"chat_id": chatID}).Post("/deleteChatStickerSet")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// AnswerCallbackQueryOptional AnswerCallbackQuery 可选参数
type AnswerCallbackQueryOptional struct {
	Text      string `json:"text,omitempty"`       // 通知文字。如果未指定，则不会显示任何内容，0-200个字符
	ShowAlert string `json:"show_alert,omitempty"` // 如果为true，则客户端将显示警报，而不是在聊天屏幕顶部显示通知。默认为false。
	URL       string `json:"url,omitempty"`        // 用户客户端将打开的URL。如果您已经创建了游戏并通过@Botfather接受了条件，请指定打开游戏的URL —请注意，仅当查询来自callback_game按钮时，此方法才有效。
	CacheTime int64  `json:"cache_time,omitempty"` // 回调查询的结果可以在客户端缓存的最长时间（以秒为单位）。电报应用程序将从版本3.14开始支持缓存。预设为0。
}

// AnswerCallbackQuery 发送答案给从嵌入式键盘发送的回调查询。答案将作为通知显示在用户的聊天屏幕顶部或作为警报
// https://core.telegram.org/bots/api#answercallbackquery
func (a API) AnswerCallbackQuery(callbackQueryID string, optional *AnswerCallbackQueryOptional) (bool, error) {
	var result bool
	err := a.handleOptional("/answerCallbackQuery", map[string]interface{}{"callback_query_id": callbackQueryID}, optional, &result)

	return result, err
}

// SetMyCommands 更改 bot 命令列表
// https://core.telegram.org/bots/api#setmycommands
func (a API) SetMyCommands(commands ...BotCommand) (bool, error) {
	if len(commands) == 0 {
		return false, nil
	}

	res, err := a.HTTPClient.SetBody(map[string]interface{}{"commands": commands}).Post("/setMyCommands")
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result bool
	err = HandleResp(res, &result)

	return result, err
}

// GetMyCommands 获取 bot 命令的当前列表
// https://core.telegram.org/bots/api#getmycommands
func (a API) GetMyCommands() ([]BotCommand, error) {
	res, err := a.HTTPClient.Get("/getMyCommands")
	if err != nil {
		return nil, err
	}

	var result []BotCommand
	err = HandleResp(res, &result)
	return result, err
}
