package hook

import (
	"log"
	"net/http"
	"strings"

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

	for _, lineevent := range events {
		if lineevent.Type == linebot.EventTypeMessage {
			message, ok := lineevent.Message.(*linebot.TextMessage) //lineevent.Message ข้อความที่ส่งมาจากผู้ใช้
			if ok {
				text := strings.TrimSpace(message.Text)
				log.Printf("Received message: %s,", text)

				// ใช้ฟังก์ชันใน pkg/event เพื่อจัดการข้อความ
				event.HandleEvent(bot, lineevent)
				// }
			}
		}
		c.Writer.WriteHeader(http.StatusOK)
		log.Println("Webhook response sent with status 200")
	}
}
