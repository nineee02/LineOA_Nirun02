package event

// import (
// 	"database/sql"
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

// var usercardidState = make(map[string]string)
// var userState = make(map[string]string)

// // HandleEvent - จัดการข้อความที่ได้รับจาก LINE
// func HandleEvent(bot *linebot.Client, event *linebot.Event) {
// 	text := event.Message.(*linebot.TextMessage).Text
// 	userID := event.Source.UserID
// 	log.Println("userID:", userID)

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
// 		case "wait status Activityrecord":
// 			handleActivityrecord(bot, event, userID)
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
// 		log.Println("Event is not a text message (handleElderlyInfo)")
// 		return
// 	}

// 	// ตรวจสอบชื่อผู้ป่วยที่ได้รับ
// 	Name := strings.TrimSpace(message.Text)
// 	log.Println("Received Card ID of Patient_info:", Name)

// 	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		// sendErrorReply(bot, event, "Unable to connect to the database. Please try again later.")
// 		return
// 	}
// 	defer db.Close()

// 	// ค้นหาข้อมูลจากฐานข้อมูล
// 	patient, err := GetPatientInfoByName(db, Name)
// 	if err != nil {
// 		log.Println("Error fetching patient info:", err)
// 		// sendErrorReply(bot, event, "No patient information found for the provided name.")
// 		return
// 	}

// 	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	replyMessage := FormatPatientInfo(patient)
// 	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
// 		log.Println("Error replying message:(handleElderlyInfo)", err)
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
// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message(handleServiceRecord):")
// 		return
// 	}

// 	// ตรวจสอบชื่อผู้ป่วยที่ได้รับ
// 	cardid := strings.TrimSpace(message.Text)
// 	log.Println("Received card id :", cardid)

// 	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		// sendErrorReply(bot, event, "Unable to connect to the database. Please try again later.")
// 		return
// 	}
// 	defer db.Close()

// 	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
// 	service, err := GetServiceInfoBycardID(db, cardid)
// 	if err != nil {
// 		log.Println("Error models.GetServiceInfoByName:", err)
// 		// sendErrorReply(bot, event, "No patient information found for the provided name.")
// 		return
// 	}

// 	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	replyMessage := FormatServiceInfo(service)
// 	quickReplyActivities := createQuickReplyActivities()

// 	if _, err := bot.ReplyMessage(
// 		event.ReplyToken,
// 		linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
// 	).Do(); err != nil {
// 		log.Println("Error replying message:", err)
// 	}

// 	usercardidState[userID] = cardid
// 	log.Printf("Saved card_id for user %s: %s", userID, cardid)

// 	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
// 	userState[userID] = "wait status Activityrecord"
// 	log.Printf("Set user state to wait status Activityrecord for user %s", userID)
// }

// func handleActivityrecord(bot *linebot.Client, event *linebot.Event, userID string) {
// 	// ตรวจสอบว่าผู้ใช้อยู่ในสถานะที่ถูกต้องหรือไม่
// 	if userState[userID] != "wait status Activityrecord" {
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

// 	// ดึง card_id จาก userState
// 	cardID := usercardidState[userID] // ปรับใช้ให้ตรงตามข้อมูลที่เก็บไว้

// 	// สร้างตัวแปร activityRecord ที่มีข้อมูลของ card_id และ activity
// 	activityRecord := &models.Activityrecord{
// 		CardID:   cardID, // ใช้ = แทน := เพราะเป็นการกำหนดค่าตัวแปรที่มีอยู่แล้ว
// 		Activity: activity,
// 	}

// 	// แสดง log สำหรับ card_id และ activity
// 	log.Printf("Using card_id: %s, Activity: %s", cardID, activity)

// 	// เชื่อมต่อกับฐานข้อมูล
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 		return
// 	}
// 	defer db.Close()

// 	// บันทึกกิจกรรมใหม่ models.ActivityRecord(db, activityRecord)
// 	if err := ActivityRecord(db, activityRecord); err != nil { // ส่ง activityRecord ไปที่ฟังก์ชัน ActivityRecord
// 		log.Printf("Error saving activity(models.ActivityRecord): %v", err)
// 		sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", activity))
// 		return
// 	}

// 	//sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรม '%s' สำเร็จสำหรับ card_id '%s'!", activity, cardID))

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

// func createQuickReplyActivities() linebot.QuickReplyItems {
// 	// รายการกิจกรรมทั้งหมด
// 	activities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// 	}

// 	quickReplyItems := linebot.QuickReplyItems{}
// 	for _, activity := range activities {
// 		quickReplyItems.Items = append(quickReplyItems.Items,
// 			linebot.NewQuickReplyButton("", linebot.NewMessageAction(activity, activity)),
// 		)
// 	}
// 	return quickReplyItems
// }

// func SaveActivity(db *sql.DB, activity string) error {
// 	if !validateActivity(activity) {
// 		return fmt.Errorf("กิจกรรม '%s' ไม่ตรงกับค่าที่อนุญาตในฐานข้อมูล", activity)
// 	}

// 	query := `INSERT INTO service_info (activity) VALUES (?)`
// 	_, err := db.Exec(query, activity)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม %s ได้: %v", activity, err)
// 	}
// 	return nil
// }

// // validateActivity ตรวจสอบว่าค่าที่ส่งมาตรงกับฐานข้อมูล หรือไม่
// func validateActivity(activity string) bool {
// 	allowedActivities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// 	}
// 	for _, allowed := range allowedActivities {
// 		if activity == allowed {
// 			return true
// 		}
// 	}
// 	return false
// }
