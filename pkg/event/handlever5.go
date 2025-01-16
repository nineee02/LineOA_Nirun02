package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/flexmessage"
	"nirun/pkg/models"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var usercardidState = make(map[string]string)
var userState = make(map[string]string)
var userActivity = make(map[string]string)        // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน
var userCheckInStatus = make(map[string]bool)     // เก็บสถานะการเช็คอินของแต่ละบัญชี LINE
var employeeLoginStatus = make(map[string]string) // เก็บสถานะล็อกอิน {employeeID: userID}

// HandleEvent - จัดการข้อความที่ได้รับจาก LINE
func HandleEvent(bot *linebot.Client, event *linebot.Event) {
	text := event.Message.(*linebot.TextMessage).Text
	log.Println("Text: ", text)
	State := event.Source.UserID
	log.Println("userstate:", State)

	switch text {
	case "ค้นหาข้อมูล":
		handleElderlyInfoStste(bot, event, event.Source.UserID)
	case "ลงเวลางาน":
		handleWorktimeStste(bot, event, event.Source.UserID)
	case "บันทึกกิจกรรม":
		handleServiceRecordStste(bot, event, event.Source.UserID)
	default:
		handleDefault(bot, event)
	}

	state, exists := userState[State]
	if exists {
		switch state {
		case "wait status worktime":
			handleWorktime(bot, event, State)
		case "wait status worktimeConfirmCheckIn":
			handleworktimeConfirmCheckIn(bot, event, State)
		case "wait status worktimeConfirmCheckOut":
			handleworktimeConfirmCheckOut(bot, event, State)
		case "wait status ElderlyInfoRequest":
			handlePateintInfo(bot, event, State)
		case "wait status ServiceRecordRequest":
			handleServiceInfo(bot, event, State)
		case "wait status Activityrecord":
			handleActivityrecord(bot, event, State)
		case "wait status ActivityStart":
			handleActivityStart(bot, event, State)
		case "wait status ActivityEnd":
			handleActivityEnd(bot, event, State)
		case "wait status Saveavtivityend":
			handleSaveavtivityend(bot, event, State)
		default:
			log.Printf("Unhandled state for user %s: %s", State, State)
		}
		return
	}

}

func setUserState(userID, state string) {
	userState[userID] = state
	log.Printf("Set user state for user %s to %s", userID, state)
}

func getUserState(userID string) (string, bool) {
	state, exists := userState[userID]
	return state, exists
}

func handleElderlyInfoStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID

	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}

	// อนุญาตให้ดำเนินการ
	setUserState(State, "wait status ElderlyInfoRequest")
}

func handleServiceRecordStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID

	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}
	setUserState(State, "wait status ServiceRecordRequest")
}

func handleWorktimeStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID

	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}

	// อนุญาตให้ดำเนินการหากยังไม่มีการเช็คอิน
	setUserState(State, "wait status worktime")
}

func handleServiceHistoryStste(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status HistoryRequest")
}

// ************************************************************************************************************************

// **********************************************************************************************************

func getUserProfile(bot *linebot.Client, userID string) (*linebot.UserProfileResponse, error) {
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return profile, nil
}
func sendCustomReply(bot *linebot.Client, replyToken string, userID string, greetingMessage string, messages ...linebot.SendingMessage) {
	if len(messages) == 0 {
		return
	}

	// ใช้ข้อความทักทายที่กำหนดเอง หรือดึงจากโปรไฟล์
	if greetingMessage == "" {
		profile, err := getUserProfile(bot, userID)
		if err == nil {
			greetingMessage = fmt.Sprintf("ยินดีต้อนรับ %s! ", profile.DisplayName)
		} else {
			greetingMessage = "ยินดีต้อนรับ!"
		}
	}

	// แทรกข้อความทักทายไปในข้อความที่ส่ง
	messages = append([]linebot.SendingMessage{linebot.NewTextMessage(greetingMessage)}, messages...)

	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
		log.Printf("Error replying message sendCustomReply: %v", err)
	}
}
func sendQRCodeForLogin(bot *linebot.Client, replyToken string) {

	flexmessage.SendRegisterLink(bot, replyToken)
}

// ตรวจสอบสถานะการเช็คอินของบัญชี
func isUserCheckedIn(userID string) bool {
	status, exists := userCheckInStatus[userID]
	return exists && status
}

// ตรวจสอบการล็อกอิน
func isEmployeeLoggedIn(employeeID, userID string) bool {
	currentUser, exists := employeeLoginStatus[employeeID]
	return exists && currentUser != userID
}

// ล็อกสถานะการใช้งาน
func lockEmployeeLogin(employeeID, userID string) bool {
	if isEmployeeLoggedIn(employeeID, userID) {
		return false // มีผู้ใช้อื่นกำลังใช้งานอยู่
	}
	employeeLoginStatus[employeeID] = userID
	return true
}

// ปลดล็อกสถานะการใช้งาน
func unlockEmployeeLogin(employeeID string) {
	delete(employeeLoginStatus, employeeID)
}

func handleWorktime(bot *linebot.Client, event *linebot.Event, userID string) {
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ดึงข้อมูลผู้ใช้ที่เชื่อมโยงกับ LINE UserID
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ที่เชื่อมโยงกับบัญชี LINE นี้.")
		} else {
			log.Println("Error fetching user info:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่.")
		}
		return
	}

	// ตรวจสอบสถานะการเข้างาน (Check-In)
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking user status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบข้อความที่ส่งมา
	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

	switch message {
	case "เช็คอิน":
		if checkedIn {
			sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่.")
			return
		}
		userState[userID] = "wait status worktimeConfirmCheckIn"
		handleworktimeConfirmCheckIn(bot, event, userID)

	case "เช็คเอ้าท์":
		if !checkedIn {
			sendReply(bot, event.ReplyToken, "คุณยังไม่ได้เช็คอิน กรุณาเช็คอินก่อนทำการเช็คเอาท์.")
			return
		}
		userState[userID] = "wait status worktimeConfirmCheckOut"
		handleworktimeConfirmCheckOut(bot, event, userID)

	default:
		quickReply := linebot.NewQuickReplyItems(
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("เช็คอิน", "เช็คอิน")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("เช็คเอ้าท์", "เช็คเอ้าท์")),
		)
		sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณาเลือก 'เช็คอิน' หรือ 'เช็คเอ้าท์'", quickReply)
	}
}

// ฟังก์ชันตรวจสอบสถานะการลงทะเบียน
func isUserRegistered(userID string) bool {
	// ตรวจสอบจากสถานะใน userState หรือฐานข้อมูล
	state, exists := userState[userID]
	return exists && state == "registered"
}

func handleworktimeConfirmCheckIn(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status worktimeConfirmCheckIn" {
		log.Printf("Unhandled state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ดึงข้อมูลผู้ใช้จาก LINE User ID
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบสถานะการเช็คอิน/เช็คเอ้าท์
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}

	// หากยังไม่ได้เช็คเอ้าท์
	if checkedIn {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอ้าท์ก่อนทำการเช็คอินใหม่.")
		userState[userID] = "wait status worktimeConfirmCheckOut"
		return
	}

	// บันทึกการเช็คอิน
	err = RecordCheckIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println(err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-in กรุณาลองใหม่.")
		return
	}

	// สร้าง worktimeRecord สำหรับ FormatworktimeCheckin
	worktimeRecord := &models.WorktimeRecord{
		UserInfo: &models.User_info{
			Name: userInfo.Name,
		},
		CheckIn: time.Now(),
	}

	flexMessage := flexmessage.FormatworktimeCheckin(worktimeRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error sending Flex Message:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
		return
	}
	log.Printf("replyMessage checkin success: %+v", flexMessage)

	userState[userID] = "wait status worktimeConfirmCheckOut"
	log.Printf("User state updated to: %s", userState[userID])
}

func handleworktimeConfirmCheckOut(bot *linebot.Client, event *linebot.Event, userID string) {
	log.Println("Starting handleworktimeConfirmCheckOut for user:", userID)

	if userState[userID] != "wait status worktimeConfirmCheckOut" {
		log.Printf("Unhandled state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	log.Println("Connecting to database...")
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ดึงข้อมูลผู้ใช้จาก LINE User ID
	log.Println("Fetching user info from database...")
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}
	log.Printf("Fetched user info: %+v", userInfo)

	// ตรวจสอบการเช็คอิน
	log.Println("Checking user check-in status...")
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
	log.Printf("Check-in status for user %s: %v", userID, checkedIn)

	// หากยังไม่ได้เช็คอิน
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "คุณยังไม่ได้เช็คอิน กรุณาเช็คอินก่อนทำการเช็คเอ้าท์.")
		return
	}

	// บันทึกการเช็คเอ้าท์
	log.Println("Recording user check-out...")
	err = RecordCheckOut(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error recording check-out:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-out กรุณาลองใหม่.")
		return
	}
	log.Println("Check-out recorded successfully for user:", userID)

	// ดึงข้อมูล WorktimeRecord
	log.Println("Fetching worktime record for user...")

	worktimeRecord, err := GetWorktimeRecordByUserID(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error fetching worktime record:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูล กรุณาลองใหม่.")
		return
	}
	log.Printf("WWWWWWWWWW:%+v", worktimeRecord)
	if worktimeRecord == nil {
		log.Println("Errrrrrrrrrr")
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}
	log.Printf("Worktime Record: %+v", worktimeRecord)

	// สร้าง WorktimeRecord สำหรับ FormatworktimeCheckout
	log.Println("Creating new worktime record for response...")
	worktimeRecord = &models.WorktimeRecord{
		UserInfo: &models.User_info{
			Name: userInfo.Name,
		},
		CheckOut: time.Now(),
		Period:   worktimeRecord.Period,
	}
	log.Printf("New worktime record: %+v", worktimeRecord)

	// ส่งข้อความตอบกลับ
	log.Println("Sending reply message...")
	flexMessage := flexmessage.FormatworktimeCheckout(worktimeRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error sending Flex Message:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
		return
	}
	log.Printf("replyMessage checkin success: %+v", flexMessage)
	// อัปเดตสถานะผู้ใช้
	userState[userID] = "wait status worktimeConfirmCheckIn"
	log.Printf("User state updated to: %s", userState[userID])
}

func resetUserState(userID string) {
	delete(userState, userID)
	delete(userActivity, userID)
	delete(usercardidState, userID)
	log.Printf("Reset state for user %s", userID)
}

func handlePateintInfo(bot *linebot.Client, event *linebot.Event, userID string) {

	state, exists := getUserState(userID)
	if !exists || state != "wait status ElderlyInfoRequest" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, state)
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	// // ป้องกันการประมวลผลซ้ำ
	// if userActivity[userID] == "processing" {
	// 	sendReply(bot, event.ReplyToken, "คำขอของคุณกำลังดำเนินการ กรุณารอสักครู่.")
	// 	return
	// }
	// userActivity[userID] = "processing"
	// defer func() { userActivity[userID] = "" }() // รีเซ็ตสถานะเมื่อเสร็จสิ้น

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := event.Message.(*linebot.TextMessage).Text

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// หากผู้ใช้ส่งข้อความ "ข้อมูลผู้สูงอายุ"
	if strings.TrimSpace(message) == "ค้นหาข้อมูล" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
	cardID := strings.TrimSpace(message)
	if cardID == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วยที่ท่านค้นหา กรุณาตรวจสอบเลขประจำตัวผู้ป่วยอีกครั้ง.")
		} else {
			log.Println("Error fetching patient info:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ป่วย กรุณาลองใหม่.")
		}
		return
	}

	replyMessage := FormatPatientInfo(patient)
	sendReply(bot, event.ReplyToken, replyMessage)

	// รีเซ็ตสถานะผู้ใช้
	userState[userID] = ""
}

func handleServiceInfo(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ServiceRecordRequest" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := event.Message.(*linebot.TextMessage).Text

	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ตรวจสอบว่าเช็คอินแล้วหรือไม่ โดยใช้ user_info_id
	userInfo, err := GetUserInfoByLINEID(db, event.Source.UserID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบสถานะการเช็คอินของผู้ใช้
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}

	// ถ้ายังไม่ได้เช็คอิน
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลางาน'")
		return
	}

	// หากผู้ใช้ส่งข้อความ "บันทึกกิจกรรม"
	if strings.TrimSpace(message) == "บันทึกกิจกรรม" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	cardID := strings.TrimSpace(message)
	if cardID == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
	service, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วยที่ท่านค้นหา กรุณาตรวจสอบเลขประจำตัวผู้ป่วยอีกครั้ง.")
		} else {
			log.Println("Error models.GetServiceInfoBycardID:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ป่วย กรุณาลองใหม่.")
		}
		return
	}

	replyMessage := FormatServiceInfo([]models.Activityrecord{*service})
	log.Println("ข้อมูลผู้สูงอายุ :", replyMessage)

	// สร้าง Quick Reply สำหรับกิจกรรม
	quickReplyActivities := createQuickReplyActivities()
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
	).Do(); err != nil {
		log.Printf("Error replying message (handleServiceInfo): %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
	}

	// บันทึก cardID สำหรับใช้ในฟังก์ชันถัดไป
	usercardidState[State] = cardID

	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
	userState[State] = "wait status Activityrecord"
	log.Printf("Set user state to wait status Activityrecord for user %s", State)
}

func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status Activityrecord:", userState)

	if userState[State] != "wait status Activityrecord" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	activity := strings.TrimSpace(message.Text)
	log.Printf("Received activity input: %s", activity)

	// ตรวจสอบว่ากิจกรรมถูกต้อง
	if !validateActivity(activity) {
		sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
		return
	}

	// เก็บกิจกรรมที่ผู้ใช้เลือก
	userActivity[State] = activity
	log.Printf("Stored activity for user %s: %s", State, activity)

	// สร้าง Quick Reply สำหรับเริ่มกิจกรรม
	quickReply := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("เริ่มกิจกรรม", "เริ่มกิจกรรม")),
	)

	sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณากดปุ่ม 'เริ่มกิจกรรม' เพื่อเริ่มบันทึกเวลา", quickReply)

	// อัปเดตสถานะผู้ใช้
	userState[State] = "wait status ActivityStart"
	log.Println("wait status ActivityStart:", userState)
}

func handleActivityStart(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status ActivityStart:", userState)
	if userState[State] != "wait status ActivityStart" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	starttime := strings.TrimSpace(message.Text)
	log.Printf("Received activity input(handleActivityStart): %s", starttime)

	if starttime != "เริ่มกิจกรรม" {
		userState[State] = "wait status Activityrecord"
		handleActivityrecord(bot, event, State)
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	cardID, exists := usercardidState[State]
	if !exists || cardID == "" {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลบัตรประชาชน กรุณากรอกใหม่")
		return
	}

	// ดึง user_info_id จาก event
	userID := event.Source.UserID

	// ดึงข้อมูล userInfo จากฐานข้อมูล
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบสถานะการเช็คอินจาก user_info_id
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}

	// ถ้ายังไม่ได้เช็คอิน
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณาเช็คอินก่อน\nที่เมนู 'ลงเวลางาน'")
		return
	}

	// ดึงข้อมูลผู้ป่วยจากฐานข้อมูล
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}

	activityRecord := &models.Activityrecord{
		PatientInfo: models.PatientInfo{
			CardID:         cardID,
			Name:           patient.PatientInfo.Name,
			PatientInfo_ID: patient.PatientInfo.PatientInfo_ID,
		},
		ServiceInfo: models.ServiceInfo{
			Activity: userActivity[State],
		},
		EmployeeInfo: models.EmployeeInfo{
			EmployeeInfo_ID: userInfo.UserInfo_ID,
		},
		StartTime: time.Now(),
		UserInfo: models.User_info{
			UserInfo_ID: userInfo.UserInfo_ID,
			Create_by:   userInfo.Name,
			Write_by:    userInfo.Name,
		},
	}

	// บันทึกกิจกรรม
	if err := SaveActivityRecord(db, activityRecord); err != nil {
		log.Printf("Error saving activity record: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกกิจกรรม กรุณาลองใหม่")
		return
	}

	activityRecord.PatientInfo.Name = patient.PatientInfo.Name
	replyMessage := FormatactivityRecordStarttime([]models.Activityrecord{*activityRecord})

	quickReply := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("เสร็จสิ้น", "เสร็จสิ้น")),
	)
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(replyMessage).WithQuickReplies(quickReply),
	).Do()
	if err != nil {
		log.Println("Error sending Quick Reply (handleActivityStart):", err)
	}

	// เปลี่ยนสถานะผู้ใช้
	userState[State] = "wait status ActivityEnd"
}

func handleActivityEnd(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status ActivityEnd" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		return
	}
	// รับข้อความกิจกรรมจากผู้ใช้
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	endtime := strings.TrimSpace(message.Text)
	log.Printf("Received activity input: %s", endtime)

	if endtime != "เสร็จสิ้น" {
		sendReply(bot, event.ReplyToken, "กรุณาพิมพ์ 'เสร็จสิ้น' เพื่อบันทึกเวลาสิ้นสุด")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	// ขอชื่อพนักงาน
	sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ:")

	// สถานะเปลี่ยนไปเป็นกรอกชื่อพนักงาน
	userState[userID] = "wait status Saveavtivityend"
}
func handleSaveavtivityend(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status Saveavtivityend" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		return
	}

	// รับชื่อพนักงานจากผู้ใช้
	employeeName := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Printf("Received employee name: %s", employeeName)

	if employeeName == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ.")
		return
	}

	// เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	// ใช้ฟังก์ชันแยกเพื่อค้นหา employee_info_id
	employeeID, err := GetEmployeeIDByName(db, employeeName)
	if err != nil {
		sendReply(bot, event.ReplyToken, err.Error())
		return
	}

	// ตรวจสอบ cardID
	cardID, exists := usercardidState[userID]
	if !exists || cardID == "" {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลบัตรประชาชน กรุณากรอกใหม่")
		return
	}

	// ดึง activityRecord_ID
	activityRecordID, err := GetActivityRecordID(db, cardID)
	if err != nil {
		sendReply(bot, event.ReplyToken, err.Error())
		return
	}

	// ดึงข้อมูล userInfo จากฐานข้อมูล
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ดึง patient_info_id โดยใช้ cardID
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}

	// บันทึกข้อมูลใน activity_record
	activityRecord := &models.Activityrecord{
		ActivityRecord_ID: activityRecordID,
		PatientInfo: models.PatientInfo{
			CardID:         cardID,
			Name:           patient.PatientInfo.Name,
			PatientInfo_ID: patient.PatientInfo.PatientInfo_ID,
		},
		ServiceInfo: models.ServiceInfo{
			Activity: userActivity[userID],
		},
		EndTime:      time.Now(),
		EmployeeInfo: models.EmployeeInfo{EmployeeInfo_ID: employeeID},
		UserInfo:     models.User_info{UserInfo_ID: userInfo.UserInfo_ID},
	}
	log.Println("ASD:", activityRecord)

	startTime, err := GetActivityStartTime(db, cardID, userActivity[userID])
	if err != nil {
		log.Printf("Error fetching StartTime: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลเวลาเริ่ม กรุณาลองใหม่")
		return
	}

	duration := activityRecord.EndTime.Sub(startTime)
	activityRecord.Period = formatDuration(duration)
	// บันทึกข้อมูลในฐานข้อมูล
	if err := UpdateActivityEndTime(db, activityRecord); err != nil {
		log.Printf("Error updating end time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
		return
	}

	// ใช้ฟังก์ชัน FormatactivityRecordEndtime เพื่อสร้างข้อความตอบกลับ
	replyMessage := FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
	sendReply(bot, event.ReplyToken, replyMessage)

	log.Printf("บันทึกกิจกรรมสำเร็จ: %s", replyMessage)
	resetUserState(userID)
}
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
}

// func handleSaveavtivityend(bot *linebot.Client, event *linebot.Event, userID string) {
// 	//รับชื่อพนักงาน
// 	//ชื่อพนักงานไปเช็คในฐานข้อมูล
// 	//บันทึกข้อมูล endtime
// 	//
// 	//
// 	if userState[userID] != "wait status Saveavtivityend" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		return
// 	}

// 	// รับชื่อพนักงานจากผู้ใช้
// 	employeeName := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Printf("handleSaveavtivityend---Received employee name: %s", employeeName)

// 	if employeeName == "" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ.")
// 		return
// 	}

// 	// เชื่อมต่อฐานข้อมูล
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 		return
// 	}
// 	defer db.Close()

// 	// ใช้ฟังก์ชันแยกเพื่อค้นหา employee_info_id
// 	employeeID, err := GetEmployeeIDByName(db, employeeName)
// 	if err != nil {
// 		sendReply(bot, event.ReplyToken, err.Error())
// 		return
// 	}
// 	//
// 	log.Println("Emm", employeeID)

// 	// // ตรวจสอบ cardID //ไว้ก่อน
// 	// cardID, exists := usercardidState[userID]
// 	// if !exists || cardID == "" {
// 	// 	sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลบัตรประชาชน กรุณากรอกใหม่")
// 	// 	return
// 	// }
// 	// // ดึง patient_info_id โดยใช้ cardID
// 	// patientinfoID, err := GetPatientInfoByName(db, cardID)
// 	// if err != nil {
// 	// 	log.Printf("Error fetching patient_info_id: %v", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
// 	// 	return
// 	// }
// 	// -----------------------------------------------------------------ตรงนี้แหละ
// 	// ดึง activityRecord_ID
// 	activityRecordID, err := GetActivityRecordByEmployeeID(db, employeeID)
// 	if err != nil {
// 		sendReply(bot, event.ReplyToken, err.Error())
// 		return
// 	}

// 	log.Printf("AAAeaA %+v", activityRecordID.ActivityRecord_ID)

// 	// บันทึกข้อมูลใน activity_record
// 	activityRecord := &models.Activityrecord{
// 		ActivityRecord_ID: activityRecordID.ActivityRecord_ID,
// 		PatientInfo: models.PatientInfo{
// 			CardID:         activityRecordID.PatientInfo.CardID,
// 			Name:           activityRecordID.PatientInfo.Name,
// 			PatientInfo_ID: activityRecordID.PatientInfo.PatientInfo_ID,
// 		},
// 		ServiceInfo: models.ServiceInfo{
// 			Activity: userActivity[userID],
// 		},
// 		EndTime:      time.Now(),
// 		EmployeeInfo: models.EmployeeInfo{EmployeeInfo_ID: employeeID},
// 		Write_by:     activityRecordID.Write_by,
// 		Period:       activityRecordID.Period,
// 	}
// 	log.Println("In::::", activityRecord)
// 	// บันทึกข้อมูลในฐานข้อมูล
// 	if err := UpdateActivityEndTimeForPatient(db, activityRecordID); err != nil {
// 		log.Printf("Error updating end time: %v", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
// 		return
// 	}
// 	// log.Println("AAAAA:", activityRecord)

// 	// ใช้ฟังก์ชัน FormatactivityRecordEndtime เพื่อสร้างข้อความตอบกลับ
// 	replyMessage := FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
// 	sendReply(bot, event.ReplyToken, replyMessage)

// 	log.Printf("บันทึกกิจกรรมสำเร็จ: %s", replyMessage)
// 	resetUserState(userID)
// }

func createQuickReplyActivities() linebot.QuickReplyItems {
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

func handleSystemManual(bot *linebot.Client, event *linebot.Event, State string) {
	sendReply(bot, event.ReplyToken, "คุณสามารถดูคู่มือการใช้งานระบบได้ที่ลิงก์: https://example.com/manual")
}

func handleDefault(bot *linebot.Client, event *linebot.Event) {
	// userID := event.Source.UserID
	// sendCustomReply(bot, event.ReplyToken, userID, linebot.NewTextMessage("Hello, World!"))
}

// ฟังก์ชันสำหรับส่งข้อความ
func sendReply(bot *linebot.Client, replyToken, message string) {
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Println("Error sending reply message:", err)
	}
}
func sendReplyWithQuickReply(bot *linebot.Client, replyToken string, message string, quickReply *linebot.QuickReplyItems) {
	textMessage := linebot.NewTextMessage(message).WithQuickReplies(quickReply)
	if _, err := bot.ReplyMessage(replyToken, textMessage).Do(); err != nil {
		log.Printf("Error sending reply with quick reply: %v", err)
	}
}
