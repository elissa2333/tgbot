package telegram

// pack

/*Games
您的机器人可以为用户提供 HTML5 游戏，使他们可以独奏或在小组和一对一聊天中相互竞争。使用 /newgame 命令通过 @BotFather 创建游戏。请注意，这种功能需要承担责任：您将需要接受机器人将提供的每款游戏的条款。
https://core.telegram.org/bots/api#games*/

//Game 代表一个游戏。使用BotFather创建和编辑游戏，它们的简称将用作唯一的标识符
//https://core.telegram.org/bots/api#game
type Game struct {
	Title        string          `json:"title"`         // 游戏标题
	Description  string          `json:"description"`   // 游戏说明
	Photo        []PhotoSize     `json:"photo"`         // 将在聊天中的游戏消息中显示的照片。
	Text         string          `json:"text"`          // 可选的。对游戏的简要说明或游戏消息中包含的高分。当漫游器调用setGameScore时，可以自动编辑以包括当前游戏的高分，或使用editMessageText进行手动编辑。0-4096个字符。
	TextEntities []MessageEntity `json:"text_entities"` // 可选的。以文本形式出现的特殊实体，例如用户名，URL，机器人命令等。
	Animation    *Animation      `json:"animation"`     // 可选的。将在聊天中显示在游戏消息中的动画。通过BotFather上传
}

// CallbackGame 占位符，当前不保存任何信息。使用BotFather来设置您的游戏
// https://core.telegram.org/bots/api#callbackgame
type CallbackGame struct{}
