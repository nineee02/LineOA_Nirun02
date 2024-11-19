package event

import (
	"database/sql"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	// 	"fmt"
	// 	"log"
	// 	"nirun/pkg/database"
	// 	"nirun/pkg/hook"
	// 	"nirun/pkg/models"
	// 	"strings"
	// 	"github.com/line/line-bot-sdk-go/linebot"
)

// / handlePatientInfo ฟังก์ชันจัดการการดึงและส่งข้อมูลผู้ป่วยผ่าน Line Bot
// public
func HandlePatientInfo(bot *linebot.Client, replyToken string, db *sql.DB, name_ string) {
	data, err := GetPatientInfoByName(db, name_) // ใช้ฟังก์ชันจาก patient_info
	if err != nil {
		log.Println("ข้อผิดพลาดในการดึงข้อมูล: ", err)
		return
	}

	replyMessage := FormatPatientInfo(data)
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("ข้อผิดพลาดในการตอบกลับ:", err)
	}
}
