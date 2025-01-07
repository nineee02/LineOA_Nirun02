package event

import (
	"fmt"
	"nirun/pkg/models"
	"strings"
	"time"
)

func FormatConfirmationCheckIn(worktimeRecord *models.WorktimeRecord) string {
	return fmt.Sprintf("--ยืนยันการเช็คอิน--\n\n%s\nรหัสพนักงาน: %s",
		worktimeRecord.EmployeeInfo.Name,
		worktimeRecord.EmployeeInfo.EmployeeCode)
}

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
func FormatHistoryAll(history []*models.Activityrecord) string {
	var result string
	var totalActivities int

	// Header
	result += fmt.Sprintf("%-6s %-6s %-10s\n", "ปี", "กิจกรรม", "จำนวนครั้ง")
	// result += fmt.Sprintf("%s\n", strings.Repeat("-", 40)) // เส้นแบ่ง

	for _, record := range history {
		detail := record.ActivityYearDetail

		// รวมจำนวนกิจกรรมทั้งหมด
		totalActivities += detail.Total

		// เพิ่มข้อมูลในแต่ละแถว
		result += fmt.Sprintf("%-6d %-6s %7d ครั้ง\n", detail.Year, detail.ActivityType, detail.Total)
	}

	// Footer รวมทั้งหมด
	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30)) // เส้นแบ่ง
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง\n", totalActivities)

	return result
}

func FormatHistoryofYear(history []*models.Activityrecord) string {
	var result string
	var totalActivities int

	// Header
	result += fmt.Sprintf("ประวัติกิจกรรมปีนี้: %d\n", time.Now().Year())
	result += fmt.Sprintf("\n%-20s %-10s\n", "กิจกรรม", "จำนวนครั้ง") // หัวตาราง
	// result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))            // เส้นแบ่ง

	for _, record := range history {
		detail := record.ActivityYearDetail
		totalActivities += detail.Total // จำนวนของกิจกรรมในแต่ละ record

		result += fmt.Sprintf("-%-20s %10d ครั้ง\n", detail.ActivityType, detail.Total)
	}

	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง", totalActivities)
	return result
}
func FormatHistoryofMonth(history []*models.Activityrecord) string {
	var result string
	var totalActivities int
	var monthName string

	monthMap := map[int]string{
		1: "มกราคม", 2: "กุมภาพันธ์", 3: "มีนาคม",
		4: "เมษายน", 5: "พฤษภาคม", 6: "มิถุนายน",
		7: "กรกฎาคม", 8: "สิงหาคม", 9: "กันยายน",
		10: "ตุลาคม", 11: "พฤศจิกายน", 12: "ธันวาคม",
	}

	// Header
	if len(history) > 0 {
		monthName = monthMap[history[0].ActivityMonthDetail.Month]
	} else {
		monthName = monthMap[int(time.Now().Month())] // แปลง time.Month เป็น int
	}
	result += fmt.Sprintf("ประวัติกิจกรรมเดือนนี้: %s\n", monthName)
	result += fmt.Sprintf("\n%-20s %-10s\n", "กิจกรรม", "จำนวนครั้ง")
	// result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))

	for _, record := range history {
		detail := record.ActivityMonthDetail
		totalActivities += detail.Total

		result += fmt.Sprintf("-%-20s %10d ครั้ง\n", detail.ActivityType, detail.Total)
	}

	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง", totalActivities)
	return result
}
func FormatHistoryofWeek(history []*models.Activityrecord) string {
	var result string
	var totalActivities int

	// คำนวณช่วงวันที่ของสัปดาห์นี้
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1) // วันจันทร์
	endOfWeek := startOfWeek.AddDate(0, 0, 6)               // วันอาทิตย์

	// Header
	result += fmt.Sprintf("ประวัติกิจกรรมสัปดาห์นี้ (%s - %s):\n",
		startOfWeek.Format("02 มกราคม 2006"), endOfWeek.Format("02 มกราคม 2006"))
	result += fmt.Sprintf("\n%-20s %-10s\n", "กิจกรรม", "จำนวนครั้ง")
	// result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))

	for _, record := range history {
		detail := record.ActivityWeekDetail
		totalActivities += detail.Total // นับจำนวนกิจกรรม

		result += fmt.Sprintf("-%-20s %10d ครั้ง\n", detail.ActivityType, detail.Total)
	}

	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง", totalActivities)
	return result
}

func FormatHistoryofDay(history []*models.Activityrecord) string {
	var result string
	var totalActivities int
	// ดึงวันที่ปัจจุบันในรูปแบบที่ต้องการ
	currentDate := time.Now().Format("02 มกราคม 2006")

	// Header
	result += fmt.Sprintf("ประวัติกิจกรรมวันนี้ (%s):\n", currentDate)
	result += fmt.Sprintf("\n%-20s %-10s\n", "กิจกรรม", "จำนวนครั้ง")
	// result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))

	for _, record := range history {
		detail := record.ActivityDayDetail
		totalActivities += detail.Total

		result += fmt.Sprintf("-%-20s %10d ครั้ง\n", detail.ActivityType, detail.Total)
	}

	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง", totalActivities)
	return result
}

func FormatHistoryofSet(history []*models.Activityrecord, startDate, endDate string) string {
	var result string
	var totalActivities int

	// Header
	result += fmt.Sprintf("ประวัติกิจกรรมระหว่างวันที่ %s ถึง %s:\n", startDate, endDate)
	result += fmt.Sprintf("\n%-20s %10s\n", "กิจกรรม", "จำนวนครั้ง")
	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))

	// Loop through records to append activity details
	for _, record := range history {
		detail := record.ActivitySetDetail
		totalActivities += detail.Total

		// Append each activity detail
		result += fmt.Sprintf("%-20s %10d ครั้ง\n", detail.ActivityType, detail.Total)
	}

	// Add total activities count
	result += fmt.Sprintf("%s\n", strings.Repeat("-", 30))
	result += fmt.Sprintf("จำนวนกิจกรรมทั้งหมด: %d ครั้ง\n", totalActivities)

	// Check for no activities case
	if totalActivities == 0 {
		result = fmt.Sprintf("ไม่มีข้อมูลกิจกรรมระหว่างวันที่ %s ถึง %s", startDate, endDate)
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
