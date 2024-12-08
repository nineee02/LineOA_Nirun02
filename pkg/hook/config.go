package hook

import (
	"log"
	"net/http"

	"nirun/pkg/event"
	linebotConfig "nirun/pkg/linebot"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

// HandleLineWebhook - จัดการ Webhook จาก LINE
func HandleLineWebhook(c *gin.Context) {
	bot := linebotConfig.GetLineBot()
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(http.StatusBadRequest)
		} else {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, lineEvent := range events {
		if lineEvent.Type == linebot.EventTypeMessage {
			message, ok := lineEvent.Message.(*linebot.TextMessage)
			if ok {
				log.Printf("Received message: %s", message.Text)
				event.HandleEvent(bot, lineEvent)
			}
		}
	}
	c.Writer.WriteHeader(http.StatusOK)
	log.Println("Webhook response sent with status 200")
}
