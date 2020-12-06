package tgbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	stdURL "net/url"

	"github.com/elissa2333/httpc"

	"github.com/elissa2333/tgbot/telegram"
	"github.com/elissa2333/tgbot/utils"
)

// MessageProcessorFunc 聊天处理器
type MessageProcessorFunc func(c *Context) error

// ActiveProcessorFunc 主动处理器函数
type ActiveProcessorFunc func(api *telegram.API) error

// Bot bot 实体
type Bot struct {
	API *telegram.API // telegram api

	timeout uint // 长时间轮询的超时时间（以秒为单位）为0即通常的短轮询。应该为正，短轮询应仅用于测试目的

	webHookEngine func() error

	MsgOffset int64 // 最后一条消息

	activeProcessorFunc []ActiveProcessorFunc

	commands       map[string]MessageProcessorFunc // 指定命令的执行方法
	defaultCommand MessageProcessorFunc            // 默认命令未指定命令时使用

	specifiedTypeMessageProcessorFunc map[string]interface{} // 指定类型消息处理器
	defaultMessageProcessorFunc       MessageProcessorFunc   // 默认消息处理器

	inlineQueryProcessorFunc InlineQueryProcessorFunc // 内联处理函数
	done                     chan struct{}            // 退出程序
	err                      chan error
}

// BotOptional bot 配置可选参数
type BotOptional struct {
	HTTPClient *http.Client
	Timeout    uint // // 长时间轮询的超时时间（以秒为单位）为0即通常的短轮询。应该为正，短轮询应仅用于测试目的
}

// New 新建 bot
func New(id int, token string, optional *BotOptional) *Bot {
	b := &Bot{
		API:      telegram.New(nil, id, token),
		timeout:  15,
		commands: map[string]MessageProcessorFunc{},
		done:     make(chan struct{}),
		err:      make(chan error),
	}

	if optional != nil {
		b.timeout = optional.Timeout

		if optional.HTTPClient != nil {
			b.API = telegram.New(optional.HTTPClient, id, token)
		}
	}

	return b
}

// SetWebhook 设置 webhook
func (b *Bot) SetWebhook(url string /*API 访问地址*/, address string /*本地监听地址*/, optional *telegram.WebhookOptional) error {
	parseURL, err := stdURL.Parse(url)
	if err != nil {
		return err
	}

	if err := b.API.SetWebhook(url, optional); err != nil {
		return err
	}

	b.webHookEngine = func() error {
		http.HandleFunc(parseURL.Path, func(writer http.ResponseWriter, request *http.Request) {
			defer request.Body.Close()
			if request.Method != http.MethodPost {
				writer.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			if request.Header.Get(httpc.ContentType) != httpc.MIMEJson {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			bodyB, err := ioutil.ReadAll(request.Body)
			if err != nil {
				b.handleError(err)
				return
			}

			m := telegram.Update{}
			if err := json.Unmarshal(bodyB, &m); err != nil {
				b.handleError(err)
				return
			}

			if m.UpdateID == 0 { // 坏请求
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			if b.inlineQueryProcessorFunc != nil {
				b.handleInlineQuery(m.InlineQuery)
			} else {
				b.handleReceivedMessages(m.Message)
			}
		})

		return http.ListenAndServe(address, nil)
	}

	return nil
}

// DeleteWebhook  删除 webhook
func (b *Bot) DeleteWebhook(optional *telegram.DeleteWebhookOptional) error {
	b.webHookEngine = nil
	return b.API.DeleteWebhook(optional)
}

// GetWebhookInfo 获取 webhook 信息
func (b *Bot) GetWebhookInfo() (*telegram.WebhookInfo, error) {
	return b.API.GetWebhookInfo()
}

// AddCommandProcessor 添加命令处理器（接收到命令后调用）
func (b *Bot) AddCommandProcessor(cmd string, execFunc MessageProcessorFunc) {
	b.commands[cmd] = execFunc
}

// SetDefaultCommandProcessor 设置默认命令处理器（在未找到命令时调用）
func (b *Bot) SetDefaultCommandProcessor(execFunc MessageProcessorFunc) {
	b.defaultCommand = execFunc
}

// MessageContextBase 基础上下文信息
type MessageContextBase struct {
	*telegram.API

	MessageID int64
	Form      *telegram.User
	Chat      *telegram.Chat

	ForwardFrom          *telegram.User // 可选的。对于转发的邮件，原始邮件的发件人
	ForwardFromChat      *telegram.Chat // 可选的。对于从频道转发的消息，有关原始频道的信息
	ForwardFromMessageID int64          // 可选的。对于从通道转发的消息，是通道中原始消息的标识符
	ForwardSignature     string         // 可选的。对于从频道转发的消息，请提供帖子作者的签名（如果有）
	ForwardSenderName    string         // 可选的。从用户转发的邮件的发件人名称，这些用户不允许在转发的邮件中添加指向其帐户的链接
	ForwardDate          int64          // 可选的。对于转发的消息，原始消息的发送日期为Unix时间
	ViaBot               *telegram.User // 可选的。发送消息的机器人
}

// GetChatID 获取聊天 ID
func (mcb *MessageContextBase) GetChatID() string {
	if mcb.Chat != nil {
		return utils.ToString(mcb.Chat.ID)
	}
	return ""
}

// TextMessageContext 文本消息上下文
type TextMessageContext struct {
	MessageContextBase

	ReplyToMessage *telegram.Message
	Text           string
}

// TextMessageProcessorFunc 文本消息处理函数
type TextMessageProcessorFunc func(c *TextMessageContext) error

// SetMessageProcessorAtText 设置消息类型为 text 的处理器
func (b *Bot) SetMessageProcessorAtText(fn TextMessageProcessorFunc) {
	b.setMessageProcessorAt(ContextTypeAtText, fn)
}

// MessageContextIncludeFile 上下文信息包含文件
type MessageContextIncludeFile struct {
	MessageContextBase
}

// GetDownloadURL 获取文件下载地址下载地址
func (mcf MessageContextIncludeFile) GetDownloadURL(filePath string) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%d:%s/%s", mcf.ID, mcf.Token, filePath)
}

// DownloadFile 下载文件
func (mcf MessageContextIncludeFile) DownloadFile(filePath string) (io.ReadCloser, error) {
	res, err := mcf.HTTPClient.DeleteBaseURL().Get(mcf.GetDownloadURL(filePath))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("http response code is not a 200")
	}

	return res.Body, err
}

// MessageContextAtPhoto 照片消息上下文
type MessageContextAtPhoto struct {
	MessageContextIncludeFile

	Photo   []telegram.PhotoSize
	Caption string
}

// ProcessorAtPhotoFunc 照片消息处理函数
type ProcessorAtPhotoFunc func(c *MessageContextAtPhoto) error

// SetMessageProcessorAtPhoto 照片处理器
func (b *Bot) SetMessageProcessorAtPhoto(fn ProcessorAtPhotoFunc) {
	b.setMessageProcessorAt(ContextTypeAtPhoto, fn)
}

// MessageContextAtVoice 语音消息上下文
type MessageContextAtVoice struct {
	MessageContextIncludeFile

	Voice *telegram.Voice
}

// ProcessorAtVoiceFunc 语音消息处理函数
type ProcessorAtVoiceFunc func(c *MessageContextAtVoice) error

// SetMessageProcessorAtVoice 语音消息处理器
func (b *Bot) SetMessageProcessorAtVoice(fn ProcessorAtVoiceFunc) {
	b.setMessageProcessorAt(ContextTypeAtVoice, fn)
}

// MessageContextAtAudio 音频消息上下文
type MessageContextAtAudio struct {
	MessageContextIncludeFile

	Audio   *telegram.Audio
	Caption string
}

// ProcessorAtAudioFunc 音频消息处理函数
type ProcessorAtAudioFunc func(c *MessageContextAtAudio) error

// SetMessageProcessorAtAudio 音频消息处理器
func (b *Bot) SetMessageProcessorAtAudio(fn ProcessorAtAudioFunc) {
	b.setMessageProcessorAt(ContextTypeAtAudio, fn)
}

// MessageContextAtVideo 视频消息上下文
type MessageContextAtVideo struct {
	MessageContextIncludeFile

	Video   *telegram.Video
	Caption string
}

// ProcessorAtVideoFunc 视频消息处理函数
type ProcessorAtVideoFunc func(c *MessageContextAtVideo) error

// SetMessageProcessorAtVideo 视频消息处理器
func (b *Bot) SetMessageProcessorAtVideo(fn ProcessorAtVideoFunc) {
	b.setMessageProcessorAt(ContextTypeAtVideo, fn)
}

// MessageContextAtAnimation 动画（gif）消息上下文
type MessageContextAtAnimation struct {
	MessageContextIncludeFile

	Animation *telegram.Animation
	Document  *telegram.Document
}

// ProcessorAtAnimationFunc 动画消息处理函数
type ProcessorAtAnimationFunc func(c *MessageContextAtAnimation) error

// SetMessageProcessorAtAnimation 动画消息处理器
func (b *Bot) SetMessageProcessorAtAnimation(fn ProcessorAtAnimationFunc) {
	b.setMessageProcessorAt(ContextTypeAtAnimation, fn)
}

// MessageContextAtDocument 文件消息上下文
type MessageContextAtDocument struct {
	MessageContextIncludeFile

	Document *telegram.Document
	Caption  string
}

// ProcessorAtDocumentFunc 文件消息处理函数
type ProcessorAtDocumentFunc func(c *MessageContextAtDocument) error

// SetMessageProcessorAtDocument 文件消息处理器
func (b *Bot) SetMessageProcessorAtDocument(fn ProcessorAtDocumentFunc) {
	b.setMessageProcessorAt(ContextTypeAtDocument, fn)
}

// MessageContextAtSticker 贴纸消息上下文
type MessageContextAtSticker struct {
	MessageContextIncludeFile

	Sticker *telegram.Sticker
}

// ProcessorAtStickerFunc 贴纸消息处理函数
type ProcessorAtStickerFunc func(c *MessageContextAtSticker) error

// SetMessageProcessorAtSticker 贴纸消息处理器
func (b *Bot) SetMessageProcessorAtSticker(fn ProcessorAtStickerFunc) {
	b.setMessageProcessorAt(ContextTypeAtSticker, fn)
}

// MessageContextAtVideoNote 视频笔记消息上下文
type MessageContextAtVideoNote struct {
	MessageContextIncludeFile

	VideoNote *telegram.VideoNote
}

// ProcessorAtVideoNoteFunc 视频笔记处理函数
type ProcessorAtVideoNoteFunc func(c *MessageContextAtVideoNote) error

// SetMessageProcessorAtVideoNote 视频笔记处理器
func (b *Bot) SetMessageProcessorAtVideoNote(fn ProcessorAtVideoNoteFunc) {
	b.setMessageProcessorAt(ContextTypeAtVideoNote, fn)
}

// MessageContextAtContact 联系人消息上下文
type MessageContextAtContact struct {
	MessageContextBase

	Contact *telegram.Contact
}

// ProcessorAtContactFunc 联系人消息处理函数
type ProcessorAtContactFunc func(c *MessageContextAtContact) error

// SetMessageProcessorAtContact 联系人消息处理器
func (b *Bot) SetMessageProcessorAtContact(fn ProcessorAtContactFunc) {
	b.setMessageProcessorAt(ContextTypeAtContact, fn)
}

// MessageContextAtDice 色子消息上下文
type MessageContextAtDice struct {
	MessageContextBase

	Dice *telegram.Dice
}

// ProcessorAtDiceFunc 色子消息处理函数
type ProcessorAtDiceFunc func(c *MessageContextAtDice) error

// SetMessageProcessorAtDice 色子消息处理器
func (b *Bot) SetMessageProcessorAtDice(fn ProcessorAtDiceFunc) {
	b.setMessageProcessorAt(ContextTypeAtDice, fn)
}

// MessageContextAtGame 游戏消息上下文
type MessageContextAtGame struct {
	MessageContextIncludeFile

	Game *telegram.Game
}

// ProcessorAtGameFunc 游戏消息处理函数
type ProcessorAtGameFunc func(c *MessageContextAtGame) error

// SetMessageProcessorAtGame 游戏消息处理函数
func (b *Bot) SetMessageProcessorAtGame(fn ProcessorAtGameFunc) {
	b.setMessageProcessorAt(ContextTypeAtGame, fn)
}

// MessageContextAtPoll 调查消息上下文
type MessageContextAtPoll struct {
	MessageContextBase

	Poll *telegram.Poll
}

// ProcessorAtPollFunc 调查消息处理函数
type ProcessorAtPollFunc func(c *MessageContextAtPoll) error

// SetMessageProcessorAtPoll 调查消息处理器
func (b *Bot) SetMessageProcessorAtPoll(fn ProcessorAtPollFunc) {
	b.setMessageProcessorAt(ContextTypeAtPoll, fn)
}

// MessageContextAtVenue 场地消息上下文
type MessageContextAtVenue struct {
	MessageContextBase

	Venue *telegram.Venue
}

// ProcessorAtVenueFunc 场地消息处理函数
type ProcessorAtVenueFunc func(c *MessageContextAtVenue) error

// SetMessageProcessorAtVenue 场地消息处理器
func (b *Bot) SetMessageProcessorAtVenue(fn ProcessorAtVenueFunc) {
	b.setMessageProcessorAt(ContextTypeAtVenue, fn)
}

// MessageContextAtLocation 共享位置消息上下文
type MessageContextAtLocation struct {
	MessageContextBase

	Location *telegram.Location
}

// ProcessorAtLocationFunc 共享位置消息处理函数
type ProcessorAtLocationFunc func(c *MessageContextAtLocation) error

// SetMessageProcessorAtLocation 共享位置消息处理器
func (b *Bot) SetMessageProcessorAtLocation(fn ProcessorAtLocationFunc) {
	b.setMessageProcessorAt(ContextTypeAtLocation, fn)
}

// setMessageProcessorAt 设置指定消息类型的处理器
func (b *Bot) setMessageProcessorAt(typeS string, fn interface{}) {
	if b.specifiedTypeMessageProcessorFunc == nil {
		b.specifiedTypeMessageProcessorFunc = map[string]interface{}{}
	}

	b.specifiedTypeMessageProcessorFunc[typeS] = fn
}

// SetMessageProcessor 消息处理器（接收到消息后调用）
func (b *Bot) SetMessageProcessor(handleMessageFunc MessageProcessorFunc) {
	b.defaultMessageProcessorFunc = handleMessageFunc
}

// InlineQueryContext 内联调用上下文
type InlineQueryContext struct {
	*telegram.API
	*telegram.InlineQuery
}

// InlineQueryProcessorFunc 内联处理函数
type InlineQueryProcessorFunc func(c *InlineQueryContext) error

// SetInlineQueryProcessor 设置内联处理器
func (b *Bot) SetInlineQueryProcessor(fn InlineQueryProcessorFunc) {
	b.inlineQueryProcessorFunc = fn
}

// AddActiveProcessor 添加主动处理器
func (b *Bot) AddActiveProcessor(activeProcessorFunc ActiveProcessorFunc) {
	if activeProcessorFunc != nil {
		b.activeProcessorFunc = append(b.activeProcessorFunc, activeProcessorFunc)
	}
}

// Run 运行 bot
// 只有再拥有 Processor 时才会正常阻塞
func (b *Bot) Run() error {
	_, err := b.API.GetMe() // check api
	if err != nil {
		return fmt.Errorf("check api call failed: %w", err)
	}

	if (b.webHookEngine) != nil { // 为了和主动处理器行为一致
		go func() {
			b.checkTask()
			if err := b.webHookEngine(); err != nil {
				b.err <- err
			}
		}()
	} else {
		if err := b.DeleteWebhook(nil); err != nil {
			return err
		}
		b.checkTask()
		go b.initiativeEngine()
	}

loop:
	for {
		select {
		case err := <-b.err:
			return err
		case <-b.done:
			break loop
		}
	}

	return nil
}

// handleError 处理错误
func (b *Bot) handleError(err error) {
	if err != nil {
		b.err <- err
	}
}

// checkTask 检查任务再合适的地方结束
func (b *Bot) checkTask() {
	totalNumberOfActiveAndPassive := 0
	cleanActiveAndPassiveCh := make(chan struct{})
	for k, vFn := range b.activeProcessorFunc {
		totalNumberOfActiveAndPassive++
		go func(index int, fn ActiveProcessorFunc) {
			if err := fn(b.API); err != nil {
				humanRead := index + 1 // 从 0开始的改成人类可读数
				b.handleError(fmt.Errorf("the %d AddActiveProcessor: %w", humanRead, err))
			}
			cleanActiveAndPassiveCh <- struct{}{}
		}(k, vFn)
	}

	if len(b.commands) != 0 || b.defaultCommand != nil || b.defaultMessageProcessorFunc != nil || len(b.specifiedTypeMessageProcessorFunc) != 0 || b.inlineQueryProcessorFunc != nil { // 统计被动
		totalNumberOfActiveAndPassive++
	}

	go func() {
		num := 0
		if num >= totalNumberOfActiveAndPassive {
			b.done <- struct{}{}
		}
		for range cleanActiveAndPassiveCh {
			num++
			if num >= totalNumberOfActiveAndPassive {
				b.done <- struct{}{}
			}
		}
	}()
}

// initiativeEngine 核心调度
func (b *Bot) initiativeEngine() {
	for {
		updates, err := b.API.GetUpdates(b.MsgOffset, 1, b.timeout)
		if err != nil {
			b.handleError(err)
			break
		}
		if len(updates) == 0 {
			continue
		}

		LastOneUpdate := updates[len(updates)-1]
		b.MsgOffset = LastOneUpdate.UpdateID + 1 // 记录消息偏量

		if LastOneUpdate.InlineQuery != nil {
			b.handleInlineQuery(LastOneUpdate.InlineQuery)
		} else {
			b.handleReceivedMessages(updates[0].Message)
		}
	}
}

// handleInlineQuery 处理内联查询
func (b *Bot) handleInlineQuery(query *telegram.InlineQuery) {
	go func() {
		if b.inlineQueryProcessorFunc != nil {
			err := b.inlineQueryProcessorFunc(&InlineQueryContext{
				API:         b.API,
				InlineQuery: query,
			})
			if err != nil {
				b.handleError(err)
			}
		}
	}()
}

// handleReceivedMessages 处理接收消息
func (b *Bot) handleReceivedMessages(message *telegram.Message) {
	ctx := &Context{
		Message: message,
		API:     b.API,
	}

	// 消息类型判断
	switch {
	case message.Text != "":
		ctx.MessageType = ContextTypeAtText
	case message.Photo != nil:
		ctx.MessageType = ContextTypeAtPhoto
	case message.Voice != nil:
		ctx.MessageType = ContextTypeAtVoice
	case message.Audio != nil:
		ctx.MessageType = ContextTypeAtAudio
	case message.Video != nil:
		ctx.MessageType = ContextTypeAtVideo
	case message.Animation != nil:
		ctx.MessageType = ContextTypeAtAnimation
	case message.Document != nil:
		ctx.MessageType = ContextTypeAtDocument
	case message.Sticker != nil:
		ctx.MessageType = ContextTypeAtSticker
	case message.VideoNote != nil:
		ctx.MessageType = ContextTypeAtVideoNote
	case message.Contact != nil:
		ctx.MessageType = ContextTypeAtContact
	case message.Dice != nil:
		ctx.MessageType = ContextTypeAtDice
	case message.Game != nil:
		ctx.MessageType = ContextTypeAtGame
	case message.Poll != nil:
		ctx.MessageType = ContextTypeAtPoll
	case message.Venue != nil:
		ctx.MessageType = ContextTypeAtVenue
	case message.Location != nil:
		ctx.MessageType = ContextTypeAtLocation
	}

	go func() {
		for _, messageEntity := range message.Entities { // 可能是 bot
			if messageEntity.Type == telegram.MessageEntityAtBotCommand {
				fn, ok := b.commands[message.Text]
				if !ok {
					if b.defaultCommand != nil { // 默认命令处理器
						if err := b.defaultCommand(ctx); err != nil {
							b.handleError(err)
						}
					}

					break // 默认命令处理器不存在则转为消息处理器进行处理
				}

				if err := fn(ctx); err != nil { // 命令处理器
					b.handleError(err)
				}
			}
		}
		if _, ok := b.specifiedTypeMessageProcessorFunc[ctx.MessageType]; ok {
			base := MessageContextBase{
				API:       ctx.API,
				MessageID: ctx.Message.MessageID,
				Form:      ctx.Message.From,
				Chat:      ctx.Message.Chat,

				ForwardFrom:          ctx.Message.ForwardFrom,
				ForwardFromChat:      ctx.Message.ForwardFromChat,
				ForwardFromMessageID: ctx.Message.ForwardFromMessageID,
				ForwardSignature:     ctx.Message.ForwardSignature,
				ForwardSenderName:    ctx.Message.ForwardSenderName,
				ForwardDate:          ctx.Message.ForwardDate,
				ViaBot:               ctx.Message.ViaBot,
			}

			switch ctx.MessageType {
			case ContextTypeAtText:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(TextMessageProcessorFunc)
				c := &TextMessageContext{
					MessageContextBase: base,
					ReplyToMessage:     ctx.Message.ReplyToMessage,
					Text:               ctx.Message.Text,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtText: %w", err))
					return
				}
			case ContextTypeAtPhoto:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtPhotoFunc)
				c := &MessageContextAtPhoto{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},
					Photo:                     ctx.Message.Photo,
					Caption:                   ctx.Message.Caption,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtPhoto: %w", err))
					return
				}
			case ContextTypeAtVoice:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtVoiceFunc)
				c := &MessageContextAtVoice{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},
					Voice:                     ctx.Message.Voice,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtVoice: %w", err))
					return
				}
			case ContextTypeAtAudio:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtAudioFunc)
				c := &MessageContextAtAudio{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Audio:   ctx.Message.Audio,
					Caption: ctx.Message.Caption,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtAudio: %w", err))
					return
				}
			case ContextTypeAtVideo:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtVideoFunc)
				c := &MessageContextAtVideo{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Video:   ctx.Message.Video,
					Caption: ctx.Message.Caption,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtVideo: %w", err))
					return
				}
			case ContextTypeAtAnimation:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtAnimationFunc)
				c := &MessageContextAtAnimation{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Animation: ctx.Message.Animation,
					Document:  ctx.Message.Document,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtAnimation: %w", err))
					return
				}
			case ContextTypeAtDocument:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtDocumentFunc)
				c := &MessageContextAtDocument{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Document: ctx.Message.Document,
					Caption:  ctx.Message.Caption,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtDocument: %w", err))
					return
				}
			case ContextTypeAtSticker:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtStickerFunc)
				c := &MessageContextAtSticker{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Sticker: ctx.Message.Sticker,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtSticker: %w", err))
					return
				}
			case ContextTypeAtVideoNote:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtVideoNoteFunc)
				c := &MessageContextAtVideoNote{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					VideoNote: ctx.Message.VideoNote,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtVideoNote: %w", err))
					return
				}
			case ContextTypeAtContact:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtContactFunc)
				c := &MessageContextAtContact{
					MessageContextBase: base,

					Contact: ctx.Message.Contact,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtContact: %w", err))
					return
				}
			case ContextTypeAtDice:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtDiceFunc)
				c := &MessageContextAtDice{
					MessageContextBase: base,

					Dice: ctx.Message.Dice,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtDice: %w", err))
					return
				}
			case ContextTypeAtGame:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtGameFunc)
				c := &MessageContextAtGame{
					MessageContextIncludeFile: MessageContextIncludeFile{MessageContextBase: base},

					Game: ctx.Message.Game,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtGame: %w", err))
					return
				}
			case ContextTypeAtPoll:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtPollFunc)
				c := &MessageContextAtPoll{
					MessageContextBase: base,

					Poll: ctx.Message.Poll,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtPoll: %w", err))
					return
				}
			case ContextTypeAtVenue:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtVenueFunc)
				c := &MessageContextAtVenue{
					MessageContextBase: base,

					Venue: ctx.Message.Venue,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtVenue: %w", err))
					return
				}
			case ContextTypeAtLocation:
				fn := b.specifiedTypeMessageProcessorFunc[ctx.MessageType].(ProcessorAtLocationFunc)
				c := &MessageContextAtLocation{
					MessageContextBase: base,

					Location: ctx.Message.Location,
				}
				if err := fn(c); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcessorAtLocation: %w", err))
					return
				}
			}
		} else {
			if b.defaultMessageProcessorFunc != nil { // 消息处理器
				if err := b.defaultMessageProcessorFunc(ctx); err != nil {
					b.handleError(fmt.Errorf("SetMessageProcesso: %w", err))
					return
				}
			}
		}
	}()
}
