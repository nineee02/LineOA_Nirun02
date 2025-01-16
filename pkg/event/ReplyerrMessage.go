package event

import (
	"fmt"
	"nirun/pkg/models"
)

//	func FormatConfirmationCheckIn(worktimeRecord *models.WorktimeRecord) string {
//		return fmt.Sprintf("--ยืนยันการเช็คอิน--\n\n%s\nรหัสพนักงาน: %s",
//			worktimeRecord.EmployeeInfo.Name,
//			worktimeRecord.EmployeeInfo.EmployeeCode)
//	}
func FormatConfirmationCheckOutNocheckin(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("%s\nคุณยังไม่ได้เช็คอิน กรุณาเช็คอินก่อนทำการเช็คเอ้าท์",
		worktimeRecord.UserInfo.Name)
}

func FormatConfirmationCheckInNocheckout(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("%s\nได้เช็คอินในระบบอยู่แล้ว กรุณาทำการเช็คเอ้าท์ก่อนทำการเช็คอินอีกครั้ง",
		worktimeRecord.UserInfo.Name)
}

func FormatworktimeCheckin(worktimeRecord *models.WorktimeRecord) string {

	return fmt.Sprintf("--ยินดีต้อนรับ--\n\n%s\nเช็คอินที่: %s",
		worktimeRecord.UserInfo.Name,
		worktimeRecord.CheckIn.Format("2006-01-02 15:04:05 PM"))
}

// func FormatConfirmationCheckOut(worktimeRecord *models.WorktimeRecord) string {
// 	return fmt.Sprintf("--ยืนยันการเช็คเอ้าท์--\n\n%s\nรหัสพนักงาน: %s",
// 		worktimeRecord.EmployeeInfo.Name,
// 		worktimeRecord.EmployeeInfo.EmployeeCode)
// }

func FormatworktimeCheckout(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("--ลาก่อน--\n\n%s\nเช็คเอ้าท์ที่: %s\n%s",
		worktimeRecord.UserInfo.Name,
		worktimeRecord.Period,
		worktimeRecord.CheckOut.Format("2006-01-02 15:04:05 PM"))

}

func FormatPatientInfo(patient *models.Activityrecord) string {
	return fmt.Sprintf(
		"ข้อมูลผู้ป่วย:\n- ชื่อ: %s\n-เลขประจำตัวประชาชน: %s\n- เบอร์โทร: %s\n- ที่อยู่: %s\n- อายุ: %s ปี\n- เพศ: %s\n- กลุ่มเลือด: %s\n- ADL %s\n- สิทธิ์การรักษา: %s",
		patient.PatientInfo.Name,
		patient.PatientInfo.CardID,
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
			"บันทึกกิจกรรม: %s\nของ %s\nที่: %s\n%s\n!!!สำเร็จ!!!",
			record.ServiceInfo.Activity,
			record.PatientInfo.Name,
			record.EndTime.Format("2006-01-02 15:04:05 PM"),
			record.Period,
		)
	}
	return result
}
