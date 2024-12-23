package event

import (
	"fmt"
	"log"
	"nirun/pkg/models"

	"github.com/line/line-bot-sdk-go/linebot"
)

func FormatGetworktimeCheckin(employee *models.EmployeeInfo) string {
	return fmt.Sprintf("CHECK_IN SUCCESS!!\nชื่อ: %s\nรหัสพนักงาน: %s\nแผนก: %s\nตำแหน่ง: %s",
		employee.Name, employee.EmployeeCode, employee.DepartmentInfo.Department, employee.JobPositionInfo.JobPosition)
}
func FormatGetworktimeCheckout(employee *models.EmployeeInfo) string {
	return fmt.Sprintf("CHECK_OUT SUCCESS!!\nชื่อ: %s\nรหัสพนักงาน: %s\nแผนก: %s\nตำแหน่ง: %s",
		employee.Name, employee.EmployeeCode, employee.DepartmentInfo.Department, employee.JobPositionInfo.JobPosition)
}

func FormatPatientInfo(patient *models.PatientInfo) string {
	return fmt.Sprintf("ข้อมูลผู้สูงอายุ:\nชื่อ: %s\nเลขประจำตัวประชาชน: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nหมายเลขโทรศัพท์: %s",
		patient.Name, patient.CardID, patient.PatientInfo_ID, patient.Age, patient.Sex, patient.Blood, patient.PhoneNumber, patient)
}

func FormatServiceInfo(activity []models.Activityrecord) string {
	// สร้างข้อความสำหรับชื่อผู้ป่วยและกิจกรรมที่สำเร็จแล้ว
	message := fmt.Sprintf("ชื่อผู้รับบริการ: %s\nกิจกรรมที่สำเร็จแล้ว:\n", activity[0].PatientInfo.Name)
	for _, info := range activity {
		message += fmt.Sprintf("- %s\n", info.ServiceInfo.Activity)
	}

	// เพิ่มรายการกิจกรรมที่สามารถเลือกเพิ่มได้
	activities := []string{
		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
	}
	message += "\nเลือกกิจกรรมที่คุณต้องการเพิ่ม:\n"
	for _, activity := range activities {
		message += fmt.Sprintf("- %s\n", activity)
		for _, activity := range activities {
			message += fmt.Sprintf("- %s\n", activity)
		}
		return message
	}
	return message
}

// func FormatEmployee(employeeInfo []models.Employee) string {
// 	message := fmt.Sprintf("ชื่อพนักงาน: %s\nกิจกรรมที่สำเร็จแล้ว:\n", employeeInfo[0].Name_Employee)
// 	for _, info := range employeeInfo {
// 		message += fmt.Sprintf("- %s\n", info.EmployeeID)
// 	}
// 	return message // คืนค่าข้อความที่สร้างขึ้น
// }

// *************ReplyError*****************************************************************************************
func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
	).Do(); err != nil {
		log.Println("ReplyErrorFormat:", err)
	}
}

// ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
		log.Println("ReplyErrorFormat:", err)
	}
}
