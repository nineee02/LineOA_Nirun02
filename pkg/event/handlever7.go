package event

// import (
// 	"database/sql"
// 	"fmt"
// 	"io"
// 	"log"
// 	"nirun/pkg/database"
// 	"nirun/pkg/flexmessage"
// 	"nirun/pkg/models"
// 	"nirun/service"
// 	"regexp"
// 	"strconv"
// 	"unicode"

// 	// "nirun/service"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/line/line-bot-sdk-go/linebot"
// )

// var usercardidState = make(map[string]string)
// var userState = make(map[string]string)                 //เก็บstate
// var userActivity = make(map[string]string)              // เก็บกิจกรรมสำหรับผู้ใช้แต่ละคน
// var userCheckInStatus = make(map[string]bool)           // เก็บสถานะการเช็คอินของแต่ละบัญชี LINE
// var userLastWorktimeAction = make(map[string]time.Time) // เก็บ timestamp การเช็คอิน/เช็คเอ้าท์ล่าสุด
// var userActivityInfoID = make(map[string]int)           // เก็บ activity_info_id ตาม userID
// var userActivityRecordID = make(map[string]int)         // เก็บ activityRecord_ID ตาม State ของผู้ใช้
// var userActivityCategory = make(map[string]string)      // เก็บมิติของกิจกรรมที่เลือก
// var userActivityStartDate = make(map[string]time.Time)  // เก็บวันที่เริ่มกิจกรรม
// var userActivityEndDate = make(map[string]time.Time)    // เก็บวันที่สิ้นสุดกิจกรรม
// var employeeLoginStatus = make(map[string]string)       // เก็บสถานะล็อกอิน {employeeID: userID}
// var userImageTimestamps = make(map[string]time.Time)    // เก็บ timestamp ของรูปภาพ

// // HandleEvent - จัดการข้อความที่ได้รับจาก LINE
// func HandleEvent(bot *linebot.Client, event *linebot.Event) {
// 	// ตรวจสอบประเภทข้อความก่อน
// 	switch message := event.Message.(type) {
// 	case *linebot.TextMessage: // อ่านข้อความจาก TextMessage
// 		text := strings.TrimSpace(message.Text)
// 		log.Println("Received TextMessage:", text)
// 		State := event.Source.UserID
// 		log.Println("User state: ", State)

// 		// ตรวจสอบคำสั่งจากข้อความ
// 		switch text {
// 		case "ค้นหาข้อมูล":
// 			handleElderlyInfoStste(bot, event, State)
// 		case "ลงเวลางาน":
// 			handleWorktimeStste(bot, event, State)
// 		case "บันทึกการบริการ":
// 			handleServiceRecordStste(bot, event, State)
// 		default:
// 			handleDefault(bot, event)
// 		}

// 		// ตรวจสอบสถานะ
// 		state, exists := userState[State]
// 		if exists {
// 			switch state {
// 			case "wait status worktime":
// 				handleWorktime(bot, event, State)
// 			case "wait status worktimeConfirmCheckIn":
// 				log.Println("✅ State matched: wait status worktimeConfirmCheckIn")

// 				// เชื่อมต่อฐานข้อมูล
// 				db, err := database.ConnectToDB()
// 				if err != nil {
// 					log.Println("❌ Database connection error:", err)
// 					sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 					return
// 				}
// 				defer db.Close()

// 				// ดึงข้อมูลผู้ใช้
// 				userInfo, err := GetUserInfoByLINEID(db, event.Source.UserID)
// 				if err != nil {
// 					log.Println("❌ Error fetching user info:", err)
// 					sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
// 					return
// 				}
// 				log.Printf("📌 Fetched user info: %+v", userInfo)

// 				// ตรวจสอบว่ายังเช็คอินอยู่หรือไม่
// 				checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
// 				if err != nil {
// 					log.Println("❌ Error checking user status:", err)
// 					sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
// 					return
// 				}
// 				log.Printf("📌 checkedIn status: %v", checkedIn)

// 				// แสดงปุ่มเช็คอิน/เช็คเอ้าท์
// 				UpdateWorktimeUI(bot, event, userInfo, checkedIn)
// 				return
// 			case "wait status ElderlyInfoRequest":
// 				handlePateintInfo(bot, event, State)
// 			case "wait status handleServiceGetCardID":
// 				handleServiceGetCardID(bot, event, State)
// 			case "wait status ServiceSelection":
// 				handleServiceSelection(bot, event, State)
// 			// case "wait status ServiceRecordRequest":
// 			// 	handleServiceInfo(bot, event, State)
// 			case "wait status ActivitySelection":
// 				handleActivitySelection(bot, event, State)
// 			case "wait status CustomActivity":
// 				handleCustomActivity(bot, event, State)
// 			case "wait status Activityrecord":
// 				handleActivityrecord(bot, event, State)
// 			// case "wait status ActivityStartDate":
// 			// 	handleActivityStartDate(bot, event, State)
// 			// case "wait status ActivityStartTime":
// 			// 	handleActivityStartTime(bot, event, State)
// 			// case "wait status ActivityEndDate":
// 			// 	handleActivityEndDate(bot, event, State)
// 			// case "wait status ActivityEndTime":
// 			// 	handleActivityEndTime(bot, event, State)
// 			case "wait status ActivityStart":
// 				handleActivityStart(bot, event, State)
// 			case "wait status ActivityEnd":
// 				handleActivityEnd(bot, event, State)

// 			// case "wait status ConfirmOrSaveEmployee":
// 			// 	handleUserChoiceForActivityRecord(bot, event, State, "")
// 			case "wait status ConfirmOrSaveEmployee":
// 				if textMessage, ok := event.Message.(*linebot.TextMessage); ok {
// 					selection := strings.TrimSpace(textMessage.Text) // รับข้อความที่ผู้ใช้ส่งมา
// 					log.Printf("Handling selection: %s", selection)  // เพิ่ม log ตรวจสอบค่า
// 					handleUserChoiceForActivityRecord(bot, event, State, selection)
// 				} else {
// 					log.Printf("Unexpected message type in ConfirmOrSaveEmployee state")
// 				}
// 			case "wait status SaveEmployeeName":
// 				handleSaveEmployeeName(bot, event, State, State, "")
// 			default:
// 				log.Printf("Unhandled state for user %s: %s", State, state)
// 			}
// 		}

// 	case *linebot.ImageMessage: // อ่านข้อความจาก ImageMessage
// 		State := event.Source.UserID
// 		state, exists := userState[State]
// 		if exists {
// 			switch state {
// 			case "wait status Saveavtivityend":
// 				log.Printf("Received ImageMessage: ID=%s", message.ID)
// 				handleSaveavtivityend(bot, event, usercardidState[event.Source.UserID], event.Source.UserID)
// 			case "wait status saveEvidenceImageafterActivity":
// 				handlesaveEvidenceImageafterActivity(bot, event, usercardidState[event.Source.UserID], event.Source.UserID)
// 			}
// 		}

// 	default:
// 		log.Printf("Unhandled message type: %T", event.Message)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถประมวลผลข้อความประเภทนี้ได้.")
// 	}
// }

// // สำหรับตั้งค่าสถานะผู้ใช้
// func setUserState(userID, state string) {
// 	userState[userID] = state
// 	log.Printf("Set user state for user %s to %s", userID, state)
// }

// // สำหรับดุงสถานะผู้ใช้
// func getUserState(userID string) (string, bool) {
// 	state, exists := userState[userID]
// 	return state, exists
// }

// // ฟังก์ชันเช็คว่าเวลากด "ลงเวลางาน" ซ้ำภายใน 40 นาทีหรือไม่
// func isRecentWorktimeAction(userID string) bool {
// 	lastActionTime, exists := userLastWorktimeAction[userID]
// 	if !exists {
// 		return false
// 	}
// 	// ตรวจสอบว่าห่างกันไม่ถึง 40 นาที
// 	if time.Since(lastActionTime) < 40*time.Minute {
// 		log.Printf("⏳ User %s pressed worktime button within 40 minutes. Locking action.", userID)
// 		return true
// 	}
// 	return false
// }

// // ฟังก์ชันอัปเดต timestamp การเช็คอิน/เช็คเอ้าท์
// func updateWorktimeAction(userID string) {
// 	userLastWorktimeAction[userID] = time.Now()
// }

// // เริ่มกระบวนการตรวจสอบสถานะการลงเวลาทำงาน
// func handleWorktimeStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	userID := event.Source.UserID

// 	// ตรวจสอบว่าผู้ใช้กดลงเวลางานซ้ำภายใน 40 นาทีหรือไม่
// 	if isRecentWorktimeAction(userID) {
// 		// ส่งข้อความถามยืนยัน พร้อมปุ่ม Quick Reply
// 		quickReply := linebot.NewQuickReplyItems(
// 			linebot.NewQuickReplyButton("", linebot.NewMessageAction("✅ ยืนยัน", "ยืนยันลงเวลางาน")),
// 			linebot.NewQuickReplyButton("", linebot.NewMessageAction("❌ ยกเลิก", "ยกเลิก")),
// 		)

// 		replyMessage := linebot.NewTextMessage("คุณลงเวลางานเมื่อไม่นานมานี้ ต้องการดำเนินการต่อหรือไม่?").
// 			WithQuickReplies(quickReply)

// 		if _, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
// 			log.Println("❌ Error sending Quick Reply message:", err)
// 		}

// 		// ตั้งสถานะให้ผู้ใช้รอการยืนยัน
// 		setUserState(userID, "wait status worktime")
// 		return
// 	}

// 	// อนุญาตให้ดำเนินการลงเวลางาน
// 	updateWorktimeAction(userID) // บันทึก timestamp
// 	setUserState(State, "wait status worktime")
// }

// // เริ่มกระบวนการขอค้นหาข้อมูล
// func handleElderlyInfoStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	// userID := event.Source.UserID
// 	// ตรวจสอบสถานะการเช็คอินของบัญชี LINE
// 	// if isUserCheckedIn(userID) {
// 	// 	return
// 	// }

// 	// อนุญาตให้ดำเนินการ
// 	setUserState(State, "wait status ElderlyInfoRequest")
// }

// // เริ่มกระบวนการขอบันทึกกิจกรรม
// func handleServiceRecordStste(bot *linebot.Client, event *linebot.Event, State string) {
// 	// userID := event.Source.UserID

// 	// // ตรวจสอบสถานะการเช็คอินของบัญชี LINE
// 	// if isUserCheckedIn(userID) {
// 	// 	return
// 	// }

// 	// ขอให้ผู้ใช้กรอกเลขบัตรประชาชน
// 	// sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")

// 	// ตั้งค่าผู้ใช้ให้อยู่ในโหมดรอเลขบัตรประชาชน
// 	setUserState(State, "wait status handleServiceGetCardID")
// }

// // **********************************************************************************************************

// func getUserProfile(bot *linebot.Client, userID string) (*linebot.UserProfileResponse, error) {
// 	profile, err := bot.GetProfile(userID).Do()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return profile, nil
// }

// // ส่งข้อความตอบกลับแบบกำหนดเอง สามารถดึงข้อมูลผู้ใช้ได้
// func sendCustomReply(bot *linebot.Client, replyToken string, userID string, greetingMessage string, messages ...linebot.SendingMessage) {
// 	if len(messages) == 0 {
// 		return
// 	}

// 	// ใช้ข้อความทักทายที่กำหนดเอง หรือดึงจากโปรไฟล์
// 	if greetingMessage == "" {
// 		profile, err := getUserProfile(bot, userID)
// 		if err == nil {
// 			greetingMessage = fmt.Sprintf("ยินดีต้อนรับ %s! ", profile.DisplayName)
// 		} else {
// 			greetingMessage = "ยินดีต้อนรับ!"
// 		}
// 	}

// 	// แทรกข้อความทักทายไปในข้อความที่ส่ง
// 	messages = append([]linebot.SendingMessage{linebot.NewTextMessage(greetingMessage)}, messages...)

// 	if _, err := bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
// 		log.Printf("Error replying message sendCustomReply: %v", err)
// 	}
// }
// func sendQRCodeForLogin(bot *linebot.Client, replyToken string) {
// 	flexmessage.SendRegisterLink(bot, replyToken)
// }

// //*************************************************************************************************************

// // // ตรวจสอบการล็อกอิน
// // func isEmployeeLoggedIn(employeeID, userID string) bool {
// // 	currentUser, exists := employeeLoginStatus[employeeID]
// // 	return exists && currentUser != userID
// // }

// // // ล็อกสถานะการใช้งาน
// // func lockEmployeeLogin(employeeID, userID string) bool {
// // 	if isEmployeeLoggedIn(employeeID, userID) {
// // 		return false // มีผู้ใช้อื่นกำลังใช้งานอยู่
// // 	}
// // 	employeeLoginStatus[employeeID] = userID
// // 	return true
// // }

// // // ปลดล็อกสถานะการใช้งาน
// // func unlockEmployeeLogin(employeeID string) {
// // 	delete(employeeLoginStatus, employeeID)
// // }

// //*************************************************************************************

// // ตรวจสอบสถานะของบัญชี
// func isUserCheckedIn(userID string) bool {
// 	status, exists := userCheckInStatus[userID]
// 	return exists && status
// }

// // ฟังก์ชันตรวจสอบสถานะการลงทะเบียน
// func isUserRegistered(userID string) bool {
// 	// ตรวจสอบจากสถานะใน userState หรือฐานข้อมูล
// 	state, exists := userState[userID]
// 	return exists && state == "registered"
// }

// // ฟังก์ชันจัดการ "ลงเวลางาน"
// func handleWorktime(bot *linebot.Client, event *linebot.Event, userID string) {
// 	log.Println("Processing worktime for user:", userID)

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("❌ Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ✅ ดึงข้อมูลผู้ใช้ตาม LINE ID
// 	userInfo, err := GetUserInfoByLINEID(db, userID)
// 	if err != nil {
// 		log.Println("❌ Error fetching user info:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ที่เชื่อมโยงกับบัญชี LINE นี้.")
// 		return
// 	}

// 	// ✅ ตรวจสอบการเช็คอินของพนักงาน
// 	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("❌ Error checking user status:", err)
// 		return
// 	}

// 	// ✅ ตรวจสอบข้อความที่ส่งมา
// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

// 	switch message {
// 	case "ยืนยันลงเวลางาน":
// 		log.Println("ยืนยันลงเวลา")
// 		log.Println("🔍 User confirmed worktime action")
// 		updateWorktimeAction(userID)
// 		setUserState(userID, "wait status worktime")

// 		// ✅ ตรวจสอบสถานะแล้วแสดงปุ่ม "เช็คอิน" หรือ "เช็คเอ้าท์"
// 		UpdateWorktimeUI(bot, event, userInfo, checkedIn)

// 	case "ยกเลิก":
// 		log.Println("ยกเลิก")
// 		sendReply(bot, event.ReplyToken, "การลงเวลางานถูกยกเลิก")

// 	case "เช็คอิน":
// 		log.Println("เช็คอิน")
// 		processCheckIn(bot, event, db, userInfo, userID)

// 	case "เช็คเอ้าท์":
// 		log.Println("เช็คเอ้าท์")
// 		processCheckOut(bot, event, db, userInfo, userID)

// 	default:
// 		UpdateWorktimeUI(bot, event, userInfo, checkedIn)
// 	}
// }

// // ฟังก์ชันยืนยันการลงเวลางาน
// func confirmWorktimeAction(bot *linebot.Client, event *linebot.Event, userID string, db *sql.DB, userInfo *models.User_info, checkedIn bool, message string) {
// 	if message == "ยืนยันลงเวลางาน" {
// 		log.Println("✅ User confirmed worktime action")
// 		// อัปเดตสถานะการลงเวลางาน
// 		updateWorktimeAction(userID)
// 		setUserState(userID, "wait status worktime") // ตั้งสถานะเป็น worktime หลังจากยืนยัน

// 		// ถ้าผู้ใช้เช็คอินอยู่แล้วให้แสดงปุ่มเช็คเอ้าท์
// 		UpdateWorktimeUI(bot, event, userInfo, checkedIn)
// 	} else if message == "ยกเลิก" {
// 		sendReply(bot, event.ReplyToken, "การลงเวลางานถูกยกเลิก")
// 		setUserState(userID, "wait status worktime") // รีเซ็ตสถานะ
// 	}
// }

// // ✅ ฟังก์ชันบันทึก "เช็คอิน"
// func processCheckIn(bot *linebot.Client, event *linebot.Event, db *sql.DB, userInfo *models.User_info, userID string) {
// 	err := RecordCheckIn(db, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("❌ Error recording check-in:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-in กรุณาลองใหม่.")
// 		return
// 	}

// 	updateWorktimeAction(userID)
// 	log.Println("updateWorktimeAction ", userID)
// 	setUserState(userID, "wait status checkOut")
// 	log.Println("setUserState CheckOut", userID)

// 	// ✅ ส่ง Flex Message
// 	worktimeRecord := &models.WorktimeRecord{
// 		UserInfo: &models.User_info{Name: userInfo.Name},
// 		CheckIn:  time.Now(),
// 	}

// 	flexMessage := flexmessage.FormatworktimeCheckin(worktimeRecord)
// 	if _, err := bot.PushMessage(userID, flexMessage).Do(); err != nil {
// 		log.Println("❌ Error flexMessage FormatworktimeCheckin:", err)
// 	}
// 	log.Println("✅ แสดงปุ่ม เช็คอินปกติ", flexMessage)
// }

// // ✅ ฟังก์ชันบันทึก "เช็คเอ้าท์"
// func processCheckOut(bot *linebot.Client, event *linebot.Event, db *sql.DB, userInfo *models.User_info, userID string) {
// 	checkedIn, err := IsEmployeeCheckedIn(db, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("❌ Error checking user status:", err)
// 		// sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะ กรุณาลองใหม่.")
// 		return
// 	}

// 	if !checkedIn {
// 		UpdateWorktimeUI(bot, event, userInfo, false)
// 		return
// 	}

// 	// ✅ บันทึกข้อมูลเช็คเอ้าท์
// 	err = RecordCheckOut(db, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("❌ Error recording check-out:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก Check-out กรุณาลองใหม่.")
// 		return
// 	}

// 	updateWorktimeAction(userID)
// 	userState[userID] = "wait status checkIn"

// 	// ✅ ดึงข้อมูลบันทึกเวลาทำงาน
// 	worktimeRecord, err := GetWorktimeRecordByUserID(db, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("❌ Error fetching worktime record:", err)
// 		// sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูล กรุณาลองใหม่.")
// 		return
// 	}

// 	// ✅ สร้าง Flex Message
// 	worktimeRecord = &models.WorktimeRecord{
// 		UserInfo: &models.User_info{Name: userInfo.Name},
// 		CheckOut: time.Now(),
// 		Period:   worktimeRecord.Period,
// 	}

// 	flexMessage := flexmessage.FormatworktimeCheckout(worktimeRecord)
// 	if _, err := bot.PushMessage(userID, flexMessage).Do(); err != nil {
// 		log.Println("Error flexMessage FormatworktimeCheckout:", err)
// 	}
// 	log.Println("แสดงปุ่ม เช็คเอ้าท์ปกติ", flexMessage)
// }

// // ฟังก์ชันสำหรับอัปเดต UI ของปุ่มเช็คอิน / เช็คเอ้าท์
// func UpdateWorktimeUI(bot *linebot.Client, event *linebot.Event, userInfo *models.User_info, checkedIn bool) {
// 	worktimeRecord := &models.WorktimeRecord{
// 		UserInfo: &models.User_info{Name: userInfo.Name},
// 		CheckIn:  time.Now(),
// 		CheckOut: time.Time{},
// 	}

// 	var flexMessage *linebot.FlexMessage
// 	if checkedIn {
// 		flexMessage = flexmessage.FormatConfirmCheckout(worktimeRecord)
// 		setUserState(event.Source.UserID, "wait status worktime") // ตั้งสถานะเป็น "รอเช็คเอ้าท์"
// 	} else {
// 		flexMessage = flexmessage.FormatConfirmCheckin(worktimeRecord)
// 		setUserState(event.Source.UserID, "wait status worktime") // ตั้งสถานะเป็น "รอเช็คอิน"
// 	}

// 	// ส่ง Flex Message ให้ผู้ใช้เลือกเช็คอิน/เช็คเอ้าท์
// 	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
// 		log.Println("❌ Error sending Flex Message:", err)
// 	}
// }

// // func resetUserState(userID string) {
// // 	delete(userState, userID)
// // 	delete(userActivity, userID)
// // 	delete(usercardidState, userID)
// // 	log.Printf("Reset state for user %s", userID)
// // }

// func sanitizeCardID(s string) string {
// 	var builder strings.Builder
// 	for _, char := range s {
// 		if unicode.IsDigit(char) { // ตรวจสอบเฉพาะตัวเลข
// 			builder.WriteRune(char)
// 		}
// 	}
// 	return builder.String()
// }

// // การค้นหาข้อมูล
// func handlePateintInfo(bot *linebot.Client, event *linebot.Event, userID string) {
// 	state, exists := getUserState(userID)
// 	if !exists || state != "wait status ElderlyInfoRequest" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, state)
// 		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
// 		return
// 	}

// 	// ดึงข้อความที่ผู้ใช้ส่งมา
// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ข้อความที่รับ = "ค้นหาข้อมูล"
// 	if message == "ค้นหาข้อมูล" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
// 		return
// 	}

// 	// ตรวจสอบเลขประจำตัวประชาชน (cardID)
// 	cardID := sanitizeCardID(message)
// 	if len(cardID) != 13 {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักที่ถูกต้อง\nตัวอย่างเช่น 1234567891234 :")
// 		return
// 	}
// 	log.Println("เลขประจำตัวประชาชน:", cardID)

// 	patient, err := service.PostRequestPatientByID(cardID)
// 	if err != nil {
// 		log.Println("ErE:", err)
// 		return
// 	}
// 	log.Println("Papatient:", patient)
// 	// ดึงข้อมูลผู้ป่วยจาก CardID

// 	flexMessage := flexmessage.FormatPatientInfo(patient)
// 	if _, err := bot.PushMessage(userID, flexMessage).Do(); err != nil {
// 		log.Println("Error sending push message:", err)
// 	}

// 	log.Println("ข้อมูลผู้ป่วย:", flexMessage)
// 	// userState[userID] = ""
// }

// func isNumeric(s string) bool {
// 	for _, c := range s {
// 		if c < '0' || c > '9' {
// 			return false
// 		}
// 	}
// 	return true
// }

// // func parseDateInput(input string) (time.Time, error) {
// // 	// ลบช่องว่างส่วนเกิน และเปลี่ยนเป็น lower case
// // 	input = strings.TrimSpace(strings.ToLower(input))

// // 	// กำหนด regex สำหรับจับวันและเดือน/ปี
// // 	re := regexp.MustCompile(`^(\d{1,2})/(\d{1,2})/(\d{4})$`)
// // 	match := re.FindStringSubmatch(input)

// // 	if len(match) == 0 {
// // 		return time.Time{}, fmt.Errorf("รูปแบบวันที่ไม่ถูกต้อง กรุณากรอกเป็น วัน/เดือน/ปี เช่น 01/01/2567")
// // 	}

// // 	// ดึงค่าจาก regex match
// // 	day, month, yearStr := match[1], match[2], match[3]

// // 	// แปลงปีเป็น int
// // 	year, _ := strconv.Atoi(yearStr)

// // 	// ตรวจสอบว่าปีเป็น พ.ศ. หรือไม่
// // 	if year > 2500 {
// // 		year -= 543 // แปลง พ.ศ. → ค.ศ.
// // 	}

// // 	// สร้างวันที่โดยไม่มีเวลา
// // 	dateStr := fmt.Sprintf("%s/%s/%d", day, month, year)
// // 	layout := "02/01/2006"
// // 	parsedDate, err := time.Parse(layout, dateStr)
// // 	if err != nil {
// // 		return time.Time{}, fmt.Errorf("ไม่สามารถแปลงวันได้")
// // 	}

// // 	return parsedDate, nil
// // }

// // func parseTimeInput(input string) (time.Time, error) {
// // 	// ลบช่องว่างส่วนเกิน และเปลี่ยนเป็น lower case
// // 	input = strings.TrimSpace(strings.ToLower(input))

// // 	// กำหนด regex สำหรับจับเวลา
// // 	re := regexp.MustCompile(`^(\d{1,2})[:.](\d{2})\s*(น\.?|น)?$`)
// // 	match := re.FindStringSubmatch(input)

// // 	if len(match) == 0 {
// // 		return time.Time{}, fmt.Errorf("รูปแบบเวลาไม่ถูกต้อง กรุณากรอกเป็น ชั่วโมง:นาที เช่น 11:20 น.")
// // 	}

// // 	// ดึงค่าจาก regex match
// // 	hour, min := match[1], match[2]

// // 	// แปลงค่าเป็นตัวเลข
// // 	hourInt, _ := strconv.Atoi(hour)
// // 	minInt, _ := strconv.Atoi(min)

// // 	// สร้างเวลา
// // 	return time.Date(0, 0, 0, hourInt, minInt, 0, 0, time.UTC), nil
// // }

// func handleServiceGetCardID(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status handleServiceGetCardID" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}

// 	// db, err := database.ConnectToDB()
// 	// if err != nil {
// 	// 	log.Println("Database connection error:", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	userID := event.Source.UserID

// 	// ตรวจสอบสถานะการเช็คอินของพนักงาน
// 	userInfo, err := service.GetUserInfoByLINEID(userID)
// 	if err != nil {
// 		log.Println("ไม่พบข้อมูลพนักงานที่เชื่อมโยงกับ LINE ID นี้.")
// 		sendReply(bot, event.ReplyToken, "คุณยังไม่ได้ลงทะเบียนในระบบ กรุณาติดต่อผู้ดูแล.")
// 		return
// 	}

// 	isCheckedIn, err := service.IsEmployeeCheckedIn(userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Println("Error checking worktime status:", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการตรวจสอบสถานะการลงเวลาทำงาน กรุณาลองใหม่.")
// 		return
// 	}

// 	if !isCheckedIn {
// 		sendReply(bot, event.ReplyToken, "คุณยังไม่ได้ลงเวลาทำงาน กรุณาเช็คอินก่อนเริ่มบันทึกการบริการ.")
// 		return
// 	}

// 	//ดึงข้อความที่ผู้ใช้ส่งมา
// 	message := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)

// 	// กรณีผู้ใช้ยังไม่ได้กรอกเลขบัตรประชาชน
// 	if message == "บันทึกการบริการ" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักของผู้เข้ารับบริการ\nตัวอย่างเช่น 1234567891234 :")
// 		return
// 	}

// 	// ตรวจสอบเลขบัตรประชาชน
// 	cardID := sanitizeCardID(message)
// 	if len(cardID) != 13 {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกเลขบัตรประชาชน 13 หลักที่ถูกต้อง\nตัวอย่างเช่น 1234567891234 :")
// 		return
// 	}
// 	log.Println("เลขประจำตัวประชาชน:", cardID)

// 	// ตรวจสอบว่ามีข้อมูลผู้ป่วยหรือไม่
// 	patient, err := service.PostRequestPatientByID(cardID)
// 	if err != nil || patient == nil {
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้สูงอายุ กรุณากรอกเลขประจำตัวประชาชนอีกครั้ง")
// 		return
// 	}

// 	// บันทึก cardID สำหรับใช้ในฟังก์ชันถัดไป
// 	usercardidState[State] = cardID
// 	setUserState(State, "wait status ActivitySelection")

// 	// ใช้ `PushMessage()` แทน `ReplyMessage()` เพื่อหลีกเลี่ยงปัญหา reply token
// 	flexMessage := flexmessage.FormatActivityCategories()
// 	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
// 		log.Println("Error sending activity category selection:", err)
// 	}
// }

// func handleServiceSelection(bot *linebot.Client, event *linebot.Event, State string) {
// 	// if userState[State] != "wait status ServiceSelection" {
// 	// 	log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 	// 	return
// 	// }

// 	// // ดึงข้อความที่ผู้ใช้เลือก
// 	// message := event.Message.(*linebot.TextMessage).Text

// 	// switch message {
// 	// case "บันทึกกิจกรรม":
// 	// 	// อัปเดตสถานะก่อนเรียก `handleServiceInfo`
// 	// 	setUserState(State, "wait status ActivitySelection")
// 	// 	log.Printf("User %s state changed to: %s", State, "wait status ActivitySelection")
// 	// 	flexMessage := flexmessage.FormatActivityCategories()
// 	// 	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
// 	// 		log.Println("Error sending activity category selection:", err)
// 	// 	}
// 	// 	// เรียกใช้ฟังก์ชันบันทึกกิจกรรม
// 	// 	handleActivitySelection(bot, event, State)

// 	// case "รายงานปัญหา":
// 	// 	// ตั้งค่าสถานะเป็น "wait status ReportIssue"
// 	// 	setUserState(State, "wait status ReportIssue")
// 	// 	log.Printf("User %s state changed to: %s", State, "wait status ReportIssue")

// 	// 	// เรียกใช้ฟังก์ชันจัดการรายงานปัญหา
// 	// 	// handleReportIssue(bot, event, State)

// 	// default:
// 	// 	sendReply(bot, event.ReplyToken, "กรุณาเลือก 'บันทึกกิจกรรม' หรือ 'รายงานปัญหา'")
// 	// }
// }

// // เลือกมิติของกิจกรรม
// func handleActivitySelection(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status ActivitySelection" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}

// 	//รับข้อความกิจกรรมที่ผู้ใช้เลือก
// 	category := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Printf("User selected category: %s", category)

// 	//ตรวจสอบว่าผู้ใช้ต้องการย้อนกลับ
// 	if category == "ย้อนกลับ" {
// 		sendActivityCategorySelection(bot, event)            //Flex Message มิติ
// 		setUserState(State, "wait status ActivitySelection") // อัปเดตสถานะกลับไปเลือกมิติ
// 		return
// 	}
// 	//ตรวจสอบมิติของกิจกรรม`category_id`
// 	categoryMapping := map[string]int{
// 		"มิติสุขภาพ":      4,
// 		"มิติสังคม":       5,
// 		"มิติเศรษฐกิจ":    6,
// 		"มิติสิ่งแวดล้อม": 7,
// 		"มิติเทคโนโลยี":   8,
// 		"มิติอื่นๆ":       9,
// 	}

// 	categoryID, exists := categoryMapping[category]
// 	if !exists {
// 		sendReply(bot, event.ReplyToken, "กรุณาเลือกมิติของกิจกรรมที่ถูกต้องจากเมนู")
// 		return
// 	}

// 	// อัปเดตหมวดหมู่กิจกรรมใน State
// 	userActivityCategory[State] = category
// 	log.Printf("Updated user activity category: %s", category)
// 	//ดึงกิจกรรมที่เกี่ยวข้องจาก API JSON-RPC และแสดงให้ผู้ใช้เลือก
// 	fetchAndShowActivities(bot, event, State,categoryID)

// 	if category == "9" {
// 		//ให้ผู้ใช้กรอกกิจกรรมเอง
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อกิจกรรมของคุณ:")
// 		userState[State] = "wait status CustomActivity"
// 	} else {
// 		//ดึงกิจกรรมที่เกี่ยวข้องจาก API JSON-RPC และแสดงให้ผู้ใช้เลือก
// 		fetchAndShowActivities(bot, event, State, categoryID)
// 	}
// }

// // ดึงกิจกรรมแต่ละมิติมาแสดงให้เลือก
// func fetchAndShowActivities(bot *linebot.Client, event *linebot.Event, State string, categoryID int) {
// 		// อัปเดตหมวดหมู่กิจกรรมใน state
// 		userActivityCategory[State] = fmt.Sprintf("%d", categoryID)
	
// 		// ใช้ API JSON-RPC เพื่อดึงรายการกิจกรรม
// 		activityList, err := service.PostActivitiesByCategory(categoryID)
// 		if err != nil {
// 			log.Printf("❌ Error fetching activities from API: %v", err)
// 			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรม กรุณาลองใหม่.")
// 			return
// 		}
	
// 		// ตรวจสอบว่ามีกิจกรรมที่ดึงมาได้หรือไม่
// 		if len(activityList) == 0 {
// 			sendReply(bot, event.ReplyToken, "❌ ไม่พบกิจกรรมในหมวดหมู่นี้ กรุณาเลือกหมวดหมู่อื่น.")
// 			return
// 		}
	
// 		// เลือกใช้ Flex Message ตาม categoryID
// 		var flexMessage *linebot.FlexMessage
// 		switch categoryID {
// 		case 4: // มิติสุขภาพ
// 			flexMessages := flexmessage.FormatActivitieshealthCarousel(activityList)
// 			// ส่ง Carousel แบบแยกเป็น 2 ชุด ถ้ามีมากกว่า 9 รายการ
// 			for _, msg := range flexMessages {
// 				if _, err := bot.PushMessage(event.Source.UserID, msg).Do(); err != nil {
// 					log.Println("❌ Error sending activity list:", err)
// 				}
// 			}
// 			return
// 		case 5: // มิติสังคม
// 			flexMessage = flexmessage.FormatActivitiessocialCarousel(activityList)
// 		case 6: // มิติเศรษฐกิจ
// 			flexMessage = flexmessage.FormatActivitieseconomicCarousel(activityList)
// 		case 7: // มิติสิ่งแวดล้อม
// 			flexMessage = flexmessage.FormatActivitiesenvironmentalCarousel(activityList)
// 		case 8: // มิติเทคโนโลยี
// 			flexMessage = flexmessage.FormatActivitiestechnologyCarousel(activityList)
// 		default:
// 			log.Printf("❌ Invalid category selection: %d", categoryID)
// 			sendReply(bot, event.ReplyToken, "❌ หมวดหมู่กิจกรรมไม่ถูกต้อง กรุณาลองใหม่.")
// 			return
// 		}
	
// 		// 🔹 ส่ง Flex Message ไปยัง LINE Bot
// 		if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
// 			log.Println("❌ Error sending activity list:", err)
// 			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการแสดงกิจกรรม กรุณาลองใหม่.")
// 			return
// 		}
	
// 		// 🔹 อัปเดตสถานะเป็น "รอเลือกกิจกรรม"
// 		setUserState(State, "wait status Activityrecord")
// 	}

// // // เลือกมิติอีกครั้งหากกดย้อนกลับ
// func sendActivityCategorySelection(bot *linebot.Client, event *linebot.Event) {
// 	// log.Println("Sending activity category selection...")

// 	//รีเซ็ตมิติที่เลือกก่อนหน้า
// 	delete(userActivityCategory, event.Source.UserID)

// 	//ส่ง Flex Message เพื่อเลือกมิติใหม่
// 	flexMessage := flexmessage.FormatActivityCategories()
// 	if _, err := bot.PushMessage(event.Source.UserID, flexMessage).Do(); err != nil {
// 		log.Println("Error sending activity category selection:", err)
// 	}
// }

// // รับกิจกรรมจากผู้ใช้หากเลือก "มิติอื่นๆ"
// func handleCustomActivity(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status CustomActivity" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		return
// 	}

// 	// รับข้อความกิจกรรมที่ผู้ใช้ป้อนเอง
// 	activity := strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Printf("User entered custom activity: %s", activity)

// 	if activity == "" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อกิจกรรมของคุณ")
// 		return
// 	}

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	//เก็บค่ากิจกรรมที่ผู้ใช้กรอกเอง
// 	userActivity[State] = activity

// 	// // ส่ง Flex Message เพื่อยืนยันเริ่มกิจกรรม
// 	// flexContainer := flexmessage.FormatStartActivity(activity)
// 	// flexMessage := linebot.NewFlexMessage("เริ่มกิจกรรม", flexContainer)
// 	// if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
// 	// 	log.Printf("Error sending Flex Message: %v", err)
// 	// }

// 	//เปลี่ยนสถานะเป็น "รอเริ่มกิจกรรม"
// 	userState[State] = "wait status ActivityStart"
// 	sendReply(bot, event.ReplyToken, "กรุณากรอกวัน/เดือน/ปี เริ่มต้นของกิจกรรม\nเช่น 01/01/2567")
// }

// // บันทึกกิจกรรม เมื่อเลือกกิจกรรมใหม่แล้ว
// func handleActivityrecord(bot *linebot.Client, event *linebot.Event, State string) {
// 	log.Println("wait status Activityrecord:", userState)

// 	if userState[State] != "wait status Activityrecord" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมใหม่:")
// 		return
// 	}

// 	// ตรวจสอบว่าเป็นข้อความ
// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	// กิจกรรมที่รับมา
// 	activity := strings.TrimSpace(strings.ToLower(message.Text))
// 	log.Printf("Received activity input: %s", activity)

// 	//ถ้าผู้ใช้กด "ย้อนกลับ" ให้กลับไปเลือกมิติใหม่
// 	if activity == "ย้อนกลับ" {
// 		log.Println("User chose to go back. Resetting state...")

// 		//รีเซ็ตค่าหมวดหมู่ที่เลือกไว้
// 		delete(userActivityCategory, State)

// 		//ตั้งค่า State ใหม่ให้กลับไปเลือกมิติ
// 		setUserState(State, "wait status ActivitySelection")

// 		//ส่ง Flex Message เพื่อเลือกมิติใหม่
// 		sendActivityCategorySelection(bot, event)
// 		return
// 	}

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Println("Database connection error:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ตรวจสอบว่าผู้ใช้เคยเลือกมิติของกิจกรรมหรือไม่
// 	category, exists := userActivityCategory[State]
// 	if !exists {
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด: ไม่พบมิติของกิจกรรม กรุณาลองใหม่.")
// 		return
// 	}
// 	log.Println("category:%s", category)

// 	// ดึงกิจกรรมจากฐานข้อมูลตามมิติที่เลือก
// 	var validActivities []string
// 	switch category {
// 	case "technology":
// 		activityList, err := GetTechnologyActivities(db)
// 		if err == nil {
// 			for _, act := range activityList {
// 				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityTechnology)))
// 			}
// 		}
// 	case "social":
// 		activityList, err := GetSocialActivities(db)
// 		if err == nil {
// 			for _, act := range activityList {
// 				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivitySocial)))
// 			}
// 		}
// 	case "health":
// 		activityList, err := GetHealthActivities(db)
// 		if err == nil {
// 			for _, act := range activityList {
// 				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityHealth)))
// 			}
// 		}
// 	case "economic":
// 		activityList, err := GetEconomicActivities(db)
// 		if err == nil {
// 			for _, act := range activityList {
// 				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityEconomic)))
// 			}
// 		}
// 	case "environmental":
// 		activityList, err := GetEnvironmentalActivities(db)
// 		if err == nil {
// 			for _, act := range activityList {
// 				validActivities = append(validActivities, strings.TrimSpace(strings.ToLower(act.ActivityEnvironmental)))
// 			}
// 		}
// 	default:
// 		log.Println("Invalid category:", category)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาด กรุณาลองใหม่.")
// 		return
// 	}
// 	// log.Println("validActivities:%s", validActivities)

// 	// ตรวจสอบว่ากิจกรรมที่ผู้ใช้เลือกอยู่ในฐานข้อมูลหรือไม่
// 	isValid := false
// 	for _, validActivity := range validActivities {
// 		if activity == validActivity {
// 			isValid = true
// 			break
// 		}
// 	}

// 	if !isValid {
// 		sendReply(bot, event.ReplyToken, fmt.Sprintf("กิจกรรม '%s' ไม่ถูกต้อง กรุณาเลือกจากรายการที่กำหนด", activity))
// 		return
// 	}

// 	// ดึง activity_info_id
// 	activityID, err := GetActivityInfoIDByType(db, category, activity)
// 	if err != nil {
// 		log.Println("Error fetching activity ID:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
// 		return
// 	}

// 	// บันทึกข้อมูล activityInfoID
// 	userActivityInfoID[State] = activityID
// 	log.Printf("Stored activityInfoID for user %s: %d", State, activityID)

// 	// บันทึกกิจกรรมที่ผู้ใช้เลือก
// 	userActivity[State] = activity
// 	log.Printf("Stored activity for user %s: %s", State, activity)

// 	// // ส่ง Flex Message เพื่อยืนยันเริ่มกิจกรรม
// 	// flexContainer := flexmessage.FormatStartActivity(activity)
// 	// flexMessage := linebot.NewFlexMessage("เริ่มกิจกรรม", flexContainer)
// 	// if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
// 	// 	log.Printf("Error sending Flex Message: %v", err)
// 	// }

// 	// อัปเดตสถานะเป็น "wait status ActivityStart"
// 	userState[State] = "wait status ActivityStart"
// 	log.Println("wait status ActivityStart: ", userState)
// 	sendReply(bot, event.ReplyToken, "กรุณากรอกวัน/เดือน/ปี เริ่มต้นของกิจกรรม\nเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น.")
// }

// func parseTimeInput(input string) (time.Time, error) {
// 	// 🔹 ลบช่องว่างส่วนเกิน และเปลี่ยนเป็น lower case
// 	input = strings.TrimSpace(strings.ToLower(input))

// 	// 🔹 กำหนด regex ให้รองรับทุกกรณี
// 	re := regexp.MustCompile(`^(\d{1,2})/(\d{1,2})/(\d{4})\s*(เวลา)?\s*(\d{1,2})[:.](\d{2})\s*(น\.?|น)?$`)
// 	match := re.FindStringSubmatch(input)

// 	if len(match) == 0 {
// 		return time.Time{}, fmt.Errorf("รูปแบบวันและเวลาไม่ถูกต้อง กรุณากรอกเป็น วัน/เดือน/ปี ชั่วโมง:นาที เช่น 01/01/2567 11:20 น.")
// 	}

// 	// 🔹 ดึงค่าจาก regex match
// 	day, month, yearStr := match[1], match[2], match[3]
// 	hour, min := match[5], match[6]

// 	// 🔹 แปลงปีเป็น int
// 	year, _ := strconv.Atoi(yearStr)

// 	// 🔹 ตรวจสอบว่าปีเป็น พ.ศ. หรือไม่
// 	if year > 2500 {
// 		year -= 543 // แปลง พ.ศ. → ค.ศ.
// 	}

// 	// 🔹 ใช้ time.Parse() ตรวจสอบความถูกต้อง
// 	dateTimeStr := fmt.Sprintf("%s/%s/%d %s:%s", day, month, year, hour, min)
// 	layout := "02/01/2006 15:04"
// 	parsedTime, err := time.Parse(layout, dateTimeStr)
// 	if err != nil {
// 		return time.Time{}, fmt.Errorf("ไม่สามารถแปลงวันและเวลาได้")
// 	}

// 	return parsedTime, nil
// }

// // กดเรื่มกิจกรรม
// func handleActivityStart(bot *linebot.Client, event *linebot.Event, State string) {
// 	log.Println("wait status ActivityStart:", userState)

// 	if userState[State] != "wait status ActivityStart" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกวันที่และเวลาเริ่มต้นของกิจกรรม\nเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น. ")
// 		return
// 	}

// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	startTimeStr := strings.TrimSpace(message.Text)
// 	log.Printf("Received start time input: %s", startTimeStr)

// 	startTime, err := parseTimeInput(startTimeStr)
// 	if err != nil {
// 		sendReply(bot, event.ReplyToken, "รูปแบบเวลาไม่ถูกต้อง ตัวอย่างเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น.")
// 		return
// 	}

// 	now := time.Now()
// 	startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 		return
// 	}
// 	defer db.Close()

// 	activityRecordID, exists := userActivityRecordID[State]
// 	if exists {
// 		log.Printf("Updating start_time for activity_record_id=%d, start_time=%v", activityRecordID, startTime)
// 		err := UpdateActivityStartTime(db, activityRecordID, startTime)
// 		if err != nil {
// 			log.Printf("Error updating start_time: %v", err)
// 			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาเริ่ม กรุณาลองใหม่")
// 			return
// 		}
// 	} else {
// 		// ดึงข้อมูลที่จำเป็น
// 		cardID := usercardidState[State]
// 		patient, err := GetPatientInfoByName(db, cardID)
// 		if err != nil {
// 			log.Printf("❌ Error fetching patient_info_id: %v", err)
// 			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
// 			return
// 		}
// 		patientInfoID := patient.PatientInfo.PatientInfo_ID
// 		if patientInfoID == 0 {
// 			sendReply(bot, event.ReplyToken, "ไม่พบรหัสข้อมูลผู้ป่วย กรุณาลองใหม่")
// 			return
// 		}

// 		category, exists := userActivityCategory[State]
// 		if !exists {
// 			sendReply(bot, event.ReplyToken, "ไม่พบหมวดหมู่กิจกรรม กรุณาลองใหม่")
// 			return
// 		}

// 		userInfo, err := GetUserInfoByLINEID(db, State)
// 		if err != nil {
// 			log.Printf("Error fetching user info: %v", err)
// 			sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่")
// 			return
// 		}
// 		userInfoID := userInfo.UserInfo_ID

// 		//ตรวจสอบว่าเป็น "มิติอื่นๆ" หรือไม่
// 		activityInfoID, exists := userActivityInfoID[State]
// 		activityOther := ""
// 		if category == "other" {
// 			activityOther = userActivity[State] // ใช้ชื่อกิจกรรมที่ผู้ใช้กรอก
// 		} else {
// 			if !exists || activityInfoID == 0 {
// 				sendReply(bot, event.ReplyToken, "กรุณาเลือกกิจกรรมก่อนเริ่ม")
// 				return
// 			}
// 		}

// 		log.Printf("Creating new activity record for patient_info_id=%d, activity_info_id=%d, activityOther=%s, start_time=%v", patientInfoID, activityInfoID, activityOther, startTime)
// 		newRecordID, err := InsertActivityStartTime(db, patientInfoID, category, activityInfoID, activityOther, startTime, userInfoID)
// 		if err != nil {
// 			log.Printf("Error inserting activity start time: %v", err)
// 			return
// 		}

// 		userActivityRecordID[State] = newRecordID
// 	}

// 	userState[State] = "wait status ActivityEnd"
// 	log.Printf("Updating userState for %s to wait status ActivityEnd", State)
// 	sendReply(bot, event.ReplyToken, "กรุณากรอกวันที่และเวลาสิ้นสุดของกิจกรรม\nเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น.")
// }

// func handleActivityEnd(bot *linebot.Client, event *linebot.Event, State string) {
// 	if userState[State] != "wait status ActivityEnd" {
// 		log.Printf("Invalid state for user %s. Current state: %s", State, userState[State])
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกวันที่และเวลาสิ้นสุดของกิจกรรม\nเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น.")
// 		return
// 	}

// 	message, ok := event.Message.(*linebot.TextMessage)
// 	if !ok {
// 		log.Println("Event is not a text message")
// 		return
// 	}

// 	// รับเวลาสิ้นสุดจากผู้ใช้
// 	endTimeStr := strings.TrimSpace(message.Text)
// 	log.Printf("Received end time input: %s", endTimeStr)

// 	// ใช้ฟังก์ชัน parseTimeInput() แปลงค่า
// 	endTime, err := parseTimeInput(endTimeStr)
// 	if err != nil {
// 		sendReply(bot, event.ReplyToken, "รูปแบบเวลาไม่ถูกต้อง ตัวอย่างเช่น 01/01/2567 11:20 น. หรือ\n01/01/2567 เวลา 11:20 น.")
// 		return
// 	}

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่")
// 		return
// 	}
// 	defer db.Close()

// 	// ดึง `activity_record_id`
// 	activityRecordID, exists := userActivityRecordID[State]
// 	if !exists {
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
// 		return
// 	}

// 	// ดึง `user_info_id` ของผู้ใช้
// 	userID, err := GetUserInfoByLINEID(db, event.Source.UserID)
// 	if err != nil {
// 		log.Printf("❌ ไม่พบ user_info_id สำหรับ LINE ID: %s", event.Source.UserID)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่")
// 		return
// 	}

// 	// บันทึก `end_time` ลงฐานข้อมูล
// 	err = UpdateActivityEndTime(db, activityRecordID, endTime, userID.UserInfo_ID)
// 	if err != nil {
// 		log.Printf("❌ Error updating end_time: %v", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่.")
// 		return
// 	}

// 	// ส่งข้อความยืนยัน และขอให้ส่งรูปหลักฐาน
// 	sendReply(bot, event.ReplyToken, fmt.Sprintf("กรุณาส่งรูปก่อนการทำกิจกรรม"))

// 	// เปลี่ยนสถานะเป็น "wait status Saveavtivityend"
// 	userState[State] = "wait status Saveavtivityend"
// }

// // รับและประมวลผลรูปการทำกิจกรรม
// func handleSaveavtivityend(bot *linebot.Client, event *linebot.Event, cardID, userID string) {
// 	if userState[userID] != "wait status Saveavtivityend" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		sendReply(bot, event.ReplyToken, "สถานะของคุณไม่ถูกต้อง กรุณาลองใหม่.")
// 		return
// 	}

// 	// ตรวจสอบประเภทข้อความว่าเป็นImage
// 	if imageMessage, ok := event.Message.(*linebot.ImageMessage); ok {
// 		log.Printf("Processing ImageMessage for user %s", userID)

// 		// อัปเดตสถานะ wait status saveEvidenceImageActivity เพื่อใช้เข้าในฟังก์ชัน saveEvidenceImageActivity
// 		userState[userID] = "wait status saveEvidenceImagebeforeActivity"
// 		log.Printf("Switching user state to 'wait status saveEvidenceImageActivity' for user %s", userID)

// 		//เข้าในฟังก์ชัน saveEvidenceImageActivity เพื่อบันทึกรูปการทำกิจกรรม
// 		if err := handlesaveEvidenceImagebeforeActivity(bot, event, cardID, userID, imageMessage); err != nil {
// 			log.Printf("Error saving image: %v", err)
// 			sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกรูปภาพ กรุณาลองใหม่.")
// 			return
// 		}

// 		// หลังจากบันทึกรูปภาพเสร็จแล้ว เปลี่ยนสถานะเป็น wait status saveEvidenceTime เพื่อเข้าในฟังก์ชัน saveEvidenceImageTime
// 		userState[userID] = "wait status saveEvidenceImageafterActivity"
// 		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปหลังการทำกิจกรรม.")
// 	} else {
// 		log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
// 		sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
// 	}
// }

// // ฟังก์ชันบันทึกรูปภาพกิจกรรม
// func handlesaveEvidenceImagebeforeActivity(bot *linebot.Client, event *linebot.Event, cardID, userID string, imageMessage *linebot.ImageMessage) error {
// 	if userState[userID] != "wait status saveEvidenceImagebeforeActivity" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
// 		return nil
// 	}
// 	//ตรวจสอบข้อความที่รับมา=Image
// 	messageID := imageMessage.ID
// 	log.Printf("Processing ImageMessage for user %s with message ID: %s", userID, messageID)

// 	// ดึงข้อมูลภาพ
// 	content, err := bot.GetMessageContent(messageID).Do()
// 	if err != nil {
// 		log.Printf("Error getting image content: %v", err)
// 		return err
// 	}
// 	defer content.Content.Close()

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		return err
// 	}
// 	defer db.Close()

// 	// ตรวจสอบข้อมูลที่เกี่ยวข้อง
// 	patientInfoID, err := GetPatientInfoIDByCardID(db, cardID)
// 	if err != nil {
// 		return err
// 	}
// 	activity, err := GetActivityNameByPatientInfoID(db, patientInfoID)
// 	if err != nil {
// 		log.Printf("Error fetching activity name: %v", err)
// 		return err
// 	}
// 	// ใช้สำหรับจัดเก็บไฟล์ชั่วคราวระหว่างการดำเนินงาน
// 	tempDir := os.TempDir()
// 	tempFilePath := fmt.Sprintf("%s\\%s.jpg", tempDir, messageID)
// 	file, err := os.Create(tempFilePath)
// 	if err != nil {
// 		log.Printf("Error creating temp file: %v", err)
// 		return err
// 	}
// 	defer file.Close()
// 	defer os.Remove(tempFilePath)

// 	// เขียนเนื้อหาภาพลงในไฟล์ (บันทึกเนื้อหาของรูปภาพหรือไฟล์ที่ได้รับจาก LINE Messaging API ลงในไฟล์ชั่วคราว)
// 	if _, err := io.Copy(file, content.Content); err != nil {
// 		log.Printf("Error writing image content to file: %v", err)
// 		return err
// 	}

// 	// เชื่อมต่อกับ MinIO
// 	minioClient, err := database.ConnectToMinio()
// 	if err != nil {
// 		log.Printf("Error connecting to MinIO: %v", err)
// 		return err
// 	}
// 	bucketName := "nirunimages"
// 	objectName := fmt.Sprintf("Image before activity/%s/%d/%s.jpg", activity, patientInfoID, messageID)

// 	// อัปโหลดไฟล์ไปยัง MinIO
// 	fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
// 	if err != nil {
// 		log.Printf("Error uploading file to MinIO: %v", err)
// 		return err
// 	}

// 	// อัปเดต URL ในฐานข้อมูล
// 	err = updateImagebeforeActivity(db, patientInfoID, fileURL)
// 	if err != nil {
// 		log.Printf("Error updating database: %v", err)
// 		return err
// 	}

// 	log.Printf("Activity Image successfully saved and URL updated: %s", fileURL)
// 	log.Printf("Last userID: %s", userID)
// 	return nil
// }

// // ฟังก์ชันบันทึกรูปหลังทำกิจกรรม
// func handlesaveEvidenceImageafterActivity(bot *linebot.Client, event *linebot.Event, cardID, userID string) error {
// 	if userState[userID] != "wait status saveEvidenceImageafterActivity" {
// 		log.Printf("Unhandled state for user %s. Current state: %s", userID, userState[userID])
// 		sendReply(bot, event.ReplyToken, "สถานะไม่ถูกต้อง กรุณาลองใหม่.")
// 		return nil
// 	}

// 	//ตรวจสอบข้อความที่รับมา=Image
// 	if imageMessage, ok := event.Message.(*linebot.ImageMessage); ok {
// 		messageID := imageMessage.ID
// 		log.Printf("Processing ImageMessage for user %s with message ID: %s", userID, messageID)

// 		// ดึงข้อมูลรูปภาพ
// 		content, err := bot.GetMessageContent(messageID).Do()
// 		if err != nil {
// 			log.Printf("Error getting image content: %v", err)
// 			return err
// 		}
// 		defer content.Content.Close()

// 		db, err := database.ConnectToDB()
// 		if err != nil {
// 			log.Printf("Database connection error: %v", err)
// 			return err
// 		}
// 		defer db.Close()

// 		// ตรวจสอบข้อมูลที่เกี่ยวข้อง
// 		patientInfoID, err := GetPatientInfoIDByCardID(db, cardID)
// 		if err != nil {
// 			return err
// 		}
// 		activity, err := GetActivityNameByPatientInfoID(db, patientInfoID)
// 		if err != nil {
// 			log.Printf("Error fetching activity name: %v", err)
// 			return err
// 		}

// 		// บันทึกไฟล์รูปภาพชั่วคราว
// 		tempDir := os.TempDir()
// 		tempFilePath := fmt.Sprintf("%s\\%s.jpg", tempDir, messageID)
// 		file, err := os.Create(tempFilePath)
// 		if err != nil {
// 			log.Printf("Error creating temp file: %v", err)
// 			return err
// 		}
// 		defer file.Close()
// 		defer os.Remove(tempFilePath)

// 		if _, err := io.Copy(file, content.Content); err != nil {
// 			log.Printf("Error writing image content to file: %v", err)
// 			return err
// 		}

// 		// อัปโหลดรูปภาพไปยัง MinIO
// 		minioClient, err := database.ConnectToMinio()
// 		if err != nil {
// 			log.Printf("Error connecting to MinIO: %v", err)
// 			return err
// 		}
// 		bucketName := "nirunimages"
// 		objectName := fmt.Sprintf("Image after activity/%s/%d/%s.jpg", activity, patientInfoID, messageID)

// 		// อัปโหลดไฟล์ไปยัง MinIO
// 		fileURL, err := UploadFileToMinIO(minioClient, bucketName, objectName, tempFilePath)
// 		if err != nil {
// 			log.Printf("Error uploading file to MinIO: %v", err)
// 			return err
// 		}

// 		// อัปเดต URL ในฐานข้อมูล
// 		err = updateImageafterActivity(db, patientInfoID, fileURL)
// 		if err != nil {
// 			log.Printf("Error updating database: %v", err)
// 			return err
// 		}

// 		log.Printf("Evidence time image successfully saved and URL updated: %s", fileURL)

// 		userState[userID] = "wait status ConfirmOrSaveEmployee"
// 		// handleUserChoiceForActivityRecord(bot, event, userID, "ยืนยันการบันทึก")
// 		log.Printf("User state updated to: %s", userState[userID])

// 		sendReply(bot, event.ReplyToken, "เลือก 'ยืนยันการบันทึก' หรือ 'บันทึกข้อมูลแทน'  ")
// 		return nil
// 	}

// 	log.Printf("Unhandled message type for user %s: %T", userID, event.Message)
// 	sendReply(bot, event.ReplyToken, "กรุณาส่งรูปภาพเท่านั้นในขั้นตอนนี้.")
// 	return nil
// }
// func handleUserChoiceForActivityRecord(bot *linebot.Client, event *linebot.Event, userID, selection string) {
// 	if userState[userID] != "wait status ConfirmOrSaveEmployee" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		return
// 	}

// 	selection = strings.TrimSpace(selection)
// 	log.Printf("Received selection: %s", selection)

// 	switch selection {
// 	case "บันทึกข้อมูลแทน", "save_employee", "บันทึกข้อมูล":
// 		log.Printf("User %s selected to save record for another employee", userID)
// 		userState[userID] = "wait status saveActivityRecordForOtherEmployee"
// 		saveActivityRecordForOtherEmployee(bot, event, userID)

// 	case "ยืนยันการบันทึก", "confirm", "ยืนยัน":
// 		log.Printf("User %s selected to confirm activity record", userID)
// 		userState[userID] = "wait status confirmActivityRecordByUser"
// 		confirmActivityRecordByUser(bot, event, userID)

// 	default:
// 		log.Printf("Invalid selection by user %s: %s", userID, selection)
// 		sendReply(bot, event.ReplyToken, "ตัวเลือกไม่ถูกต้อง กรุณาลองใหม่.\nกรุณาเลือก:\n- ยืนยันการบันทึก\n- บันทึกข้อมูลแทน")
// 	}
// }

// func confirmActivityRecordByUser(bot *linebot.Client, event *linebot.Event, userID string) {
// 	if userState[userID] != "wait status confirmActivityRecordByUser" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		return
// 	}

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ดึง activity_record_id ที่เพิ่งบันทึก
// 	activityRecord, err := GetLatestActivityRecord(db, userID)
// 	if err != nil {
// 		log.Printf("ไม่พบข้อมูลกิจกรรมสำหรับ UserID: %s", userID)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลกิจกรรม กรุณาลองใหม่.")
// 		return
// 	}

// 	// ตรวจสอบสิทธิ์พนักงาน
// 	userInfo, err := GetUserInfoByLINEID(db, userID)
// 	if err != nil || userInfo.EmployeeInfo.EmployeeInfo_ID == 0 {
// 		log.Printf("User %s is not an employee", userID)
// 		sendReply(bot, event.ReplyToken, "คุณไม่มีสิทธิ์ยืนยัน กรุณาเลือก 'บันทึกข้อมูลแทน'")
// 		return
// 	}

// 	// บันทึก employee_info_id ลง activity_record
// 	err = UpdateActivityEmployeeID(db, activityRecord.ActivityRecord_ID, userInfo.EmployeeInfo.EmployeeInfo_ID, userInfo.UserInfo_ID)
// 	if err != nil {
// 		log.Printf("Error updating employee_info_id: %v", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก กรุณาลองใหม่.")
// 		return
// 	}

// 	sendReply(bot, event.ReplyToken, "บันทึกกิจกรรมสำเร็จ!")
// 	userState[userID] = ""
// }

// // ฟังก์ชันให้ผู้ใช้กรอกชื่อพนักงานแทน
// func saveActivityRecordForOtherEmployee(bot *linebot.Client, event *linebot.Event, userID string) {
// 	if userState[userID] != "wait status saveActivityRecordForOtherEmployee" {
// 		log.Printf("Invalid state for user %s. Current state: %s", userID, userState[userID])
// 		return
// 	}

// 	sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการแทน:")
// 	userState[userID] = "wait status SaveEmployeeName"
// }

// // ฟังก์ชันบันทึกชื่อพนักงานที่ทำการบริการแทน
// func handleSaveEmployeeName(bot *linebot.Client, event *linebot.Event, userID, State, employeeName string) {
// 	// ตรวจสอบค่าที่ได้รับ
// 	// log.Printf("User %s entered employee name (before trim): '%s'", userID, employeeName)

// 	employeeName = strings.TrimSpace(event.Message.(*linebot.TextMessage).Text)
// 	log.Printf("Received employee name: %s", employeeName)
// 	if employeeName == "" {
// 		sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการ.")
// 		return
// 	}
// 	// //ตรวจสอบว่าค่าที่รับมาเป็นค่าว่างหรือไม่
// 	// if employeeName == "" {
// 	// 	sendReply(bot, event.ReplyToken, "กรุณากรอกชื่อพนักงานที่ให้บริการแทน")
// 	// 	return
// 	// }

// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		log.Printf("Database connection error: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่.")
// 		return
// 	}
// 	defer db.Close()

// 	// ค้นหา employeeID
// 	employeeID, err := GetEmployeeIDByName(db, employeeName)
// 	if err != nil {
// 		log.Printf("ไม่พบข้อมูล employee_info_id สำหรับพนักงาน: '%s'", employeeName)
// 		sendReply(bot, event.ReplyToken, fmt.Sprintf("ไม่พบพนักงานชื่อ %s\nกรุณากรอกชื่อใหม่", employeeName))
// 		return
// 	}

// 	log.Printf("Employee ID found: %d for name: %s", employeeID, employeeName)

// 	//อัปเดตข้อมูลใน activity_record
// 	cardID := usercardidState[State]
// 	patient, err := GetPatientInfoByName(db, cardID)
// 	if err != nil {
// 		log.Printf("Error fetching patient_info_id: %v", err)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่")
// 		return
// 	}
// 	//ดึงข้อมูลผู้ใช้ตาม LINE ID
// 	userInfo, err := GetUserInfoByLINEID(db, userID)
// 	if err != nil {
// 		log.Println("Error fetching user info:", err)
// 		sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่.")
// 		return
// 	}
// 	// ดึง Activity Record ID
// 	activityRecordID, err := GetActivityRecordID(db, cardID)
// 	if err != nil {
// 		sendReply(bot, event.ReplyToken, err.Error())
// 		return
// 	}
// 	// เตรียมข้อมูล activityRecord
// 	activityRecord := &models.Activityrecord{
// 		ActivityRecord_ID: activityRecordID.ActivityRecord_ID,
// 		PatientInfo: models.PatientInfo{
// 			CardID:         cardID,
// 			Name:           patient.PatientInfo.Name,
// 			PatientInfo_ID: patient.PatientInfo.PatientInfo_ID,
// 		},
// 		EmployeeInfo: models.EmployeeInfo{EmployeeInfo_ID: employeeID},
// 		UserInfo:     models.User_info{UserInfo_ID: userInfo.UserInfo_ID},
// 	}
// 	//เหลือแก้การดึงการคำนวณ
// 	// startTime, err := GetActivityStartTime(db, cardID, userActivity[userID])
// 	// if err != nil {
// 	// 	log.Printf("Error fetching StartTime: %v", err)
// 	// 	sendReply(bot, event.ReplyToken, "ไม่พบข้อมูลเวลาเริ่ม กรุณาลองใหม่")
// 	// 	return
// 	// }
// 	// duration := activityRecord.EndTime.Sub(startTime)
// 	// activityRecord.Period = formatDuration(duration)

// 	log.Printf("Updating employee_info_id=%d for ActivityRecord_ID=%d", employeeID, activityRecord.ActivityRecord_ID)

// 	if err := UpdateActivityEmployee(db, activityRecord); err != nil {
// 		log.Printf("Error updating end time: %v", err)
// 		sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึกเวลาสิ้นสุด กรุณาลองใหม่")
// 		return
// 	}
// 	flexMessage := flexmessage.FormatactivityRecordEndtime([]models.Activityrecord{*activityRecord})
// 	if _, err := bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
// 		log.Printf("Error sending reply message: %v", err)
// 		// sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการส่งข้อความ กรุณาลองใหม่.")
// 		return
// 	}
// 	// err = UpdateActivityEndTimeForPatient(db, activityRecord.ActivityRecord_ID, employeeID, 0)
// 	// if err != nil {
// 	// 	log.Printf("Error updating employee_info_id: %v", err)
// 	// 	sendReply(bot, event.ReplyToken, "เกิดข้อผิดพลาดในการบันทึก กรุณาลองใหม่.")
// 	// 	return
// 	// }

// 	sendReply(bot, event.ReplyToken, "บันทึกข้อมูลพนักงานสำเร็จ!")
// 	userState[userID] = ""
// }

// // แปลงระยะเวลาของกิจกรรมเป็นชั่วโมงและนาที
// // func formatDuration(d time.Duration) string {
// // 	hours := int(d.Hours())
// // 	minutes := int(d.Minutes()) % 60
// // 	return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
// // }

// // ตรวจสอบกิจกรรมที่รับมาตรงกับฐานข้อมูลไหม
// func validateActivity(activity string) bool {
// 	allowedActivities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กันฟัง", "ซุโดกุ", "จับคู่ภาพ",
// 	}
// 	for _, allowed := range allowedActivities {
// 		if activity == allowed {
// 			return true
// 		}
// 	}
// 	return false
// }

// func handleDefault(bot *linebot.Client, event *linebot.Event) {
// 	sendCustomReply(bot, event.ReplyToken, event.Source.UserID, "กรุณาเลือกเมนู")
// }

// // ส่งข้อความตอบกลับแบบธรรมดา
// func sendReply(bot *linebot.Client, replyToken, message string) {
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
// 		log.Println("Error sending reply message funcsendReply:", err)
// 	}
// }

// // ส่งข้อความตอบกลับพร้อมปุ่ม
// func sendReplyWithQuickReply(bot *linebot.Client, replyToken string, message string, quickReply *linebot.QuickReplyItems) {
// 	textMessage := linebot.NewTextMessage(message).WithQuickReplies(quickReply)
// 	if _, err := bot.ReplyMessage(replyToken, textMessage).Do(); err != nil {
// 		log.Printf("Error sending reply with quick reply: %v", err)
// 	}
// }
