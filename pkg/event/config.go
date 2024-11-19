package event

import (
	"fmt"
	"log"
	"nirun/pkg/database"              // แก้ไขตาม path ที่ถูกต้องของ database package
	linebotConfig "nirun/pkg/linebot" // แก้ไขตาม path ที่ถูกต้องของ linebotConfig package
	"nirun/pkg/models"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func HandleEvent(event *linebot.Event) {
	bot := linebotConfig.GetLineBot()

	if event.Type != linebot.EventTypeMessage {
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return
	}

	text := message.Text
	log.Println("Received message:", text)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		models.ReplyDatabaseError(bot, event.ReplyToken)
		return
	}
	defer db.Close()

	if strings.HasPrefix(text, "ข้อมูลผู้ป่วย") {
		// ค้นหาผู้ป่วย
		name := strings.TrimSpace(strings.TrimPrefix(text, "ข้อมูลผู้ป่วย "))
		log.Println("Searching for patient with name:", name)

		patientInfo, err := models.GetPatientInfoByName(db, name)
		if err != nil {
			log.Println("Error fetching patient info:", err)
			models.ReplyDataNotFound(bot, event.ReplyToken)
			return
		}

		// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
		replyMessage := fmt.Sprintf(
			"ข้อมูลผู้ป่วย:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nเบอร์โทรศัพท์: %s",
			patientInfo.Name, patientInfo.PatientID, patientInfo.Age, patientInfo.Sex, patientInfo.Blood, patientInfo.PhoneNumber,
		)
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Println("Error replying patient info message:", err)
		}
	}
}
