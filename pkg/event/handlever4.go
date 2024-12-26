package event

// import (
// 	"database/sql"
// 	"log"
// 	"nirun/pkg/database"
// 	"nirun/pkg/models"
// 	"strconv"
// 	"strings"

// 	"github.com/line/line-bot-sdk-go/linebot"
// )

// func sendCustomReply(bot *linebot.Client, replyToken string, messages ...linebot.SendingMessage) {
// 	if len(messages) == 0 {
// 		return // ไม่ส่งข้อความถ้าไม่มี
// 	}

// 	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
// 		log.Printf("Error replying message sendCustomReply: %v", err)
// 	}
// }

// var usercardidState = make(map[string]string)
// var userState = make(map[string]string)
// var userActivity = make(map[string]string) // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน

// // HandleEvent - จัดการข้อความที่ได้รับจาก LINE
// func HandleEvent(bot *linebot.Client, event *linebot.Event) {
// 	text := event.Message.(*linebot.TextMessage).Text
// 	log.Println("Text: ", text)
// 	State := event.Source.UserID
// 	// log.Println("userstate:", State)

// 	switch text {
// 	case "NIRUN":
// 		handleNirunStste(bot, event, event.Source.UserID)
// 	case "ข้อมูลผู้สูงอายุ":
// 		handleElderlyInfoStste(bot, event, event.Source.UserID)
// 	case "ลงเวลาการทำงาน":
// 		handleWorktimeStste(bot, event, event.Source.UserID)
// 	// case "ประวัติการเข้ารับบริการ":
// 	// 	handleServiceHistoryStste(bot, event, event.Source.UserID)
// 	case "บันทึกการเข้ารับบริการ":
// 		handleServiceRecordStste(bot, event, event.Source.UserID)
// 	// case "คู่มือการใช้งานระบบ":
// 	// 	handleSystemManualStste(bot, event, event.Source.UserID)
// 	default:
// 		handleDefault(bot, event)
// 	}

// 	state, exists := userState[State]
// 	if exists {
// 		switch state {
// 		case "wait status NirunRequest":
// 			handleNIRUN(bot, event, State)
// 		case "wait status worktimeCheckIn":
// 			handleworktimeCheckIn(bot, event, State)
// 		case "wait status worktimeCheckOut":
// 			handleworktimeCheckOut(bot, event, State)
// 		case "wait status worktime":
// 			handleWorktime(bot, event, State)
// 		case "wait status ElderlyInfoRequest":
// 			handlePateintInfo(bot, event, State)
// 		case "wait status ServiceRecordRequest":
// 			handleServiceInfo(bot, event, State)
// 		case "wait status Activityrecord":
// 			handleActivityrecord(bot, event, State)
// 		case "wait status ActivityPeriodRecord":
// 			handleActivityPeriodRecord(bot, event, State)

// 		default:
// 			log.Printf("Unhandled state for user %s: %s", State, State)
// 		}
// 		return
// 	}

// }
// func setUserState(State, state string) {
// 	userState[State] = state
// 	// log.Printf("Set user state for user %s to %s", State, state)

// }

// func handleNirunStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status NirunRequest")
// }

// func handleElderlyInfoStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status ElderlyInfoRequest")
// }

// func handleServiceRecordStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status ServiceRecordRequest")
// }

// func handleWorktimeStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	setUserState(State, "wait status worktime")
// }

// // func handleServiceHistoryStste(bot *linebot.Client, event *linebot.Event, State string) {
// // 	setUserState(State, "wait status HistoryRequest")
// // }

// // func handleSystemManualStste(bot *linebot.Client, event *linebot.Event, State string) {
// // 	setUserState(State, "wait status ManualRequest")
// // }

// // ************************************************************************************************************************

// // **********************************************************************************************************
// func handleNIRUN(bot *linebot.Client, event *linebot.Event, State string) {
// 	sendReply(bot, event.ReplyToken, "ยินดีต้อนรับสู่ระบบ NIRUN! กรุณาเลือกเมนูที่ต้องการ.")
// }

// func handleWorktime(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status worktime" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}
// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Println("Message received:", message)

// 	if message == "ลงเวลาการทำงาน" {
// 		// sendReply(bot, event.ReplyToken, "กรุณากรอกรหัสพนักงาน :")
// 		return
// 	}

// 	employeeCode := strings.TrimSpace(message)
// 	if employeeCode == "" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกรหัสพนักงาน:")
// 		return
// 	}
// 	log.Println("รหัสพนักงาน:", employeeCode)

// 	if len(employeeCode) < 6 {
// 		sendReply(bot, event.ReplyToken, "รหัสพนักงานไม่ถูกต้อง กรุณากรอกใหม่:")
// 		return
// 	}
// 	// เชื่อมต่อฐานข้อมูล
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ดึง employee_info_id จาก employee_info
// 	employeeID, err := GetEmployeeID(db, employeeCode)
// 	if err != nil {
// 		log.Println("Error fetching employee ID:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลสำหรับรหัสพนักงาน.")
// 		return
// 	}
// 	log.Printf("Found employee_info_id: %d for employeeCode: %s", employeeID, employeeCode)

// 	// ดึงข้อมูล worktimeRecord
// 	worktimeRecord, err := GetWorktime(db, employeeCode)
// 	if err != nil {
// 		log.Println("Error fetching worktime record:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
// 		return
// 	}

// 	var replyMessage string
// 	var quickReply *linebot.QuickReplyItems
// 	userState[State+"_employeeID"] = strconv.Itoa(employeeID)
// 	userState[State+"_employeeCode"] = employeeCode
// 	log.Printf("Updated userState for %s: %+v", State, userState)
// 	// ตรวจสอบสถานะการทำงาน
// 	if worktimeRecord == nil {
// 		// ดึงข้อมูลพนักงานจาก employee_info
// 		employeeInfo, err := GetEmployeeInfo(db, employeeCode)
// 		if err != nil {
// 			log.Println("Error fetching employee info:", err)
// 			sendReply(bot, event.ReplyToken, "ไม่สามารถดึงข้อมูลพนักงานได้ กรุณาลองใหม่.")
// 			return
// 		}

// 		// สร้าง worktimeRecord ชั่วคราวสำหรับการฟอร์แมต
// 		tempWorktimeRecord := &models.WorktimeRecord{
// 			EmployeeInfo: *employeeInfo, // ใช้ข้อมูลพนักงานที่ดึงมาได้
// 		}

// 		// ฟอร์แมตข้อความ
// 		replyMessage = FormatConfirmationCheckIn(tempWorktimeRecord)

// 		confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยัน Check-in", "ยืนยัน Check-in"))
// 		cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
// 		quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

// 		// อัปเดต userState
// 		userState[State+"_employeeCode"] = employeeCode
// 		userState[State+"_employeeID"] = strconv.Itoa(employeeID)
// 		userState[State] = "wait status worktimeCheckIn"
// 	} else if worktimeRecord.CheckOut == "" {
// 		// มี Check-in แต่ยังไม่มี Check-out -> ให้ทำ Check-out
// 		replyMessage = FormatConfirmationCheckOut(worktimeRecord)

// 		confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยัน Check-out", "ยืนยัน Check-out"))
// 		cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
// 		quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

// 		userState[State] = "wait status worktimeCheckOut"
// 	} else {
// 		// มี Check-out -> ให้ทำ Check-in ใหม่
// 		replyMessage = FormatConfirmationCheckIn(worktimeRecord)

// 		confirmButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยืนยัน Check-in", "ยืนยัน Check-in"))
// 		cancelButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก"))
// 		quickReply = linebot.NewQuickReplyItems(confirmButton, cancelButton)

// 		userState[State] = "wait status worktimeCheckIn"
// 	}

// 	sendReplyWithQuickReply(bot, event.ReplyToken, replyMessage, quickReply)
// }

// // ลงเวลาเข้างาน
// func handleworktimeCheckIn(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status worktimeCheckIn" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}

// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	if message != "ยืนยัน Check-in" {
// 		sendReply(bot, event.ReplyToken, "กรุณายืนยัน Check-in หรือพิมพ์ 'ยกเลิก'")
// 		return
// 	}

// 	employeeIDStr, ok := userState[State+"_employeeID"]
// 	if !ok {
// 		log.Printf("Employee ID not found in userState for state: %s, userState: %+v", State, userState)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
// 		return
// 	}

// 	// แปลงค่า employeeID เป็น int
// 	employeeID, err := strconv.Atoi(employeeIDStr)
// 	if err != nil {
// 		log.Printf("Error converting employeeID: %s", employeeIDStr)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการประมวลผลข้อมูล กรุณาลองใหม่.")
// 		return
// 	}

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	err = RecordCheckIn(db, employeeID)
// 	if err != nil {
// 		log.Println("Error recording Check-in:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการ Check-in กรุณาลองใหม่.")
// 		return
// 	}

// 	worktimeRecord, err := GetWorktime(db, userState[State+"_employeeCode"])
// 	if err != nil {
// 		log.Println("Error fetching worktime record:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
// 		return
// 	}
// 	if worktimeRecord == nil {
// 		sendReply(bot, event.ReplyToken, "ไม่มีข้อมูลการทำงาน กรุณาลองใหม่.")
// 		return
// 	}

// 	replyMessage := FormatworktimeCheckin(worktimeRecord)
// 	sendReply(bot, event.ReplyToken, replyMessage)

// 	userState[State] = "wait status worktimeCheckOut"
// }

// // ลงเวลาออกงาน
// func handleworktimeCheckOut(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status worktimeCheckOut" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		sendReply(bot, event.ReplyToken, "สถานะปัจจุบันไม่ถูกต้อง กรุณาเริ่มใหม่.")
// 		return
// 	}

// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Println("Message received:", message)

// 	if message == "ยกเลิก" {
// 		sendReply(bot, event.ReplyToken, "การเช็คเอาท์ถูกยกเลิก.")
// 		userState[State] = "wait status worktimeCheckIn"
// 		return
// 	}

// 	if message != "ยืนยัน Check-out" {
// 		sendReply(bot, event.ReplyToken, "กรุณายืนยัน Check-out หรือพิมพ์ 'ยกเลิก'")
// 		return
// 	}

// 	employeeIDStr, ok := userState[State+"_employeeID"]
// 	if !ok {
// 		log.Printf("Employee ID not found in userState for state: %s, userState: %+v", State, userState)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณาลองใหม่.")
// 		return
// 	}

// 	// แปลงค่า employeeID เป็น int
// 	employeeID, err := strconv.Atoi(employeeIDStr)
// 	if err != nil {
// 		log.Printf("Error converting employeeID: %s", employeeIDStr)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการประมวลผลข้อมูล กรุณาลองใหม่.")
// 		return
// 	}


// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	err = RecordCheckOut(db, employeeID)
// 	if err != nil {
// 		log.Println("Error recording Check-out:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการ Check-out กรุณาลองใหม่.")
// 		return
// 	}

// 	worktimeRecord, err := GetWorktime(db, userState[State+"_employeeCode"])
// 	if err != nil {
// 		log.Println("Error fetching worktime record:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน กรุณาลองใหม่.")
// 		return
// 	}
// 	if worktimeRecord == nil {
// 		sendReply(bot, event.ReplyToken, "ไม่มีข้อมูลการทำงาน กรุณาลองใหม่.")
// 		return
// 	}
// 	replyMessage := FormatworktimeCheckout(worktimeRecord)
// 	sendReply(bot, event.ReplyToken, replyMessage)

// 	userState[State] = "wait status worktimeCheckIn"

// }

// func handlePateintInfo(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status ElderlyInfoRequest" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}

// 	// ดึงข้อความที่ผู้ใช้ส่งมา
// 	message := event.Message.(*linebot.TextMessage).Text

// 	// ตรวจสอบ employeeID ใน state
// 	employeeIDStr, exists := userState[State+"_employeeID"]
// 	if !exists {
// 		log.Println("Employee ID not found in state")
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลพนักงาน กรุณา Check-in ก่อน.")
// 		return
// 	}

// 	// แปลง employeeID เป็นตัวเลข
// 	employeeID, err := strconv.Atoi(employeeIDStr)
// 	if err != nil {
// 		log.Println("Error parsing employee ID:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลพนักงาน.")
// 		return
// 	}

// 	// ตรวจสอบสถานะ Check-in ของพนักงาน
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	checkedIn, err := IsEmployeeCheckedIn(db, employeeID)
// 	if err != nil {
// 		log.Println("Error checking employee status:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
// 		return
// 	}
// 	if !checkedIn {
// 		sendReply(bot, event.ReplyToken, "ไม่พบสถานะ Check-in สำหรับพนักงานคนนี้ กรุณา Check-in ก่อน.")
// 		return
// 	}

// 	// หากผู้ใช้ส่งข้อความ "ข้อมูลผู้สูงอายุ"
// 	if strings.TrimSpace(message) == "ข้อมูลผู้สูงอายุ" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน :")
// 		return
// 	}

// 	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
// 	cardID := strings.TrimSpace(message)
// 	if cardID == "" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน:")
// 		return
// 	}
// 	log.Println("เลขประจำตัวประชาชน:", cardID)

// 	// ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
// 	patient, err := GetPatientInfoByName(db, cardID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วยที่ท่านค้นหา กรุณาตรวจสอบเลขประจำตัวผู้ป่วยอีกครั้ง.")
// 		} else {
// 			log.Println("Error fetching patient info:", err)
// 			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการค้นหาข้อมูลผู้ป่วย กรุณาลองใหม่.")
// 		}
// 		return
// 	}

// 	// ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	replyMessage := FormatPatientInfo(patient)
// 	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
// 		log.Println("Error replying message (handlePateintInfo):", err)
// 	}
// 	log.Println("ข้อมูลผู้สูงอายุ :", replyMessage)

// 	// รีเซ็ตสถานะผู้ใช้
// 	userState[State] = ""
// }

// func handleServiceHistory(bot *linebot.Client, event *linebot.Event) {
// 	sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชนเพื่อดูประวัติการเข้ารับบริการ:")
// }

// func handleServiceInfo(bot *linebot.Client, event *linebot.Event, State string) {
// 	// message := event.Message.(*linebot.TextMessage).Text
// 	// // log.Println("Message pateint:", message)

// 	// if message == "บันทึกการเข้ารับบริการ" {
// 	// 	sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน :")
// 	// 	return
// 	// }

// 	// cardID := strings.TrimSpace(message)
// 	// if cardID == "" {
// 	// 	sendReply(bot, event.ReplyToken, "กรุณากรอกเลขประจำตัวประชาชน:")
// 	// 	return
// 	// }
// 	// log.Println("เลขประจำตัวประชาชน:", cardID)

// 	// // เชื่อมต่อกับฐานข้อมูลและค้นหาข้อมูลผู้ป่วย
// 	// db, err := database.ConnectToDB()
// 	// if err != nil {
// 	// 	log.Println("Database connection error:", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	// // ค้นหาข้อมูลผู้ป่วยจากฐานข้อมูล
// 	// service, err := GetServiceInfoBycardID(db, cardID)
// 	// if err != nil {
// 	// 	log.Println("Error models.GetServiceInfoBycardID:", err)
// 	// 	// sendErrorReply(bot, event, "No patient information found for the provided name.")
// 	// 	return
// 	// }

// 	// // ส่งข้อมูลผู้ป่วยกลับไปยังผู้ใช้
// 	// replyMessage := FormatServiceInfo(service)
// 	// // log.Println("reply Message Format: ", replyMessage)
// 	// quickReplyActivities := createQuickReplyActivities()

// 	// if _, err := bot.ReplyMessage(
// 	// 	event.ReplyToken,
// 	// 	linebot.NewTextMessage(replyMessage).WithQuickReplies(&quickReplyActivities),
// 	// ).Do(); err != nil {
// 	// 	log.Println("Error replying message:", err)
// 	// }

// 	// usercardidState[State] = cardID
// 	// // log.Printf("Saved card_id for user %s: %s", State, cardID)

// 	// // เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
// 	// userState[State] = "wait status Activityrecord"
// 	// //log.Printf("Set user state to wait status Activityrecord for user %s", State)
// }
// func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
// 	// // ตรวจสอบว่าผู้ใช้อยู่ในสถานะที่ถูกต้องหรือไม่
// 	// if userState[State] != "wait status Activityrecord" {
// 	// 	log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 	// 	sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
// 	// 	return
// 	// }

// 	// // รับข้อความกิจกรรมจากผู้ใช้
// 	// message, ok := event.Message.(*linebot.TextMessage)
// 	// if !ok {
// 	// 	log.Println("Event is not a text message")
// 	// 	return
// 	// }

// 	// activity := strings.TrimSpace(message.Text)
// 	// log.Printf("Received activity input: %s", activity)

// 	// if !validateActivity(activity) {
// 	// 	sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
// 	// 	return
// 	// }

// 	// userActivity[State] = activity // เก็บกิจกรรมที่ผู้ใช้กรอก
// 	// log.Printf("Stored activity for user %s: %s", State, activity)
// 	// // ดึง card_id จาก userState
// 	// cardID := usercardidState[State] // ปรับใช้ให้ตรงตามข้อมูลที่เก็บไว้

// 	// // เชื่อมต่อกับฐานข้อมูล
// 	// db, err := database.ConnectToDB()
// 	// if err != nil {
// 	// 	log.Printf("Database connection error: %v", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	// // สร้างตัวแปร activityRecord
// 	// activityRecord := &models.Activityrecord{
// 	// 	CardID:   usercardidState[State],
// 	// 	Activity: activity,
// 	// }

// 	// // แสดง log card_id และ activity
// 	// log.Printf("Using card_id: %s, Activity: %s", cardID, activity)

// 	// // บันทึกกิจกรรมใหม่ GetActivityRecord(db, activityRecord)
// 	// if err := GetActivityRecord(db, activityRecord); err != nil { // ส่ง activityRecord ไปที่ฟังก์ชัน ActivityRecord
// 	// 	log.Printf("Error saving activity(models.ActivityRecord): %v", err)
// 	// 	sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", activity))
// 	// 	return
// 	// }

// 	// //sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรม '%s' สำเร็จสำหรับ card_id '%s'!", activity, cardID))

// 	// // รีเซ็ตสถานะผู้ใช้
// 	// userState[State] = "wait status ActivityPeriodRecord"
// 	// sendReply(bot, event.ReplyToken, "กรุณากรอกระยะเวลาในการทำกิจกรรม (เช่น 30นาที หรือ 1ชั่วโมง 10นาที):")
// }

// func handleActivityPeriodRecord(bot *linebot.Client, event *linebot.Event, State string) {
// 	// if userState[State] != "wait status ActivityPeriodRecord" {
// 	// 	log.Printf("สถานะไม่ถูกต้อง %s. สถานะปัจจุบัน: %s", State, userState[State])
// 	// 	sendReply(bot, event.ReplyToken, "กรอกระยะเวลาในการทำกิจกรรม:")
// 	// 	return
// 	// }

// 	// message, ok := event.Message.(*linebot.TextMessage)
// 	// if !ok {
// 	// 	log.Println("Event is not a text message")
// 	// 	return
// 	// }

// 	// period := strings.TrimSpace(message.Text)
// 	// log.Printf("Received activity input: %s", period)

// 	// if period == "" {
// 	// 	sendReply(bot, event.ReplyToken, "กรุณากรอกระยะเวลาในการทำกิจกรรม (เช่น 30นาที หรือ 1ชั่วโมง 10นาที):")
// 	// 	return
// 	// }
// 	// //ตรวจสอบว่าในuserActivity[State] มีชื่อกิจกรรมหรือยัง
// 	// activity, exists := userActivity[State]
// 	// if !exists || activity == "" {
// 	// 	log.Printf("No activity found for user %s", State)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่พบกิจกรรม กรุณากรอกกิจกรรมใหม่")
// 	// 	return
// 	// }

// 	// db, err := database.ConnectToDB()
// 	// if err != nil {
// 	// 	log.Printf("Database connection error: %v", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	// activityRecord := &models.Activityrecord{
// 	// 	Activity: activity,
// 	// 	Period:   period,
// 	// 	EndTime:  time.Now(),
// 	// }

// 	// //บันทึกระยะเวลา
// 	// if err := GetActivityPeriodRecord(db, activityRecord); err != nil {
// 	// 	log.Printf("Error saving period(models.ActivityPeriodRecord): %v", err)
// 	// 	sendReply(bot, event.ReplyToken, fmt.Sprintf("เกิดข้อผิดพลาดในการบันทึกกิจกรรม '%s' กรุณาลองใหม่", period))
// 	// 	return
// 	// }

// 	// sendReply(bot, event.ReplyToken, fmt.Sprintf("บันทึกกิจกรรมสำเร็จ !!!"))

// 	// // รีเซ็ตสถานะผู้ใช้
// 	// userState[State] = ""
// }

// func handleSystemManual(bot *linebot.Client, event *linebot.Event, State string) {
// 	sendReply(bot, event.ReplyToken, "คุณสามารถดูคู่มือการใช้งานระบบได้ที่ลิงก์: https://example.com/manual")
// }

// func handleDefault(bot *linebot.Client, event *linebot.Event) {
// 	sendCustomReply(bot, event.ReplyToken)
// }

// // ฟังก์ชันสำหรับส่งข้อความ
// func sendReply(bot *linebot.Client, replyToken, message string) {
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
// 		log.Println("Error sending reply message:", err)
// 	}
// }
// func sendReplyWithQuickReply(bot *linebot.Client, replyToken string, message string, quickReply *linebot.QuickReplyItems) {
// 	textMessage := linebot.NewTextMessage(message).WithQuickReplies(quickReply)
// 	if _, err := bot.ReplyMessage(replyToken, textMessage).Do(); err != nil {
// 		log.Printf("Error sending reply with quick reply: %v", err)
// 	}
// }
