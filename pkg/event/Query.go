package event

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/models"
	"time"

	"github.com/minio/minio-go/v7"
)

// ดึงข้อมูลผู้ใช้ตาม LINE ID
func GetUserInfoByLINEID(db *sql.DB, lineUserID string) (*models.User_info, error) {
	query := `SELECT user_info_id, line_user_id, sex, name_surname, email, phone_number, create_date, update_date
	          FROM user_info WHERE line_user_id = ?`
	row := db.QueryRow(query, lineUserID)

	user := &models.User_info{}
	err := row.Scan(
		&user.UserInfo_ID,
		&user.Line_user_id,
		&user.Sex,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.CreateDate,
		&user.UpdateDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ LINE User ID: %s", lineUserID)
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้: %v", err)
	}

	return user, nil
}

// ตรวจสอบการเช็คอินของพนักงาน
func IsEmployeeCheckedIn(db *sql.DB, userInfoID int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM worktime_record
		WHERE user_info_id = ? AND check_out IS NULL`
	var count int
	err := db.QueryRow(query, userInfoID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking check-in status: %v", err)
	}
	return count > 0, nil
}

// บันทึกการเช็คอิน
func RecordCheckIn(db *sql.DB, userID int) error {
	// ตรวจสอบว่าผู้ใช้เช็คอินอยู่หรือไม่
	var existingCheckIn sql.NullTime
	queryCheck := `SELECT check_in FROM worktime_record WHERE user_info_id = ? AND check_out IS NULL`
	err := db.QueryRow(queryCheck, userID).Scan(&existingCheckIn)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing check-in: %v", err)
	}

	// ถ้าพบว่ามีการเช็คอินอยู่แล้ว
	if existingCheckIn.Valid {
		return fmt.Errorf("ผู้ใช้ได้ทำการเช็คอินไปแล้ว")
	}

	// ดำเนินการเช็คอิน
	queryInsert := `
		INSERT INTO worktime_record (
			check_in,
			create_by,
			update_by,
			user_info_id,
			create_date,
			update_date
		)
		VALUES (NOW(), ?, ?, ?, NOW(), NOW())`
	_, err = db.Exec(queryInsert, userID, userID, userID)
	if err != nil {
		return fmt.Errorf("error recording check-in: %v", err)
	}

	return nil
}

// บันทึกการเช็คเอ้าท์
func RecordCheckOut(db *sql.DB, userID int) error {
	query := `
		UPDATE worktime_record
		SET check_out = NOW(),
			period = TIMEDIFF(NOW(), check_in),
			update_date = NOW(),
			update_by = ?
		WHERE user_info_id = ? AND check_out IS NULL`
	_, err := db.Exec(query, userID, userID)
	if err != nil {
		return fmt.Errorf("error recording check-out: %v", err)
	}
	return nil
}

// AutoCheckOut - อัปเดต Check-out เป็นเที่ยงคืนอัตโนมัติ
func AutoCheckOut(db *sql.DB) error {
	query := `
		UPDATE worktime_record
		SET check_out = NOW(),
			update_date = NOW()
		WHERE check_out IS NULL`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error performing auto check-out: %v", err)
	}
	log.Println("Auto check-out completed successfully.")
	return nil
}

// รันเซิร์ฟใหม่ทุกๆเที่ยงคืน
func StartAutoCheckOutScheduler(db *sql.DB) {
	ticker := time.NewTicker(24 * time.Hour) // รันทุก 24 ชั่วโมง
	go func() {
		for {
			<-ticker.C
			log.Println("Running auto check-out at midnight...")
			err := AutoCheckOut(db)
			if err != nil {
				log.Println("Auto check-out error:", err)
			}
		}
	}()
}

// ดึงข้อมูลบันทึกเวลาทำงานสำหรับผู้ใช้ตาม ID
func GetWorktimeRecordByUserID(db *sql.DB, UserInfo_ID int) (*models.WorktimeRecord, error) {
	query := `SELECT
    	w.worktime_record_id,
   		w.check_in,
    	w.check_out,
    	w.period,
    	u.user_name
		FROM worktime_record w
		INNER JOIN user_info u ON w.user_info_id = u.user_info_id
		WHERE u.user_info_id = ?
		ORDER BY w.check_in DESC
		LIMIT 1;`

	row := db.QueryRow(query, UserInfo_ID)
	log.Println("Executing query:", query)
	log.Printf("Query parameter userID: %d", UserInfo_ID)

	worktimeRecord := &models.WorktimeRecord{
		UserInfo: &models.User_info{},
	}

	// ดึงค่า period เป็น string โดยตรง
	err := row.Scan(
		&worktimeRecord.WorktimeRecord_ID,
		&worktimeRecord.CheckIn,
		&worktimeRecord.CheckOut,
		&worktimeRecord.Period,
		&worktimeRecord.UserInfo.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// กรณีไม่มีข้อมูล
			log.Println("No worktime record found for userID:", UserInfo_ID)
			return nil, nil
		}
		log.Printf("Error scanning row for userID %d: %v", UserInfo_ID, err)
		return nil, fmt.Errorf("error fetching worktime record: %v", err)
	}

	log.Printf("Fetched worktime record: %+v", worktimeRecord)
	return worktimeRecord, nil
}

//**********************************************************************************************************************************************

// ดึง Service Info ID สำหรับกิจกรรม
func GetServiceInfoIDByActivity(db *sql.DB, activity string) (int, error) {
	query := "SELECT service_info_id FROM service_info WHERE activity = ?"
	var serviceInfoID int
	err := db.QueryRow(query, activity).Scan(&serviceInfoID)
	if err != nil {
		return 0, fmt.Errorf("ไม่พบ service_info_id สำหรับกิจกรรม: %s, error: %v", activity, err)
	}
	return serviceInfoID, nil
}
func GetPatientInfoByName(db *sql.DB, cardID string) (*models.Activityrecord, error) {
	query := `SELECT
				p.card_id,
				p.patient_info_id,
				p.name_surname,
				p.phone_number,
				p.email,
				p.address,
				p.date_of_birth,
				p.age,
				p.sex,
				p.blood,
				p.ADL,

				c.country_info_id,
				c.country,
				c.create_date,
				c.update_date,

				r.religion_info_id,
				r.religion,
				r.create_date,
				r.update_date,

				rtt.right_to_treatment_info_id,
				rtt.right_to_treatment,
				rtt.create_date,
				rtt.update_date
			FROM patient_info p
			LEFT JOIN country_info c ON p.country_info_id = c.country_info_id
			LEFT JOIN religion_info r ON p.religion_info_id = r.religion_info_id
			LEFT JOIN right_to_treatment_info rtt ON p.right_to_treatment_info_id = rtt.right_to_treatment_info_id
			WHERE p.card_id = ?`

	// สร้างโครงสร้างเพื่อเก็บผลลัพธ์
	patient := &models.Activityrecord{}
	// var imagePath []byte

	err := db.QueryRow(query, cardID).Scan(
		&patient.PatientInfo.CardID,
		&patient.PatientInfo.PatientInfo_ID,
		&patient.PatientInfo.Name,
		&patient.PatientInfo.PhoneNumber,
		&patient.PatientInfo.Email,
		&patient.PatientInfo.Address,
		&patient.PatientInfo.DateOfBirth,
		&patient.PatientInfo.Age,
		&patient.PatientInfo.Sex,
		&patient.PatientInfo.Blood,
		&patient.PatientInfo.ADL,

		&patient.PatientInfo.CountryInfo.CountryInfo_ID,
		&patient.PatientInfo.CountryInfo.Country,
		&patient.PatientInfo.CountryInfo.CreateDate,
		&patient.PatientInfo.CountryInfo.UpdateDate,

		&patient.PatientInfo.Religion.ReligionInfo_ID,
		&patient.PatientInfo.Religion.Religion,
		&patient.PatientInfo.Religion.CreateDate,
		&patient.PatientInfo.Religion.UpdateDate,

		&patient.PatientInfo.RightToTreatmentInfo.RightToTreatmentInfo_ID,
		&patient.PatientInfo.RightToTreatmentInfo.Right_to_treatment,
		&patient.PatientInfo.RightToTreatmentInfo.CreateDate,
		&patient.PatientInfo.RightToTreatmentInfo.UpdateDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลผู้ป่วยที่มี CardID: %s", cardID)
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ป่วย: %v", err)
	}

	// // แปลง imagePath จาก []byte เป็น string และเก็บในโครงสร้าง
	// patient.PatientInfo.ImagePath = string(imagePath)

	return patient, nil
}

// กิจกรรมมิติเทคโนโลยี
func GetTechnologyActivities(db *sql.DB) ([]models.ActivityTechnologyInfo, error) {
	query := `SELECT activity_technology_info_id, activity, service_type, create_date FROM activity_technology_info`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivityTechnologyInfo
	for rows.Next() {
		var activity models.ActivityTechnologyInfo
		if err := rows.Scan(&activity.ActivityTechnologyInfo_ID, &activity.ActivityTechnology, &activity.ServiceType, &activity.CreateDate); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// กิจกรรมมิติสังคม
func GetSocialActivities(db *sql.DB) ([]models.ActivitySocialInfo, error) {
	query := `SELECT activity_social_info_id, activity, service_type, create_date FROM activity_social_info`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivitySocialInfo
	for rows.Next() {
		var activity models.ActivitySocialInfo
		if err := rows.Scan(&activity.ActivitySocialInfo_ID, &activity.ActivitySocial, &activity.ServiceType, &activity.CreateDate); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// กิจกรรมมิติสุขภาพ
func GetHealthActivities(db *sql.DB) ([]models.ActivityHealthInfo, error) {
	query := `SELECT activity_health_info_id, activity, service_type, create_date FROM activity_health_info`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivityHealthInfo
	for rows.Next() {
		var activity models.ActivityHealthInfo
		if err := rows.Scan(&activity.ActivityHealthInfo_ID, &activity.ActivityHealth, &activity.ServiceType, &activity.CreateDate); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// กิจกรรมมิติเศรษฐกิจ
func GetEconomicActivities(db *sql.DB) ([]models.ActivityEconomicInfo, error) {
	query := `SELECT activity_economic_info_id, activity, service_type, create_date FROM activity_economic_info`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivityEconomicInfo
	for rows.Next() {
		var activity models.ActivityEconomicInfo
		if err := rows.Scan(&activity.ActivityEconomicInfo_ID, &activity.ActivityEconomic, &activity.ServiceType, &activity.CreateDate); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// กิจกรรมมิติสภาพแวดล้อม
func GetEnvironmentalActivities(db *sql.DB) ([]models.ActivityEnvironmentalInfo, error) {
	query := `SELECT activity_environmental_info_id, activity, service_type, create_date FROM activity_environmental_info`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivityEnvironmentalInfo
	for rows.Next() {
		var activity models.ActivityEnvironmentalInfo
		if err := rows.Scan(&activity.ActivityEnvironmentalInfo_ID, &activity.ActivityEnvironmental, &activity.ServiceType, &activity.CreateDate); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
func GetActivityInfoIDByType(db *sql.DB, category string, activityName string) (int, error) {
	var query string
	switch category {
	case "technology":
		query = `SELECT activity_technology_info_id FROM activity_technology_info WHERE activity = ?`
	case "social":
		query = `SELECT activity_social_info_id FROM activity_social_info WHERE activity = ?`
	case "health":
		query = `SELECT activity_health_info_id FROM activity_health_info WHERE activity = ?`
	case "economic":
		query = `SELECT activity_economic_info_id FROM activity_economic_info WHERE activity = ?`
	case "environmental":
		query = `SELECT activity_environmental_info_id FROM activity_environmental_info WHERE activity = ?`
	default:
		return 0, fmt.Errorf("invalid activity category")
	}

	var activityInfoID int
	err := db.QueryRow(query, activityName).Scan(&activityInfoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("ไม่พบข้อมูลกิจกรรม: %s", activityName)
		}
		return 0, err
	}
	return activityInfoID, nil
}

// บันทึกระยะเวลา
func SaveActivityDuration(db *sql.DB, cardID string, period string) error {
	query := `
		UPDATE activity_record
		SET period = ?
		WHERE patient_info_id = (SELECT patient_info_id FROM patient_info WHERE card_id = ?)
		AND end_time IS NULL
		LIMIT 1
	`
	_, err := db.Exec(query, period, cardID)
	if err != nil {
		log.Printf("Error updating activity period for cardID %s: %v", cardID, err)
		return fmt.Errorf("error updating activity period: %v", err)
	}

	log.Printf("Successfully updated period for patient with cardID: %s", cardID)
	return nil
}

// [บันทึกกกิจกรรม]
func SaveActivityRecord(db *sql.DB, activity *models.Activityrecord, category string) error {
	// ✅ ดึง patient_info_id
	patient, err := GetPatientInfoByName(db, activity.PatientInfo.CardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		return fmt.Errorf("error fetching patient_info_id: %v", err)
	}
	activity.PatientInfo.PatientInfo_ID = patient.PatientInfo.PatientInfo_ID

	// ✅ ตรวจสอบ patient_info_id
	if activity.PatientInfo.PatientInfo_ID == 0 {
		log.Println("Invalid patient_info_id")
		return fmt.Errorf("patient_info_id is missing or invalid")
	}

	// ✅ ตรวจสอบ activity_info_id และเลือกคอลัมน์ที่ถูกต้อง
	var activityInfoColumn string
	var activityInfoID interface{} // ใช้ `interface{}` เพื่อรองรับ `NULL`

	switch category {
	case "technology":
		activityInfoColumn = "activity_technology_info_id"
		activityInfoID = activity.ActivityRecord_ID
	case "social":
		activityInfoColumn = "activity_social_info_id"
		activityInfoID = activity.ActivityRecord_ID
	case "health":
		activityInfoColumn = "activity_health_info_id"
		activityInfoID = activity.ActivityRecord_ID
	case "economic":
		activityInfoColumn = "activity_economic_info_id"
		activityInfoID = activity.ActivityRecord_ID
	case "environmental":
		activityInfoColumn = "activity_environmental_info_id"
		activityInfoID = activity.ActivityRecord_ID
	case "other":
		activityInfoColumn = "activity_other"
		activityInfoID = activity.ActivityOther
	default:
		log.Println("Invalid category:", category)
		return fmt.Errorf("invalid category selected")
	}

	// ✅ ใช้ Dynamic SQL Query เพื่อเลือกคอลัมน์ที่ถูกต้อง
	query := fmt.Sprintf(`
		INSERT INTO activity_record (
			patient_info_id,
			%s, 
			start_time,
			create_by,
			update_by
		) VALUES (?, ?, ?, ?, ?)
	`, activityInfoColumn)

	_, err = db.Exec(query,
		activity.PatientInfo.PatientInfo_ID,
		activityInfoID,
		activity.StartTime,
		activity.UserInfo.UserInfo_ID,
		activity.UserInfo.UserInfo_ID,
	)
	if err != nil {
		log.Printf("Error inserting activity record: %v", err)
		return fmt.Errorf("error inserting activity record: %v", err)
	}

	log.Printf("✅ ActivityRecord saved successfully!")
	return nil
}

// บันทึกกิจกรรม
// func SaveActivityRecord(db *sql.DB, activity *models.Activityrecord) error {
// 	// ดึง patient_info_id
// 	patient, err := GetPatientInfoByName(db, activity.PatientInfo.CardID)
// 	if err != nil {
// 		log.Printf("Error fetching patient_info_id: %v", err)
// 		return fmt.Errorf("error fetching patient_info_id: %v", err)
// 	}
// 	activity.PatientInfo.PatientInfo_ID = patient.PatientInfo.PatientInfo_ID
// 	// ตรวจสอบว่าค่า patient_info_id ถูกต้อง
// 	if activity.PatientInfo.PatientInfo_ID == 0 {
// 		log.Println("Invalid patient_info_id")
// 		return fmt.Errorf("patient_info_id is missing or invalid")
// 	}

// 	// ดึง service_info_id
// 	serviceInfoID, err := GetServiceInfoIDByActivity(db, activity.ServiceInfo.Activity)
// 	if err != nil {
// 		log.Printf("Error fetching service_info_id: %v", err)
// 		return fmt.Errorf("error fetching service_info_id: %v", err)
// 	}

// 	// ตรวจสอบว่า service_info_id มีอยู่จริง
// 	if serviceInfoID == 0 {
// 		log.Println("Invalid service_info_id")
// 		return fmt.Errorf("service_info_id is invalid or missing")
// 	}

// 	query := `
// 		INSERT INTO activity_record (
// 			patient_info_id,
// 			service_info_id,
// 			start_time,
// 			create_by,
// 			write_by
// 		)
// 		VALUES (?, ?, ?, ?, ?)
// 		`
// 	result, err := db.Exec(query,
// 		activity.PatientInfo.PatientInfo_ID,
// 		serviceInfoID,
// 		time.Now(),
// 		activity.UserInfo.UserInfo_ID, // ใช้ UserInfo_ID สำหรับ create_by
// 		activity.UserInfo.UserInfo_ID, // ใช้ UserInfo_ID สำหรับ write_by
// 	)
// 	if err != nil {
// 		log.Printf("Error inserting activity record: %v", err)
// 		return fmt.Errorf("error inserting activity record: %v", err)
// 	}

// 	// ดึง activity_record_id ที่ถูกเพิ่มขึ้นมา
// 	activityRecordID, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("Error retrieving last insert id: %v", err)
// 		return fmt.Errorf("error retrieving last insert id: %v", err)
// 	}

// 	// บันทึก activityRecordID ลงใน activity
// 	activity.ActivityRecord_ID = int(activityRecordID)

// 	log.Printf("activity.ActivityRecord_ID:%d", activity.ActivityRecord_ID)
// 	log.Println("Activity record saved successfully")
// 	return nil

// }

// ดึงข้อมูลactivity_recordผ่าน card_idของผู้สูงอายุ
func GetActivityRecordID(db *sql.DB, cardID string) (*models.Activityrecord, error) {
	query := `
        SELECT a.activity_record_id, a.patient_info_id, a.start_time
        FROM activity_record a
        INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
        WHERE p.card_id = ? AND a.end_time IS NULL
        ORDER BY a.start_time DESC
        LIMIT 1
    `

	activityRecord := &models.Activityrecord{}
	err := db.QueryRow(query, cardID).Scan(
		&activityRecord.ActivityRecord_ID,
		&activityRecord.PatientInfo.PatientInfo_ID,
		&activityRecord.StartTime,
	)

	if err != nil {
		log.Printf("ไม่พบข้อมูลกิจกรรมที่ยังไม่เสร็จสิ้นสำหรับ CardID: %s", cardID)
		return nil, fmt.Errorf("ไม่มีกิจกรรมที่ยังไม่เสร็จสิ้นสำหรับผู้ใช้นี้")
	}

	log.Printf("พบกิจกรรม: ActivityRecord_ID=%d, PatientInfo_ID=%d, StartTime=%v",
		activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID, activityRecord.StartTime)

	return activityRecord, nil
}
func GetLatestActivityRecord(db *sql.DB, userID string) (*models.Activityrecord, error) {
	// ดึง patient_info_id จาก userID ก่อน
	patientInfoID, err := GetPatientInfoByName(db, userID)
	if err != nil {
		log.Printf("ไม่พบข้อมูล patient_info_id สำหรับ UserID: %s", userID)
		return nil, fmt.Errorf("ไม่พบข้อมูลผู้ป่วยสำหรับบัญชีนี้")
	}

	// ค้นหากิจกรรมล่าสุดของ patient_info_id
	query := `
        SELECT a.activity_record_id, a.patient_info_id, a.start_time, a.end_time, a.period
        FROM activity_record a
        WHERE a.patient_info_id = ?
        ORDER BY a.start_time DESC
        LIMIT 1
    `

	activityRecord := &models.Activityrecord{}
	err = db.QueryRow(query, patientInfoID).Scan(
		&activityRecord.ActivityRecord_ID,
		&activityRecord.PatientInfo.PatientInfo_ID,
		&activityRecord.StartTime,
		&activityRecord.EndTime,
		&activityRecord.Period,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ไม่พบข้อมูลกิจกรรมล่าสุดสำหรับ PatientInfo_ID: %d", patientInfoID)
			return nil, fmt.Errorf("ไม่พบกิจกรรมที่เกี่ยวข้องกับผู้ป่วยนี้")
		}
		log.Printf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรมล่าสุด: %v", err)
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรมล่าสุด")
	}

	log.Printf("พบกิจกรรมล่าสุด: ActivityRecord_ID=%d, PatientInfo_ID=%d, StartTime=%v, EndTime=%v, Period=%s",
		activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID, activityRecord.StartTime, activityRecord.EndTime, activityRecord.Period)

	return activityRecord, nil
}
func GetLatestActivityRecordByPatientID(db *sql.DB, patientInfoID int) (*models.Activityrecord, error) {
	query := `
        SELECT a.activity_record_id, a.patient_info_id, a.start_time, a.end_time, a.period
        FROM activity_record a
        WHERE a.patient_info_id = ?
        ORDER BY a.start_time DESC
        LIMIT 1
    `

	activityRecord := &models.Activityrecord{}
	err := db.QueryRow(query, patientInfoID).Scan(
		&activityRecord.ActivityRecord_ID,
		&activityRecord.PatientInfo.PatientInfo_ID,
		&activityRecord.StartTime,
		&activityRecord.EndTime,
		&activityRecord.Period,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ไม่พบข้อมูลกิจกรรมล่าสุดสำหรับ PatientInfo_ID: %d", patientInfoID)
			return nil, fmt.Errorf("ไม่พบกิจกรรมที่เกี่ยวข้องกับผู้ป่วยนี้")
		}
		log.Printf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรมล่าสุด: %v", err)
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรมล่าสุด")
	}

	log.Printf("พบกิจกรรมล่าสุด: ActivityRecord_ID=%d, PatientInfo_ID=%d, StartTime=%v, EndTime=%v, Period=%s",
		activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID, activityRecord.StartTime, activityRecord.EndTime, activityRecord.Period)

	return activityRecord, nil
}

// func GetActivityRecordID(db *sql.DB, cardID string) (*models.Activityrecord, error) {
// 	query := `
//         SELECT a.activity_record_id, a.start_time
//         FROM activity_record a
//         INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
//         WHERE p.card_id = ? AND a.end_time IS NULL
//         ORDER BY a.start_time DESC
//         LIMIT 1
//     `

// 	activityRecord := &models.Activityrecord{}
// 	err := db.QueryRow(query, cardID).Scan(
// 		&activityRecord.ActivityRecord_ID,
// 		&activityRecord.StartTime,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Printf("❌ ไม่พบข้อมูลกิจกรรมที่ยังไม่เสร็จสิ้นสำหรับ CardID: %s", cardID)
// 			return nil, fmt.Errorf("❌ ไม่มีกิจกรรมที่ยังไม่เสร็จสิ้นสำหรับผู้ใช้นี้")
// 		}
// 		log.Printf("❌ Error fetching activity record for CardID %s: %v", cardID, err)
// 		return nil, fmt.Errorf("❌ เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรม")
// 	}

// 	log.Printf("✅ พบกิจกรรมที่ยังไม่เสร็จสิ้น: ActivityRecord_ID=%d, StartTime=%v", activityRecord.ActivityRecord_ID, activityRecord.StartTime)
// 	return activityRecord, nil
// }

// ข้อมูล start_time ของกิจกรรมผ่าน cardID
func GetActivityStartTime(db *sql.DB, cardID string, activity string) (time.Time, error) {
	query := `
		SELECT a.start_time
		FROM activity_record a
		INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
		INNER JOIN service_info s ON a.service_info_id = s.service_info_id
		WHERE p.card_id = ? AND s.activity = ? AND a.end_time IS NULL
		LIMIT 1
	`

	var startTime time.Time
	err := db.QueryRow(query, cardID, activity).Scan(&startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, fmt.Errorf("ไม่พบเวลาเริ่มสำหรับ card_id: %s และกิจกรรม: %s", cardID, activity)
		}
		return time.Time{}, fmt.Errorf("เกิดข้อผิดพลาดในการดึงเวลาเริ่ม: %v", err)
	}

	return startTime, nil
}

func UpdateActivityStartTime(db *sql.DB, activityRecordID int, startTime time.Time) error {
	query := `
		UPDATE activity_record
		SET start_time = ?
		WHERE activity_record_id = ?`
	_, err := db.Exec(query, startTime, activityRecordID)
	if err != nil {
		return fmt.Errorf("error updating activity start time: %v", err)
	}
	return nil
}
func InsertActivityStartTime(db *sql.DB, patientInfoID int, category string, activityInfoID interface{}, activityOther string, startTime time.Time, createdBy int) (int, error) {
	// ตรวจสอบประเภทกิจกรรม และเลือกคอลัมน์ที่เหมาะสม
	var activityColumn string
	var activityValue interface{}

	if category == "other" {
		activityColumn = "activity_other"
		activityValue = activityOther
	} else {
		switch category {
		case "technology":
			activityColumn = "activity_technology_info_id"
		case "social":
			activityColumn = "activity_social_info_id"
		case "health":
			activityColumn = "activity_health_info_id"
		case "economic":
			activityColumn = "activity_economic_info_id"
		case "environmental":
			activityColumn = "activity_environmental_info_id"
		default:
			return 0, fmt.Errorf("Invalid activity category: %s", category)
		}
		activityValue = activityInfoID
	}

	// ใช้ `activityColumn` ที่เลือกมาแทนที่ `activity_info_id`
	query := fmt.Sprintf(`
		INSERT INTO activity_record (patient_info_id, %s, start_time, create_date, update_date, create_by, update_by)
		VALUES (?, ?, ?, NOW(), NOW(), ?, ?)`, activityColumn)

	result, err := db.Exec(query, patientInfoID, activityValue, startTime, createdBy, createdBy)
	if err != nil {
		return 0, fmt.Errorf("Error inserting activity start time: %v", err)
	}

	// ดึง `activity_record_id` ที่เพิ่งสร้าง
	recordID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(" Error retrieving last insert id: %v", err)
	}

	log.Printf("Successfully inserted activity record. Record ID: %d", recordID)
	return int(recordID), nil
}

func SaveActivityCompletion(db *sql.DB, cardID string, startTime, endTime time.Time, period string, lineUserID string) error {
	// ดึง user_info_id จาก line_user_id
	userID, err := GetUserInfoByLINEID(db, lineUserID)
	if err != nil {
		log.Printf("ไม่พบ user_info_id สำหรับ LINE ID: %s", lineUserID)
		return fmt.Errorf("ไม่พบข้อมูลผู้ใช้ กรุณาลองใหม่")
	}

	// ดึง `activity_record_id` และ `patient_info_id`
	activityRecord, err := GetActivityRecordID(db, cardID)
	if err != nil {
		return fmt.Errorf("ไม่พบข้อมูล activity record สำหรับ CardID: %s", cardID)
	}

	log.Printf("ตรวจสอบค่าที่ได้จาก GetActivityRecordID: ActivityRecord_ID=%d, PatientInfo_ID=%d",
		activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID)

	//อัปเดต `end_time` และ `period` โดยใช้ `activity_record_id` และ `patient_info_id`
	query := `
        UPDATE activity_record
        SET end_time = ?, period = ?, update_by = ?
        WHERE activity_record_id = ? 
        AND patient_info_id = ?
        LIMIT 1
    `

	log.Printf("🛠️ Preparing to update: activity_record_id=%d, patient_info_id=%d, end_time=%v, period=%s, update_by=%d",
		activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID,
		endTime, period, userID.UserInfo_ID)

	result, err := db.Exec(query, endTime, period, userID.UserInfo_ID, activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID)
	if err != nil {
		log.Printf("Error updating activity completion data: %v", err)
		return fmt.Errorf("error updating activity completion: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("ไม่มีแถวถูกอัปเดต! ตรวจสอบ activity_record_id=%d, patient_info_id=%d", activityRecord.ActivityRecord_ID, activityRecord.PatientInfo.PatientInfo_ID)
		return fmt.Errorf("ไม่มีข้อมูลที่ถูกอัปเดต โปรดตรวจสอบว่ามี record นี้ในฐานข้อมูลหรือไม่")
	}

	log.Printf("Activity completion updated successfully: %d rows affected", rowsAffected)
	return nil
}

// อัปเดต end_time, update_date และ update_by ในฐานข้อมูล
func UpdateActivityEndTime(db *sql.DB, activityRecordID int, endTime time.Time, userInfoID int) error {
	query := `
		UPDATE activity_record
		SET end_time = ?, update_date = NOW(), update_by = ?
		WHERE activity_record_id = ?
	`

	_, err := db.Exec(query, endTime, userInfoID, activityRecordID)
	if err != nil {
		return fmt.Errorf("Error updating end_time: %v", err)
	}
	return nil
}

// อัปเดตเวลาสิ้นสุดของกิจกรรม
// func UpdateActivityRecordWithoutEndTime(db *sql.DB, activity *models.Activityrecord) error {
// 	// ตรวจสอบค่าข้อมูลพื้นฐานก่อนอัปเดต
// 	if activity.PatientInfo.PatientInfo_ID == 0 {
// 		return fmt.Errorf("❌ patient_info_id is invalid")
// 	}
// 	if activity.ActivityRecord_ID == 0 {
// 		log.Println("❌ Invalid ActivityRecord_ID")
// 		return fmt.Errorf("❌ activity record ID is invalid")
// 	}

// 	// ตรวจสอบว่า activity_record_id ตรงกับ patient_info_id หรือไม่ และยังไม่มี end_time
// 	checkQuery := `
// 		SELECT COUNT(*)
// 		FROM activity_record
// 		WHERE activity_record_id = ? AND patient_info_id = ?
// 	`
// 	var count int
// 	err := db.QueryRow(checkQuery, activity.ActivityRecord_ID, activity.PatientInfo.PatientInfo_ID).Scan(&count)
// 	if err != nil {
// 		log.Printf("⚠️ SQL Execution error (checking record match): %v", err)
// 		return fmt.Errorf("⚠️ error verifying activity record: %v", err)
// 	}
// 	if count == 0 {
// 		return fmt.Errorf("⚠️ activity_record_id does not match with patient_info_id")
// 	}

// 	// อัปเดตข้อมูล **โดยไม่แตะต้อง `end_time`**
// 	updateQuery := `
// 	        UPDATE activity_record
// 	        SET
// 	            employee_info_id = ?,
// 	            write_by = ?,
// 	            write_date = NOW()
// 	        WHERE activity_record_id = ?
// 			LIMIT 1;`

// 	log.Printf("🔄 Updating activity_record (without end_time) with:\n➡️ EmployeeInfo_ID: %d\n➡️ WriteBy: %d\n➡️ ActivityRecord_ID: %d",
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 	)

// 	result, err := db.Exec(updateQuery,
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 	)

// 	if err != nil {
// 		log.Printf("❌ SQL Execution error: %v", err)
// 		return fmt.Errorf("❌ error updating activity record: %v", err)
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	log.Printf("✅ Rows affected: %d", rowsAffected)

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("⚠️ no rows were updated - check your WHERE conditions")
// 	}

// 	return nil
// }

// func UpdateActivityEndTime(db *sql.DB, activity *models.Activityrecord) error {
// 	// ตรวจสอบข้อมูลพื้นฐาน
// 	if activity.PatientInfo.PatientInfo_ID == 0 {
// 		return fmt.Errorf("patient_info_id is invalid")
// 	}

// 	if activity.ActivityRecord_ID == 0 {
// 		log.Println("Invalid ActivityRecord_ID")
// 		return fmt.Errorf("activity record ID is invalid")
// 	}

// 	// ตรวจสอบว่า activity_record_id ตรงกับ patient_info_id หรือไม่
// 	checkQuery := `
// 		SELECT COUNT(*)
// 		FROM activity_record
// 		WHERE activity_record_id = ? AND patient_info_id = ? AND end_time IS NULL
// 	`
// 	var count int
// 	err := db.QueryRow(checkQuery, activity.ActivityRecord_ID, activity.PatientInfo.PatientInfo_ID).Scan(&count)
// 	if err != nil {
// 		log.Printf("SQL Execution error (checking record match): %v", err)
// 		return fmt.Errorf("error verifying activity record: %v", err)
// 	}
// 	if count == 0 {
// 		return fmt.Errorf("activity_record_id does not match with patient_info_id or record already has end_time")
// 	}

// 	// อัปเดตข้อมูลเมื่อการตรวจสอบผ่าน
// 	updateQuery := `
// 	        UPDATE activity_record
// 	        SET
// 	            end_time = ?,
// 	            employee_info_id = ?,
// 	            write_by = ?,
// 	            write_date = NOW()
// 	        WHERE activity_record_id = ? AND end_time IS NULL
// 			LIMIT 1;`

// 	log.Printf("Updating activity_record with: EndTime: %v, EmployeeInfo_ID: %d, WriteBy: %d, ActivityRecord_ID: %d",
// 		activity.EndTime,
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 	)

// 	result, err := db.Exec(updateQuery,
// 		activity.EndTime,
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 	)

// 	if err != nil {
// 		log.Printf("SQL Execution error: %v", err)
// 		return fmt.Errorf("error updating end time: %v", err)
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	log.Printf("Rows affected: %d", rowsAffected)

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no rows were updated - check your WHERE conditions")
// 	}

// 	return nil
// }

// บันทึกกิจกรรม
// func SaveActivity(db *sql.DB, activity string) error {
// 	if !validateActivity(activity) {
// 		return fmt.Errorf("กิจกรรม '%s' ไม่ตรงกับค่าที่อนุญาตในฐานข้อมูล", activity)
// 	}

// 	query := `INSERT INTO service_info (activity) VALUES (?)`
// 	_, err := db.Exec(query, activity)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม %s ได้: %v", activity, err)
// 	}
// 	return nil
// }

// ***************************************************************************************************************************
// อัปโหลดไฟล์ไปยัง MinIO
func UploadFileToMinIO(client *minio.Client, bucketName, objectName, filePath string) (string, error) {
	// อัปโหลดไฟล์ไปยัง MinIO
	_, err := client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// สร้าง Public URL สำหรับไฟล์
	fileURL := fmt.Sprintf("http://10.221.43.191:9000/%s/%s", bucketName, objectName)
	return fileURL, nil
}

// func DownloadFileFromMinIO(minioClient *minio.Client, bucketName, objectName, filePath string) error {
// 	//ตั้งค่า Context
// 	ctx := context.Background()

// 	//ดึงไฟล์จาก MinIO (แบบ Stream)
// 	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
// 	if err != nil {
// 		return fmt.Errorf("❌ Error getting object from MinIO: %v", err)
// 	}
// 	defer object.Close()

// 	//สร้างไฟล์ใหม่เพื่อบันทึกข้อมูล
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return fmt.Errorf("❌ Error creating file: %v", err)
// 	}
// 	defer file.Close()

// 	//คัดลอกข้อมูลจาก Stream ไปยังไฟล์
// 	_, err = io.Copy(file, object)
// 	if err != nil {
// 		return fmt.Errorf("❌ Error writing file: %v", err)
// 	}

// 	log.Printf("✅ Successfully downloaded file from MinIO: %s", filePath)
// 	return nil
// }

// อัปเดต URL ของรูปการทำกิจกรรมในฐานข้อมูล
func updateImagebeforeActivity(db *sql.DB, patientInfoID int, fileURL string) error {
	query := `
		UPDATE activity_record
		SET image_before = ?
		WHERE patient_info_id = ?
		ORDER BY create_date DESC
		LIMIT 1
	`
	_, err := db.Exec(query, fileURL, patientInfoID)
	if err != nil {
		return fmt.Errorf("error updating image before activity URL: %v", err)
	}
	return nil
}

// เลือกกิจกรรมเพื่อเข้าถึงฐานข้อมูลหลักฐานในminio ผ่านpatient_info_id
func GetActivityNameByPatientInfoID(db *sql.DB, patientInfoID int) (string, error) {
	var activity string
	query := `
    SELECT 
        COALESCE(
            at.activity_technology_info_id, 
            aso.activity_social_info_id, 
            ah.activity_health_info_id, 
            aec.activity_economic_info_id, 
            aenv.activity_environmental_info_id, 
            ar.activity_other
        ) AS activity
    FROM activity_record ar
    LEFT JOIN activity_technology_info at ON at.activity_technology_info_id = ar.activity_technology_info_id
    LEFT JOIN activity_social_info aso ON aso.activity_social_info_id = ar.activity_social_info_id
    LEFT JOIN activity_health_info ah ON ah.activity_health_info_id = ar.activity_health_info_id
    LEFT JOIN activity_economic_info aec ON aec.activity_economic_info_id = ar.activity_economic_info_id
    LEFT JOIN activity_environmental_info aenv ON aenv.activity_environmental_info_id = ar.activity_environmental_info_id
    WHERE ar.patient_info_id = ?
    ORDER BY ar.create_date DESC
    LIMIT 1
`
	err := db.QueryRow(query, patientInfoID).Scan(&activity)
	if err != nil {
		return "", fmt.Errorf("error fetching activity name: %v", err)
	}
	return activity, nil
}

// ข้อมูลpatient_info_idผ่านcard_id การเก็บรูปหลักฐาน
func GetPatientInfoIDByCardID(db *sql.DB, cardID string) (int, error) {
	var patientInfoID int
	err := db.QueryRow("SELECT patient_info_id FROM patient_info WHERE card_id = ?", cardID).Scan(&patientInfoID)
	if err != nil {
		log.Printf("Error fetching patientInfoID for cardID %s: %v", cardID, err)
		return 0, err
	}
	return patientInfoID, nil
}

// อัปเดต URL ของรูปการจับเวลาทำกิจกรรมในฐานข้อมูล
func updateImageafterActivity(db *sql.DB, patientInfoID int, fileURL string) error {
	query := `
		UPDATE activity_record
		SET image_after = ?
		WHERE patient_info_id = ?
		ORDER BY create_date DESC
		LIMIT 1
	`
	_, err := db.Exec(query, fileURL, patientInfoID)
	if err != nil {
		return fmt.Errorf("error updating image after activity URL: %v", err)
	}
	return nil
}

//**************************************************************************************************************************

// ข้อมูลพนักงานผ่านชื่อ
func GetEmployeeIDByName(db *sql.DB, employeeName string) (int, error) {
	var employeeID int
	query := `SELECT employee_info_id 
		FROM employee_info 
		WHERE name_surname = ?
		ORDER BY create_date DESC
		LIMIT 1`
	err := db.QueryRow(query, employeeName).Scan(&employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("ไม่พบพนักงานที่ชื่อ %s\nกรุณากรอกชื่อพนักงานอีกครั้ง", employeeName)
		}
		log.Printf("Error fetching employee info: %v", err)
		return 0, fmt.Errorf("เกิดข้อผิดพลาดในการค้นหาข้อมูลพนักงาน")
	}
	return employeeID, nil
}

// func GetEmployeeIDByName(db *sql.DB, employeeName string) (int, error) {
//     var employeeID int
//     query := `SELECT employee_info_id FROM employee_info WHERE LOWER(TRIM(name_surname)) = LOWER(TRIM(?))`
//     err := db.QueryRow(query, employeeName).Scan(&employeeID)
//     if err != nil {
//         if err == sql.ErrNoRows {
//             log.Printf("No employee found with name: %s", employeeName)
//             return 0, fmt.Errorf("ไม่พบพนักงานชื่อ: %s", employeeName)
//         }
//         log.Printf("Error querying employee_info_id: %v", err)
//         return 0, err
//     }
//     log.Printf("Found employee_info_id: %d for name: %s", employeeID, employeeName)
//     return employeeID, nil
// }

func UpdateActivityEmployeeID(db *sql.DB, activityRecordID int, employeeID int, userInfoID int) error {
	query := `UPDATE activity_record SET employee_info_id = ? WHERE activity_record_id = ?`
	_, err := db.Exec(query, employeeID, userInfoID, activityRecordID)
	return err
}
func UpdateActivityEmployee(db *sql.DB, activity *models.Activityrecord) error {
	// ตรวจสอบข้อมูลพื้นฐาน
	if activity.PatientInfo.PatientInfo_ID == 0 {
		return fmt.Errorf("patient_info_id is invalid")
	}

	if activity.EmployeeInfo.EmployeeInfo_ID == 0 {
		return fmt.Errorf("employee_info_id is invalid")
	}

	// ค้นหา activity_record_id ล่าสุดของ patient_info_id ที่ยังไม่มี employee_info_id และยังไม่สิ้นสุด (end_time IS NULL)
	var activityRecordID int
	findQuery := `
		SELECT activity_record_id 
		FROM activity_record 
		WHERE patient_info_id = ? 
			AND employee_info_id IS NULL 
		ORDER BY create_date DESC 
		LIMIT 1;
	`

	err := db.QueryRow(findQuery, activity.PatientInfo.PatientInfo_ID).Scan(&activityRecordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("⚠️ ไม่พบกิจกรรมที่ต้องอัปเดต หรืออาจถูกอัปเดตไปแล้ว")
		}
		log.Printf("❌ SQL Execution error (finding latest activity record): %v", err)
		return fmt.Errorf("❌ error fetching latest activity record: %v", err)
	}

	log.Printf("พบ ActivityRecord_ID ล่าสุดที่ต้องอัปเดต: %d", activityRecordID)

	// อัปเดต employee_info_id ลง activity_record ล่าสุด
	updateQuery := `
		UPDATE activity_record 
		SET 
			employee_info_id = ?, 
			update_by = ?, 
			update_date = NOW()
		WHERE activity_record_id = ?;
	`

	result, err := db.Exec(updateQuery,
		activity.EmployeeInfo.EmployeeInfo_ID, // ใส่รหัสพนักงาน
		activity.UserInfo.UserInfo_ID,         // ใส่รหัสผู้ที่อัปเดต
		activityRecordID,                      // ใส่รหัสกิจกรรมที่ต้องการอัปเดต
	)

	if err != nil {
		log.Printf("❌ SQL Execution error (updating activity record): %v", err)
		return fmt.Errorf("❌ error updating employee_info_id: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		return fmt.Errorf("ไม่มีแถวที่ถูกอัปเดต ตรวจสอบเงื่อนไข WHERE")
	}

	log.Println("อัปเดต employee_info_id สำเร็จ!")
	return nil
}

// //บันทึกชื่อพนักงานที่บริการ
// func GetActivityRecordByEmployeeID(db *sql.DB, patientinfoID int) (*models.Activityrecord, error) {
// 	query := `SELECT
// 		a.activity_record_id,
// 		a.start_time,
// 		a.end_time,
// 		a.period,
// 		a.write_by,

// 		p.patient_info_id,
// 		p.card_id,
// 		p.username,

// 		s.service_info_id,
// 		s.activity,

// 		e.employee_info_id

// 	FROM activity_record a
// 	INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
// 	INNER JOIN service_info s ON a.service_info_id = s.service_info_id
// 	INNER JOIN employee_info e ON a.employee_info_id = e.employee_info_id
// 	WHERE e.employee_info_id = ?
// 	ORDER BY a.start_time DESC
// 	LIMIT 1;`

// 	row := db.QueryRow(query, patientinfoID)

// 	var record models.Activityrecord
// 	var patientInfo models.PatientInfo
// 	var serviceInfo models.ServiceInfo
// 	var employeeInfo models.EmployeeInfo

// 	err := row.Scan(
// 		&record.ActivityRecord_ID,
// 		&record.StartTime,
// 		&record.EndTime,
// 		&record.Period,
// 		&record.Write_by,

// 		&patientInfo.PatientInfo_ID,
// 		&patientInfo.CardID,
// 		&patientInfo.Name,

// 		&serviceInfo.ServiceInfo_ID,
// 		&serviceInfo.Activity,

// 		&employeeInfo.EmployeeInfo_ID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// กรณีไม่พบข้อมูล
// 			return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรมที่ยังไม่เสร็จสำหรับ employeeID: %s", patientinfoID)
// 		}
// 		// กรณีเกิดข้อผิดพลาด
// 		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรม: %v", err)
// 	}

// 	// Assign ข้อมูลให้กับ Activityrecord
// 	record.PatientInfo = patientInfo
// 	record.ServiceInfo = serviceInfo
// 	record.EmployeeInfo = employeeInfo

// 	return &record, nil
// }

// //บันทึกการอัปเดตข้อมูลการทำกิจกรรมผ่านid ผู้สูงอายุ
// func UpdateActivityEndTimeForPatient(db *sql.DB, activity *models.Activityrecord) error {
// 	if activity.PatientInfo.PatientInfo_ID == 0 {
// 		return fmt.Errorf("invalid patient_info_id")
// 	}

// 	if activity.ActivityRecord_ID == 0 {
// 		log.Println("Invalid ActivityRecord_ID")
// 		return fmt.Errorf("invalid activity record ID")
// 	}

// 	query := `
//         UPDATE activity_record
//         SET
//             end_time = ?,
//             employee_info_id = ?,
//             write_by = ?,
//             write_date = NOW()
//         WHERE activity_record_id = ?
//           AND patient_info_id = ?
//           AND end_time IS NULL
//     `
// 	log.Printf("Updating activity_record with: EndTime=%v, EmployeeInfo_ID=%d, WriteBy=%d, ActivityRecord_ID=%d, PatientInfo_ID=%d",
// 		activity.EndTime,
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 		activity.PatientInfo.PatientInfo_ID,
// 	)

// 	result, err := db.Exec(query,
// 		activity.EndTime,
// 		activity.EmployeeInfo.EmployeeInfo_ID,
// 		activity.UserInfo.UserInfo_ID,
// 		activity.ActivityRecord_ID,
// 		activity.PatientInfo.PatientInfo_ID,
// 	)

// 	if err != nil {
// 		log.Printf("SQL Execution error: %v", err)
// 		return fmt.Errorf("error updating activity record: %v", err)
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	log.Printf("Rows affected: %d", rowsAffected)

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no rows were updated - check your WHERE conditions")
// 	}

// 	return nil
// }

// ข้อมูลimamge ของ patient_info
// func GetImageFromDatabase(db *sql.DB, cardID string) ([]byte, error) {
// 	var imageData []byte
// 	query := "SELECT image FROM patient_info WHERE card_id = ?"
// 	err := db.QueryRow(query, cardID).Scan(&imageData)
// 	if err != nil {
// 		return nil, fmt.Errorf("ไม่พบข้อมูลรูปภาพสำหรับ CardID: %s, Error: %v", cardID, err)
// 	}

// 	return imageData, nil
// }
