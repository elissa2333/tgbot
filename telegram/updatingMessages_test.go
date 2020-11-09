package telegram

import (
	"os"
	"testing"
	"time"

	"github.com/elissa2333/tgbot/utils"
)

//go:generate go test -v -test.run TestAPI_EditMessageText
func TestAPI_EditMessageText(t *testing.T) {
	msg, err := tAPI.SendMessage(tChatID, "original", nil)
	if err != nil {
		t.Error(err)
	}

	_, err = tAPI.EditMessageText("modify", EditMessageTextOptional{ChatID: tChatID, MessageID: msg.MessageID})
	if err != nil {
		t.Error(err)
	}
}

//go:generate go test -v -test.run TestAPI_EditMessageCaption
func TestAPI_EditMessageCaption(t *testing.T) {
	file, err := os.Open("./testdata/eso1907a.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	msg, err := tAPI.SendPhoto(tChatID, file, &SendPhotoOptional{Caption: "original"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tAPI.EditMessageCaption(EditMessageCaptionOptional{ChatID: tChatID, MessageID: msg.MessageID})
	if err != nil {
		t.Error(err)
	}
}

//go:generate go test -v -test.run TestAPI_EditMessageMedia
func TestAPI_EditMessageMedia(t *testing.T) {
	file, err := os.Open("./testdata/eso1907a.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	msg, err := tAPI.SendPhoto(tChatID, file, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	_, err = tAPI.EditMessageMedia(InputMedia{
		InputMediaPhoto: &InputMediaPhoto{
			Type:  InputMediaPhotoType,
			Media: "https://http.cat/302",
		},
	}, EditMessageMediaOptional{
		ChatID:    utils.ToString(msg.Chat.ID),
		MessageID: msg.MessageID,
	})
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_StopPoll
func TestAPI_StopPoll(t *testing.T) {
	msg, err := tAPI.SendPoll(tChatID, "select a and b", []string{"a", "b"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tAPI.StopPoll(utils.ToString(msg.Chat.ID), msg.MessageID, nil)
	if err != nil {
		t.Fatal(err)
	}
}

//go:generate go test -v -test.run TestAPI_DeleteMessage
func TestAPI_DeleteMessage(t *testing.T) {
	msg, err := tAPI.SendMessage(tChatID, "delete after three seconds", nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	_, err = tAPI.DeleteMessage(tChatID, msg.MessageID)
	if err != nil {
		t.Fatal(err)
	}
}
