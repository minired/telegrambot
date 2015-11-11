package main

import (
	"container/list"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"time"
)

type ReceivedMessage struct {
	message         tgbotapi.MessageConfig
	replyChatIdList *list.List
}

var (
	MessageList list.List
	MessageChan = make(chan tgbotapi.Message)
)

func AddNewMessage(message tgbotapi.MessageConfig) {
	MessageList.PushBack(ReceivedMessage{message, new(list.List)})
}
func (r ReceivedMessage) IsSentMessage(ChatId int) bool {
	x := r.replyChatIdList
	for e := x.Front(); e != nil; e = e.Next() {
		b := e.Value.(int)
		if b == ChatId {
			return true
		}
	}
	return false

}

func GetReplyMessage(ChatId int) tgbotapi.MessageConfig {
	for e := MessageList.Front(); e != nil; e = e.Next() {
		b := e.Value.(ReceivedMessage)
		if b.message.ChatID != ChatId && b.IsSentMessage(ChatId) != true {
			b.replyChatIdList.PushBack(ChatId)
			return b.message
		}
	}

	var x tgbotapi.MessageConfig
	x.ChatID = -1
	return x
}
func SendReplyMessage(bot *tgbotapi.BotAPI) {
	for {
		updatemsg := <-MessageChan
		msg := tgbotapi.NewMessage(updatemsg.Chat.ID, updatemsg.From.LastName+updatemsg.From.FirstName+": "+updatemsg.Text)
		AddNewMessage(msg)
		replyMsg := GetReplyMessage(updatemsg.Chat.ID)
		if replyMsg.ChatID != -1 {
			bot.SendMessage(replyMsg)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
func main() {

	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_API_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	err = bot.UpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	go SendReplyMessage(bot)

	for update := range bot.Updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		MessageChan <- update.Message
		time.Sleep(100 * time.Millisecond)
	}
}
