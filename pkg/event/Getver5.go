package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/models"
	"time"
)

func GetUserInfoByLINEID(db *sql.DB, lineUserID string) (*models.User_info, error) {
	query := `SELECT user_info_id, line_user_id, sex, name, email, phone_number, create_date, write_date 
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
		&user.WriteDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ LINE User ID: %s", lineUserID)
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้: %v", err)
	}

	return user, nil
}

func GetWorktime(db *sql.DB, employeeCode string) (*models.WorktimeRecord, error) {
	query := `
		SELECT wr.worktime_record_id, wr.check_in, wr.check_out, 
		       e.employee_code, e.username,
		       d.department, j.job_position
		FROM worktime_record wr
		LEFT JOIN employee_info e ON wr.employee_info_id = e.employee_info_id
		LEFT JOIN department_info d ON e.department_info_id = d.department_info_id
		LEFT JOIN job_position_info j ON e.job_position_info_id = j.job_position_info_id
		WHERE e.employee_code = ?
		ORDER BY wr.check_in DESC
		LIMIT 1`

	row := db.QueryRow(query, employeeCode)

	var record models.WorktimeRecord
	var checkOut sql.NullTime
	var department, jobPosition sql.NullString

	err := row.Scan(
		&record.WorktimeRecord_ID,
		&record.CheckIn,
		&checkOut,
		&record.EmployeeInfo.EmployeeCode,
		&record.EmployeeInfo.Name,
		&department,
		&jobPosition,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // คืนค่า nil หากไม่มีข้อมูล
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลการทำงาน: %v", err)
	}

	// จัดการค่า NULL
	if checkOut.Valid {
		record.CheckOut = checkOut.Time // กรณีที่มีค่า
	} else {
		record.CheckOut = time.Time{} // กรณีไม่มีค่า
	}

	record.EmployeeInfo.DepartmentInfo.Department = department.String
	record.EmployeeInfo.JobPositionInfo.JobPosition = jobPosition.String

	return &record, nil
}

// ตรวจสอบการเช็คอิน
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

// บันทึก Check-in
func RecordCheckIn(db *sql.DB, userID int) error {
	// ดึง employee_info_id ที่สัมพันธ์กับ user_info_id
	var employeeID sql.NullInt64 // ใช้ sql.NullInt64 เพื่อจัดการค่า NULL
	query := `SELECT employee_info_id FROM user_info WHERE user_info_id = ?`
	err := db.QueryRow(query, userID).Scan(&employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("ไม่พบ employee_info_id สำหรับ user_info_id: %d", userID)
		}
		return fmt.Errorf("error fetching employee_info_id: %v", err)
	}

	// เตรียม Query สำหรับ INSERT
	var insertQuery string
	var args []interface{}

	if employeeID.Valid { // ตรวจสอบว่ามีค่า employee_info_id หรือไม่
		insertQuery = `
			INSERT INTO worktime_record (
				check_in, 
				create_by, 
				write_by, 
				user_info_id, 
				employee_info_id, 
				create_date, 
				write_date
			) 
			VALUES (NOW(), ?, ?, ?, ?, NOW(), NOW())`
		args = []interface{}{userID, userID, userID, employeeID.Int64}
	} else {
		insertQuery = `
			INSERT INTO worktime_record (
				check_in, 
				create_by, 
				write_by, 
				user_info_id, 
				create_date, 
				write_date
			) 
			VALUES (NOW(), ?, ?, ?, NOW(), NOW())`
		args = []interface{}{userID, userID, userID}
	}

	// บันทึกข้อมูลการเช็คอิน
	_, err = db.Exec(insertQuery, args...)
	if err != nil {
		return fmt.Errorf("error recording check-in: %v", err)
	}
	return nil
}

func RecordCheckOut(db *sql.DB, userID int) error {
	query := `
		UPDATE worktime_record 
		SET 
			check_out = NOW(), 
			write_date = NOW(), 
			write_by = ?
		WHERE user_info_id = ? AND check_out IS NULL`
	_, err := db.Exec(query, userID, userID)
	if err != nil {
		return fmt.Errorf("error recording check-out: %v", err)
	}
	return nil
}

// GetPatientInfoByName ค้นหาข้อมูลผู้ป่วยจากชื่อ
func GetPatientInfoByName(db *sql.DB, cardID string) (*models.Activityrecord, error) {
	query := `SELECT 
				p.card_id,
				p.patient_info_id,
				p.username, 
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
				c.write_date,

				r.religion_info_id, 
				r.religion,
				r.create_date,
				r.write_date,
				 
				rtt.right_to_treatment_info_id, 
				rtt.right_to_treatment, 
				rtt.create_date,
				rtt.write_date
			FROM patient_info p
			LEFT JOIN country_info c ON p.country_info_id = c.country_info_id
			LEFT JOIN religion_info r ON p.religion_info_id = r.religion_info_id
			LEFT JOIN right_to_treatment_info rtt ON p.right_to_treatment_info_id = rtt.right_to_treatment_info_id 
			WHERE p.card_id LIKE ?`

	patient := &models.Activityrecord{}
	err := db.QueryRow(query, "%"+cardID+"%").Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลผู้สูงอายุ")
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาด: %v", err)
	}

	return patient, nil
}

func GetServiceInfoBycardID(db *sql.DB, cardID string) ([]models.Activityrecord, error) {
	query := `SELECT a.avtivity_info_id,
					a.start_time,
					a.end_time,
					a.period,
					p.card_id, 
					p.username, 
					p.patient_info_id,
					s.service_info_id,
					s.activity
					e.employee_info_id
					
			  FROM activity_record a
			  INNER JOIN patient_info p ON a.patient_info_id = p.patient_info_id
			  INNER JOIN service_info s ON a.service_info_id = s.service_info_id
			  INNER JOIN employee_info e ON a.employee_info_id = e.employee_info_id
			  WHERE p.card_id =?`

	rows, err := db.Query(query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activityrecord []models.Activityrecord

	for rows.Next() {
		var record models.Activityrecord
		var patientInfo models.PatientInfo
		var serviceInfo models.ServiceInfo
		var employeeInfo models.EmployeeInfo

		err := rows.Scan(
			&record.ActivityRecord_ID,
			&record.StartTime,
			&record.EndTime,
			&record.Period,
			&patientInfo.CardID,
			&patientInfo.Name,
			&patientInfo.PatientInfo_ID,
			&serviceInfo.ServiceInfo_ID,
			&serviceInfo.Activity,
			&employeeInfo.EmployeeInfo_ID)
		if err != nil {
			return nil, err
		}
		// Assign ข้อมูลให้กับ Activityrecord
		record.PatientInfo = patientInfo
		record.ServiceInfo = serviceInfo
		record.EmployeeInfo = employeeInfo

		activityrecord = append(activityrecord, record)
	}

	if len(activityrecord) == 0 {
		return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรม: %s", cardID)
	}

	return activityrecord, nil
}

func GetServiceInfoIDByActivity(db *sql.DB, activity string) (int, error) {
	query := "SELECT service_info_id FROM service_info WHERE activity = ?"
	var serviceInfoID int
	err := db.QueryRow(query, activity).Scan(&serviceInfoID)
	if err != nil {
		return 0, fmt.Errorf("ไม่พบ service_info_id สำหรับกิจกรรม: %s, error: %v", activity, err)
	}
	return serviceInfoID, nil
}
func GetEmployeeIDByName(db *sql.DB, employeeName string) (int, error) {
    var employeeID int
    query := `SELECT employee_info_id FROM employee_info WHERE username = ?`
    err := db.QueryRow(query, employeeName).Scan(&employeeID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("ไม่พบพนักงานที่ชื่อ %s", employeeName)
        }
        log.Printf("Error fetching employee info: %v", err)
        return 0, fmt.Errorf("เกิดข้อผิดพลาดในการค้นหาข้อมูลพนักงาน")
    }
    return employeeID, nil
}

func SaveActivityRecord(db *sql.DB, activity *models.Activityrecord) error {
	// ดึง patient_info_id
	patient, err := GetPatientInfoByName(db, activity.PatientInfo.CardID)
	if err != nil {
		log.Printf("Error fetching patient_info_id: %v", err)
		return fmt.Errorf("error fetching patient_info_id: %v", err)
	}
	activity.PatientInfo.PatientInfo_ID = patient.PatientInfo.PatientInfo_ID

	// ตรวจสอบว่าค่า patient_info_id ถูกต้อง
	if activity.PatientInfo.PatientInfo_ID == 0 {
		log.Println("Invalid patient_info_id")
		return fmt.Errorf("patient_info_id is missing or invalid")
	}

	// ดึง service_info_id
	serviceInfoID, err := GetServiceInfoIDByActivity(db, activity.ServiceInfo.Activity)
	if err != nil {
		log.Printf("Error fetching service_info_id: %v", err)
		return fmt.Errorf("error fetching service_info_id: %v", err)
	}

	// ตรวจสอบว่า service_info_id มีอยู่จริง
	if serviceInfoID == 0 {
		log.Println("Invalid service_info_id")
		return fmt.Errorf("service_info_id is invalid or missing")
	}

	query := `
		INSERT INTO activity_record (
			patient_info_id, 
			service_info_id, 
			start_time, 
			create_by, 
			write_by
		) 
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query,
		activity.PatientInfo.PatientInfo_ID,
		serviceInfoID,
		time.Now(),
		activity.UserInfo.UserInfo_ID, // ใช้ UserInfo_ID สำหรับ create_by
		activity.UserInfo.UserInfo_ID, // ใช้ UserInfo_ID สำหรับ write_by
	)
	if err != nil {
		log.Printf("Error inserting activity record: %v", err)
		return fmt.Errorf("error inserting activity record: %v", err)
	}

	// ดึง activity_record_id ที่ถูกเพิ่มขึ้นมา
	activityRecordID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert id: %v", err)
		return fmt.Errorf("error retrieving last insert id: %v", err)
	}

	// บันทึก activityRecordID ลงใน activity
	activity.ActivityRecord_ID = int(activityRecordID)

	log.Println("Activity record saved successfully")
	return nil

}
func GetActivityRecordID(db *sql.DB, cardID string) (int, error) {
    var activityRecordID int
    query := `
        SELECT activity_record_id 
        FROM activity_record 
        WHERE patient_info_id = (SELECT patient_info_id FROM patient_info WHERE card_id = ?) 
          AND end_time IS NULL
    `
    err := db.QueryRow(query, cardID).Scan(&activityRecordID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("ไม่มีกิจกรรมที่ยังไม่เสร็จสิ้นสำหรับผู้ใช้นี้")
        }
        log.Printf("Error fetching activityRecord_ID: %v", err)
        return 0, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรม")
    }
    return activityRecordID, nil
}

func UpdateActivityEndTime(db *sql.DB, activity *models.Activityrecord) error {
    if activity.PatientInfo.PatientInfo_ID == 0 {
        return fmt.Errorf("patient_info_id is invalid")
    }

    if activity.ActivityRecord_ID == 0 {
        log.Println("Invalid ActivityRecord_ID")
        return fmt.Errorf("activity record ID is invalid")
    }

    query := `
        UPDATE activity_record 
        SET 
            end_time = ?, 
            employee_info_id = ?, 
            write_by = ?, 
            write_date = NOW()
        WHERE activity_record_id = ? AND end_time IS NULL
    `

    log.Printf("Updating activity_record with: EndTime: %v, EmployeeInfo_ID: %d, WriteBy: %d, ActivityRecord_ID: %d",
        activity.EndTime,
        activity.EmployeeInfo.EmployeeInfo_ID,
        activity.UserInfo.UserInfo_ID,
        activity.ActivityRecord_ID,
    )

    result, err := db.Exec(query,
        activity.EndTime,
        activity.EmployeeInfo.EmployeeInfo_ID,
        activity.UserInfo.UserInfo_ID,
        activity.ActivityRecord_ID,
    )

    if err != nil {
        log.Printf("SQL Execution error: %v", err)
        return fmt.Errorf("error updating end time: %v", err)
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Rows affected: %d", rowsAffected)

    if rowsAffected == 0 {
        return fmt.Errorf("no rows were updated - check your WHERE conditions")
    }

    return nil
}


