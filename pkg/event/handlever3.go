package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/models"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func sendCustomReply(bot *linebot.Client, replyToken string, messages ...linebot.SendingMessage) {
	if len(messages) == 0 {
		return // ไม่ส่งข้อความถ้าไม่มี
	}

	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
		log.Printf("Error replying message sendCustomReply: %v", err)
	}
}

var usercardidState = make(map[string]string)
var userState = make(map[string]string)
var userActivity = make(map[string]string) // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน

// HandleEvent - จัดการข้อความที่ได้รับจาก LINE
func HandleEvent(bot *linebot.Client, event *linebot.Event) {
	text := event.Message.(*linebot.TextMessage).Text
	log.Println("Text: ", text)
	State := event.Source.UserID
	// log.Println("userstate:", State)

	switch text {
	case "NIRUN":
		handleNIRUN(bot, event)
	case "ข้อมูลผู้สูงอายุ":
		handleElderlyInfoRequest(bot, event, event.Source.UserID)
	case "ลงเวลาการทำงานสำหรับเจ้าหน้าที":
		// handleWorkTime(bot, event)
	case "ประวัติการเข้ารับบริการ":
		handleServiceHistory(bot, event)
	case "บันทึกการเข้ารับบริการ":
		handleServiceRecordRequest(bot, event, event.Source.UserID)
	case "คู่มือการใช้งานระบบ":
		handleSystemManual(bot, event)
	default:
		handleDefault(bot, event)
	}

	state, exists := userState[State]
	if exists {
		switch state {
		case "wait status ElderlyInfoRequest":
			handlePateintInfo(bot, event, State)
		case "wait status ServiceRecordRequest":
			handleServiceInfo(bot, event, State)
		case "wait status Activityrecord":
			handleActivityrecord(bot, event, State)
		case "wait status ActivityPeriodRecord":
			handleActivityPeriodRecord(bot, event, State)
		default:
			log.Printf("Unhandled state for user %s: %s", State, state)
		}
		return
	}

}
func setUserState(State, state string) {
	userState[State] = state
	// log.Printf("Set user state for user %s to %s", State, state)

}

func handleElderlyInfoRequest(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status ElderlyInfoRequest")

}

func handleServiceRecordRequest(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status ServiceRecordRequest")
}

// ************************************************************************************************************************

// **********************************************************************************************************
func handleNIRUN(bot *linebot.Client, event *linebot.Event) {
	sendReply(bot, event.ReplyToken, "ยินดีต้อนรับสู่ระบบ NIRUN! กรุณาเลือกเมนูที่ต้องการ.")
}

func handlePateintInfo(bot *linebot.Client, event *linebot.Event, State string) {
	message := event.Message.(*linebot.TextMessage).Text
	log.Println("Message pateint:", message)

	if message == "ข้อมูลผู้สูงอายุ" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน :")
		return
	}

	cardID := strings.TrimSpace(message)
	if cardID == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน:")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ค้นหาข้อมูลจากฐานข้อมูล
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		// แทนที่จะส่ง error ในกรณีไม่พบข้อมูลผู้ป่วย, ให้ส่งข้อความที่เหมาะสม
		if err == sql.ErrNoRows {
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วยที่ท่านค้นหา กรุณาตรวจสอบเลขประจำตัวผู้ป่วยอีกครั้ง.")
		} else {
			log.Println("Error fetching patient info:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ป่วย กรุณาลองใหม่.")
		}
		return
	}

	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
	replyMessage := FormatPatientInfo(patient)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message:(handleElderlyInfo)", err)
	}
	log.Println("ข้อมูลผู้สูงอายุ :", replyMessage)

	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
	userState[State] = "" // เปลี่ยนสถานะเพื่อให้พร้อมรับข้อมูลใหม่
}

func handleServiceHistory(bot *linebot.Client, event *linebot.Event) {
	sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชนเพื่อดูประวัติการเข้ารับบริการ:")
}

func handleServiceInfo(bot *linebot.Client, event *linebot.Event, State string) {
	message := event.Message.(*linebot.TextMessage).Text
	// log.Println("Message pateint:", message)

	if message == "บันทึกการเข้ารับบริการ" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน :")
		return
	}

	cardID := strings.TrimSpace(message)
	if cardID == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน:")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
	service, err := GetServiceInfoBycardID(db, cardID)
	if err != nil {
		log.Println("Error models.GetServiceInfoBycardID:", err)
		// sendErrorReply(bot, event, "No patient information found for the provided name.")
		return
	}

	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
	replyMessage := FormatServiceInfo(service)
	// log.Println("reply Message Format: ", replyMessage)
	quickReplyActivities := createQuickReplyActivities()

	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
	).Do(); err != nil {
		log.Println("Error replying message:", err)
	}

	usercardidState[State] = cardID
	// log.Printf("Saved card_id for user %s: %s", State, cardID)

	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
	userState[State] = "wait status Activityrecord"
	//log.Printf("Set user state to wait status Activityrecord for user %s", State)
}
func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
	// ตรวจสอบว่าผู้ใช้อยู่ในสถานะที่ถูกต้องหรือไม่
	if userState[State] != "wait status Activityrecord" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
		return
	}

	// รับข้อความกิจกรรมจากผู้ใช้
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	activity := strings.TrimSpace(message.Text)
	log.Printf("Received activity input: %s", activity)

	if !validateActivity(activity) {
		sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
		return
	}

	userActivity[State] = activity // เก็บกิจกรรมที่ผู้ใช้กรอก
	log.Printf("Stored activity for user %s: %s", State, activity)
	// ดึง card_id จาก userState
	cardID := usercardidState[State] // ปรับใช้ให้ตรงตามข้อมูลที่เก็บไว้

	// เชื่อมต่อกับฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	// สร้างตัวแปร activityRecord
	activityRecord := &models.Activityrecord{
		CardID:   usercardidState[State],
		Activity: activity,
	}

	// แสดง log card_id และ activity
	log.Printf("Using card_id: %s, Activity: %s", cardID, activity)

	// บันทึกกิจกรรมใหม่ GetActivityRecord(db, activityRecord)
	if err := GetActivityRecord(db, activityRecord); err != nil { // ส่ง activityRecord ไปที่ฟังก์ชัน ActivityRecord
		log.Printf("Error saving activity(models.ActivityRecord): %v", err)
		sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", activity))
		return
	}

	//sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรม '%s' สำเร็จสำหรับ card_id '%s'!", activity, cardID))

	// รีเซ็ตสถานะผู้ใช้
	userState[State] = "wait status ActivityPeriodRecord"
	sendReply(bot, event.ReplyToken, "กรุณากรอกระยะเวลาในการทำกิจกรรม (เช่น 30นาที หรือ 1ชั่วโมง 10นาที):")
}

func handleActivityPeriodRecord(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ActivityPeriodRecord" {
		log.Printf("สถานะไม่ถูกต้อง %s. สถานะปัจจุบัน: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรอกระยะเวลาในการทำกิจกรรม:")
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	period := strings.TrimSpace(message.Text)
	log.Printf("Received activity input: %s", period)

	if period == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกระยะเวลาในการทำกิจกรรม (เช่น 30นาที หรือ 1ชั่วโมง 10นาที):")
		return
	}
	//ตรวจสอบว่าในuserActivity[State] มีชื่อกิจกรรมหรือยัง
	activity, exists := userActivity[State]
	if !exists || activity == "" {
		log.Printf("No activity found for user %s", State)
		sendReply(bot, event.ReplyToken, "ไม่พบกิจกรรม กรุณากรอกกิจกรรมใหม่")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	activityRecord := &models.Activityrecord{
		Activity: activity,
		Period:   period,
		EndTime:  time.Now(),
	}

	//บันทึกระยะเวลา
	if err := ActivityPeriodRecord(db, activityRecord); err != nil {
		log.Printf("Error saving period(models.ActivityPeriodRecord): %v", err)
		sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", period))
		return
	}

	sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรมสำเร็จ !!!"))

	// รีเซ็ตสถานะผู้ใช้
	userState[State] = ""
}

func handleSystemManual(bot *linebot.Client, event *linebot.Event) {
	sendReply(bot, event.ReplyToken, "คุณสามารถดูคู่มือการใช้งานระบบได้ที่ลิงก์: https://example.com/manual")
}

func handleDefault(bot *linebot.Client, event *linebot.Event) {
	sendCustomReply(bot, event.ReplyToken)
}

// ฟังก์ชันสำหรับส่งข้อความ
func sendReply(bot *linebot.Client, replyToken, message string) {
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Printf("Error replying message sendReply: %v", err)
	}
}

func createQuickReplyActivities() linebot.QuickReplyItems {
	// รายการกิจกรรมทั้งหมด
	activities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
	}

	quickReplyItems := linebot.QuickReplyItems{}
	for _, activity := range activities {
		quickReplyItems.Items = append(quickReplyItems.Items,
			linebot.NewQuickReplyButton("", linebot.NewMessageAction(activity, activity)),
		)
	}
	return quickReplyItems
}

func SaveActivity(db *sql.DB, activity string) error {
	if !validateActivity(activity) {
		return fmt.Errorf("กิจกรรม '%s' ไม่ตรงกับค่าที่อนุญาตในฐานข้อมูล", activity)
	}

	query := `INSERT INTO service_info (activity) VALUES (?)`
	_, err := db.Exec(query, activity)
	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม %s ได้: %v", activity, err)
	}
	return nil
}

// validateActivity ตรวจสอบว่าค่าที่ส่งมาตรงกับฐานข้อมูล หรือไม่
func validateActivity(activity string) bool {
	allowedActivities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
	}
	for _, allowed := range allowedActivities {
		if activity == allowed {
			return true
		}
	}
	return false
}
