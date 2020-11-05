package tgbot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	stdURL "net/url"

	"github.com/elissa2333/httpc"

	"github.com/elissa2333/tgbot/telegram"
)

// MessageProcessorFunc 聊天处理器
type MessageProcessorFunc func(c *Context) error

// ActiveProcessorFunc 主动处理器函数
type ActiveProcessorFunc func(api *telegram.API) error

// Bot bot 实体
type Bot struct {
	api     *telegram.API // telegram API
	timeout uint          // 长时间轮询的超时时间（以秒为单位）为0即通常的短轮询。应该为正，短轮询应仅用于测试目的

	webHookEngine func() error

	MsgOffset int64 // 最后一条消息

	activeProcessorFunc ActiveProcessorFunc

	commands       map[string]MessageProcessorFunc // 指定命令的执行方法
	defaultCommand MessageProcessorFunc            // 默认命令未指定命令时使用

	messageProcessorFunc MessageProcessorFunc

	done chan struct{} // 退出程序
	err  chan error
}

// BotOptional bot 配置可选参数
type BotOptional struct {
	HTTPClient *http.Client
	Timeout    uint // // 长时间轮询的超时时间（以秒为单位）为0即通常的短轮询。应该为正，短轮询应仅用于测试目的
}

// New 新建 bot
func New(id int, token string, optional *BotOptional) *Bot {
	b := &Bot{
		api:      telegram.New(nil, id, token),
		timeout:  15,
		commands: map[string]MessageProcessorFunc{},
		done:     make(chan struct{}),
		err:      make(chan error),
	}

	if optional != nil {
		b.timeout = optional.Timeout

		if optional.HTTPClient != nil {
			b.api = telegram.New(optional.HTTPClient, id, token)
		}
	}

	return b
}

// SetWebhook 设置 webhook
func (b *Bot) SetWebhook(url string /*api 访问地址*/, address string /*本地监听地址*/, optional *telegram.OptionalWebhook) error {
	parseURL, err := stdURL.Parse(url)
	if err != nil {
		return err
	}

	if err := b.api.SetWebhook(url, optional); err != nil {
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
				b.err <- err
				return
			}

			m := telegram.Update{}
			if err := json.Unmarshal(bodyB, &m); err != nil {
				b.err <- err
				return
			}

			if m.UpdateID == 0 || m.Message == nil { // 坏请求
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			b.handleReceivedMessages(m.Message)
		})

		return http.ListenAndServe(address, nil)
	}

	return nil
}

// DeleteWebhook  删除 webhook
func (b *Bot) DeleteWebhook() error {
	b.webHookEngine = nil
	return b.api.DeleteWebhook()
}

// GetWebhookInfo 获取 webhook 信息
func (b *Bot) GetWebhookInfo() (*telegram.WebhookInfo, error) {
	return b.api.GetWebhookInfo()
}

// handleError 处理错误
func (b *Bot) handleError(err error) {
	if err != nil {
		b.err <- err
	}
}

// initiativeEngine 核心调度
func (b *Bot) initiativeEngine() {
	TotalNumberOfActiveAndPassive := 0
	cleanActiveAndPassiveCh := make(chan struct{})
	if b.activeProcessorFunc != nil { // 激活主动处理器
		TotalNumberOfActiveAndPassive++ // 统计被动
		go func() {
			if err := b.activeProcessorFunc(b.api); err != nil {
				b.handleError(err)
			}
			cleanActiveAndPassiveCh <- struct{}{}
		}()
	}

	if len(b.commands) != 0 || b.defaultCommand != nil || b.messageProcessorFunc != nil { // 统计被动
		TotalNumberOfActiveAndPassive++
	}

	go func() {
		num := 0
		for range cleanActiveAndPassiveCh {
			num++
			if num >= TotalNumberOfActiveAndPassive {
				b.done <- struct{}{}
			}
		}
	}()

	for {
		updates, err := b.api.GetUpdates(b.MsgOffset, 1, b.timeout)
		if err != nil {
			b.handleError(err)
			break
		}
		if len(updates) == 0 {
			continue
		}

		LastOneUpdate := updates[len(updates)-1]
		b.MsgOffset = LastOneUpdate.UpdateID + 1 // 记录消息偏量

		b.handleReceivedMessages(updates[0].Message)
	}
}

// handleReceivedMessages 处理接收消息
func (b *Bot) handleReceivedMessages(message *telegram.Message) {
	ctx := &Context{
		Message: message,
		API:     b.api,
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

		if b.messageProcessorFunc != nil { // 消息处理器
			if err := b.messageProcessorFunc(ctx); err != nil {
				b.handleError(err)
			}
		}
	}()
}

// AddCommandProcessor 添加命令处理器（接收到命令后调用）
func (b *Bot) AddCommandProcessor(cmd string, execFunc MessageProcessorFunc) {
	b.commands[cmd] = execFunc
}

// SetDefaultCommandProcessor 设置默认命令处理器（在未找到命令时调用）
func (b *Bot) SetDefaultCommandProcessor(execFunc MessageProcessorFunc) {
	b.defaultCommand = execFunc
}

// SetMessageProcessor 消息处理器（接收到消息后调用）
func (b *Bot) SetMessageProcessor(handleMessageFunc MessageProcessorFunc) {
	b.messageProcessorFunc = handleMessageFunc
}

// SetActiveProcessor 设置主动处理器
func (b *Bot) SetActiveProcessor(activeProcessorFunc ActiveProcessorFunc) {
	b.activeProcessorFunc = activeProcessorFunc
}

// Run 运行 bot
func (b *Bot) Run() error {
	_, err := b.api.GetMe() // check api
	if err != nil {
		return err
	}

	if (b.webHookEngine) != nil { // 为了和主动处理器行为一致
		go func() {
			if err := b.webHookEngine(); err != nil {
				b.err <- err
			}
		}()
	} else {
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
