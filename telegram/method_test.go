package telegram

import (
	"fmt"
	"os"
	"testing"

	"github.com/elissa2333/tgbot/utils"
)

//go:generate go test -v -test.run TestAPI_GetMe
func TestAPI_GetMe(t *testing.T) {
	user, err := tAPI.GetMe()
	if err != nil {
		t.Fatal(err)
	}

	if utils.ToInt(user.ID) != tID {
		t.Fatal("userID != setting id")
	}
}

func TestAPI_LogOut(t *testing.T) {
	// TODO
}

func TestAPI_Close(t *testing.T) {
	// TODO
}

//go:generate go test -v -test.run TestAPI_SendMessage
func TestAPI_SendMessage(t *testing.T) {
	// 你会收到一条消息
	content := "niconiconi"
	msg, err := tAPI.SendMessage(tChatID, content, nil)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Text != content {
		t.Fatal("回应内容与发送内容不一致")
	}
}

//go:generate go test -v -test.run TestAPI_ForwardMessage
func TestAPI_ForwardMessage(t *testing.T) {
	// 你的 bot 会收到两条消息 一条原始消息 一条转发消息
	msg, err := tAPI.SendMessage(tChatID, "forward test", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tAPI.ForwardMessage(tChatID, utils.ToString(msg.Chat.ID), msg.MessageID, nil)
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_CopyMessage
func TestAPI_CopyMessage(t *testing.T) {
	// 两条完全一致的消息
	msg, err := tAPI.SendMessage(tChatID, "copy", nil)
	if err != nil {
		t.Fatal(err)
	}
	messageID, err := tAPI.CopyMessage(tChatID, utils.ToString(msg.Chat.ID), msg.MessageID, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(messageID)

	if messageID == 0 {
		t.Fatal(" 消息拷贝失败")
	}
}

//go:generate go test -v -test.run TestAPI_SendPhoto
func TestAPI_SendPhoto(t *testing.T) {
	// 收到一张图片

	f, err := os.Open("./testdata/eso1907a.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = tAPI.SendPhoto(tChatID, f, nil)
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_SendAudio
func TestAPI_SendAudio(t *testing.T) {
	f, err := os.Open("./testdata/the-wires.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = tAPI.SendAudio(tChatID, f, nil)
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_SendDocument
func TestAPI_SendDocument(t *testing.T) {
	f, err := os.Open("./testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = tAPI.SendDocument(tChatID, f, nil)
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_SendPoll
func TestAPI_SendPoll(t *testing.T) {
	_, err := tAPI.SendPoll(tChatID, "select a and b", []string{"a", "b"}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
