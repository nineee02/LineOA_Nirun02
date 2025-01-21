package event

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"nirun/pkg/database"
	"nirun/pkg/flexmessage"
	"nirun/pkg/models"
	"os"
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
	// ตรวจสอบประเภทข้อความก่อน
	switch message := event.Message.(type) {
	case *linebot.TextMessage: // อ่านข้อความจาก TextMessage
		text := strings.TrimSpace(message.Text)
		log.Println("Received TextMessage:", text)
		State := event.Source.UserID
		log.Println("User state:", State)

		// ตรวจสอบคำสั่งจากข้อความ
		switch text {
		case "ค้นหาข้อมูล":
			handleElderlyInfoStste(bot, event, State)
		case "ลงเวลางาน":
			handleWorktimeStste(bot, event, State)
		case "บันทึกกิจกรรม":
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
			case "wait status SaveEmployeeName":
				handleSaveEmployeeName(bot, event, State)
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
			case "wait status saveEvidenceTime":
				handlesaveEvidenceImageTime(bot, event, usercardidState[event.Source.UserID], event.Source.UserID)
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
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
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
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
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
		sendReply(bot, event.ReplyToken, "คุณได้เช็คอินอยู่แล้ว กรุณาเช็คเอาท์ก่อนทำการเช็คอินใหม่")
		return
	}
	setUserState(State, "wait status ServiceRecordRequest")
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
		worktimeRecord := &models.WorktimeRecord{
			UserInfo: &models.User_info{
				Name: userInfo.Name,
			},
			CheckIn:  time.Now(),  // สามารถปรับข้อมูลจริงจากฐานข้อมูลได้
			CheckOut: time.Time{}, // ค่า CheckOut จะเป็นเวลาเริ่มต้น
		}

		// ใช้ Flex Message
		flexMessage := flexmessage.FormatConfirmationWorktime(worktimeRecord)
		if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
			log.Printf("Error sending Flex Message: %v", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
		}
	}
}

// ฟังก์ชันตรวจสอบสถานะการลงทะเบียน
func isUserRegistered(userID string) bool {
	// ตรวจสอบจากสถานะใน userState หรือฐานข้อมูล
	state, exists := userState[userID]
	return exists && state == "registered"
}

// เลือก เช็คอิน
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

	// ดึงข้อมูลผู้ใช้ตาม LINE ID
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบการเช็คอินของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
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
}

// เลือก เช็คเอ้าท์
func handleworktimeConfirmCheckOut(bot *linebot.Client, event *linebot.Event, userID string) {
	log.Println("Starting handleworktimeConfirmCheckOut for user:", userID)

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

	// ดึงข้อมูลผู้ใช้ตาม LINE ID
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}
	log.Printf("Fetched user info: %+v", userInfo)

	// ตรวจสอบการเช็คอินของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
	log.Printf("Check-in status for user %s: %v", userID, checkedIn)
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "คุณยังไม่ได้เช็คอิน กรุณาเช็คอินก่อนทำการเช็คเอ้าท์.")
		return
	}

	// บันทึกการเช็คเอ้าท์
	err = RecordCheckOut(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error recording check-out:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-out กรุณาลองใหม่.")
		return
	}
	log.Println("Check-out recorded successfully for user:", userID)

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
}

// func resetUserState(userID string) {
// 	delete(userState, userID)
// 	delete(userActivity, userID)
// 	delete(usercardidState, userID)
// 	log.Printf("Reset state for user %s", userID)
// }

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

	// หข้อความที่รับ = "ค้นหาข้อมูล"
	if message == "ค้นหาข้อมูล" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
		return
	}

	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
	cardID := message
	if len(cardID) != 13 {
		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักที่ถูกต้อง\nตัวอย่างเช่น 1234567891234 :")
		return
	}
	log.Println("เลขประจำตัวประชาชน:", cardID)

	// ดึงข้อมูลผู้ป่วยจาก CardID
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Println("Error fetching patient info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาตรวจสอบเลขประจำตัวประชาชนอีกครั้ง")
		return
	}

	//**************************************************************ดึงรูปมาแสดง***************************************
	// // ดึงข้อมูลรูปภาพจากฐานข้อมูล
	// imageData, err := GetImageFromDatabase(db, cardID)
	// if err != nil {
	// 	log.Println("Error fetching image from database:", err)
	// 	sendReply(bot, event.ReplyToken, "ไม่พบรูปภาพสำหรับผู้ป่วย กรุณาลองใหม่.")
	// 	return
	// }
	// log.Printf("inmageData: %+v", imageData)

	// // บันทึกภาพเป็นไฟล์ชั่วคราว
	// tempDir := os.TempDir() // ดึงตำแหน่ง temporary directory
	// if _, err := os.Stat(tempDir); os.IsNotExist(err) {
	// 	err = os.MkdirAll(tempDir, os.ModePerm)
	// 	if err != nil {
	// 		log.Println("Error creating temp directory:", err)
	// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการสร้างไดเรกทอรีชั่วคราว")
	// 		return
	// 	}
	// }
	// tempFilePath := fmt.Sprintf("%s\\%s.jpg", tempDir, cardID)
	// err = os.WriteFile(tempFilePath, imageData, 0644)
	// if err != nil {
	// 	log.Println("Error writing image file:", err)
	// 	sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการจัดการรูปภาพ")
	// 	return
	// }
	// defer os.Remove(tempFilePath) // ลบไฟล์หลังใช้งาน

	// // เชื่อมต่อ MinIO และอัปโหลดรูปภาพ
	// minioClient, err := database.ConnectToMinio()
	// if err != nil {
	// 	log.Println("Error connecting to MinIO:", err)
	// 	sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อ MinIO ได้")
	// 	return
	// }
	// objectName := fmt.Sprintf("patient_info/%d/%s.jpg", patient.PatientInfo.PatientInfo_ID, cardID)
	// bucketName := "nirunimages" // แทนที่ด้วยชื่อ bucket ของคุณ

	// fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
	// if err != nil {
	// 	log.Println("Error uploading file to MinIO:", err)
	// 	sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการอัปโหลดรูปภาพไปยัง MinIO")
	// 	return
	// }
	// log.Printf("fileURL: %s", fileURL)

	flexMessage := flexmessage.FormatPatientInfo(patient)
	if _, err := bot.PushMessage(userID, flexMessage).Do(); err != nil {
		log.Println("Error sending push message:", err)
	}

	log.Println("ข้อมูลผู้ป่วยและรูปภาพส่งสำเร็จ:", flexMessage)
	userState[userID] = ""
}

// บันทึกกิจกรรม
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

	// ดึงข้อมูลผู้ใช้ตาม LINE ID
	userInfo, err := GetUserInfoByLINEID(db, event.Source.UserID)
	if err != nil {
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบการเช็คอินของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณา Check-in ก่อน\nที่เมนู 'ลงเวลางาน'")
		return
	}

	// ข้อความที่รับ = "บันทึกกิจกรรม"
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

	// ตรวจสอบข้อมูลผู้ป่วยจากฐานข้อมูล
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

	flexMessage := flexmessage.FormatServiceInfo([]models.Activityrecord{*service})
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Printf("Error sending Flex Message (handleServiceInfo): %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
		return
	}

	// บันทึก cardID สำหรับใช้ในฟังก์ชันถัดไป
	usercardidState[State] = cardID

	// เปลี่ยนสถานะผู้ใช้หลังจากได้รับข้อมูล
	userState[State] = "wait status Activityrecord"
	log.Printf("Set user state to wait status Activityrecord for user %s", State)
}

// จัดการการเลือกกิจกรรม เมื่อเลือกกิจกรรมใหม่แล้ว
func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
	log.Println("wait status Activityrecord:", userState)

	if userState[State] != "wait status Activityrecord" {
		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
		return
	}
	//ตรวจสอบว่าเป็นtext
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		log.Println("Event is not a text message")
		return
	}
	//กิจกรรมที่รับมา
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

	// ใช้ Flex Message แทน Quick Reply
	flexContainer := flexmessage.FormatStartActivity(activity)
	flexMessage := linebot.NewFlexMessage("เริ่มกิจกรรม", flexContainer)

	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Printf("Error sending Flex Message: %v", err)
	}

	userState[State] = "wait status ActivityStart"
	log.Println("wait status ActivityStart:", userState)
}

// เริ่มกิจกรรมที่เลือก เมื่อกดเรื่มกิจกรรม
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
	//ข้อความที่รับ 'เริ่มกิจกรรม'
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

	//ดึงข้อมูลผู้ใช้จากLINE
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบการเช็คอินของพนักงาน
	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
	if err != nil {
		log.Println("Error checking check-in status:", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
		return
	}
	if !checkedIn {
		sendReply(bot, event.ReplyToken, "กรุณาเช็คอินก่อน\nที่เมนู 'ลงเวลางาน'")
		return
	}

	// ตรวจสอบข้อมูลผู้ป่วยจากฐานข้อมูล
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}
	//เตรียมข้อมูล activityRecord
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

	// activityRecord.PatientInfo.Name = patient.PatientInfo.Name
	flexMessage := flexmessage.FormatactivityRecordStarttime(activityRecord)
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Printf("Error sending Flex Message: %v", err)
	}

	userState[State] = "wait status ActivityEnd"
}

// จบกิจกรรมและรอรับหลักฐานภาพ เมื่อกด เสร็จสิ้น
func handleActivityEnd(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status ActivityEnd" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะของคุณไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	// ตรวจสอบประเภทข้อความ
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		endtime := strings.TrimSpace(message.Text)
		if endtime != "เสร็จสิ้น" {
			sendReply(bot, event.ReplyToken, "กรุณาพิมพ์ 'เสร็จสิ้น' เพื่อบันทึกเวลาสิ้นสุด.")
			return
		}

		// เปลี่ยนสถานะเพื่อใช้ในฟังก์ชันอื่น
		userState[userID] = "wait status Saveavtivityend"
		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปการทำกิจกรรมเพื่อบันทึกการทำกิจกรรม.")

	default:
		log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
		sendReply(bot, event.ReplyToken, "กรุณาส่งข้อความ 'เสร็จสิ้น' ในขั้นตอนนี้.")
	}
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
		userState[userID] = "wait status saveEvidenceImageActivity"
		log.Printf("Switching user state to 'wait status saveEvidenceImageActivity' for user %s", userID)

		//เข้าในฟังก์ชัน saveEvidenceImageActivity เพื่อบันทึกรูปการทำกิจกรรม
		if err := handlesaveEvidenceImageActivity(bot, event, cardID, userID, imageMessage); err != nil {
			log.Printf("Error saving image: %v", err)
			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกรูปภาพ กรุณาลองใหม่.")
			return
		}

		// หลังจากบันทึกรูปภาพเสร็จแล้ว เปลี่ยนสถานะเป็น wait status saveEvidenceTime เพื่อเข้าในฟังก์ชัน saveEvidenceImageTime
		userState[userID] = "wait status saveEvidenceTime"
		sendReply(bot, event.ReplyToken, "รูปการทำกิจกรรมถูกบันทึกแล้ว\nกรุณาส่งรูปการจับเวลาการทำกิจกรรม.")
	} else {
		log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
	}
}

// ฟังก์ชันบันทึกรูปภาพกิจกรรม
func handlesaveEvidenceImageActivity(bot *linebot.Client, event *linebot.Event, cardID, userID string, imageMessage *linebot.ImageMessage) error {
	if userState[userID] != "wait status saveEvidenceImageActivity" {
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

	// เขียนเนื้อหาภาพลงในไฟล์
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
	objectName := fmt.Sprintf("evidence_activity/%s/%d/%s.jpg", activity, patientInfoID, messageID)

	// อัปโหลดไฟล์ไปยัง MinIO
	fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
	if err != nil {
		log.Printf("Error uploading file to MinIO: %v", err)
		return err
	}

	// อัปเดต URL ในฐานข้อมูล
	err = updateEvidenceImageActivity(db, patientInfoID, fileURL)
	if err != nil {
		log.Printf("Error updating database: %v", err)
		return err
	}

	log.Printf("Activity Image successfully saved and URL updated: %s", fileURL)
	log.Printf("Last userID: %s", userID)
	return nil
}

// ฟังก์ชันบันทึกรูปภาพเวลาการทำกิจกรรม
func handlesaveEvidenceImageTime(bot *linebot.Client, event *linebot.Event, cardID, userID string) error {
	if userState[userID] != "wait status saveEvidenceTime" {
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
		objectName := fmt.Sprintf("evidence_time/%s/%d/%s.jpg", activity, patientInfoID, messageID)

		// อัปโหลดไฟล์ไปยัง MinIO
		fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
		if err != nil {
			log.Printf("Error uploading file to MinIO: %v", err)
			return err
		}

		// อัปเดต URL ในฐานข้อมูล
		err = updateEvidenceImageTime(db, patientInfoID, fileURL)
		if err != nil {
			log.Printf("Error updating database: %v", err)
			return err
		}

		log.Printf("Evidence time image successfully saved and URL updated: %s", fileURL)

		userState[userID] = "wait status SaveEmployeeName"
		log.Printf("User state updated to: %s", userState[userID])

		sendReply(bot, event.ReplyToken, "บันทึกรูปภาพเวลาสำเร็จ\nกรุณากรอกชื่อพนักงาน.")
		return nil
	}

	log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
	sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
	return nil
}

// บันทึกชื่อพนักงานที่ทำการบริการ
func handleSaveEmployeeName(bot *linebot.Client, event *linebot.Event, userID string) {
	if userState[userID] != "wait status SaveEmployeeName" {
		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
		sendReply(bot, event.ReplyToken, "สถานะของคุณไม่ถูกต้อง กรุณาลองใหม่.")
		return
	}

	// ชื่อพนักงานที่รับจากผู้ใช้
	employeeName := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
	log.Printf("Received employee name: %s", employeeName)
	if employeeName == "" {
		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ.")
		return
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
		return
	}
	defer db.Close()

	// ตรวจสอบข้อมูลEmployee
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

	// ดึง Activity Record ID
	activityRecordID, err := GetActivityRecordID(db, cardID)
	if err != nil {
		sendReply(bot, event.ReplyToken, err.Error())
		return
	}
	//ดึงข้อมูลผู้ใช้ตาม LINE ID
	userInfo, err := GetUserInfoByLINEID(db, userID)
	if err != nil {
		log.Println("Error fetching user info:", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
		return
	}
	patient, err := GetPatientInfoByName(db, cardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
		return
	}

	// เตรียมข้อมูล activityRecord
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
	log.Println("Activity Record to be updated:", activityRecord)

	// คำนวณระยะเวลา
	startTime, err := GetActivityStartTime(db, cardID, userActivity[userID])
	if err != nil {
		log.Printf("Error fetching StartTime: %v", err)
		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลเวลาเริ่ม กรุณาลองใหม่")
		return
	}
	duration := activityRecord.EndTime.Sub(startTime)
	activityRecord.Period = formatDuration(duration)

	// อัปเดตข้อมูลในฐานข้อมูล
	if err := UpdateActivityEndTime(db, activityRecord); err != nil {
		log.Printf("Error updating end time: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
		return
	}

	flexMessage := flexmessage.FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
		log.Printf("Error sending reply message: %v", err)
		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
		return
	}
	log.Printf("บันทึกกิจกรรมสำเร็จ: %s", flexMessage)
	// resetUserState(userID)
	userState[userID] = ""

}

// แปลงระยะเวลาของกิจกรรมเป็นชั่วโมงและนาที
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
}

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
