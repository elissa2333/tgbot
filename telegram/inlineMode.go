package telegram

import (
	"errors"

	"github.com/elissa2333/tgbot/utils"
)

// pack

/*Inline mode
以下方法和对象使您的机器人可以以串联模式工作。有关更多详细信息，
请参见我们的内联机器人简介。
要启用此选项，请将 /setinline 命令发送到 @BotFather 并提供用户在键入您的机器人名称后将在输入字段中看到的占位符文本。
https://core.telegram.org/bots/api#inline-mode*/

// InlineQuery 该对象表示传入的内联查询。当用户发送空查询时，您的漫游器可能会返回一些默认或趋势结果。
// https://core.telegram.org/bots/api#inlinequery
type InlineQuery struct {
	ID       string    `json:"id"`       // 此查询的唯一标识符
	From     *User     `json:"from"`     // 发件人
	Location *Location `json:"location"` // 可选的。发件人位置，仅适用于请求用户位置的机器人
	Query    string    `json:"query"`    // 查询文字（最多256个字符）
	Offset   string    `json:"offset"`   // 要返回的结果的偏移量，可以由机器人控制
}

// AnswerInlineQueryOptional AnswerInlineQuery 可选参数
type AnswerInlineQueryOptional struct {
	CacheTime         int64  `json:"cache_time,omitempty"`          // 内联查询的结果可以在服务器上缓存的最长时间（以秒为单位）。默认为300。
	IsPersonal        bool   `json:"is_personal,omitempty"`         // 如果将结果仅在发送查询的用户的服务器端缓存，则传递True。默认情况下，结果可能会返回给发送相同查询的任何用户
	NextOffset        string `json:"next_offset,omitempty"`         // 在下一个查询中传递带有相同文本的客户端应发送的偏移量，以接收更多结果。如果没有更多结果或不支持分页，则传递一个空字符串。偏移长度不能超过64个字节。
	SwitchPmText      string `json:"switch_pm_text,omitempty"`      // 如果通过，客户端将显示带有指定文本的按钮，该按钮会将用户切换到与机器人的私人聊天，并向机器人发送带有参数switch_pm_parameter的启动消息
	SwitchPmParameter string `json:"switch_pm_parameter,omitempty"` // 示例：发送YouTube视频的嵌入式漫游器可以要求用户将漫游器连接到其YouTube帐户，以相应地调整搜索结果。为此，它会在结果上方甚至显示任何结果之前显示一个“连接您的YouTube帐户”按钮。用户按下按钮，切换到与机器人的私人聊天，并在此过程中传递一个开始参数，该参数指示机器人返回oauth链接。完成后，该机器人可以提供一个switch_inline按钮，以便用户可以轻松地返回到他们想使用该机器人的内联功能的聊天室。
}

// AnswerInlineQuery 使用此方法将答案发送给内联查询。成功时，返回True。每个查询的结果不得超过50个。
// https://core.telegram.org/bots/api#answerinlinequery
func (a API) AnswerInlineQuery(inlineQueryID string, results []InlineQueryResult, optional *AnswerInlineQueryOptional) (bool, error) {
	m := map[string]interface{}{"inline_query_id": inlineQueryID}
	if optional != nil {
		om, err := utils.StructToMap(optional)
		if err != nil {
			return false, err
		}
		for k, v := range om {
			m[k] = v
		}
	}
	var merge []map[string]interface{}
	for _, v := range results {
		convM := map[string]interface{}{}
		var err error
		switch {
		default:
			return false, errors.New("no parameters specified")
		case v.InlineQueryResultCachedAudio != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedAudio)
		case v.InlineQueryResultCachedDocument != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedDocument)
		case v.InlineQueryResultCachedGif != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedGif)
		case v.InlineQueryResultCachedMpeg4Gif != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedMpeg4Gif)
		case v.InlineQueryResultCachedPhoto != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedPhoto)
		case v.InlineQueryResultCachedSticker != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedSticker)
		case v.InlineQueryResultCachedVideo != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedVideo)
		case v.InlineQueryResultCachedVoice != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultCachedVoice)
		case v.InlineQueryResultArticle != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultArticle)
		case v.InlineQueryResultAudio != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultAudio)
		case v.InlineQueryResultContact != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultContact)
		case v.InlineQueryResultGame != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultGame)
		case v.InlineQueryResultDocument != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultDocument)
		case v.InlineQueryResultGif != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultGif)
		case v.InlineQueryResultLocation != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultLocation)
		case v.InlineQueryResultMpeg4Gif != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultMpeg4Gif)
		case v.InlineQueryResultPhoto != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultPhoto)
		case v.InlineQueryResultVenue != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultVenue)
		case v.InlineQueryResultVideo != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultVideo)
		case v.InlineQueryResultVoice != nil:
			convM, err = utils.StructToMap(v.InlineQueryResultVoice)
		}
		if err != nil {
			return false, err
		}
		merge = append(merge, convM)
	}

	m["results"] = merge
	var result bool
	err := a.handleOptional("/answerInlineQuery", m, nil, &result)
	return result, err
}

// InlineQueryResult 该对象表示内联查询的一个结果。电报客户端当前支持以下20种类型的结果：
// https://core.telegram.org/bots/api#inlinequeryresult
type InlineQueryResult struct {
	*InlineQueryResultCachedAudio
	*InlineQueryResultCachedDocument
	*InlineQueryResultCachedGif
	*InlineQueryResultCachedMpeg4Gif
	*InlineQueryResultCachedPhoto
	*InlineQueryResultCachedSticker
	*InlineQueryResultCachedVideo
	*InlineQueryResultCachedVoice
	*InlineQueryResultArticle
	*InlineQueryResultAudio
	*InlineQueryResultContact
	*InlineQueryResultGame
	*InlineQueryResultDocument
	*InlineQueryResultGif
	*InlineQueryResultLocation
	*InlineQueryResultMpeg4Gif
	*InlineQueryResultPhoto
	*InlineQueryResultVenue
	*InlineQueryResultVideo
	*InlineQueryResultVoice
}

// InlineQueryResultArticleType InlineQueryResultArticle Type 字段值
const InlineQueryResultArticleType = "article"

// InlineQueryResultArticle 表示文章或网页的链接。
// https://core.telegram.org/bots/api#inlinequeryresultarticle
type InlineQueryResultArticle struct {
	Type                string                `json:"type,omitempty"`                  // 结果类型，必须为 article
	ID                  string                `json:"id,omitempty"`                    // 此结果的唯一标识符，1-64字节
	Title               string                `json:"title,omitempty"`                 // 结果标题
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"` // 要发送的消息内容
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // 可选的。消息附带的嵌入式键盘
	URL                 string                `json:"url,omitempty"`                   // 可选的。结果网址
	HideURL             bool                  `json:"hide_url,omitempty"`              // 可选的。如果您不希望在邮件中显示URL，则传递True
	Description         string                `json:"description,omitempty"`           // 可选的。结果简短说明
	ThumbURL            string                `json:"thumb_url,omitempty"`             // 可选的。结果缩略图的网址
	ThumbWidth          int64                 `json:"thumb_width,omitempty"`           // 可选的。缩图宽度
	ThumbHeight         int64                 `json:"thumb_height,omitempty"`          // 可选的。缩略图高度
}

// InlineQueryResultPhotoType InlineQueryResultPhoto 类型
const InlineQueryResultPhotoType = "photo"

// InlineQueryResultPhoto 表示照片的链接。默认情况下，该照片将由用户发送并带有可选标题。或者，您可以使用input_message_content发送具有指定内容（而不是照片）的消息。
// https://core.telegram.org/bots/api#inlinequeryresultphoto
type InlineQueryResultPhoto struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 photo
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	PhotoURL            string                `json:"photo_url"`             // 照片的有效网址。照片必须为jpeg格式。相片大小不得超过5MB
	ThumbURL            string                `json:"thumb_url"`             // 照片缩略图的网址
	PhotoWidth          int64                 `json:"photo_width"`           // 可选的。照片的宽度
	PhotoHeight         int64                 `json:"photo_height"`          // 可选的。照片的宽度
	Title               string                `json:"title"`                 // 可选的。结果标题
	Description         string                `json:"description"`           // 可选的。结果简短说明
	Caption             string                `json:"caption"`               // 可选的。要发送的照片的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析照片标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是照片的内容
}

// InlineQueryResultGifType InlineQueryResultGif 类型
const InlineQueryResultGifType = "gif"

// InlineQueryResultGif 表示指向动画GIF文件的链接。默认情况下，该动画GIF文件将由用户带有可选的标题发送。或者，您可以使用input_message_content而不是动画发送具有指定内容的消息。
// https://core.telegram.org/bots/api#inlinequeryresultgif
type InlineQueryResultGif struct {
	Type                string                `json:"type"`                  // 结果类型，必须为gif
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	GifURL              string                `json:"gif_url"`               // GIF文件的有效URL。档案大小不得超过1MB
	GifWidth            int64                 `json:"gif_width"`             // 可选的。GIF宽度
	GifHeight           int64                 `json:"gif_height"`            // 可选的。GIF的高度
	GifDuration         int64                 `json:"gif_duration"`          // 可选的。GIF的持续时间
	ThumbURL            string                `json:"thumb_url"`             // 结果的静态（JPEG或GIF）或动画（MPEG4）缩略图的URL
	ThumbMimeType       string                `json:"thumb_mime_type"`       // 可选的。缩略图的MIME类型必须是 “image/jpeg”, “image/gif”, or “video/mp4”. 默认 “image/jpeg”
	Title               string                `json:"title"`                 // 可选的。结果标题
	Caption             string                `json:"caption"`               // 可选的。要发送的GIF文件的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容，而不是GIF动画
}

// InlineQueryResultMpeg4GifType InlineQueryResultMpeg4Gif 类型
const InlineQueryResultMpeg4GifType = "mpeg4_gif"

// InlineQueryResultMpeg4Gif 表示视频动画的链接（无声音的H.264 / MPEG-4 AVC视频）。默认情况下，该动画MPEG-4文件将由用户带有可选的标题发送。或者，您可以使用input_message_content而不是动画发送具有指定内容的消息。
// https://core.telegram.org/bots/api#inlinequeryresultmpeg4gif
type InlineQueryResultMpeg4Gif struct {
	Type                string                `json:"type"`                  // 结果类型，必须为 mpeg4_gif
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	Mpeg4URL            string                `json:"mpeg_4_url"`            // MP4文件的有效URL。档案大小不得超过1MB
	Mpeg4Width          int64                 `json:"mpeg_4_width"`          // 可选的。影片宽度
	Mpeg4Height         int64                 `json:"mpeg_4_height"`         // 可选的。影片高度
	Mpeg4Duration       int64                 `json:"mpeg_4_duration"`       // 可选的。影片时长
	ThumbURL            string                `json:"thumb_url"`             // 结果的静态（JPEG或GIF）或动画（MPEG4）缩略图的URL
	ThumbMimeType       string                `json:"thumb_mime_type"`       // 可选的。缩略图的MIME类型必须是 “image/jpeg”, “image/gif”, or “video/mp4”. 默认 “image/jpeg”
	Title               string                `json:"title"`                 // 可选的。结果标题
	Caption             string                `json:"caption"`               // 可选的。待发送的MPEG-4文件的标题，实体解析后为0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容，而不是视频动画
}

// InlineQueryResultVideoType InlineQueryResultVideo 类型
const InlineQueryResultVideoType = "video"

// InlineQueryResultVideo 表示指向包含嵌入式视频播放器或视频文件的页面的链接。默认情况下，该视频文件将由用户发送并带有可选的标题。或者，您可以使用input_message_content发送具有指定内容而不是视频的消息。
// 如果InlineQueryResultVideo消息包含嵌入式视频（例如YouTube），则必须使用input_message_content替换其内容
// https://core.telegram.org/bots/api#inlinequeryresultvideo
type InlineQueryResultVideo struct {
	Type                string               `json:"type"`                  // 结果类型，必须是 video
	ID                  string               `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	VideoURL            string               `json:"video_url"`             // 嵌入式视频播放器或视频文件的有效URL
	MimeType            string               `json:"mime_type"`             // 视频网址的内容的MIME类型，“text/html”或“ video/mp4”
	ThumbURL            string               `json:"thumb_url"`             // 视频缩略图的网址（仅jpeg）
	Title               string               `json:"title"`                 // 结果标题
	Caption             string               `json:"caption"`               // 可选的。要发送的视频的标题，实体解析后为0-1024个字符
	ParseMode           string               `json:"parse_mode"`            // 可选的。解析视频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity      `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	VideoWidth          int64                `json:"video_width"`           // 可选的。影片宽度
	VideoHeight         int64                `json:"video_height"`          // 可选的。影片高度
	VideoDuration       int64                `json:"video_duration"`        // 可选的。视频时长，以秒为单位
	Description         string               `json:"description"`           // 可选的。结果简短说明
	ReplyMarkup         InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent `json:"input_message_content"` // 可选的。要发送的消息内容而不是视频内容。如果使用InlineQueryResultVideo作为结果发送HTML页面（例如YouTube视频），则此字段为必填字段
}

// InlineQueryResultAudioType InlineQueryResultAudio 类型
const InlineQueryResultAudioType = "audio"

// InlineQueryResultAudio 表示指向MP3音频文件的链接。默认情况下，该音频文件将由用户发送。或者，您可以使用input_message_content发送具有指定内容而不是音频的消息。
// https://core.telegram.org/bots/api#inlinequeryresultaudio
type InlineQueryResultAudio struct {
	Type                string                `json:"type"`                  // 结果类型，必须为 audio
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	AudioURL            string                `json:"audio_url"`             // 音频文件的有效URL
	Title               string                `json:"title"`                 // 标题
	Caption             string                `json:"caption"`               // 可选的。标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析音频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	Performer           string                `json:"performer"`             // 可选的。演员
	AudioDuration       int64                 `json:"audio_duration"`        // 可选的。音频持续时间（以秒为单位）
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是音频
}

// InlineQueryResultVoiceType InlineQueryResultVoice 类型
const InlineQueryResultVoiceType = "voice"

// InlineQueryResultVoice 表示指向以OPUS编码的.OGG容器中的语音记录的链接。默认情况下，此录音将由用户发送。或者，您可以使用input_message_content来发送具有指定内容的消息，而不是语音消息。
// https://core.telegram.org/bots/api#inlinequeryresultvoice
type InlineQueryResultVoice struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 voice
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	VoiceURL            string                `json:"voice_url"`             // 录音的有效URL
	Title               string                `json:"title"`                 // 录音标题
	Caption             string                `json:"caption"`               // 可选的。标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。标题，实体解析后0-1024个字符
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	VoiceDuration       int64                 `json:"voice_duration"`        // 可选的。记录持续时间（秒）
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是语音记录
}

// InlineQueryResultDocumentType InlineQueryResultDocument 类型
const InlineQueryResultDocumentType = "document"

// InlineQueryResultDocument 表示文件的链接。默认情况下，该文件将由用户发送并带有可选的标题。或者，您可以使用input_message_content发送具有指定内容而不是文件的消息。当前，使用此方法只能发送.PDF和.ZIP文件。
// https://core.telegram.org/bots/api#inlinequeryresultdocument
type InlineQueryResultDocument struct {
	Type                string                `json:"type"`                  // 结果类型，必须为 document
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	Title               string                `json:"title"`                 // 结果标题
	Caption             string                `json:"caption"`               // 可选的。待发送文档的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析文档标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DocumentURL         string                `json:"document_url"`          // 文件的有效URL
	MimeType            string                `json:"mime_type"`             // 文件内容的MIME类型，“application/pdf”或“application/zip”
	Description         string                `json:"description"`           // 可选的。结果简短说明
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是文件的内容
	ThumbURL            string                `json:"thumb_url"`             // 可选的。文件缩略图的URL（仅jpeg）
	ThumbWidth          int64                 `json:"thumb_width"`           // 可选的。缩图宽度
	ThumbHeight         int64                 `json:"thumb_height"`          // 可选的。缩略图高度
}

// InlineQueryResultLocationType InlineQueryResultLocation 类型
const InlineQueryResultLocationType = "location"

// InlineQueryResultLocation 表示地图上的位置。默认情况下，该位置将由用户发送。或者，您可以使用input_message_content发送具有指定内容而不是位置的消息。
// https://core.telegram.org/bots/api#inlinequeryresultlocation
type InlineQueryResultLocation struct {
	Type                 string                `json:"type"`                   // 结果类型，必须是 location
	ID                   string                `json:"id"`                     // 此结果的唯一标识符，1-64字节
	Latitude             float64               `json:"latitude"`               // 位置纬度
	Longitude            float64               `json:"longitude"`              // 位置经度
	Title                string                `json:"title"`                  // 位置标题
	HorizontalAccuracy   float64               `json:"horizontal_accuracy"`    // 可选的。位置的不确定性半径，以米为单位；0-1500
	LivePeriod           int64                 `json:"live_period"`            // 可选的。可以更新位置的时间段（以秒为单位）应该在60到86400之间。
	Heading              int64                 `json:"heading"`                // 可选的。对于现场位置，用户移动的方向（以度为单位）。如果指定，则必须介于1到360之间。
	ProximityAlertRadius int64                 `json:"proximity_alert_radius"` // 可选的。对于实时位置，有关接近另一个聊天成员的接近警报的最大距离（以米为单位）。如果指定，则必须介于1到100000之间。
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup"`           // 可选的。消息附带的嵌入式键盘
	InputMessageContent  *InputMessageContent  `json:"input_message_content"`  // 可选的。要发送的消息内容而不是位置
	ThumbURL             string                `json:"thumb_url"`              // 可选的。结果缩略图的网址
	ThumbWidth           int64                 `json:"thumb_width"`            // 可选的。缩图宽度
	ThumbHeight          int64                 `json:"thumb_height"`           // 可选的。缩略图高度
}

// InlineQueryResultVenueType InlineQueryResultVenue 类型
const InlineQueryResultVenueType = "venue"

// InlineQueryResultVenue 代表场地。默认情况下，场地将由用户发送。或者，您可以使用input_message_content来发送具有指定内容的消息，而不是发送场地。
// https://core.telegram.org/bots/api#inlinequeryresultvenue
type InlineQueryResultVenue struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 venue
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64字节
	Latitude            float64               `json:"latitude"`              // 会场位置的纬度
	Longitude           float64               `json:"longitude"`             // 会场位置的经度
	Title               string                `json:"title"`                 // 会场名称
	Address             string                `json:"address"`               // 会场地址
	FoursquareID        string                `json:"foursquare_id"`         // 可选的。场地的Foursquare标识符（如果已知）
	FoursquareType      string                `json:"foursquare_type"`       // 可选的。场地的Foursquare类型（如果已知）。（例如，“ arts_entertainment /默认”，“ arts_entertainment /水族馆”或“食品/冰淇淋”。）
	GooglePlaceID       string                `json:"google_place_id"`       // 可选的。场地的Google地方信息标识符
	GooglePlaceType     string                `json:"google_place_type"`     // 可选的。场所的Google地方信息类型。（请参阅支持的类型。
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是地点
	ThumbURL            string                `json:"thumb_url"`             // 可选的。结果缩略图的网址
	ThumbWidth          int64                 `json:"thumb_width"`           // 可选的。缩图宽度
	ThumbHeight         int64                 `json:"thumb_height"`          // 可选的。缩略图高度
}

// InlineQueryResultContactType InlineQueryResultContact 类型
const InlineQueryResultContactType = "contact"

// InlineQueryResultContact 代表具有电话号码的联系人。默认情况下，此联系人将由用户发送。或者，您可以使用input_message_content而不是联系人发送具有指定内容的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcontact
type InlineQueryResultContact struct {
	Type                string                `json:"type"`                  // 结果类型，必须为 contact
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64字节
	PhoneNumber         string                `json:"phone_number"`          // 联系人的电话号码
	FirstName           string                `json:"first_name"`            // 联系人的名字
	LastName            string                `json:"last_name"`             // 可选的。联系人的姓氏
	Vcard               string                `json:"vcard"`                 // 可选的。vCard形式的有关联系人的其他数据，0-2048字节
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容，而不是联系人的内容
	ThumbURL            string                `json:"thumb_url"`             // 可选的。结果缩略图的网址
	ThumbWidth          int64                 `json:"thumb_width"`           // 可选的。缩图宽度
	ThumbHeight         int64                 `json:"thumb_height"`          // 可选的。缩略图高度
}

// InlineQueryResultGameType InlineQueryResultGame 类型
const InlineQueryResultGameType = "game"

// InlineQueryResultGame 代表一个 Game。
// https://core.telegram.org/bots/api#inlinequeryresultgame
type InlineQueryResultGame struct {
	Type          string                `json:"type"`            // 结果类型，必须是 game
	ID            string                `json:"id"`              // 此结果的唯一标识符，1-64个字节
	GameShortName string                `json:"game_short_name"` // 游戏简称
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup"`    // 可选的。消息附带的嵌入式键盘
}

// InlineQueryResultCachedPhotoType InlineQueryResultCachedPhoto 类型
const InlineQueryResultCachedPhotoType = "photo"

// InlineQueryResultCachedPhoto 表示指向存储在电报服务器上的照片的链接。默认情况下，该照片将由用户发送并带有可选的标题。或者，您可以使用input_message_content发送具有指定内容（而不是照片）的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedphoto
type InlineQueryResultCachedPhoto struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 photo
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	PhotoFileID         string                `json:"photo_file_id"`         // 照片的有效文件标识符
	Title               string                `json:"title"`                 // 可选的。结果标题
	Description         string                `json:"description"`           // 可选的。结果简短说明
	Caption             string                `json:"caption"`               // 可选的。要发送的照片的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析照片标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是照片的内容
}

// InlineQueryResultCachedGifType InlineQueryResultCachedGif 类型
const InlineQueryResultCachedGifType = " gif"

// InlineQueryResultCachedGif 表示指向存储在电报服务器上的GIF动画文件的链接。默认情况下，该动画GIF文件将由用户带有可选的标题发送。或者，您可以使用input_message_content而不是动画发送具有指定内容的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedgif
type InlineQueryResultCachedGif struct {
	Type                string                `json:"type"`                  // 结果类型，必须为gif
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	GifFileID           string                `json:"gif_file_id"`           // GIF文件的有效文件标识符
	Title               string                `json:"title"`                 // 可选的。结果标题
	Caption             string                `json:"caption"`               // 可选的。要发送的GIF文件的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容，而不是GIF动画
}

// InlineQueryResultCachedMpeg4GifType InlineQueryResultCachedMpeg4Gif 类型
const InlineQueryResultCachedMpeg4GifType = "mpeg4_gif"

// InlineQueryResultCachedMpeg4Gif 表示指向存储在电报服务器上的视频动画（无声音的H.264 / MPEG-4 AVC视频）的链接。默认情况下，该动画MPEG-4文件将由用户发送，并带有可选的标题。或者，您可以使用input_message_content而不是动画发送具有指定内容的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedmpeg4gif
type InlineQueryResultCachedMpeg4Gif struct {
	Type                string                `json:"type"`                  // 结果类型，必须为mpeg4_gif
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	Mpeg4FileID         string                `json:"mpeg_4_file_id"`        // MP4文件的有效文件标识符
	Title               string                `json:"title"`                 // 可选的。结果标题
	Caption             string                `json:"caption"`               // 可选的。待发送的MPEG-4文件的标题，实体解析后为0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容，而不是视频动画
}

// InlineQueryResultCachedStickerType InlineQueryResultCachedSticker 类型
const InlineQueryResultCachedStickerType = "sticker"

// InlineQueryResultCachedSticker 表示指向存储在电报服务器上的标签的链接。默认情况下，该标签将由用户发送。或者，您可以使用input_message_content发送具有指定内容而不是标签的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedsticker
type InlineQueryResultCachedSticker struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 sticker
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	StickerFileID       string                `json:"sticker_file_id"`       // 贴纸的有效文件标识符
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是标签
}

// InlineQueryResultCachedDocumentType InlineQueryResultCachedDocument 类型
const InlineQueryResultCachedDocumentType = "document"

// InlineQueryResultCachedDocument 表示指向存储在电报服务器上的文件的链接。默认情况下，该文件将由用户发送并带有可选的标题。或者，您可以使用input_message_content发送具有指定内容而不是文件的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcacheddocument
type InlineQueryResultCachedDocument struct {
	Type                string                `json:"type"`                  // 结果类型，必须为 document
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	Title               string                `json:"title"`                 // 结果标题
	DocumentFileID      string                `json:"document_file_id"`      // 该文件的有效文件标识符
	Description         string                `json:"description"`           // 可选的。结果简短说明
	Caption             string                `json:"caption"`               // 可选的。待发送文档的标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析文档标题中的实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是文件的内容
}

// InlineQueryResultCachedVideoType InlineQueryResultCachedVideo 类型
const InlineQueryResultCachedVideoType = "video"

// InlineQueryResultCachedVideo 表示指向存储在电报服务器上的视频文件的链接。默认情况下，该视频文件将由用户发送并带有可选的标题。或者，您可以使用input_message_content发送具有指定内容而不是视频的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedvideo
type InlineQueryResultCachedVideo struct {
	Type                string                `json:"type"`                  // 结果类型，必须是 video
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	VideoFileID         string                `json:"video_file_id"`         // 视频文件的有效文件标识符
	Title               string                `json:"title"`                 // 结果标题
	Description         string                `json:"description"`           // 可选的。结果简短说明
	Caption             string                `json:"caption"`               // 可选的。要发送的视频的标题，实体解析后为0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析视频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是视频内容
}

// InlineQueryResultCachedVoiceType InlineQueryResultCachedVoice类型
const InlineQueryResultCachedVoiceType = "voice"

// InlineQueryResultCachedVoice 表示指向存储在电报服务器上的语音消息的链接。默认情况下，此语音消息将由用户发送。或者，您可以使用input_message_content发送具有指定内容的消息，而不是语音消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedvoice
type InlineQueryResultCachedVoice struct {
	Type                string                `json:"type"`                  // 结果类型，必须是语音
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	VideoFileID         string                `json:"video_file_id"`         // 语音留言的有效文件标识符
	Title               string                `json:"title"`                 // 语音留言标题
	Caption             string                `json:"caption"`               // 可选的。标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。语音消息标题中的实体解析模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是语音消息
}

// InlineQueryResultCachedAudioType InlineQueryResultCachedAudio 类型
const InlineQueryResultCachedAudioType = "audio"

// InlineQueryResultCachedAudio 表示指向存储在电报服务器上的MP3音频文件的链接。默认情况下，该音频文件将由用户发送。或者，您可以使用input_message_content发送具有指定内容而不是音频的消息。
// https://core.telegram.org/bots/api#inlinequeryresultcachedaudio
type InlineQueryResultCachedAudio struct {
	Type                string                `json:"type"`                  // 结果类型，必须为音频
	ID                  string                `json:"id"`                    // 此结果的唯一标识符，1-64个字节
	AudioFileID         string                `json:"audio_file_id"`         // 音频文件的有效文件标识符
	Caption             string                `json:"caption"`               // 可选的。标题，实体解析后0-1024个字符
	ParseMode           string                `json:"parse_mode"`            // 可选的。解析音频字幕中实体的模式。有关更多详细信息，请参见格式化选项。
	CaptionEntities     []MessageEntity       `json:"caption_entities"`      // 可选的。标题中显示的特殊实体的列表，可以指定这些实体，而不是parse_mode
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`          // 可选的。消息附带的嵌入式键盘
	InputMessageContent *InputMessageContent  `json:"input_message_content"` // 可选的。要发送的消息内容而不是音频
}

// InputMessageContent 该对象表示作为内联查询结果要发送的消息的内容。电报客户端当前支持以下4种类型：
// https://core.telegram.org/bots/api#inputmessagecontent
type InputMessageContent struct {
	*InputTextMessageContent
	*InputLocationMessageContent
	*InputVenueMessageContent
	*InputContactMessageContent
}

// InputTextMessageContent 表示要作为内联查询结果发送的文本消息的内容。
// https://core.telegram.org/bots/api#inputtextmessagecontent
type InputTextMessageContent struct {
	MessageText           string          `json:"message_text"`             // 要发送的消息文本，1-4096个字符
	ParseMode             string          `json:"parse_mode"`               // 可选的。消息文本中的实体解析模式。有关更多详细信息，请参见格式化选项。
	Entities              []MessageEntity `json:"entities"`                 // 可选的。出现在消息文本中的特殊实体的列表，可以指定这些实体，而不是parse_mode
	DisableWebPagePreview bool            `json:"disable_web_page_preview"` // 可选的。禁用已发送消息中链接的链接预览
}

// InputLocationMessageContent 表示要作为内联查询结果发送的位置消息的内容。
// https://core.telegram.org/bots/api#inputlocationmessagecontent
type InputLocationMessageContent struct {
	Latitude             float64 `json:"latitude"`               // 位置的纬度
	Longitude            float64 `json:"longitude"`              // 位置的经度
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`    // 可选的。位置的不确定性半径，以米为单位；0-1500
	LivePeriod           int64   `json:"live_period"`            // 可选的。可以更新位置的时间段（以秒为单位）应该在60到86400之间。
	Heading              int64   `json:"heading"`                // 可选的。对于现场位置，用户移动的方向（以度为单位）。如果指定，则必须介于1到360之间。
	ProximityAlertRadius int64   `json:"proximity_alert_radius"` // 可选的。对于实时位置，有关接近另一个聊天成员的接近警报的最大距离（以米为单位）。如果指定，则必须介于1到100000之间。
}

// InputVenueMessageContent 表示作为内联查询结果发送的场所消息的内容。
// https://core.telegram.org/bots/api#inputvenuemessagecontent
type InputVenueMessageContent struct {
	Latitude        float64 `json:"latitude"`          // 场地的纬度
	Longitude       float64 `json:"longitude"`         // 场地经度
	Title           string  `json:"title"`             // 会场名称
	Address         string  `json:"address"`           // 	会场地址
	FoursquareID    string  `json:"foursquare_id"`     // 可选的。场地的Foursquare标识符（如果已知）
	FoursquareType  string  `json:"foursquare_type"`   // 可选的。场地的Foursquare类型（如果已知）。（例如，“ arts_entertainment /默认”，“ arts_entertainment /水族馆”或“食品/冰淇淋”。）
	GooglePlaceID   string  `json:"google_place_id"`   // 可选的。场地的Google地方信息标识符
	GooglePlaceType string  `json:"google_place_type"` // 可选的。场所的Google地方信息类型。（请参阅支持的类型。）
}

// InputContactMessageContent 表示作为内联查询结果发送的联系消息的内容。
// https://core.telegram.org/bots/api#inputcontactmessagecontent
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"` // 联系人的电话号码
	FirstName   string `json:"first_name"`   // 联系人的名字
	LastName    string `json:"last_name"`    // 可选的。联系人的姓氏
	Vcard       string `json:"vcard"`        // 可选的。vCard形式的有关联系人的其他数据，0-2048字节
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
