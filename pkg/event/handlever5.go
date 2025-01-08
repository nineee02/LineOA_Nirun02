package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/flex"
	"nirun/pkg/models"
	"strconv"
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
var userActivity = make(map[string]string)    // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน
var userCheckInStatus = make(map[string]bool) // เก็บสถานะการเช็คอินของแต่ละบัญชี LINE

// ตรวจสอบสถานะการเช็คอินของบัญชี
func isUserCheckedIn(userID string) bool {
	status, exists := userCheckInStatus[userID]
	return exists && status
}

// อัปเดตสถานะการเช็คอินของบัญชี
func setUserCheckInStatus(userID string, status bool) {
	userCheckInStatus[userID] = status
	log.Printf("User %s check-in status set to %v", userID, status)
}
func setUserCheckOutStatus(userID string, status bool) {
	userCheckInStatus[userID] = status
	log.Printf("User %s check-in status set to %v", userID, status)
}

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
		case "wait status NirunRequest":
			handleNIRUN(bot, event, State)
		case "wait status worktimeCheckIn":
			handleworktimeCheckIn(bot, event, State)
		case "wait status worktimeCheckOut":
			handleworktimeCheckOut(bot, event, State)
		case "wait status worktime":
			handleWorktime(bot, event, State)
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
func setUserState(State, state string) {
	userState[State] = state
	// log.Printf("Set user state for user %s to %s", State, state)

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
	// sendReply(bot, event.ReplyToken, "กรุณากรอกรหัสพนักงาน:")
}

func handleServiceHistoryStste(bot *linebot.Client, event *linebot.Event, State string) {
	setUserState(State, "wait status HistoryRequest")
}

// func handleSystemManualStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status ManualRequest")
// }

// ************************************************************************************************************************

// **********************************************************************************************************
func handleNIRUN(bot *linebot.Client, event *linebot.Event, State string) {
	sendReply(bot, event.ReplyToken, "https://community.app.nirun.life/web/login#action=348&model=ni.patient&view_type=kanban&menu_id=257")
}

func handleWorktime(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status worktime" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}
	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Println("Message received:", message)

	if message == "ลงเวลาเข้าและออกงาน" {
		// sendReply(bot, event.ReplyToken, "กรุณากรอกรหัสพนักงาน :")
		return
	}
	log.Println("worktime", message)

	employeeCode := strings.TrimSpace(message)
	if employeeCode == "" {
		// sendReply(bot, event.ReplyToken, "กรุณากรอกรหัสพนักงาน:")
		return
	}
	log.Println("รหัสพนักงาน:", employeeCode)

	if len(employeeCode) < 6 {
		sendReply(bot, event.ReplyToken, "รหัสพนักงานไม่ถูกต้อง กรุณากรอกใหม่:")
		return
	}
	// เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	// ดึง employee_info_id จาก employee_info
	employeeID, err := GetEmployeeID(db, employeeCode)
	if err != nil {
		log.Println("Error fetching employee ID:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลสำหรับรหัสพนักงาน.")
		return
	}
	log.Printf("Found employee_info_id: %d for employeeCode: %s", employeeID, employeeCode)

	// ดึงข้อมูล worktimeRecord
	worktimeRecord, err := GetWorktime(db, employeeCode)
	if err != nil {
		log.Println("Error fetching worktime record:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}

	var replyMessage string
	var quickReply *linebot.QuickReplyItems
	userState[State+"_employeeID"] = strconv.Itoa(employeeID)
	userState[State+"_employeeCode"] = employeeCode
	log.Printf("Updated userState for %s: %+v", State, userState)
	// ตรวจสอบสถานะการทำงาน
	if worktimeRecord == nil {
		// ดึงข้อมูลพนักงานจาก employee_info
		employeeInfo, err := GetEmployeeInfo(db, employeeCode)
		if err != nil {
			log.Println("Error fetching employee info:", err)
			sendReply(bot, event.ReplyToken, "ไม่สามารถดึงข้อมูลพนักงานได้ กรุณาลองใหม่.")
			return
		}

		// สร้าง worktimeRecord ชั่วคราวสำหรับการฟอร์แมต
		tempWorktimeRecord := &models.WorktimeRecord{
			EmployeeInfo: *employeeInfo, // ใช้ข้อมูลพนักงานที่ดึงมาได้
		}

		// ฟอร์แมตข้อความ
		replyMessage = FormatConfirmationCheckIn1(tempWorktimeRecord)

		// confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยัน Check-in", "ยืนยัน Check-in"))
		// cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
		// quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

		// อัปเดต userState
		userState[State+"_employeeCode"] = employeeCode
		userState[State+"_employeeID"] = strconv.Itoa(employeeID)
		userState[State] = "wait status worktimeCheckIn"
	} else if worktimeRecord.CheckOut == (time.Time{}) {
		// มี Check-in แต่ยังไม่มี Check-out -> ให้ทำ Check-out
		flexMessage := flex.FormatConfirmationCheckOut(worktimeRecord)
		if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
			log.Println("Error replying message:", err)
		}
		log.Println("ยืนยันเช็คเอ้าท์", flexMessage)

		// confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยันการออกงาน", "ยืนยันการออกงาน"))
		// cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
		// quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

		userState[State] = "wait status worktimeCheckOut"
	} else {
		// มี Check-out -> ให้ทำ Check-in ใหม่
		flexMessage := flex.FormatConfirmationCheckIn(worktimeRecord)
		if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
			log.Println("Error replying message:", err)
		}
		log.Println("ยืนยันเช็คอิน", flexMessage)

		// confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยัน Check-in", "ยืนยัน Check-in"))
		// cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
		// quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

		userState[State] = "wait status worktimeCheckIn"
	}

	sendReplyWithQuickReply(bot, event.ReplyToken, replyMessage, quickReply)
}

// ลงเวลาเข้างาน
func handleworktimeCheckIn(bot *linebot.Client, event *linebot.Event, State string) {
	userID := event.Source.UserID

	// ตรวจสอบสถานะการเช็คอิน
	if isUserCheckedIn(userID) {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว\nกรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}

	if userState[State] != "wait status worktimeCheckIn" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	if message != "เช็คอิน" {
		sendReply(bot, event.ReplyToken, "กรุณาพิมพ์ 'เช็คอิน' เพื่อยืนยันการเข้างาน")
		return
	}

	employeeIDStr, ok := userState[State+"_employeeID"]
	if !ok {
		log.Printf("Employee ID not found in userState for state: %s, userState: %+v", State, userState)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
		return
	}

	// แปลงค่า employeeID เป็น int
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Printf("Error converting employeeID: %s", employeeIDStr)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการประมวลผลข้อมูล กรุณาลองใหม่.")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	err = RecordCheckIn(db, employeeID)
	if err != nil {
		log.Println("Error recording Check-in:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการ Check-in กรุณาลองใหม่.")
		return
	}

	worktimeRecord, err := GetWorktime(db, userState[State+"_employeeCode"])
	if err != nil {
		log.Println("Error fetching worktime record:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}
	if worktimeRecord == nil {
		sendReply(bot, event.ReplyToken, "ไม่มีข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}
	flexMessage := flex.FormatworktimeCheckin(worktimeRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	}
	log.Println("เช็คอินสำเร็จ :", flexMessage)
	// replyMessage := FormatworktimeCheckin(worktimeRecord)
	// sendReply(bot, event.ReplyToken, replyMessage)
	setUserCheckInStatus(userID, true)
	sendReply(bot, event.ReplyToken, "เช็คอินสำเร็จ")
	userState[State] = "wait status worktimeCheckOut"
}

// ลงเวลาออกงาน
func handleworktimeCheckOut(bot *linebot.Client, event *linebot.Event, State string) {

	if userState[State] != "wait status worktimeCheckOut" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "สถานะปัจจุบันไม่ถูกต้อง กรุณาเริ่มใหม่.")
		return
	}
	userID := event.Source.UserID
	// ตรวจสอบสถานะการเช็คอิน
	if isUserCheckedIn(userID) {
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว\nกรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}
	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Println("Message received:", message)

	if message == "ยกเลิก" {
		sendReply(bot, event.ReplyToken, "การเช็คเอาท์ถูกยกเลิก.")
		userState[State] = "wait status worktimeCheckIn"
		return
	}

	if message != "เช็คเอ้าท์" {
		sendReply(bot, event.ReplyToken, "กรุณายืนยันการออกงานก่อนทำการเช็คอินอีกครั้ง")
		return
	}

	employeeIDStr, ok := userState[State+"_employeeID"]
	if !ok {
		log.Printf("Employee ID not found in userState for state: %s, userState: %+v", State, userState)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
		return
	}

	// แปลงค่า employeeID เป็น int
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Printf("Error converting employeeID: %s", employeeIDStr)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการประมวลผลข้อมูล กรุณาลองใหม่.")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	err = RecordCheckOut(db, employeeID)
	if err != nil {
		log.Println("Error recording Check-out:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการ Check-out กรุณาลองใหม่.")
		return
	}

	worktimeRecord, err := GetWorktime(db, userState[State+"_employeeCode"])
	if err != nil {
		log.Println("Error fetching worktime record:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}
	if worktimeRecord == nil {
		sendReply(bot, event.ReplyToken, "ไม่มีข้อมูลการทำงาน กรุณาลองใหม่.")
		return
	}
	flexMessage := flex.FormatworktimeCheckout(worktimeRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	}
	log.Println("เช็คเอ้าท์สำเร็จ :", flexMessage)
	userCheckInStatus[event.Source.UserID] = false
	userState[State] = "wait status worktimeCheckIn"

}

func handlePateintInfo(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ElderlyInfoRequest" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	// ดึงข้อความที่ผู้ใช้ส่งมา
	message := event.Message.(*linebot.TextMessage).Text

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
		// sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
	cardID := strings.TrimSpace(message)
	if cardID == "" {
		// sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
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

	flexMessage := flex.FormatPatientInfo(patient)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	}

	// รีเซ็ตสถานะผู้ใช้
	userState[State] = ""
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
		// sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	cardID := strings.TrimSpace(message)
	if cardID == "" {
		// sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
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

	// replyMessage := FormatServiceInfo([]models.Activityrecord{*service})
	// log.Println("ข้อมูลผู้สูงอายุ :", replyMessage)
	// // log.Println("reply Message Format: ", replyMessage)
	// quickReplyActivities := createQuickReplyActivities()
	// if _, err := bot.ReplyMessage(
	// 	event.ReplyToken,
	// 	linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
	// ).Do(); err != nil {
	// 	log.Printf("Error replying message (handleServiceInfo): %v", err)
	// 	sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
	// }

	// สร้าง Flex Message
	flexMessage := flex.FormatServiceInfo([]models.Activityrecord{*service})

	// // สร้าง Quick Replies
	// quickReplyActivities := createQuickReplyActivities()

	// // แนบ Quick Replies กับ Flex Message
	// flexMessage.WithQuickReplies(&quickReplyActivities)

	// ส่งข้อความ
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		flexMessage,
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

	bubble := flex.FormatStartActivity(activity)
	flexMessage := linebot.NewFlexMessage("เลือกกิจกรรม", bubble)

	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	} else {
		log.Println("เลือกกิจกรรมสำเร็จ:", flexMessage)
	}

	// // สร้าง Quick Reply สำหรับเริ่มกิจกรรม
	// quickReply := linebot.NewQuickReplyItems(
	// 	linebot.NewQuickReplyButton("", linebot.NewMessageAction("เริ่มกิจกรรม", "เริ่มกิจกรรม")),
	// )

	// // ส่งข้อความให้ผู้ใช้กด "เริ่มกิจกรรม"
	// sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณากดปุ่ม 'เริ่มกิจกรรม' เพื่อเริ่มบันทึกเวลา", quickReply)

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

	flexMessage := flex.FormatactivityRecordStarttime(activityRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	}
	log.Println("บันทึกกิจกรรมสำเร็จ :", flexMessage)

	quickReply := linebot.NewQuickReplyItems(
	// linebot.NewQuickReplyButton("", linebot.NewMessageAction("เสร็จสิ้น", "เสร็จสิ้น")),
	)
	_, err = bot.ReplyMessage(
		event.ReplyToken,
		flexMessage.WithQuickReplies(quickReply),
	).Do()
	if err != nil {
		log.Println("Error sending Quick Reply (handleActivityStart):", err)
	}

	// sendReplyWithQuickReply(bot, event.ReplyToken, "กรุณากดปุ่ม 'เสร็จสิ้น' เมื่อทำกิจกรรมเสร็จ", quickReply)
	userState[State] = "wait status ActivityEnd"
}
func handleActivityEnd(bot *linebot.Client, event *linebot.Event, State string) {
	if userState[State] != "wait status ActivityEnd" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
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
		// sendReply(bot, event.ReplyToken, "กรุณาพิมพ์ 'เสร็จสิ้น' เพื่อบันทึกเวลาสิ้นสุด")
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
			Activity: userActivity[State],
		},
		EndTime: time.Now(),
	}
	startTime, err := GetActivityStartTime(db, cardID, userActivity[State])
	if err != nil {
		log.Printf("Error fetching StartTime: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลเวลาเริ่ม กรุณาลองใหม่")
		return
	}

	duration := activityRecord.EndTime.Sub(startTime)
	activityRecord.Period = formatDuration(duration)

	if err := UpdateActivityEndTime(db, activityRecord); err != nil {
		log.Printf("Error updating end time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
		return
	}

	// สร้างและส่ง Flex Message
	flexMessage := flex.FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Println("Error replying message:", err)
	}
	log.Println("บันทึกกิจกรรมสำเร็จ", flexMessage)

	if err := UpdateActivityEndTime(db, activityRecord); err != nil {
		log.Printf("Error updating end time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
		return
	}
	userCheckInStatus[event.Source.UserID] = true
	sendReply(bot, event.ReplyToken, "บันทึกกิจกรรมสำเร็จ! ขอบคุณที่ใช้บริการ")
	userState[State] = ""
}
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
}

func createQuickReplyActivities() linebot.QuickReplyItems {
	activities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กันฟัง", "ซุโดกุ", "จับคู่ภาพ",
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
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กันฟัง", "ซุโดกุ", "จับคู่ภาพ",
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
    flexMessage :=flex.FormatHistoryType()
    if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
        log.Printf("Failed to send Flex Message: %v", err)
    }
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
		// sendReply(bot, event.ReplyToken, "กรุณาระบุช่วงวันที่ในรูปแบบ YYYY-MM-DD ถึง YYYY-MM-DD \nเช่น2025-01-01 ถึง 2025-02-01:")
		return
	} else {
		sendReply(bot, event.ReplyToken, "กรุณาเลือกประเภทการดูประวัติอีกครั้ง:")
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

	sendReply(bot, event.ReplyToken, replyMessage)
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

	sendReply(bot, event.ReplyToken, replyMessage)
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

	sendReply(bot, event.ReplyToken, replyMessage)
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

	sendReply(bot, event.ReplyToken, replyMessage)
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
		// sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
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

	sendReply(bot, event.ReplyToken, replyMessage)
}
func handleActivityHistoryofSet(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status HistoryofSet:", userState)
	if userState[State] != "wait status HistoryofSet" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		return
	}

	message := event.Message.(*linebot.TextMessage).Text
	log.Println("Received message:", message)

	// แยกข้อความเพื่อดึงช่วงวันที่
	dates := strings.Split(message, "ถึง")
	log.Println("Split dates:", dates)

	if len(dates) != 2 {
		sendReply(bot, event.ReplyToken, "กรุณาระบุช่วงวันที่ในรูปแบบ YYYY-MM-DD ถึง YYYY-MM-DD\nตัวอย่าง: 2025-01-01 ถึง 2025-02-01")
		return
	}

	startDate := strings.TrimSpace(dates[0])
	endDate := strings.TrimSpace(dates[1])
	log.Printf("Parsed startDate: %s, endDate: %s", startDate, endDate)

	// ตรวจสอบรูปแบบวันที่
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบวันที่เริ่มต้นไม่ถูกต้อง กรุณาระบุในรูปแบบ YYYY-MM-DD เช่น 2025-01-01")
		return
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		sendReply(bot, event.ReplyToken, "รูปแบบวันที่สิ้นสุดไม่ถูกต้อง กรุณาระบุในรูปแบบ YYYY-MM-DD เช่น 2025-02-01")
		return
	}

	startDateWithTime := startDate + " 00:00:00"
	endDateWithTime := endDate + " 23:59:59"
	log.Printf("StartDateWithTime: %s, EndDateWithTime: %s", startDateWithTime, endDateWithTime)

	// เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Println("Database connection error:", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
		return
	}
	defer db.Close()

	historyRecords, err := historyOfSet(db, startDateWithTime, endDateWithTime)
	if err != nil {
		log.Println("Error fetching activity history:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูล กรุณาลองใหม่.")
		return
	}
	log.Printf("Fetched history records: %+v", historyRecords)

	if len(historyRecords) == 0 {
		sendReply(bot, event.ReplyToken, "ไม่มีข้อมูลกิจกรรมในช่วงวันที่ที่คุณระบุ")
		return
	}

	replyMessage := FormatHistoryofSet(historyRecords, startDate, endDate)
	log.Printf("Formatted reply message: %s", replyMessage)

	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Println("Error replying message:", err)
	}
}

func handleSystemManual(bot *linebot.Client, event *linebot.Event, State string) {
	sendReply(bot, event.ReplyToken, "คุณสามารถดูคู่มือการใช้งานระบบได้ที่ลิงก์: https://example.com/manual")
}

func handleDefault(bot *linebot.Client, event *linebot.Event) {
	sendCustomReply(bot, event.ReplyToken)
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
