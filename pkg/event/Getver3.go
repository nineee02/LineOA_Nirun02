package event

// import (
// 	"database/sql"
// 	"fmt"
// 	"nirun/pkg/models"
// 	"time"
// )

// // GetPatientInfoByName ค้นหาข้อมูลผู้ป่วยจากชื่อ
// func GetPatientInfoByName(db *sql.DB, card_id string) (*models.PatientInfo, error) {
// 	query := `SELECT card_id, patient_id, username, phone_number, email, 
// 	                 address, country,  religion, sex, 
// 	                 blood, date_of_birth, age 
// 	          FROM patient_info 
// 	          WHERE card_id LIKE ?`

// 	patient := &models.PatientInfo{}
// 	err := db.QueryRow(query, "%"+card_id+"%").Scan(
// 		&patient.CardID,
// 		&patient.PatientID,
// 		&patient.Name,
// 		&patient.PhoneNumber,
// 		&patient.Email,
// 		&patient.Address,
// 		&patient.Country,
// 		&patient.Religion,
// 		&patient.Sex,
// 		&patient.Blood,
// 		&patient.DateOfBirth,
// 		&patient.Age,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("ไม่พบข้อมูลผู้ป่วย: %v", err)
// 	}
// 	return patient, nil
// }

// func GetServiceInfoBycardID(db *sql.DB, card_id string) ([]models.PatientInfo, error) {
// 	query := `SELECT patient_info.card_id, patient_info.username, activity_record.activity
// 			  FROM patient_info
// 			  INNER JOIN activity_record ON patient_info.card_id = activity_record .card_id
// 			  WHERE patient_info.card_id = ?`

// 	rows, err := db.Query(query, card_id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var activity []models.PatientInfo
// 	for rows.Next() {
// 		var patientinfo models.PatientInfo
// 		err := rows.Scan(&patientinfo.Activityrecord.CardID, &patientinfo.Name,
// 			&patientinfo.Activityrecord.Activity)
// 		if err != nil {
// 			return nil, err
// 		}
// 		//log.Printf("ดึงข้อมูลกิจกรรม: %+v\n", serviceInfo) // ตรวจสอบข้อมูล
// 		activity = append(activity, patientinfo)
// 	}

// 	if len(activity) == 0 {
// 		return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรม: %s", card_id)
// 	}

// 	return activity, nil
// }

// func GetActivityRecord(db *sql.DB, activity *models.Activityrecord) error {
// 	// query สำหรับการบันทึกข้อมูลกิจกรรมลงในฐานข้อมูล
// 	query := `INSERT INTO activity_record
// 	  				(card_id, activity)
// 					VALUES (?, ?)`

// 	// ใช้ข้อมูลจาก activity เพื่อบันทึก
// 	_, err := db.Exec(query, activity.CardID, activity.Activity)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม: %v", err)
// 	}
// 	return nil
// }
// func ActivityPeriodRecord(db *sql.DB, activity *models.Activityrecord) error {
// 	// query สำหรับการบันทึกข้อมูลระยะเวลาและเวลาที่กรอก
// 	query := `UPDATE activity_record
// 			  SET period = ?, end_time = ?
// 			  WHERE activity = ?`

// 	// ใช้ time.Now() สำหรับ end_time
// 	endTime := time.Now()

// 	_, err := db.Exec(query, activity.Period, endTime, activity.Activity)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถบันทึกระยะเวลา: %v", err)
// 	}
// 	return nil
// }

// // **********************************************************************************************************************
// // FormatPatientInfo จัดรูปแบบข้อมูลผู้ป่วยให้อยู่ในรูปแบบข้อความที่เหมาะสมสำหรับการแสดงผลหรือส่งกลับไปยังผู้ใช้
// // func FormatPatientInfo(patient *models.PatientInfo) string {
// // 	return fmt.Sprintf("ข้อมูลผู้ป่วย:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nหมายเลขโทรศัพท์: %s",
// // 		patient.Name, patient.PatientID, patient.Age, patient.Sex, patient.Blood, patient.PhoneNumber)
// // }

// // // formatServiceInfo จัดรูปแบบข้อมูลกิจกรรมของผู้สูงอายุให้เหมาะสมสำหรับการแสดงผล
// // func FormatServiceInfo(activity []models.PatientInfo) string {
// // 	// สร้างข้อความสำหรับชื่อผู้ป่วยและกิจกรรมที่สำเร็จแล้ว
// // 	message := fmt.Sprintf("ชื่อผู้ป่วย: %s\nกิจกรรมที่สำเร็จแล้ว:\n", activity[0].Name)
// // 	for _, info := range activity {
// // 		message += fmt.Sprintf("- %s\n", info.Activityrecord)
// // 	}

// // 	// เพิ่มรายการกิจกรรมที่สามารถเลือกเพิ่มได้
// // 	activities := []string{
// // 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// // 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// // 	}
// // 	message += "\nเลือกกิจกรรมที่คุณต้องการเพิ่ม:\n"
// // 	for _, activity := range activities {
// // 		message += fmt.Sprintf("- %s\n", activity)
// // 		for _, activity := range activities {
// // 			message += fmt.Sprintf("- %s\n", activity)
// // 		}
// // 		return message
// // 	}
// // 	return message
// // }

// // ******************************************************************************************************************************************
// // replyErrorFormat ส่งข้อความตัวอย่างการใช้งานที่ถูกต้องกลับไปยังผู้ใช้ เมื่อรูปแบบคำสั่งที่ได้รับไม่ถูกต้อง
// // func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
// // 	if _, err := bot.PushMessage(
// // 		replyToken,
// // 		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
// // 		//linebot.NewTextMessage("กรุณากรอกรูปแบบ 'ข้อมูลผู้กิจกรรม []'"),
// // 	).Do(); err != nil {
// // 		log.Println("เกิดข้อผิดพลาดในการส่งข้อความ:", err)
// // 	}
// // }

// // // ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
// // func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
// // 	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
// // 	if _, err := bot.PushMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
// // 		log.Println("Error sending not found message:", err)
// // 	}
// // }

// // ฟังก์ชัน replyDatabaseError ข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล
// // func ReplyDatabaseError(bot *linebot.Client, replyToken string) {
// // 	dbErrorMessage := "เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล กรุณาลองใหม่อีกครั้งภายหลัง"
// // 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(dbErrorMessage)).Do(); err != nil {
// // 		log.Println("Error sending database error message:", err)
// // 	}
// // }

// // // func GetEmployee(db *sql.DB, NameEmployee string) (*Employee, error) {
// // // 	query := "INSERT INTO employee (username, start_time) VALUES (?, ?)"
// // // 	startTime := time.Now().Format("2006-01-02 15:04:05")

// // // 	log.Printf("Executing query: %s, Values: %s, %s", query, NameEmployee, startTime)

// // // 	_, err := db.Exec(query, NameEmployee, startTime_ServiceInfo)
// // // 	if err != nil {
// // // 		return nil, fmt.Errorf("ไม่สามารถบันทึกเวลาเข้างานได้: %v", err)
// // // 	}

// // // 	return &Employee{Name: NameEmployee, Starttime_ServiceInfo: startTime}, nil
// // // }

// // // **********************************************************************************************************************
