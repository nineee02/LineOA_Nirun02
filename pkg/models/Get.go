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
func GetServiceInfoByName(db *sql.DB, name string) (*ServiceInfo, error) {
	query := `SELECT id, patiet_id,name_,activity 
	FROM patient_info
	INNER JOIN service_info 
	ON patient_info.id = service_info.service_id WHERE name_ =?`
	row := db.QueryRow(query, name)

	var serviceInfo ServiceInfo
	err := row.Scan(&serviceInfo.ID, &serviceInfo.PatientInfo.ID, serviceInfo, &serviceInfo.PatientInfo.Name, &serviceInfo.Activity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรม %s", name)
		}
		return nil, err
	}
	log.Println("ข้อมูลกิจกรรมที่ดึงมา:", &serviceInfo)
	return &serviceInfo, nil
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

func SaveActivity(db *sql.DB, patientID string, ID int) error {
	// อัปเดตกิจกรรมที่เลือก
	query := `UPDATE service_info SET selected = 1 WHERE patient_id = ? AND id = ?`
	_, err := db.Exec(query, patientID, ID)
	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรมได้: %v", err)
	}
	return nil
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
func FormatServiceInfo(serviceInfo *ServiceInfo) string {
	return fmt.Sprintf("กรอกข้อมูล:")
}

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
func ReplyDatabaseError(bot *linebot.Client, replyToken string) {
	dbErrorMessage := "เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล กรุณาลองใหม่อีกครั้งภายหลัง"
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(dbErrorMessage)).Do(); err != nil {
		log.Println("Error sending database error message:", err)
	}
}
