package event

// import (
// 	"fmt"
// 	"log"
// 	"nirun/pkg/database"
// 	"nirun/pkg/models"
// 	"strings"

// 	"github.com/line/line-bot-sdk-go/linebot"
// )

// func sendCustomReply(bot *linebot.Client, replyToken string, messages ...linebot.SendingMessage) {
// 	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
// 		log.Printf("Error replying message: %v", err)
// 	}
// }

// var userState = make(map[string]string)

// // HandleEvent - จัดการข้อความที่ได้รับจาก LINE
// func HandleEvent(bot *linebot.Client, event *linebot.Event) {
// 	text := event.Message.(*linebot.TextMessage).Text
// 	userID := event.Source.UserID
// 	log.Println("Error replying message:", userID)

// 	switch text {
// 	case "NIRUN":
// 		handleNIRUN(bot, event)
// 	case "ข้อมูลผู้สูงอายุ":
// 		handleElderlyInfoRequest(bot, event, event.Source.UserID)
// 	case "ลงเวลาการทำงานสำหรับเจ้าหน้าที":
// 		handleWorkTime(bot, event)
// 	case "ประวัติการเข้ารับบริการ":
// 		handleServiceHistory(bot, event)
// 	case "บันทึกการเข้ารับบริการ":
// 		handleServiceRecordRequest(bot, event, event.Source.UserID)
// 	case "คู่มือการใช้งานระบบ":
// 		handleSystemManual(bot, event)
// 	default:
// 		handleDefault(bot, event)
// 	}

// 	state, exists := userState[userID]
// 	if exists {
// 		switch state {
// 		case "wait status ElderlyInfoRequest":
// 			handleElderlyInfo(bot, event, userID)
// 		case "wait status ServiceRecordRequest":
// 			handleServiceRecord(bot, event, userID)
// 		case "wait status SaveInput":
// 			handleActivityInput(bot, event, userID)
// 		default:
// 			log.Printf("Unhandled state for user %s: %s", userID, state)
// 		}
// 		return
// 	}

// }
// func setUserState(userID, state string) {
// 	userState[userID] = state
// 	log.Printf("Set user state for user %s to %s", userID, state)
// }

// func handleElderlyInfoRequest(bot *linebot.Client, event *linebot.Event, userID string) {
// 	setUserState(userID, "wait status ElderlyInfoRequest")
// }

// func handleServiceRecordRequest(bot *linebot.Client, event *linebot.Event, userID string) {
// 	setUserState(userID, "wait status ServiceRecordRequest")
// }

// // ************************************************************************************************************************

// // **********************************************************************************************************
// func handleNIRUN(bot *linebot.Client, event *linebot.Event) {
// 	sendReply(bot, event.ReplyToken, "ยินดีต้อนรับสู่ระบบ NIRUN! กรุณาเลือกเมนูที่ต้องการ.")
// }

// func handleElderlyInfo(bot *linebot.Client, event *linebot.Event, userID string) {
// 	// รับข้อความจากผู้ใช้ (ชื่อผู้ป่วย)
// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	// ตรวจสอบชื่อผู้ป่วยที่ได้รับ
// 	patientName := strings.TrimSpace(message.Text)
// 	log.Println("Received patient name:", patientName)

// 	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		// sendErrorReply(bot, event, "Unable to connect to the database. Please try again later.")
// 		return
// 	}
// 	defer db.Close()

// 	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
// 	patientInfo, err := models.GetPatientInfoByName(db, patientName)
// 	if err != nil {
// 		log.Println("Error fetching patient info:", err)
// 		// sendErrorReply(bot, event, "No patient information found for the provided name.")
// 		return
// 	}

// 	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	replyMessage := models.FormatPatientInfo(patientInfo)
// 	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
// 		log.Println("Error replying message:", err)
// 	}

// 	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
// 	userState[userID] = "" // เปลี่ยนสถานะเพื่อให้พร้อมรับข้อมูลใหม่
// }

// func handleWorkTime(bot *linebot.Client, event *linebot.Event) {
// 	sendReply(bot, event.ReplyToken, "กรุณาลงเวลาการทำงานโดยกรอกชื่อและเวลาที่ต้องการ:")
// }

// func handleServiceHistory(bot *linebot.Client, event *linebot.Event) {
// 	sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อเพื่อดูประวัติการเข้ารับบริการ:")
// }

// func handleServiceRecord(bot *linebot.Client, event *linebot.Event, userID string) {
// 	//sendReply(bot, event.ReplyToken, "กรุณากรอกข้อมูลเพื่อบันทึกการเข้ารับบริการ:")
// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	// ตรวจสอบชื่อผู้ป่วยที่ได้รับ
// 	patientName := strings.TrimSpace(message.Text)
// 	log.Println("Received patient name:", patientName)

// 	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		// sendErrorReply(bot, event, "Unable to connect to the database. Please try again later.")
// 		return
// 	}
// 	defer db.Close()

// 	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
// 	serviceInfo, err := models.GetServiceInfoByName(db, patientName)
// 	if err != nil {
// 		log.Println("Error fetching patient info:", err)
// 		// sendErrorReply(bot, event, "No patient information found for the provided name.")
// 		return
// 	}

// 	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	replyMessage := models.FormatServiceInfo(serviceInfo)
// 	quickReplyActivities := createQuickReplyActivities()

// 	if _, err := bot.ReplyMessage(
// 		event.ReplyToken,
// 		linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
// 	).Do(); err != nil {
// 		log.Println("Error replying message:", err)
// 	}

// 	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
// 	userState[userID] = "wait status SaveInput"
// 	log.Printf("Set user state to wait status SaveInput for user %s", userID)
// }

// func handleActivityInput(bot *linebot.Client, event *linebot.Event, userID string) {
// 	// ตรวจสอบว่าผู้ใช้อยู่ในสถานะที่ถูกต้องหรือไม่
// 	if userState[userID] != "wait status SaveInput" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
// 		return
// 	}

// 	// รับข้อความกิจกรรมจากผู้ใช้
// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	activity := strings.TrimSpace(message.Text)
// 	log.Printf("Received activity input: %s", activity)

// 	if !validateActivity(activity) {
// 		sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
// 		return
// 	}
// 	// เชื่อมต่อกับฐานข้อมูล
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 		return
// 	}
// 	defer db.Close()

// 	// บันทึกกิจกรรมใหม่
// 	if err := SaveActivity(db, activity); err != nil {
// 		log.Printf("Error saving activity: %v", err)
// 		sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", activity))
// 		return
// 	}

// 	sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรม '%s' สำเร็จ!", activity))

// 	// รีเซ็ตสถานะผู้ใช้
// 	userState[userID] = ""
// 	log.Printf("Reset user state for user %s", userID)
// }

// func handleSystemManual(bot *linebot.Client, event *linebot.Event) {
// 	sendReply(bot, event.ReplyToken, "คุณสามารถดูคู่มือการใช้งานระบบได้ที่ลิงก์: https://example.com/manual")
// }

// func handleDefault(bot *linebot.Client, event *linebot.Event) {
// 	sendCustomReply(bot, event.ReplyToken)
// }

// // ฟังก์ชันสำหรับส่งข้อความ
// func sendReply(bot *linebot.Client, replyToken, message string) {
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
// 		log.Printf("Error replying message: %v", err)
// 	}
// }
