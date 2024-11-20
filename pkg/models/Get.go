package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

// ฟังก์ชันสำหรับดึงข้อมูลผู้ป่วยตาม Name
func GetPatientInfoByName(db *sql.DB, name_ string) (*PatientInfo, error) {
	query := `SELECT name_, patiet_id, age, sex, blood, phone_numbers FROM patient_info WHERE name_ = ?`
	row := db.QueryRow(query, name_)

	var patientInfo PatientInfo
	err := row.Scan(&patientInfo.Name, &patientInfo.PatientID, &patientInfo.Age, &patientInfo.Sex, &patientInfo.Blood, &patientInfo.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลผู้ป่วยที่มีชื่อ %s", name_)
		}
		return nil, err
	}
	log.Println("ข้อมูลผู้ป่วยที่ดึงมา:", &patientInfo)
	return &patientInfo, nil
}

// ฟังก์ชันสำหรับดึงข้อมูลกิจกรรม
func GetServiceInfoByName(db *sql.DB, name string) ([]ServiceInfo, error) {
	query := `SELECT patient_info.id, patient_info.patiet_id, patient_info.name_, service_info.activity
	FROM patient_info
	INNER JOIN service_info 
	ON patient_info.id = service_info.service_id WHERE patient_info.name_ = ?`
	rows, err := db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serviceInfos []ServiceInfo
	for rows.Next() {
		var serviceInfo ServiceInfo
		err := rows.Scan(&serviceInfo.PatientInfo.ID, &serviceInfo.PatientInfo.PatientID, &serviceInfo.PatientInfo.Name, &serviceInfo.Activity)
		if err != nil {
			return nil, err
		}
		//log.Printf("ดึงข้อมูลกิจกรรม: %+v\n", serviceInfo) // ตรวจสอบข้อมูล
		serviceInfos = append(serviceInfos, serviceInfo)
	}

	if len(serviceInfos) == 0 {
		return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรมสำหรับผู้ป่วย: %s", name)
	}

	return serviceInfos, nil

}

func GetAllActivity(db *sql.DB, patientID string) ([]ServiceInfo, error) {
	// ดึงข้อมูลกิจกรรมทั้งหมดจากฐานข้อมูล
	query := `SELECT activity FROM service_info WHERE patient_info = ?`
	rows, err := db.Query(query, patientID)

	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลกิจกรรมได้: %v", err)
	}
	defer rows.Close()

	//activities คือ Slice ว่างสำหรับเก็บข้อมูลกิจกรรม (Activity) ทั้งหมดที่ดึงมาจากฐานข้อมูล
	var activities []ServiceInfo

	//วนลูปในผลลัพธ์ที่ได้จากการ Query ฐานข้อมูล
	for rows.Next() {
		var activity ServiceInfo //activity เก็บข้อมูลของกิจกรรมที่ได้จากผลลัพธ์
		err := rows.Scan(&activity.Activity)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil

}

// **********************************************************************************************************************
// FormatPatientInfo จัดรูปแบบข้อมูลผู้ป่วยให้อยู่ในรูปแบบข้อความที่เหมาะสมสำหรับการแสดงผลหรือส่งกลับไปยังผู้ใช้
func FormatPatientInfo(patient *PatientInfo) string {
	return fmt.Sprintf("ข้อมูลผู้ป่วย:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nหมายเลขโทรศัพท์: %s",
		patient.Name, patient.PatientID, patient.Age, patient.Sex, patient.Blood, patient.PhoneNumber)
}

// func FormatAllActivity(activities []ServiceInfo) string {
//     var message string
//     message = "กิจกรรมทั้งหมดที่คุณสามารถเลือกได้:\n"
//     for _, activity := range activities {
//         message += fmt.Sprintf("- %s\n", activity.Activity)
//     }
//     return message
// }

// formatServiceInfo จัดรูปแบบข้อมูลกิจกรรมของผู้สูงอายุให้เหมาะสมสำหรับการแสดงผล
func FormatServiceInfo(serviceInfo []ServiceInfo) string {
	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(serviceInfo) == 0 {
		return "ไม่พบกิจกรรมสำหรับผู้ป่วยนี้"
	}

	// สร้างข้อความสำหรับชื่อผู้ป่วยและกิจกรรมที่สำเร็จแล้ว
	message := fmt.Sprintf("ชื่อผู้ป่วย: %s\nกิจกรรมที่สำเร็จแล้ว:\n", serviceInfo[0].PatientInfo.Name)
	for _, info := range serviceInfo {
		message += fmt.Sprintf("- %s\n", info.Activity)
	}

	// เพิ่มรายการกิจกรรมที่สามารถเลือกเพิ่มได้
	activities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
	}
	message += "\nเลือกกิจกรรมที่คุณต้องการเพิ่ม:\n"
	for _, activity := range activities {
		message += fmt.Sprintf("- %s\n", activity)
	}
	return message
}

// func FormatSaveactivitysuccess(activity string) string {
// 	return fmt.Sprintf("บันทึกกิจกรรม '%s' สำเร็จ!\n", activity)

// }

// ******************************************************************************************************************************************
// replyErrorFormat ส่งข้อความตัวอย่างการใช้งานที่ถูกต้องกลับไปยังผู้ใช้ เมื่อรูปแบบคำสั่งที่ได้รับไม่ถูกต้อง
func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
		//linebot.NewTextMessage("กรุณากรอกรูปแบบ 'ข้อมูลผู้กิจกรรม []'"),
	).Do(); err != nil {
		log.Println("เกิดข้อผิดพลาดในการส่งข้อความ:", err)
	}
}

// ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
		log.Println("Error sending not found message:", err)
	}
}

// ฟังก์ชัน replyDatabaseError ข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล
// func ReplyDatabaseError(bot *linebot.Client, replyToken string) {
// 	dbErrorMessage := "เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล กรุณาลองใหม่อีกครั้งภายหลัง"
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(dbErrorMessage)).Do(); err != nil {
// 		log.Println("Error sending database error message:", err)
// 	}
// }
