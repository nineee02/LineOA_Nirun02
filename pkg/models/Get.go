package models

import (
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

// FormatPatientInfo จัดรูปแบบข้อมูลผู้ป่วยให้อยู่ในรูปแบบข้อความที่เหมาะสมสำหรับการแสดงผลหรือส่งกลับไปยังผู้ใช้
func FormatPatientInfo(patient *PatientInfo) string {
	return fmt.Sprintf("ข้อมูลผู้ป่วย:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nหมายเลขโทรศัพท์: %s",
		patient.Name, patient.PatientID, patient.Age, patient.Sex, patient.Blood, patient.PhoneNumber)
}

// // formatServiceInfo จัดรูปแบบข้อมูลกิจกรรมของผู้สูงอายุให้เหมาะสมสำหรับการแสดงผล
// func FormatServiceInfo(serviceInfo *ServiceInfo) string {
// 	return fmt.Sprintf("ข้อมูลผู้สูงอายุ:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nกิจกรรม: %s\n\nกรุณาพิมพ์ชื่อกิจกรรมที่ต้องการ:",
// 		serviceInfo.PatientInfo.Name,
// 		serviceInfo.PatientInfo.PatientID,
// 		serviceInfo.Activity)

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
func ReplyDatabaseError(bot *linebot.Client, replyToken string) {
	dbErrorMessage := "เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล กรุณาลองใหม่อีกครั้งภายหลัง"
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(dbErrorMessage)).Do(); err != nil {
		log.Println("Error sending database error message:", err)
	}
}
