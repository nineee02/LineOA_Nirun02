package event

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/flexmessage"
	"nirun/pkg/models"
	"regexp"
	"unicode"

	// "nirun/service"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var usercardidState = make(map[string]string)
var userState = make(map[string]string)              //เก็บstate
var userActivity = make(map[string]string)           // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน
var userCheckInStatus = make(map[string]bool)        // เก็บสถานะการเช็คอินของแต่ละบัญชี LINE
var userActivityInfoID = make(map[string]int)        // เก็บ activity_info_id ตาม userID
var userActivityRecordID = make(map[string]int)      // เก็บ activityRecord_ID ตาม State ของผู้ใช้
var userActivityCategory = make(map[string]string)   // เก็บมิติของกิจกรรมที่เลือก
var employeeLoginStatus = make(map[string]string)    // เก็บสถานะล็อกอิน {employeeID: userID}
var userImageTimestamps = make(map[string]time.Time) // เก็บ timestamp ของรูปภาพ

// HandleEvent - จัดการข้อความที่ได้รับจาก LINE
func HandleEvent(bot *linebot.Client, event *linebot.Event) {
	// ตรวจสอบประเภทข้อความก่อน
	switch message := event.Message.(type) {
	case *linebot.TextMessage: // อ่านข้อความจาก TextMessage
		text := strings.TrimSpace(message.Text)
		log.Println("Received TextMessage:", text)
		State := event.Source.UserID
		log.Println("User state: ", State)

		// ตรวจสอบคำสั่งจากข้อความ
		switch text {
		case "ค้นหาข้อมูล":
			handleElderlyInfoStste(bot, event, State)
		case "ลงเวลางาน":
			handleWorktimeStste(bot, event, State)
		case "บันทึกการบริการ":
			handleServiceRecordStste(bot, event, State)
		default:
			handleDefault(bot, event)
		}

		// ตรวจสอบสถานะ
		state, exists := userState[State]
		if exists {
			switch state {
			case "wait status worktime":
				handleWorktime(bot, event, State)
			case "wait status ElderlyInfoRequest":
				handlePateintInfo(bot, event, State)
			case "wait status handleServiceGetCardID":
				handleServiceGetCardID(bot, event, State)
			case "wait status ServiceSelection":
				handleServiceSelection(bot, event, State)
			// case "wait status ServiceRecordRequest":
			// 	handleServiceInfo(bot, event, State)
			case "wait status ActivitySelection":
				handleActivitySelection(bot, event, State)
			case "wait status CustomActivity":
				handleCustomActivity(bot, event, State)
			case "wait status Activityrecord":
				handleActivityrecord(bot, event, State)
			case "wait status ActivityStart":
				handleActivityStart(bot, event, State)
			case "wait status ActivityEnd":
				handleActivityEnd(bot, event, State)
			// case "wait status ConfirmOrSaveEmployee":
			// 	handleUserChoiceForActivityRecord(bot, event, State, "")
			case "wait status ConfirmOrSaveEmployee":
				if textMessage, ok := event.Message.(*linebot.TextMessage); ok {
					selection := strings.TrimSpace(textMessage.Text)  // รับข้อความที่ผู้ใช้ส่งมา
					log.Printf("📌 Handling selection: %s", selection) // เพิ่ม log ตรวจสอบค่า
					handleUserChoiceForActivityRecord(bot, event, State, selection)
				} else {
					log.Printf("❌ Unexpected message type in ConfirmOrSaveEmployee state")
					sendReply(bot, event.ReplyToken, "กรุณาเลือก 'ยืนยันการบันทึก' หรือ 'บันทึกข้อมูลแทน'")
				}
			// case "wait status saveActivityRecordForOtherEmployee":
			// 	saveActivityRecordForOtherEmployee(bot, event, State)
			case "wait status SaveEmployeeName":
				handleSaveEmployeeName(bot, event, State, State, "")
			default:
				log.Printf("Unhandled state for user %s: %s", State, state)
			}
		}

	case *linebot.ImageMessage: // อ่านข้อความจาก ImageMessage
		State := event.Source.UserID
		state, exists := userState[State]
		if exists {
			switch state {
			case "wait status Saveavtivityend":
				log.Printf("Received ImageMessage: ID=%s", message.ID)
				handleSaveavtivityend(bot, event, usercardidState[event.Source.UserID], event.Source.UserID)
			case "wait status saveEvidenceImageafterActivity":
				handlesaveEvidenceImageafterActivity(bot, event, usercardidState[event.Source.UserID], event.Source.UserID)
			}
		}

	default:
		log.Printf("Unhandled message type: %T", event.Message)
		sendReply(bot, event.ReplyToken, "ไม่สามารถประมวลผลข้อความประเภทนี้ได้.")
	}
}

// สำหรับตั้งค่าสถานะผู้ใช้
func setUserState(userID, state string) {
	userState[userID] = state
	log.Printf("Set user state for user %s to %s", userID, state)
}

// สำหรับดุงสถานะผู้ใช้
func getUserState(userID string) (string, bool) {
	state, exists := userState[userID]
	return state, exists
}

// เริ่มกระบวนการตรวจสอบสถานะการลงเวลาทำงาน
func handleWorktimeStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID
	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		return
	}

	// อนุญาตให้ดำเนินการหากยังไม่มีการเช็คอิน
	setUserState(State, "wait status worktime")
}

// เริ่มกระบวนการขอค้นหาข้อมูล
func handleElderlyInfoStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID
	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		return
	}

	// อนุญาตให้ดำเนินการ
	setUserState(State, "wait status ElderlyInfoRequest")
}

// เริ่มกระบวนการขอบันทึกกิจกรรม
func handleServiceRecordStste(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID

	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
	if isUserCheckedIn(userID) {
		return
	}

	// ขอให้ผู้ใช้กรอกเลขบัตรประชาชน
	sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")

	// ตั้งค่าผู้ใช้ให้อยู่ในโหมดรอเลขบัตรประชาชน
	setUserState(State, "wait status handleServiceGetCardID")
}

// **********************************************************************************************************

func getUserProfile(bot *linebot.Client, userID string) (*linebot.UserProfileResponse, error) {
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// ส่งข้อความตอบกลับแบบกำหนดเอง สามารถดึงข้อมูลผู้ใช้ได้
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

//*************************************************************************************************************

// // ตรวจสอบการล็อกอิน
// func isEmployeeLoggedIn(employeeID, userID string) bool {
// 	currentUser, exists := employeeLoginStatus[employeeID]
// 	return exists && currentUser != userID
// }

// // ล็อกสถานะการใช้งาน
// func lockEmployeeLogin(employeeID, userID string) bool {
// 	if isEmployeeLoggedIn(employeeID, userID) {
// 		return false // มีผู้ใช้อื่นกำลังใช้งานอยู่
// 	}
// 	employeeLoginStatus[employeeID] = userID
// 	return true
// }

// // ปลดล็อกสถานะการใช้งาน
// func unlockEmployeeLogin(employeeID string) {
// 	delete(employeeLoginStatus, employeeID)
// }

//*************************************************************************************

// ตรวจสอบสถานะของบัญชี
func isUserCheckedIn(userID string) bool {
	status, exists := userCheckInStatus[userID]
	return exists && status
}

// ฟังก์ชันตรวจสอบสถานะการลงทะเบียน
func isUserRegistered(userID string) bool {
	// ตรวจสอบจากสถานะใน userState หรือฐานข้อมูล
	state, exists := userState[userID]
	return exists && state == "registered"
}

// ฟังก์ชันลงเวลาเข้าและออกงาน
func handleWorktime(bot *linebot.Client, event *linebot.Event, userID string) {
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ดึงข้อมูลผู้ใช้ตาม LINE ID
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

	// ตรวจสอบการเช็คอินของพนักงาน
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
			// ถ้าผู้ใช้เช็คอินแล้ว ให้แสดงปุ่ม "เช็คเอ้าท์"
			UpdateWorktimeUI(bot, event, userInfo, true)
			return
		}

		// บันทึกข้อมูลเช็คอินลงฐานข้อมูล
		err = RecordCheckIn(db, userInfo.UserInfo_ID)
		if err != nil {
			log.Println("Error recording check-in:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-in กรุณาลองใหม่.")
			return
		}
		// เตรียมข้อมูล worktimeRecord
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
		// อัปเดตปุ่มเป็น "เช็คเอ้าท์"
		UpdateWorktimeUI(bot, event, userInfo, true)

	case "เช็คเอ้าท์":
		if !checkedIn {
			// ถ้าผู้ใช้ยังไม่ได้เช็คอิน ให้แสดงปุ่ม "เช็คอิน"
			UpdateWorktimeUI(bot, event, userInfo, false)
			return
		}

		// บันทึกข้อมูลเช็คเอ้าท์ลงฐานข้อมูล
		err = RecordCheckOut(db, userInfo.UserInfo_ID)
		if err != nil {
			log.Println("Error recording check-out:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-out กรุณาลองใหม่.")
			return
		}
		// ดึงข้อมูลบันทึกเวลาทำงาน
		worktimeRecord, err := GetWorktimeRecordByUserID(db, userInfo.UserInfo_ID)
		if err != nil {
			log.Println("Error fetching worktime record:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูล กรุณาลองใหม่.")
			return
		}
		log.Printf("worktimeRecor(check out):%+v", worktimeRecord)
		if worktimeRecord == nil {
			log.Println("Error worktimeRecord check out")
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลการทำงาน กรุณาลองใหม่.")
			return
		}
		log.Printf("Worktime Record: %+v", worktimeRecord)

		// เตรียมข้อมูล WorktimeRecord
		worktimeRecord = &models.WorktimeRecord{
			UserInfo: &models.User_info{
				Name: userInfo.Name,
			},
			CheckOut: time.Now(),
			Period:   worktimeRecord.Period,
		}
		log.Printf("New worktime record: %+v", worktimeRecord)

		flexMessage := flexmessage.FormatworktimeCheckout(worktimeRecord)
		if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
			log.Println("Error sending Flex Message:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
			return
		}
		log.Printf("replyMessage checkin success: %+v", flexMessage)

		userState[userID] = "wait status worktimeConfirmCheckIn"
		log.Printf("User state updated to: %s", userState[userID])
		// อัปเดตปุ่มเป็น "เช็คอิน"
		UpdateWorktimeUI(bot, event, userInfo, false)

	default:
		// ถ้าไม่ใช่ "เช็คอิน" หรือ "เช็คเอ้าท์" ให้แสดงปุ่มตามสถานะปัจจุบัน
		UpdateWorktimeUI(bot, event, userInfo, checkedIn)
	}
}

// ฟังก์ชันสำหรับอัปเดต UI ของปุ่มเช็คอิน / เช็คเอ้าท์
func UpdateWorktimeUI(bot *linebot.Client, event *linebot.Event, userInfo *models.User_info, checkedIn bool) {
	worktimeRecord := &models.WorktimeRecord{
		UserInfo: &models.User_info{
			Name: userInfo.Name,
		},
		CheckIn:  time.Now(),
		CheckOut: time.Time{},
	}

	var flexMessage *linebot.FlexMessage
	if checkedIn {
		// แสดงปุ่ม "เช็คเอ้าท์"
		flexMessage = flexmessage.FormatConfirmCheckout(worktimeRecord)
	} else {
		// แสดงปุ่ม "เช็คอิน"
		flexMessage = flexmessage.FormatConfirmCheckin(worktimeRecord)
	}

	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error sending Flex Message:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
	}
}

// func resetUserState(userID string) {
// 	delete(userState, userID)
// 	delete(userActivity, userID)
// 	delete(usercardidState, userID)
// 	log.Printf("Reset state for user %s", userID)
// }

func sanitizeCardID(s string) string {
	var builder strings.Builder
	for _, char := range s {
		if unicode.IsDigit(char) { // ตรวจสอบเฉพาะตัวเลข
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

// การค้นหาข้อมูล
func handlePateintInfo(bot *linebot.Client, event *linebot.Event, userID string) {
	state, exists := getUserState(userID)
	if !exists || state != "wait status ElderlyInfoRequest" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, state)
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ข้อความที่รับ = "ค้นหาข้อมูล"
	if message == "ค้นหาข้อมูล" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
	cardID := sanitizeCardID(message)
	if len(cardID) != 13 {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักที่ถูกต้อง\nตัวอย่างเช่น 1234567891234 :")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// patient, err := service.PostRequestByID(cardID)
	// if err != nil {
	// 	log.Println("ErE:", err)
	// 	return
	// }
	// log.Println("Papatient:", patient)
	// ดึงข้อมูลผู้ป่วยจาก CardID
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Println("Error fetching patient info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้สูงอายุ กรุณากรอกเลขประจำตัวประชาชนอีกครั้ง")
		return
	}
	flexMessage := flexmessage.FormatPatientInfo(&patient.PatientInfo)
	if _, err := bot.PushMessage(userID, flexMessage).Do(); err != nil {
		log.Println("Error sending push message:", err)
	}

	log.Println("ข้อมูลผู้ป่วย:", flexMessage)
	userState[userID] = ""
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func parseTimeInput(input string) (time.Time, error) {
	// 🔹 ลบช่องว่างส่วนเกิน และทำให้เป็น lower case
	input = strings.TrimSpace(strings.ToLower(input))

	// 🔹 กำหนด regex ให้รับค่าหลายรูปแบบ
	re := regexp.MustCompile(`(\d{1,2})[:.](\d{2})`) // รองรับ "11:00", "11.00"
	match := re.FindStringSubmatch(input)

	if len(match) < 3 {
		return time.Time{}, fmt.Errorf("รูปแบบเวลาไม่ถูกต้อง")
	}

	// 🔹 แปลงชั่วโมงและนาทีเป็นตัวเลข
	hour, min := match[1], match[2]
	parsedTime, err := time.Parse("15:04", fmt.Sprintf("%s:%s", hour, min))
	if err != nil {
		return time.Time{}, fmt.Errorf("เกิดข้อผิดพลาดในการแปลงเวลา")
	}

	// 🔹 ใช้วันปัจจุบัน และนำเวลาเข้าไปใช้
	now := time.Now()
	finalTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

	return finalTime, nil
}
func handleServiceGetCardID(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status handleServiceGetCardID" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ข้อความที่รับ = "บันทึกการบริการ"
	if message == "บันทึกการบริการ" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
	cardID := sanitizeCardID(message)
	if len(cardID) != 13 {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักที่ถูกต้อง\nตัวอย่างเช่น 1234567891234 :")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	//ตรวจสอบว่ามีข้อมูลผู้ป่วยหรือไม่
	if _, err := GetPatientInfoByName(db, cardID); err != nil {
		if err == sql.ErrNoRows {
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้สูงอายุ กรุณากรอกเลขประจำตัวประชาชนอีกครั้ง")
		} else {
			log.Println("Error models.GetServiceInfoBycardID:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการค้นหาข้อมูล กรุณาลองใหม่.")
		}
		return
	}

	//บันทึก cardID สำหรับใช้ในฟังก์ชันถัดไป
	usercardidState[State] = cardID
	setUserState(State, "wait status ActivitySelection")

	//ใช้ `PushMessage()` แทน `ReplyMessage()` เพื่อหลีกเลี่ยงปัญหา reply token
	flexMessage := flexmessage.FormatActivityCategories()
	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
		log.Println("Error sending activity category selection:", err)
	}
}

func handleServiceSelection(bot *linebot.Client, event *linebot.Event, State string) {
	// if userState[State] != "wait status ServiceSelection" {
	// 	log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
	// 	return
	// }

	// // ดึงข้อความที่ผู้ใช้เลือก
	// message := event.Message.(*linebot.TextMessage).Text

	// switch message {
	// case "บันทึกกิจกรรม":
	// 	// อัปเดตสถานะก่อนเรียก `handleServiceInfo`
	// 	setUserState(State, "wait status ActivitySelection")
	// 	log.Printf("User %s state changed to: %s", State, "wait status ActivitySelection")
	// 	flexMessage := flexmessage.FormatActivityCategories()
	// 	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
	// 		log.Println("Error sending activity category selection:", err)
	// 	}
	// 	// เรียกใช้ฟังก์ชันบันทึกกิจกรรม
	// 	handleActivitySelection(bot, event, State)

	// case "รายงานปัญหา":
	// 	// ตั้งค่าสถานะเป็น "wait status ReportIssue"
	// 	setUserState(State, "wait status ReportIssue")
	// 	log.Printf("User %s state changed to: %s", State, "wait status ReportIssue")

	// 	// เรียกใช้ฟังก์ชันจัดการรายงานปัญหา
	// 	// handleReportIssue(bot, event, State)

	// default:
	// 	sendReply(bot, event.ReplyToken, "กรุณาเลือก 'บันทึกกิจกรรม' หรือ 'รายงานปัญหา'")
	// }
}

// เลือกมิติของกิจกรรม
func handleActivitySelection(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ActivitySelection" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	//รับข้อความกิจกรรมที่ผู้ใช้เลือก
	category := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Printf("User selected category: %s", category)

	//ตรวจสอบว่าผู้ใช้ต้องการย้อนกลับ
	if category == "ย้อนกลับ" {
		sendActivityCategorySelection(bot, event)            //Flex Message มิติ
		setUserState(State, "wait status ActivitySelection") // อัปเดตสถานะกลับไปเลือกมิติ
		return
	}
	//ตรวจสอบมิติของกิจกรรม
	validCategories := map[string]string{
		"มิติเทคโนโลยี":   "technology",
		"มิติสังคม":       "social",
		"มิติสุขภาพ":      "health",
		"มิติเศรษฐกิจ":    "economic",
		"มิติสิ่งแวดล้อม": "environmental",
		"มิติอื่นๆ":       "other",
	}

	categoryKey, exists := validCategories[category]
	log.Printf("categoryKey:%s", categoryKey)

	if !exists {
		sendReply(bot, event.ReplyToken, "กรุณาเลือกมิติของกิจกรรมที่ถูกต้องจากเมนู")
		return
	}

	//อัปเดตหมวดหมู่ของกิจกรรมใน State
	userActivityCategory[State] = categoryKey
	log.Printf("userActivityCategory: %s", userActivityCategory)

	if categoryKey == "other" {
		//ให้ผู้ใช้กรอกกิจกรรมเอง
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อกิจกรรมของคุณ:")
		userState[State] = "wait status CustomActivity"
	} else {
		//ดึงกิจกรรมที่เกี่ยวข้องจากฐานข้อมูลและแสดงให้ผู้ใช้เลือก
		fetchAndShowActivities(bot, event, State, categoryKey)
	}
}

// ดึงกิจกรรมแต่ละมิติมาแสดงให้เลือก
func fetchAndShowActivities(bot *linebot.Client, event *linebot.Event, State string, category string) {
	//อัปเดตหมวดหมู่กิจกรรมใน state
	userActivityCategory[State] = category

	//เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	var activities []string

	//ดึงข้อมูลกิจกรรมตามหมวดหมู่จากฐานข้อมูล
	switch category {
	case "technology":
		activityList, err := GetTechnologyActivities(db)
		if err == nil {
			for _, activity := range activityList {
				activities = append(activities, strings.TrimSpace(activity.ActivityTechnology))
			}
		}
	case "social":
		activityList, err := GetSocialActivities(db)
		if err == nil {
			for _, activity := range activityList {
				activities = append(activities, strings.TrimSpace(activity.ActivitySocial))
			}
		}
	case "health":
		activityList, err := GetHealthActivities(db)
		if err == nil {
			for _, activity := range activityList {
				activities = append(activities, strings.TrimSpace(activity.ActivityHealth))
			}
		}
	case "economic":
		activityList, err := GetEconomicActivities(db)
		if err == nil {
			for _, activity := range activityList {
				activities = append(activities, strings.TrimSpace(activity.ActivityEconomic))
			}
		}
	case "environmental":
		activityList, err := GetEnvironmentalActivities(db)
		if err == nil {
			for _, activity := range activityList {
				activities = append(activities, strings.TrimSpace(activity.ActivityEnvironmental))
			}
		}
	default:
		log.Println("Invalid category selection:", category)
		return
	}

	//แสดง Flex Message ให้ผู้ใช้เลือกกิจกรรม
	if len(activities) > 0 {
		flexMessage := flexmessage.FormatActivities(activities)
		if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
			log.Println("Error sending activity list:", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการแสดงกิจกรรม กรุณาลองใหม่.")
		}
	} else {
		sendReply(bot, event.ReplyToken, "ไม่พบกิจกรรมในหมวดหมู่นี้ กรุณาเลือกหมวดหมู่อื่น.")
	}

	//อัปเดตสถานะเป็น "รอเลือกกิจกรรม"
	setUserState(State, "wait status Activityrecord")
}

// เลือกมิติอีกครั้งหากกดย้อนกลับ
func sendActivityCategorySelection(bot *linebot.Client, event *linebot.Event) {
	// log.Println("Sending activity category selection...")

	//รีเซ็ตมิติที่เลือกก่อนหน้า
	delete(userActivityCategory, event.Source.UserID)

	//ส่ง Flex Message เพื่อเลือกมิติใหม่
	flexMessage := flexmessage.FormatActivityCategories()
	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
		log.Println("Error sending activity category selection:", err)
	}
}

// รับกิจกรรมจากผู้ใช้หากเลือก "มิติอื่นๆ"
func handleCustomActivity(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status CustomActivity" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	// รับข้อความกิจกรรมที่ผู้ใช้ป้อนเอง
	activity := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Printf("User entered custom activity: %s", activity)

	if activity == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อกิจกรรมของคุณ")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	//เก็บค่ากิจกรรมที่ผู้ใช้กรอกเอง
	userActivity[State] = activity

	// // ส่ง Flex Message เพื่อยืนยันเริ่มกิจกรรม
	// flexContainer := flexmessage.FormatStartActivity(activity)
	// flexMessage := linebot.NewFlexMessage("เริ่มกิจกรรม", flexContainer)
	// if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
	// 	log.Printf("Error sending Flex Message: %v", err)
	// }

	//เปลี่ยนสถานะเป็น "รอเริ่มกิจกรรม"
	userState[State] = "wait status ActivityStart"
	sendReply(bot, event.ReplyToken, "กรุณากรอกเวลาเริ่มทำกิจกรรม เช่น 11:20 ")
}

// บันทึกกิจกรรม เมื่อเลือกกิจกรรมใหม่แล้ว
func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status Activityrecord:", userState)

	if userState[State] != "wait status Activityrecord" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
		return
	}

	// ตรวจสอบว่าเป็นข้อความ
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	// กิจกรรมที่รับมา
	activity := strings.TrimSpace(strings.ToLower(message.Text))
	log.Printf("Received activity input: %s", activity)

	//ถ้าผู้ใช้กด "ย้อนกลับ" ให้กลับไปเลือกมิติใหม่
	if activity == "ย้อนกลับ" {
		log.Println("User chose to go back. Resetting state...")

		//รีเซ็ตค่าหมวดหมู่ที่เลือกไว้
		delete(userActivityCategory, State)

		//ตั้งค่า State ใหม่ให้กลับไปเลือกมิติ
		setUserState(State, "wait status ActivitySelection")

		//ส่ง Flex Message เพื่อเลือกมิติใหม่
		sendActivityCategorySelection(bot, event)
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ตรวจสอบว่าผู้ใช้เคยเลือกมิติของกิจกรรมหรือไม่
	category, exists := userActivityCategory[State]
	if !exists {
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด: ไม่พบมิติของกิจกรรม กรุณาลองใหม่.")
		return
	}
	log.Println("category:%s", category)

	// ดึงกิจกรรมจากฐานข้อมูลตามมิติที่เลือก
	var validActivities []string
	switch category {
	case "technology":
		activityList, err := GetTechnologyActivities(db)
		if err == nil {
			for _, act := range activityList {
				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityTechnology)))
			}
		}
	case "social":
		activityList, err := GetSocialActivities(db)
		if err == nil {
			for _, act := range activityList {
				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivitySocial)))
			}
		}
	case "health":
		activityList, err := GetHealthActivities(db)
		if err == nil {
			for _, act := range activityList {
				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityHealth)))
			}
		}
	case "economic":
		activityList, err := GetEconomicActivities(db)
		if err == nil {
			for _, act := range activityList {
				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityEconomic)))
			}
		}
	case "environmental":
		activityList, err := GetEnvironmentalActivities(db)
		if err == nil {
			for _, act := range activityList {
				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityEnvironmental)))
			}
		}
	default:
		log.Println("Invalid category:", category)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่.")
		return
	}
	// log.Println("validActivities:%s", validActivities)

	// ตรวจสอบว่ากิจกรรมที่ผู้ใช้เลือกอยู่ในฐานข้อมูลหรือไม่
	isValid := false
	for _, validActivity := range validActivities {
		if activity == validActivity {
			isValid = true
			break
		}
	}

	if !isValid {
		sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
		return
	}

	// ดึง activity_info_id
	activityID, err := GetActivityInfoIDByType(db, category, activity)
	if err != nil {
		log.Println("Error fetching activity ID:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
		return
	}

	// บันทึกข้อมูล activityInfoID
	userActivityInfoID[State] = activityID
	log.Printf("Stored activityInfoID for user %s: %d", State, activityID)

	// บันทึกกิจกรรมที่ผู้ใช้เลือก
	userActivity[State] = activity
	log.Printf("Stored activity for user %s: %s", State, activity)

	// // ส่ง Flex Message เพื่อยืนยันเริ่มกิจกรรม
	// flexContainer := flexmessage.FormatStartActivity(activity)
	// flexMessage := linebot.NewFlexMessage("เริ่มกิจกรรม", flexContainer)
	// if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
	// 	log.Printf("Error sending Flex Message: %v", err)
	// }

	// อัปเดตสถานะเป็น "wait status ActivityStart"
	userState[State] = "wait status ActivityStart"
	log.Println("wait status ActivityStart: ", userState)
	sendReply(bot, event.ReplyToken, "กรุณากรอกเวลาเริ่มทำกิจกรรม เช่น 11:20 ")
}

// กดเรื่มกิจกรรม
func handleActivityStart(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status ActivityStart:", userState)

	if userState[State] != "wait status ActivityStart" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณากรอกเวลาเริ่มทำกิจกรรม:")
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	startTimeStr := strings.TrimSpace(message.Text)
	log.Printf("Received start time input: %s", startTimeStr)

	startTime, err := parseTimeInput(startTimeStr)
	if err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบเวลาไม่ถูกต้อง ตัวอย่างเช่น '13:30', '13.30', '13:30นาที'")
		return
	}

	now := time.Now()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	activityRecordID, exists := userActivityRecordID[State]
	if exists {
		log.Printf("Updating start_time for activity_record_id=%d, start_time=%v", activityRecordID, startTime)
		err := UpdateActivityStartTime(db, activityRecordID, startTime)
		if err != nil {
			log.Printf("Error updating start_time: %v", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาเริ่ม กรุณาลองใหม่")
			return
		}
	} else {
		// ✅ ดึงข้อมูลที่จำเป็น
		cardID := usercardidState[State]
		patient, err := GetPatientInfoByName(db, cardID)
		if err != nil {
			log.Printf("❌ Error fetching patient_info_id: %v", err)
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
			return
		}
		patientInfoID := patient.PatientInfo.PatientInfo_ID
		if patientInfoID == 0 {
			sendReply(bot, event.ReplyToken, "ไม่พบรหัสข้อมูลผู้ป่วย กรุณาลองใหม่")
			return
		}

		category, exists := userActivityCategory[State]
		if !exists {
			sendReply(bot, event.ReplyToken, "ไม่พบหมวดหมู่กิจกรรม กรุณาลองใหม่")
			return
		}

		userInfo, err := GetUserInfoByLINEID(db, State)
		if err != nil {
			log.Printf("Error fetching user info: %v", err)
			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่")
			return
		}
		userInfoID := userInfo.UserInfo_ID

		// ✅ ตรวจสอบว่าเป็น "มิติอื่นๆ" หรือไม่
		activityInfoID, exists := userActivityInfoID[State]
		activityOther := ""
		if category == "other" {
			activityOther = userActivity[State] // ใช้ชื่อกิจกรรมที่ผู้ใช้กรอก
		} else {
			if !exists || activityInfoID == 0 {
				sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมก่อนเริ่ม")
				return
			}
		}

		log.Printf("Creating new activity record for patient_info_id=%d, activity_info_id=%d, activityOther=%s, start_time=%v", patientInfoID, activityInfoID, activityOther, startTime)
		newRecordID, err := InsertActivityStartTime(db, patientInfoID, category, activityInfoID, activityOther, startTime, userInfoID)
		if err != nil {
			log.Printf("Error inserting activity start time: %v", err)
			return
		}

		userActivityRecordID[State] = newRecordID
	}

	userState[State] = "wait status ActivityEnd"
	log.Printf("Updating userState for %s to wait status ActivityEnd", State)
	sendReply(bot, event.ReplyToken, "กรุณากรอกเวลาสิ้นสุดของการทำกิจกรรม เช่น 11:20 ")
}

func handleActivityEnd(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ActivityEnd" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณากรอกเวลาสิ้นสุดกิจกรรม:")
		return
	}

	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	// รับเวลาสิ้นสุดจากผู้ใช้
	endTimeStr := strings.TrimSpace(message.Text)
	log.Printf("Received end time input: %s", endTimeStr)

	// ✅ ใช้ฟังก์ชัน parseTimeInput() แปลงค่า
	endTime, err := parseTimeInput(endTimeStr)
	if err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบเวลาไม่ถูกต้อง ตัวอย่างเช่น '11:05', '11.05', '11:05นาที'")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("❌ Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	// ดึง `activity_record_id`
	activityRecordID, exists := userActivityRecordID[State]
	if !exists {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
		return
	}

	// ดึง `user_info_id` ของผู้ใช้
	userID, err := GetUserInfoByLINEID(db, event.Source.UserID)
	if err != nil {
		log.Printf("❌ ไม่พบ user_info_id สำหรับ LINE ID: %s", event.Source.UserID)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่")
		return
	}

	// บันทึก `end_time` ลงฐานข้อมูล
	err = UpdateActivityEndTime(db, activityRecordID, endTime, userID.UserInfo_ID)
	if err != nil {
		log.Printf("❌ Error updating end_time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่.")
		return
	}

	// ส่งข้อความยืนยัน และขอให้ส่งรูปหลักฐาน
	sendReply(bot, event.ReplyToken, fmt.Sprintf("กรุณาส่งรูปก่อนการทำกิจกรรม"))

	// เปลี่ยนสถานะเป็น "wait status Saveavtivityend"
	userState[State] = "wait status Saveavtivityend"
}

// รับและประมวลผลรูปการทำกิจกรรม
func handleSaveavtivityend(bot *linebot.Client, event *linebot.Event, cardID, userID string) {
	if userState[userID] != "wait status Saveavtivityend" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะของคุณไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบประเภทข้อความว่าเป็นImage
	if imageMessage, ok := event.Message.(*linebot.ImageMessage); ok {
		log.Printf("Processing ImageMessage for user %s", userID)

		// อัปเดตสถานะ wait status saveEvidenceImageActivity เพื่อใช้เข้าในฟังก์ชัน saveEvidenceImageActivity
		userState[userID] = "wait status saveEvidenceImagebeforeActivity"
		log.Printf("Switching user state to 'wait status saveEvidenceImageActivity' for user %s", userID)

		//เข้าในฟังก์ชัน saveEvidenceImageActivity เพื่อบันทึกรูปการทำกิจกรรม
		if err := handlesaveEvidenceImagebeforeActivity(bot, event, cardID, userID, imageMessage); err != nil {
			log.Printf("Error saving image: %v", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกรูปภาพ กรุณาลองใหม่.")
			return
		}

		// หลังจากบันทึกรูปภาพเสร็จแล้ว เปลี่ยนสถานะเป็น wait status saveEvidenceTime เพื่อเข้าในฟังก์ชัน saveEvidenceImageTime
		userState[userID] = "wait status saveEvidenceImageafterActivity"
		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปหลังการทำกิจกรรม.")
	} else {
		log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
	}
}

// ฟังก์ชันบันทึกรูปภาพกิจกรรม
func handlesaveEvidenceImagebeforeActivity(bot *linebot.Client, event *linebot.Event, cardID, userID string, imageMessage *linebot.ImageMessage) error {
	if userState[userID] != "wait status saveEvidenceImagebeforeActivity" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return nil
	}
	//ตรวจสอบข้อความที่รับมา=Image
	messageID := imageMessage.ID
	log.Printf("Processing ImageMessage for user %s with message ID: %s", userID, messageID)

	// ดึงข้อมูลภาพ
	content, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		log.Printf("Error getting image content: %v", err)
		return err
	}
	defer content.Content.Close()

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return err
	}
	defer db.Close()

	// ตรวจสอบข้อมูลที่เกี่ยวข้อง
	patientInfoID, err := GetPatientInfoIDByCardID(db, cardID)
	if err != nil {
		return err
	}
	activity, err := GetActivityNameByPatientInfoID(db, patientInfoID)
	if err != nil {
		log.Printf("Error fetching activity name: %v", err)
		return err
	}
	// ใช้สำหรับจัดเก็บไฟล์ชั่วคราวระหว่างการดำเนินงาน
	tempDir := os.TempDir()
	tempFilePath := fmt.Sprintf("%s\\%s.jpg", tempDir, messageID)
	file, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		return err
	}
	defer file.Close()
	defer os.Remove(tempFilePath)

	// เขียนเนื้อหาภาพลงในไฟล์ (บันทึกเนื้อหาของรูปภาพหรือไฟล์ที่ได้รับจาก LINE Messaging API ลงในไฟล์ชั่วคราว)
	if _, err := io.Copy(file, content.Content); err != nil {
		log.Printf("Error writing image content to file: %v", err)
		return err
	}

	// เชื่อมต่อกับ MinIO
	minioClient, err := database.ConnectToMinio()
	if err != nil {
		log.Printf("Error connecting to MinIO: %v", err)
		return err
	}
	bucketName := "nirunimages"
	objectName := fmt.Sprintf("Imgage before activity/%s/%d/%s.jpg", activity, patientInfoID, messageID)

	// อัปโหลดไฟล์ไปยัง MinIO
	fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
	if err != nil {
		log.Printf("Error uploading file to MinIO: %v", err)
		return err
	}

	// อัปเดต URL ในฐานข้อมูล
	err = updateImagebeforeActivity(db, patientInfoID, fileURL)
	if err != nil {
		log.Printf("Error updating database: %v", err)
		return err
	}

	log.Printf("Activity Image successfully saved and URL updated: %s", fileURL)
	log.Printf("Last userID: %s", userID)
	return nil
}

// ฟังก์ชันบันทึกรูปหลังทำกิจกรรม
func handlesaveEvidenceImageafterActivity(bot *linebot.Client, event *linebot.Event, cardID, userID string) error {
	if userState[userID] != "wait status saveEvidenceImageafterActivity" {
		log.Printf("Unhandled state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
		return nil
	}

	//ตรวจสอบข้อความที่รับมา=Image
	if imageMessage, ok := event.Message.(*linebot.ImageMessage); ok {
		messageID := imageMessage.ID
		log.Printf("Processing ImageMessage for user %s with message ID: %s", userID, messageID)

		// ดึงข้อมูลรูปภาพ
		content, err := bot.GetMessageContent(messageID).Do()
		if err != nil {
			log.Printf("Error getting image content: %v", err)
			return err
		}
		defer content.Content.Close()

		db, err := database.ConnectToDB()
		if err != nil {
			log.Printf("Database connection error: %v", err)
			return err
		}
		defer db.Close()

		// ตรวจสอบข้อมูลที่เกี่ยวข้อง
		patientInfoID, err := GetPatientInfoIDByCardID(db, cardID)
		if err != nil {
			return err
		}
		activity, err := GetActivityNameByPatientInfoID(db, patientInfoID)
		if err != nil {
			log.Printf("Error fetching activity name: %v", err)
			return err
		}

		// บันทึกไฟล์รูปภาพชั่วคราว
		tempDir := os.TempDir()
		tempFilePath := fmt.Sprintf("%s\\%s.jpg", tempDir, messageID)
		file, err := os.Create(tempFilePath)
		if err != nil {
			log.Printf("Error creating temp file: %v", err)
			return err
		}
		defer file.Close()
		defer os.Remove(tempFilePath)

		if _, err := io.Copy(file, content.Content); err != nil {
			log.Printf("Error writing image content to file: %v", err)
			return err
		}

		// อัปโหลดรูปภาพไปยัง MinIO
		minioClient, err := database.ConnectToMinio()
		if err != nil {
			log.Printf("Error connecting to MinIO: %v", err)
			return err
		}
		bucketName := "nirunimages"
		objectName := fmt.Sprintf("Imgage after activity/%s/%d/%s.jpg", activity, patientInfoID, messageID)

		// อัปโหลดไฟล์ไปยัง MinIO
		fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
		if err != nil {
			log.Printf("Error uploading file to MinIO: %v", err)
			return err
		}

		// อัปเดต URL ในฐานข้อมูล
		err = updateImageafterActivity(db, patientInfoID, fileURL)
		if err != nil {
			log.Printf("Error updating database: %v", err)
			return err
		}

		log.Printf("Evidence time image successfully saved and URL updated: %s", fileURL)

		userState[userID] = "wait status ConfirmOrSaveEmployee"
		// handleUserChoiceForActivityRecord(bot, event, userID, "ยืนยันการบันทึก")
		log.Printf("User state updated to: %s", userState[userID])

		sendReply(bot, event.ReplyToken, "เลือก 'ยืนยันการบันทึก' หรือ 'บันทึกข้อมูลแทน'  ")
		return nil
	}

	log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
	sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
	return nil
}
func handleUserChoiceForActivityRecord(bot *linebot.Client, event *linebot.Event, userID, selection string) {
	if userState[userID] != "wait status ConfirmOrSaveEmployee" {
		log.Printf("❌ Invalid state for user %s. Current state: %s", userID, userState[userID])
		return
	}

	selection = strings.TrimSpace(selection)
	log.Printf("📌 Received selection: %s", selection)

	switch selection {
	case "บันทึกข้อมูลแทน", "save_employee", "บันทึกข้อมูล":
		log.Printf("✅ User %s selected to save record for another employee", userID)
		userState[userID] = "wait status saveActivityRecordForOtherEmployee"
		saveActivityRecordForOtherEmployee(bot, event, userID)

	case "ยืนยันการบันทึก", "confirm", "ยืนยัน":
		log.Printf("✅ User %s selected to confirm activity record", userID)
		userState[userID] = "wait status confirmActivityRecordByUser"
		confirmActivityRecordByUser(bot, event, userID)

	default:
		log.Printf("❌ Invalid selection by user %s: %s", userID, selection)
		sendReply(bot, event.ReplyToken, "❌ ตัวเลือกไม่ถูกต้อง กรุณาลองใหม่.\nกรุณาเลือก:\n- ยืนยันการบันทึก\n- บันทึกข้อมูลแทน")
	}
}

func confirmActivityRecordByUser(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status confirmActivityRecordByUser" {
		log.Printf("❌ Invalid state for user %s. Current state: %s", userID, userState[userID])
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("❌ Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ✅ ดึง activity_record_id ที่เพิ่งบันทึก
	activityRecord, err := GetLatestActivityRecord(db, userID)
	if err != nil {
		log.Printf("❌ ไม่พบข้อมูลกิจกรรมสำหรับ UserID: %s", userID)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
		return
	}

	// ✅ ตรวจสอบสิทธิ์พนักงาน
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil || userInfo.EmployeeInfo.EmployeeInfo_ID == 0 {
		log.Printf("❌ User %s is not an employee", userID)
		sendReply(bot, event.ReplyToken, "คุณไม่มีสิทธิ์ยืนยัน กรุณาเลือก 'บันทึกข้อมูลแทน'")
		return
	}

	// ✅ บันทึก employee_info_id ลง activity_record
	err = UpdateActivityEmployeeID(db, activityRecord.ActivityRecord_ID, userInfo.EmployeeInfo.EmployeeInfo_ID, userInfo.UserInfo_ID)
	if err != nil {
		log.Printf("❌ Error updating employee_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก กรุณาลองใหม่.")
		return
	}

	sendReply(bot, event.ReplyToken, "✅ บันทึกกิจกรรมสำเร็จ!")
	userState[userID] = ""
}

// ✅ ฟังก์ชันให้ผู้ใช้กรอกชื่อพนักงานแทน
func saveActivityRecordForOtherEmployee(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status saveActivityRecordForOtherEmployee" {
		log.Printf("❌ Invalid state for user %s. Current state: %s", userID, userState[userID])
		return
	}

	sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการแทน:")
	userState[userID] = "wait status SaveEmployeeName"
}

// ✅ ฟังก์ชันบันทึกชื่อพนักงานที่ทำการบริการแทน

func handleSaveEmployeeName(bot *linebot.Client, event *linebot.Event, userID, State, employeeName string) {
	// 🔍 ตรวจสอบค่าที่ได้รับ
	// log.Printf("📌 User %s entered employee name (before trim): '%s'", userID, employeeName)

	employeeName = strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Printf("Received employee name: %s", employeeName)
	if employeeName == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ.")
		return
	}
	// 🛑 ตรวจสอบว่าค่าที่รับมาเป็นค่าว่างหรือไม่
	if employeeName == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการแทน")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("❌ Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ✅ ค้นหา employeeID จากชื่อที่ถูกลบคำนำหน้าแล้ว
	employeeID, err := GetEmployeeIDByName(db, employeeName)
	if err != nil {
		log.Printf("❌ ไม่พบข้อมูล employee_info_id สำหรับพนักงาน: '%s'", employeeName)
		sendReply(bot, event.ReplyToken, fmt.Sprintf("ไม่พบพนักงานชื่อ %s\nกรุณากรอกชื่อใหม่", employeeName))
		return
	}

	log.Printf("✅ Employee ID found: %d for name: %s", employeeID, employeeName)

	// ✅ อัปเดตข้อมูลใน activity_record
	cardID := usercardidState[State]
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("❌ Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}
	patientInfoID := patient.PatientInfo.PatientInfo_ID
	if patientInfoID == 0 {
		sendReply(bot, event.ReplyToken, "ไม่พบรหัสข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}

	activityRecord, err := GetLatestActivityRecordByPatientID(db, patientInfoID)
	if err != nil {
		log.Printf("⚠️ ไม่พบข้อมูลกิจกรรมล่าสุดสำหรับ PatientInfo_ID: %d", patientInfoID)
		sendReply(bot, event.ReplyToken, "⚠️ ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่ หรือ ตรวจสอบว่ากิจกรรมถูกบันทึกเสร็จสิ้นแล้ว")
		return
	}

	log.Printf("✅ Updating employee_info_id=%d for ActivityRecord_ID=%d", employeeID, activityRecord.ActivityRecord_ID)

	err = UpdateActivityEmployeeID(db, activityRecord.ActivityRecord_ID, employeeID, 0)
	if err != nil {
		log.Printf("❌ Error updating employee_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก กรุณาลองใหม่.")
		return
	}

	sendReply(bot, event.ReplyToken, "✅ บันทึกข้อมูลพนักงานสำเร็จ!")
	userState[userID] = ""
}

// แปลงระยะเวลาของกิจกรรมเป็นชั่วโมงและนาที
// func formatDuration(d time.Duration) string {
// 	hours := int(d.Hours())
// 	minutes := int(d.Minutes()) % 60
// 	return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
// }

// ตรวจสอบกิจกรรมที่รับมาตรงกับฐานข้อมูลไหม
func validateActivity(activity string) bool {
	allowedActivities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กันฟัง", "ซุโดกุ", "จับคู่ภาพ",
	}
	for _, allowed := range allowedActivities {
		if activity == allowed {
			return true
		}
	}
	return false
}

func handleDefault(bot *linebot.Client, event *linebot.Event) {
	sendCustomReply(bot, event.ReplyToken, event.Source.UserID, "กรุณาเลือกเมนู")
}

// ส่งข้อความตอบกลับแบบธรรมดา
func sendReply(bot *linebot.Client, replyToken, message string) {
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Println("Error sending reply message:", err)
	}
}

// ส่งข้อความตอบกลับพร้อมปุ่ม
func sendReplyWithQuickReply(bot *linebot.Client, replyToken string, message string, quickReply *linebot.QuickReplyItems) {
	textMessage := linebot.NewTextMessage(message).WithQuickReplies(quickReply)
	if _, err := bot.ReplyMessage(replyToken, textMessage).Do(); err != nil {
		log.Printf("Error sending reply with quick reply: %v", err)
	}
}
