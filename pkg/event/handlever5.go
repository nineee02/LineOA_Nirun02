package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/models"
	"strconv"
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
	// log.Println("userstate:", State)

	switch text {
	case "NIRUN":
		handleNirunStste(bot, event, event.Source.UserID)
	case "ค้นหาข้อมูลผู้สูงอายุ":
		handleElderlyInfoStste(bot, event, event.Source.UserID)
	case "ลงเวลาเข้าและออกงาน":
		handleWorktimeStste(bot, event, event.Source.UserID)
	case "ประวัติการเข้ารับบริการ":
		handleServiceHistoryStste(bot, event, event.Source.UserID)
	case "บันทึกการเข้ารับบริการ":
		handleServiceRecordStste(bot, event, event.Source.UserID)
	// case "คู่มือการใช้งานระบบ":
	// 	handleSystemManualStste(bot, event, event.Source.UserID)
	default:
		handleDefault(bot, event)
	}

	state, exists := userState[State]
	if exists {
		switch state {
		// case "wait status NirunRequest":
		// 	handleNIRUN(bot, event, State)

		case "wait status worktime":
			handleWorktime(bot, event, State)
		// case "wait status worktimeConfirm":
		// 	handleworktimeConfirm(bot, event, State)
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
		case "wait status HistoryRequest":
			handleServiceHistory(bot, event, State)
		case "wait status HistoryAll":
			handleActivityHistoryAll(bot, event, State)
		case "wait status HistoryofYear":
			handleActivityHistoryofYear(bot, event, State)
		case "wait status HistoryofMonth":
			handleActivityHistoryofMonth(bot, event, State)
		case "wait status HistoryofWeek":
			handleActivityHistoryofWeek(bot, event, State)
		case "wait status HistoryofDay":
			handleActivityHistoryofDay(bot, event, State)
		case "wait status HistoryofSet":
			handleActivityHistoryofSet(bot, event, State)
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

func handleNirunStste(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status NirunRequest")
}

func handleElderlyInfoStste(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status ElderlyInfoRequest")
}

func handleServiceRecordStste(bot *linebot.Client, event *linebot.Event, State string) {
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

// func handleSystemManualStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status ManualRequest")
// }

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
	sendRegisterLink(bot, replyToken)
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
    // ตรวจสอบสถานะการลงทะเบียน
    if !isUserRegistered(userID) {
        sendQRCodeForLogin(bot, event.ReplyToken) // ส่ง QR Code สำหรับลงทะเบียน
        return
    }

    db, err := database.ConnectToDB()
    if err != nil {
        log.Println("Database connection error:", err)
        sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
        return
    }
    defer db.Close()

    // ดึงข้อมูลพนักงานที่เชื่อมโยงกับ LINE UserID
    employeeInfo, err := GetEmployeeByLINEID(db, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงานที่เชื่อมโยงกับบัญชี LINE นี้.")
        } else {
            log.Println("Error fetching employee info:", err)
            sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่.")
        }
        return
    }

    // ตรวจสอบสถานะเช็คอิน
    checkedIn, err := IsEmployeeCheckedIn(db, employeeInfo.EmployeeInfo_ID)
    if err != nil {
        log.Println("Error checking employee status:", err)
        sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
        return
    }

    // ตรวจสอบข้อความที่ส่งมา
    message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

    switch message {
    case "เช็คอิน":
        if checkedIn {
            sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
            return
        }
        userState[userID] = "wait status worktimeConfirmCheckIn"
        handleworktimeConfirmCheckIn(bot, event, userID)

    case "เช็คเอ้าท์":
        if !checkedIn {
            sendReply(bot, event.ReplyToken, "คุณยังไม่ได้เช็คอิน กรุณาเช็คอินก่อน")
            return
        }
        userState[userID] = "wait status worktimeConfirmCheckOut"
        handleworktimeConfirmCheckOut(bot, event, userID)

        var quickReply *linebot.QuickReplyItems

        sendReply(bot, event.ReplyToken, "กรุณาเลือก 'เช็คอิน' หรือ 'เช็คเอ้าท์'")
        quickReply = linebot.NewQuickReplyItems(
            linebot.NewQuickReplyButton("", linebot.NewMessageAction("เช็คอิน", "เช็คอิน")),
            linebot.NewQuickReplyButton("", linebot.NewMessageAction("เช็คเอ้าท์", "เช็คเอ้าท์")),
        )

        sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณาเลือกตัวเลือก:", quickReply)
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

	// ดึงข้อมูล employee_info_id จาก LINE User ID
	employeeInfo, err := GetEmployeeByLINEID(db, userID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
		return
	}

	// บันทึกการเช็คอิน
	err = RecordCheckIn(db, employeeInfo.EmployeeInfo_ID)
	if err != nil {
		log.Println("Error recording Check-in:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-in กรุณาลองใหม่.")
		return
	}

	sendReply(bot, event.ReplyToken, fmt.Sprintf("เช็คอินสำเร็จสำหรับ %s", employeeInfo.Name))
	userState[userID] = "wait status worktimeConfirmCheckOut"
}


func handleworktimeConfirmCheckOut(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status worktimeConfirmCheckOut" {
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

	// ดึงข้อมูล employee_info_id จาก LINE User ID
	employeeInfo, err := GetEmployeeByLINEID(db, userID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
		return
	}

	// บันทึกการเช็คเอาท์
	err = RecordCheckOut(db, employeeInfo.EmployeeInfo_ID)
	if err != nil {
		log.Println("Error recording Check-out:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-out กรุณาลองใหม่.")
		return
	}

	sendReply(bot, event.ReplyToken, fmt.Sprintf("เช็คเอาท์สำเร็จสำหรับ %s", employeeInfo.Name))
	userState[userID] = "wait status worktime"
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

	// ป้องกันการประมวลผลซ้ำ
	if userActivity[userID] == "processing" {
		sendReply(bot, event.ReplyToken, "คำขอของคุณกำลังดำเนินการ กรุณารอสักครู่.")
		return
	}
	userActivity[userID] = "processing"
	defer func() { userActivity[userID] = "" }() // รีเซ็ตสถานะเมื่อเสร็จสิ้น

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := event.Message.(*linebot.TextMessage).Text

	// ตรวจสอบ employeeID ใน state
	employeeIDStr, exists := userState[userID+"_employeeID"]
	if !exists {
		log.Println("Employee ID not found in state")
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลาเข้าและออกงาน'")
		return
	}

	// แปลง employeeID เป็นตัวเลข
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Println("Error parsing employee ID:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลพนักงาน.")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ตรวจสอบสถานะ Check-in ของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, employeeID)
	if err != nil {
		log.Println("Error checking employee status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}

	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลาเข้าและออกงาน'")
		return
	}

	// หากผู้ใช้ส่งข้อความ "ข้อมูลผู้สูงอายุ"
	if strings.TrimSpace(message) == "ค้นหาข้อมูลผู้สูงอายุ" {
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
	// log.Println("Message pateint:", message)

	// ตรวจสอบ employeeID ใน state
	employeeIDStr, exists := userState[State+"_employeeID"]
	if !exists {
		log.Println("Employee ID not found in state")
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลาเข้าและออกงาน'")
		return
	}

	// แปลง employeeID เป็นตัวเลข
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Println("Error parsing employee ID:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลพนักงาน.")
		return
	}

	// เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ตรวจสอบสถานะ Check-in ของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, employeeID)
	if err != nil {
		log.Println("Error checking employee status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลาเข้าและออกงาน'")
		return
	}

	// หากผู้ใช้ส่งข้อความ "บันทึกการเข้ารับบริการ"
	if strings.TrimSpace(message) == "บันทึกการเข้ารับบริการ" {
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
	// log.Println("reply Message Format: ", replyMessage)
	quickReplyActivities := createQuickReplyActivities()
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
	).Do(); err != nil {
		log.Printf("Error replying message (handleServiceInfo): %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
	}

	usercardidState[State] = cardID
	// log.Printf("Saved card_id for user %s: %s", State, cardID)

	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
	userState[State] = "wait status Activityrecord"
	log.Printf("Set user state to wait status Activityrecord for user %s", State)

}

func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status Activityrecord:", userState)
	// ตรวจสอบสถานะของผู้ใช้
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

	// ส่งข้อความให้ผู้ใช้กด "เริ่มกิจกรรม"
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

	// รับข้อความกิจกรรมจากผู้ใช้
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}

	starttime := strings.TrimSpace(message.Text)
	log.Printf("Received activity input: %s", starttime)

	if starttime != "เริ่มกิจกรรม" {
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
	employeeIDStr, exists := userState[State+"_employeeID"]
	if !exists || employeeIDStr == "" {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่")
		return
	}
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Printf("Error parsing employeeID: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการแปลงข้อมูลพนักงาน กรุณาลองใหม่")
		return
	}
	// ดึง patient_info_id โดยใช้ cardID
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}
	activityRecord := &models.Activityrecord{
		PatientInfo: models.PatientInfo{
			CardID:         cardID,
			Name:           patient.PatientInfo.Name, // กำหนดค่า Name จากฐานข้อมูล
			PatientInfo_ID: patient.PatientInfo.PatientInfo_ID,
		},
		ServiceInfo: models.ServiceInfo{
			Activity: userActivity[State],
		},
		EmployeeInfo: models.EmployeeInfo{
			EmployeeInfo_ID: employeeID,
		},
		StartTime: time.Now(),
	}
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

	// sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณากดปุ่ม 'เสร็จสิ้น' เมื่อทำกิจกรรมเสร็จ", quickReply)
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

	cardID, exists := usercardidState[userID]
	if !exists || cardID == "" {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลบัตรประชาชน กรุณากรอกใหม่")
		return
	}
	// ดึง patient_info_id โดยใช้ cardID
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}

	// บันทึก end_time และ employee_info_id
	activityRecord := &models.Activityrecord{
		PatientInfo: models.PatientInfo{
			CardID:         cardID,
			Name:           patient.PatientInfo.Name, // กำหนดค่า Name จากฐานข้อมูล
			PatientInfo_ID: patient.PatientInfo.PatientInfo_ID,
		},
		ServiceInfo: models.ServiceInfo{
			Activity: userActivity[userID],
		},
		EndTime: time.Now(),
	}
	if err := UpdateActivityEndTime(db, activityRecord); err != nil {
		log.Printf("Error updating end time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
		return
	}
	activityRecord.PatientInfo.Name = patient.PatientInfo.Name
	replyMessage := FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityEnd):", err)
	}
	log.Println("บันทึกกิจกรรมสำเร็จ :", replyMessage)

	resetUserState(userID)
	userState[userID] = ""
}

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

func handleServiceHistory(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status HistoryRequest" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	employeeIDStr, exists := userState[State+"_employeeID"]
	if !exists {
		log.Println("Employee ID not found in state")
		return
	}

	// แปลง employeeID เป็นตัวเลข
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Println("Error parsing employee ID:", err)
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	// ตรวจสอบสถานะ Check-in ของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, employeeID)
	if err != nil {
		log.Println("Error checking employee status:", err)
		return
	}
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลาเข้าและออกงาน'")
		return
	}

	// รับข้อความจากผู้ใช้
	message := event.Message.(*linebot.TextMessage).Text
	if message == "ประวัติการเข้ารับบริการ" {
		quickReply := linebot.NewQuickReplyItems(
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("ทั้งหมด", "ทั้งหมด")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("ปีนี้", "ปีนี้")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("เดือนนี้", "เดือนนี้")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("สัปดาห์นี้", "สัปดาห์นี้")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("วันนี้", "วันนี้")),
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("ระบุช่วงเวลา", "ระบุช่วงเวลา")),
		)

		sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณาเลือกประเภทการดูประวัติ:", quickReply)
		return
	}
	if message == "ทั้งหมด" {
		userState[State] = "wait status HistoryAll"
		handleActivityHistoryAll(bot, event, State)
		return
	} else if message == "ปีนี้" {
		userState[State] = "wait status HistoryofYear"
		handleActivityHistoryofYear(bot, event, State)
		return
	} else if message == "เดือนนี้" {
		userState[State] = "wait status HistoryofMonth"
		handleActivityHistoryofMonth(bot, event, State)
	} else if message == "สัปดาห์นี้" {
		userState[State] = "wait status HistoryofWeek"
		handleActivityHistoryofWeek(bot, event, State)
	} else if message == "วันนี้" {
		userState[State] = "wait status HistoryofDay"
		handleActivityHistoryofDay(bot, event, State)
	} else if message == "ระบุช่วงเวลา" {
		userState[State] = "wait status HistoryofSet"
		handleActivityHistoryofSet(bot, event, State)
		return
	} else {
		// sendReply(bot, event.ReplyToken, "กรุณาเลือกประเภแทการดูประวัติอีกครั้ง:")
		return
	}
	log.Printf("Set user state to %s for user %s", userState[State], State)
}
func handleActivityHistoryAll(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryoAll:", userState)
	if userState[State] != "wait status HistoryAll" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		// sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()
	historyRecords, err := historyALL(db)
	if err != nil {
		log.Println("Error fetching yearly activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลปีนี้.")
		return
	}
	replyMessage := FormatHistoryAll(historyRecords)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryAll):", err)
	}
	log.Println("ประวัติกิจกรรมทั้งหมด:", replyMessage)

}
func handleActivityHistoryofYear(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofYear:", userState)
	if userState[State] != "wait status HistoryofYear" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		// sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()
	historyRecords, err := historyOfYear(db)
	if err != nil {
		log.Println("Error fetching yearly activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลปีนี้.")
		return
	}
	replyMessage := FormatHistoryofYear(historyRecords)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryofYear):", err)
	}
	log.Println("ประวัติกิจกรรมปีนี้:", replyMessage)

}
func handleActivityHistoryofMonth(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofMonth:", userState)
	if userState[State] != "wait status HistoryofMonth" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		// sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()
	historyRecords, err := historyOfMonth(db)
	if err != nil {
		log.Println("Error fetching yearly activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลสัปดาห์นี้.")
		return
	}
	replyMessage := FormatHistoryofMonth(historyRecords)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryofMonth):", err)
	}
	log.Println("ประวัติกิจกรรมเดือนนี้:", replyMessage)

}
func handleActivityHistoryofWeek(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofWeek:", userState)
	if userState[State] != "wait status HistoryofWeek" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		// sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()
	historyRecords, err := historyOfWeek(db)
	if err != nil {
		log.Println("Error fetching yearly activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลสัปดาห์นี้.")
		return
	}
	replyMessage := FormatHistoryofWeek(historyRecords)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryofWeek):", err)
	}
	log.Println("ประวัติกิจกรรมสัปดาห์นี้:", replyMessage)

}
func handleActivityHistoryofDay(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofDay:", userState)
	if userState[State] != "wait status HistoryofDay" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()
	historyRecords, err := historyOfDay(db)
	if err != nil {
		log.Println("Error fetching yearly activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลสัปดาห์นี้.")
		return
	}
	replyMessage := FormatHistoryofDay(historyRecords)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryofDay):", err)
	}
	log.Println("ประวัติกิจกรรมวันนี้:", replyMessage)

}
func handleActivityHistoryofSet(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofSet:", userState)
	if userState[State] != "wait status HistoryofSet" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}
	message := event.Message.(*linebot.TextMessage).Text

	// แยกข้อความเพื่อดึงช่วงวันที่
	dates := strings.Split(message, "ถึง")
	if len(dates) != 2 {
		sendReply(bot, event.ReplyToken, "กรุณาระบุช่วงวันที่ในรูปแบบ YYYY-MM-DD ถึง YYYY-MM-DD\nตัวอย่าง: 2025-01-01 ถึง 2025-02-01")
		return
	}

	startDate := strings.TrimSpace(dates[0])
	endDate := strings.TrimSpace(dates[1])

	// ตรวจสอบรูปแบบวันที่
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบวันที่เริ่มต้นไม่ถูกต้อง กรุณาระบุในรูปแบบ YYYY-MM-DD เช่น 2025-01-01")
		return
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบวันที่สิ้นสุดไม่ถูกต้อง กรุณาระบุในรูปแบบ YYYY-MM-DD เช่น 2025-02-01")
		return
	}

	// เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	// ดึงข้อมูลจากฐานข้อมูล
	historyRecords, err := historyOfSet(db, startDate, endDate)
	if err != nil {
		log.Println("Error fetching activity history for set dates:", err)
		return
	}

	// ฟอร์แมตข้อความผลลัพธ์
	replyMessage := FormatHistoryofSet(historyRecords, startDate, endDate)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message (handleActivityHistoryofSet):", err)
	}
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
