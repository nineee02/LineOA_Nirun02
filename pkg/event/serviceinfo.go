package event

// import (
// 	"log"
// 	"nirun/pkg/database"

// 	"nirun/pkg/hook"

// 	"github.com/line/line-bot-sdk-go/linebot"
// )

// // handleServiceData จัดการการบันทึกข้อมูลการเข้ารับบริการ
// func handleServiceData(bot *linebot.Client, event *linebot.Event, text string) {
// 	//name_ := strings.TrimSpace(strings.TrimPrefix(text, "บันทึกการเข้ารับบริการ"))

// 	// บันทึกข้อมูลลงในฐานข้อมูล
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		ReplyDatabaseError(bot, event.ReplyToken)
// 		return
// 	}
// 	defer db.Close()

// }
