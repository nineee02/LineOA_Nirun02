package event

import (
	"fmt"
	"nirun/pkg/models"
)

func FormatConfirmationCheckIn(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("--ยืนยันการเช็คอิน--\n\n%s\nรหัสพนักงาน: %s",
		worktimeRecord.EmployeeInfo.Name,
		worktimeRecord.EmployeeInfo.EmployeeCode)
}

// func FormatConfirmationCheckIn2(worktimeRecord *models.WorktimeRecord) string {
// 	return fmt.Sprintf("--ยืนยันการเช็คอิน--\n\n%s\nรหัสพนักงาน: %s",
// 		worktimeRecord.EmployeeInfo.Name,
// 		worktimeRecord.EmployeeInfo.EmployeeCode)
// }

func FormatworktimeCheckin(worktimeRecord *models.WorktimeRecord) string {
	if worktimeRecord == nil {
		return "ไม่พบข้อมูลการทำงาน กรุณาลองใหม่."
	}
	return fmt.Sprintf("--ยินดีต้อนรับ--\n\nชื่อ: %s\nรหัสพนักงาน: %s\nเช็คอินที่: %s",
		worktimeRecord.EmployeeInfo.Name,
		worktimeRecord.EmployeeInfo.EmployeeCode,
		worktimeRecord.CheckIn.Format("2006-01-02 15:04:05 PM"))
}

func FormatConfirmationCheckOut(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("--ยืนยันการเช็คเอ้าท์--\n\n%s\nรหัสพนักงาน: %s",
		worktimeRecord.EmployeeInfo.Name,
		worktimeRecord.EmployeeInfo.EmployeeCode)
}

func FormatworktimeCheckout(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("--ลาก่อน--\n\nชื่อ: %s\nรหัสพนักงาน: %s\nเช็คเอ้าท์ที่: %s",
		worktimeRecord.EmployeeInfo.Name,
		worktimeRecord.EmployeeInfo.EmployeeCode,
		worktimeRecord.CheckOut.Format("2006-01-02 15:04:05 PM"))
}

func FormatPatientInfo(patient *models.Activityrecord) string {
	return fmt.Sprintf(
		"ข้อมูลผู้ป่วย:\n- ชื่อ: %s\n- เบอร์โทร: %s\n- ที่อยู่: %s\n- อายุ: %s ปี\n- เพศ: %s\n- กลุ่มเลือด: %s\n- ADL %s\n- สิทธิ์การรักษา: %s",
		patient.PatientInfo.Name,
		patient.PatientInfo.PhoneNumber,
		patient.PatientInfo.Address,
		patient.PatientInfo.Age,
		patient.PatientInfo.Sex,
		patient.PatientInfo.Blood,
		patient.PatientInfo.ADL,
		patient.PatientInfo.RightToTreatmentInfo.Right_to_treatment,
	)
}

func FormatServiceInfo(activity []models.Activityrecord) string {
	// สร้างข้อความสำหรับชื่อผู้ป่วยและกิจกรรมที่สำเร็จแล้ว
	message := fmt.Sprintf("ชื่อผู้รับบริการ: %s\n", activity[0].PatientInfo.Name)
	for _, info := range activity {
		message += fmt.Sprintf(" %s\n", info.ServiceInfo.Activity)
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
func FormatactivityRecordStarttime(starttime []models.Activityrecord) string {
	var result string
	for _, record := range starttime {
		result += fmt.Sprintf(
			"เริ่มบันทึกกิจกรรม: %s\nของ %s\nที่: %s\n\nกรุณาพิมพ์ 'เสร็จสิ้น' เมื่อทำกิจกรรมเสร็จ",
			record.ServiceInfo.Activity,
			record.PatientInfo.Name,
			record.StartTime.Format("2006-01-02 15:04:05 PM"),
		)
	}
	return result
}
func FormatactivityRecordEndtime(endtime []models.Activityrecord) string {
	var result string
	for _, record := range endtime {
		result += fmt.Sprintf(
			"บันทึกกิจกรรม: %s\nของ %s\nที่: %s\n\n!!!สำเร็จ!!!",
			record.ServiceInfo.Activity,
			record.PatientInfo.Name,
			record.EndTime.Format("2006-01-02 15:04:05 PM"),
		)
	}
	return result
}

// // *************ReplyError*****************************************************************************************
// func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
// 	if _, err := bot.ReplyMessage(
// 		replyToken,
// 		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
// 	).Do(); err != nil {
// 		log.Println("ReplyErrorFormat:", err)
// 	}
// }

// // ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
// func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
// 	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
// 		log.Println("ReplyErrorFormat:", err)
// 	}
// }
