package event

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

// func ReplyTextMessage(bot *linebot.Client, replyToken string, textMessage string) {
// 	message := linebot.NewTextMessage(textMessage)
// 	_, err := bot.ReplyMessage(replyToken, message).Do()
// 	if err != nil {
// 		log.Println("ErrorReplyMessage:") // ถ้ามีข้อผิดพลาด
// 	}
// 	log.Printf("ReplyTextMessage success!!!")
// }

func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
	).Do(); err != nil {
		log.Println("ReplyErrorFormat:", err)
	}
}

// ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
		log.Println("ReplyErrorFormat:", err)
	}
}
