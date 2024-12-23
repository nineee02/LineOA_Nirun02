package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/models"
	"time"
)

func GetWorktime(db *sql.DB, employeeCode string) (*models.EmployeeInfo, error) {
	query := `SELECT e.employee_info_id, e.employee_code, e.username, e.phone_number, e.email, 
	d.department, j.job_position
	FROM employee_info e
	LEFT JOIN department_info d ON e.department_info_id = d.department_info_id
	LEFT JOIN job_position_info j ON e.job_position_info_id = j.job_position_info_id
	WHERE e.employee_code = ?`

	row := db.QueryRow(query, employeeCode)

	var employee models.EmployeeInfo
	var department sql.NullString
	var jobPosition sql.NullString

	// เพิ่ม log เพื่อช่วย debug
	// log.Println("Executing query for employee code:", employeeCode)
	// log.Println("Executing query for employee code:", employeeCode)

	err := row.Scan(
		&employee.EmployeeInfo_ID,
		&employee.EmployeeCode,
		&employee.Name,
		&employee.PhoneNumber,
		&employee.Email,
		&department,
		&jobPosition,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No employee found with code: %s", employeeCode)
			return nil, fmt.Errorf("ไม่มีข้อมูลของพนักงาน %s", employeeCode)
		}
		log.Println("Error executing query:", err)
		return nil, err
	}

	// จัดการค่าที่อาจเป็น NULL
	employee.DepartmentInfo.Department = department.String
	employee.JobPositionInfo.JobPosition = jobPosition.String

	// Log ข้อมูลที่ได้จากการ Query
	// log.Println("Query result for employee:", employee)

	return &employee, nil
}

func InsertWorktime(db *sql.DB, employee_code *models.WorktimeRecord) error {
	query := `INSERT INTO worktime_record
			(employee_id, check_in, check_out)
			VALUES (?, ?, ?)`
	_, err := db.Exec(query, employee_code.EmployeeInfo.EmployeeCode, employee_code.CheckIn, employee_code.CheckOut)
	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม: %v", err)
	}
	return nil
}

// GetPatientInfoByName ค้นหาข้อมูลผู้ป่วยจากชื่อ
func GetPatientInfoByName(db *sql.DB, card_id string) (*models.Activityrecord, error) {
	query := `SELECT p.card_id,
					p.patient__info_id,
					p.username, 
					p.phone_number, 
					p.email, 
					p.address,
					p.date_of_birth, 
					p.age,
					p.sex, 

					c.country_info.id, 
					c.country, 
					c.create_date,
					c.write_date,

					r.religion_info_id, 
					r.religion,
					r.create_date,
					r.write_date,

	                b.blood_info_id,
					b.blood,
					 
					rtt.right_to_treatment_info_id, 
    				rtt.right_to_treatment, 
					rtt.create_date,
					rtt.write_date
					
	        FROM patient_info p
			LEFT JOIN blood_info  b ON p.blood_info_id = b.blood_info_id
			LEFT JOIN country_info  c ON p.country_info_id = c.country_info_id
			LEFT JOIN religion  r ON p.religion_info_id = r.religion_info_id
			LEFT JOIN right_to_treatment_info  rtt ON p.right_to_treatment_info_id = r.right_to_treatment_info_id 
	    	WHERE p.card_id LIKE ?`

	patient := &models.Activityrecord{}
	err := db.QueryRow(query, "%"+card_id+"%").Scan(
		&patient.PatientInfo.PatientInfo_ID,
		&patient.PatientInfo.CardID,
		&patient.PatientInfo.Name,
		&patient.PatientInfo.PhoneNumber,
		&patient.PatientInfo.Email,
		&patient.PatientInfo.Address,
		&patient.PatientInfo.DateOfBirth,
		&patient.PatientInfo.Age,
		&patient.PatientInfo.Sex,
		&patient.PatientInfo.CreateDate,
		&patient.PatientInfo.WriteDate,

		&patient.PatientInfo.BloodInfo.BloodInfo_ID,
		&patient.PatientInfo.BloodInfo.Blood,
		&patient.PatientInfo.CountryInfo.CountryInfo_ID,
		&patient.PatientInfo.CountryInfo.Country,
		&patient.PatientInfo.CountryInfo.CreateDate,
		&patient.PatientInfo.CountryInfo.WriteDate,
		&patient.PatientInfo.Religion.ReligionInfo_ID,
		&patient.PatientInfo.Religion.Religion,
		&patient.PatientInfo.Religion.CreateDate,
		&patient.PatientInfo.Religion.WriteDate,
		&patient.PatientInfo.RightToTreatmentInfo.RightToTreatmentInfo_ID,
		&patient.PatientInfo.RightToTreatmentInfo.Right_to_treatment,
		&patient.PatientInfo.RightToTreatmentInfo.CreateDate,
		&patient.PatientInfo.RightToTreatmentInfo.WriteDate,
	)
	if err != nil {
		return nil, fmt.Errorf("ไม่พบข้อมูลผู้สูงอายุ: %v", err)
	}
	return patient, nil
}

func GetServiceInfoBycardID(db *sql.DB, card_id string) ([]models.Activityrecord, error) {
	query := `SELECT p.card_id, 
					p.username, 
					s.service_info_id,
					s.activity,
			  FROM activity_record a
			  INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
			  INNER JOIN service_info s ON a.service_info_id = s.service_info_id
			  WHERE p.crad_id =?`

	rows, err := db.Query(query, card_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activityrecord []models.Activityrecord

	for rows.Next() {
		var record models.Activityrecord
		var patientInfo models.PatientInfo
		var serviceInfo models.ServiceInfo

		err := rows.Scan(
			&patientInfo.CardID,
			&patientInfo.Name,
			&serviceInfo.ServiceInfo_Id)
		if err != nil {
			return nil, err
		}
		// Assign ข้อมูลให้กับ Activityrecord
		record.PatientInfo = patientInfo
		record.ServiceInfo = serviceInfo

		activityrecord = append(activityrecord, record)
	}

	if len(activityrecord) == 0 {
		return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรม: %s", card_id)
	}

	return activityrecord, nil
}

func GetActivityRecord(db *sql.DB, activity *models.Activityrecord) error {
	// query สำหรับการบันทึกข้อมูลกิจกรรมลงในฐานข้อมูล
	query := `INSERT INTO activity_record
	  				(card_id, activity)
					VALUES (?, ?)`

	// ใช้ข้อมูลจาก activity เพื่อบันทึก
	_, err := db.Exec(query, activity.PatientInfo.CardID, activity.ServiceInfo.Activity)
	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม: %v", err)
	}
	return nil
}
func GetActivityPeriodRecord(db *sql.DB, activity *models.Activityrecord) error {
	// query สำหรับการบันทึกข้อมูลระยะเวลาและเวลาที่กรอก
	query := `UPDATE activity_record
			  SET period = ?, end_time = ?
			  WHERE activity = ?`

	// ใช้ time.Now() สำหรับ end_time
	endTime := time.Now()

	_, err := db.Exec(query, activity.Period, endTime, activity.ActivityRecord_ID)
	if err != nil {
		return fmt.Errorf("ไม่สามารถบันทึกระยะเวลา: %v", err)
	}
	return nil
}

// **********************************************************************************************************************
// FormatPatientInfo จัดรูปแบบข้อมูลผู้ป่วยให้อยู่ในรูปแบบข้อความที่เหมาะสมสำหรับการแสดงผลหรือส่งกลับไปยังผู้ใช้
// func FormatPatientInfo(patient *models.PatientInfo) string {
// 	return fmt.Sprintf("ข้อมูลผู้ป่วย:\nชื่อ: %s\nรหัสผู้ป่วย: %s\nอายุ: %d\nเพศ: %s\nหมู่เลือด: %s\nหมายเลขโทรศัพท์: %s",
// 		patient.Name, patient.PatientID, patient.Age, patient.Sex, patient.Blood, patient.PhoneNumber)
// }

// // formatServiceInfo จัดรูปแบบข้อมูลกิจกรรมของผู้สูงอายุให้เหมาะสมสำหรับการแสดงผล
// func FormatServiceInfo(activity []models.PatientInfo) string {
// 	// สร้างข้อความสำหรับชื่อผู้ป่วยและกิจกรรมที่สำเร็จแล้ว
// 	message := fmt.Sprintf("ชื่อผู้ป่วย: %s\nกิจกรรมที่สำเร็จแล้ว:\n", activity[0].Name)
// 	for _, info := range activity {
// 		message += fmt.Sprintf("- %s\n", info.Activityrecord)
// 	}

// 	// เพิ่มรายการกิจกรรมที่สามารถเลือกเพิ่มได้
// 	activities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// 	}
// 	message += "\nเลือกกิจกรรมที่คุณต้องการเพิ่ม:\n"
// 	for _, activity := range activities {
// 		message += fmt.Sprintf("- %s\n", activity)
// 		for _, activity := range activities {
// 			message += fmt.Sprintf("- %s\n", activity)
// 		}
// 		return message
// 	}
// 	return message
// }

// ******************************************************************************************************************************************
// replyErrorFormat ส่งข้อความตัวอย่างการใช้งานที่ถูกต้องกลับไปยังผู้ใช้ เมื่อรูปแบบคำสั่งที่ได้รับไม่ถูกต้อง
// func ReplyErrorFormat(bot *linebot.Client, replyToken string) {
// 	if _, err := bot.PushMessage(
// 		replyToken,
// 		linebot.NewTextMessage("กรุณากรอกรูปแบบข้อความให้ถูกต้อง เช่น 'นางสมหวัง สดใส'"),
// 		//linebot.NewTextMessage("กรุณากรอกรูปแบบ 'ข้อมูลผู้กิจกรรม []'"),
// 	).Do(); err != nil {
// 		log.Println("เกิดข้อผิดพลาดในการส่งข้อความ:", err)
// 	}
// }

// // ฟังก์ชัน replyDataNotFound แจ้งผู้ใช้เมื่อไม่พบข้อมูลผู้สูงอายุ
// func ReplyDataNotFound(bot *linebot.Client, replyToken string) {
// 	notFoundMessage := "ไม่พบข้อมูลผู้สูงอายุตามชื่อ กรุณาลองใหม่อีกครั้ง"
// 	if _, err := bot.PushMessage(replyToken, linebot.NewTextMessage(notFoundMessage)).Do(); err != nil {
// 		log.Println("Error sending not found message:", err)
// 	}
// }

// ฟังก์ชัน replyDatabaseError ข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล
// func ReplyDatabaseError(bot *linebot.Client, replyToken string) {
// 	dbErrorMessage := "เกิดข้อผิดพลาดในการเชื่อมต่อฐานข้อมูล กรุณาลองใหม่อีกครั้งภายหลัง"
// 	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(dbErrorMessage)).Do(); err != nil {
// 		log.Println("Error sending database error message:", err)
// 	}
// }

// // func GetEmployee(db *sql.DB, NameEmployee string) (*Employee, error) {
// // 	query := "INSERT INTO employee (username, start_time) VALUES (?, ?)"
// // 	startTime := time.Now().Format("2006-01-02 15:04:05")

// // 	log.Printf("Executing query: %s, Values: %s, %s", query, NameEmployee, startTime)

// // 	_, err := db.Exec(query, NameEmployee, startTime_ServiceInfo)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("ไม่สามารถบันทึกเวลาเข้างานได้: %v", err)
// // 	}

// // 	return &Employee{Name: NameEmployee, Starttime_ServiceInfo: startTime}, nil
// // }

// // **********************************************************************************************************************
