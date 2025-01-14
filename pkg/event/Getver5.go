package event

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/models"
	"strconv"
	"time"
)
// func GetEmployeeByLINEID(db *sql.DB, lineUserID string) (*models.EmployeeInfo, error) {
//     query := `
//         SELECT employee_info_id, employee_code, username, phone_number, email 
//         FROM employee_info 
//         WHERE line_user_id = ?`
//     row := db.QueryRow(query, lineUserID)
// 	log.Printf("Fetching employee data for LINE User ID: %s", lineUserID)
//     var employeeInfo models.EmployeeInfo
//     err := row.Scan(
//         &employeeInfo.EmployeeInfo_ID,
//         &employeeInfo.EmployeeCode,
//         &employeeInfo.Name,
//         &employeeInfo.PhoneNumber,
//         &employeeInfo.Email,
//     )
//     if err != nil {
//         return nil, err
//     }
//     return &employeeInfo, nil
// }
func GetEmployeeByLINEID(db *sql.DB, lineUserID string) (*models.EmployeeInfo, error) {
	query := `SELECT employee_info_id, employee_code, username, phone_number, email FROM employee_info WHERE line_user_id = ?`
	row := db.QueryRow(query, lineUserID)

	employee := &models.EmployeeInfo{}
	err := row.Scan(&employee.EmployeeInfo_ID, &employee.EmployeeCode, &employee.Name, &employee.PhoneNumber, &employee.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลสำหรับ LINE User ID: %s", lineUserID)
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลพนักงาน: %v", err)
	}

	return employee, nil
}


func GetEmployeeInfo(db *sql.DB, employeeCode string) (*models.EmployeeInfo, error) {
	query := `
		SELECT employee_info_id, employee_code, username, phone_number, email, create_date, write_date
		FROM employee_info
		WHERE employee_code = ?`

	row := db.QueryRow(query, employeeCode)

	var employeeInfo models.EmployeeInfo
	err := row.Scan(
		&employeeInfo.EmployeeInfo_ID,
		&employeeInfo.EmployeeCode,
		&employeeInfo.Name,
		&employeeInfo.PhoneNumber,
		&employeeInfo.Email,
		&employeeInfo.CreateDate,
		&employeeInfo.WriteDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ไม่พบข้อมูลพนักงานสำหรับรหัสพนักงาน: %s", employeeCode)
		}
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลพนักงาน: %v", err)
	}

	return &employeeInfo, nil
}

func GetEmployeeID(db *sql.DB, employeeCode string) (int, error) {
	var employeeID int
	query := "SELECT employee_info_id FROM employee_info WHERE employee_code = ?"
	err := db.QueryRow(query, employeeCode).Scan(&employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("ไม่พบข้อมูลสำหรับรหัสพนักงาน: %s", employeeCode)
		}
		return 0, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูล employee_info_id: %v", err)
	}
	return employeeID, nil
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

// ตรวจสอบว่าพนักงานคนนี้มีการเช็คอินอยู่แล้วหรือยัง
func IsEmployeeCheckedIn(db *sql.DB, employeeID int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM worktime_record WHERE employee_info_id = ? AND check_out IS NULL"
	err := db.QueryRow(query, employeeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("เกิดข้อผิดพลาดในการตรวจสอบสถานะ: %v", err)
	}
	return count > 0, nil
}

// บันทึก Check-in
func RecordCheckIn(db *sql.DB, employeeID int) error {
	query := "INSERT INTO worktime_record (employee_info_id, check_in) VALUES (?, ?)"
	_, err := db.Exec(query, employeeID, time.Now())
	if err != nil {
		return fmt.Errorf("เกิดข้อผิดพลาดในการบันทึก Check-in: %v", err)
	}
	return nil
}

// บันทึก Check-out
func RecordCheckOut(db *sql.DB, employeeID int) error {
	query := "UPDATE worktime_record SET check_out = ? WHERE employee_info_id = ? AND check_out IS NULL"
	_, err := db.Exec(query, time.Now(), employeeID)
	if err != nil {
		return fmt.Errorf("เกิดข้อผิดพลาดในการบันทึก Check-out: %v", err)
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

func SaveActivityRecord(db *sql.DB, activity *models.Activityrecord) error {
	// ดึง patient_info_id
	patient, err := GetPatientInfoByName(db, activity.PatientInfo.CardID)
	if err != nil {
		return fmt.Errorf("error fetching patient_info_id: %v", err)
	}
	activity.PatientInfo.PatientInfo_ID = patient.PatientInfo.PatientInfo_ID

	// ตรวจสอบว่าค่า patient_info_id ถูกต้อง
	if activity.PatientInfo.PatientInfo_ID == 0 {
		return fmt.Errorf("patient_info_id is missing or invalid")
	}

	// ดึง service_info_id
	serviceInfoID, err := GetServiceInfoIDByActivity(db, activity.ServiceInfo.Activity)
	if err != nil {
		return fmt.Errorf("error fetching service_info_id: %v", err)
	}

	// ตรวจสอบว่า service_info_id มีอยู่จริง
	if serviceInfoID == 0 {
		return fmt.Errorf("service_info_id is invalid or missing")
	}

	// บันทึกลง activity_record
	query := "INSERT INTO activity_record (patient_info_id, service_info_id, employee_info_id, start_time) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, activity.PatientInfo.PatientInfo_ID, serviceInfoID, activity.EmployeeInfo.EmployeeInfo_ID, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting activity record: %v", err)
	}

	return nil
}
func GetActivityStartTime(db *sql.DB, cardID string, activity string) (time.Time, error) {
	// คิวรีเพื่อดึง start_time จาก activity_record โดยใช้ cardID และ activity
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

func UpdateActivityEndTime(db *sql.DB, activity *models.Activityrecord) error {
	if activity.PatientInfo.PatientInfo_ID == 0 {
		return fmt.Errorf("patient_info_id is invalid")
	}

	query := `
    UPDATE activity_record 
    SET end_time = ? 
    WHERE patient_info_id = ? AND end_time IS NULL`

	_, err := db.Exec(query, activity.EndTime, activity.PatientInfo.PatientInfo_ID)
	if err != nil {
		return fmt.Errorf("error updating end time: %v", err)
	}

	return nil
}
func historyALL(db *sql.DB) ([]*models.Activityrecord, error) {
	query := `SELECT 
			YEAR(ar.create_date) AS year,
			s.activity AS activity_type,
			COUNT(*) AS total
		FROM 
			activity_record ar
		JOIN 
			service_info s ON ar.service_info_id = s.service_info_id
		GROUP BY 
			YEAR(ar.create_date), s.activity
		ORDER BY 
			year DESC, activity_type;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{}
		detail := &models.ActivityYearDetail{}

		if err := rows.Scan(&detail.Year, &detail.ActivityType, &detail.Total); err != nil {
			return nil, err
		}
		record.ActivityYearDetail = *detail
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func historyOfYear(db *sql.DB) ([]*models.Activityrecord, error) {
	query := `
			SELECT 
				YEAR(ar.create_date) AS year,
				s.activity AS activity_type,
				COUNT(*) AS total
			FROM 
				activity_record ar
			JOIN 
				service_info s ON ar.service_info_id = s.service_info_id
			WHERE 
				YEAR(ar.create_date) = YEAR(CURDATE()) -- เงื่อนไขสำหรับปีปัจจุบัน
			GROUP BY 
				YEAR(ar.create_date), s.activity
			ORDER BY 
				year DESC, activity_type;
		`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{}
		detail := &models.ActivityYearDetail{}

		if err := rows.Scan(&detail.Year, &detail.ActivityType, &detail.Total); err != nil {
			return nil, err
		}

		record.ActivityYearDetail = *detail
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func historyOfMonth(db *sql.DB) ([]*models.Activityrecord, error) {
	query := `SELECT 
		MONTH(ar.create_date) AS month,
		s.activity AS activity_type,
		COUNT(*) AS total
	FROM 
		activity_record ar
	JOIN 
		service_info s ON ar.service_info_id = s.service_info_id
	WHERE 
		YEAR(ar.create_date) = YEAR(CURDATE()) -- ตรวจสอบว่าปีปัจจุบัน
		AND MONTH(ar.create_date) = MONTH(CURDATE()) -- ตรวจสอบว่าเป็นเดือนปัจจุบัน
	GROUP BY 
		MONTH(ar.create_date), s.activity
	ORDER BY 
		month DESC, activity_type;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{}
		detail := &models.ActivityMonthDetail{}

		if err := rows.Scan(&detail.Month, &detail.ActivityType, &detail.Total); err != nil {
			return nil, err
		}

		record.ActivityMonthDetail = *detail
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func historyOfWeek(db *sql.DB) ([]*models.Activityrecord, error) {
	query := `SELECT 
                s.activity AS activity_type,
                COUNT(*) AS total
          FROM 
                activity_record ar
          JOIN 
                service_info s ON ar.service_info_id = s.service_info_id
          WHERE 
                ar.create_date BETWEEN 
                DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY) 
                AND 
                DATE_ADD(CURDATE(), INTERVAL (6 - WEEKDAY(CURDATE())) DAY)
          GROUP BY 
                s.activity
          ORDER BY 
                total DESC;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{}
		detail := &models.ActivityWeekDetail{}

		if err := rows.Scan(&detail.ActivityType, &detail.Total); err != nil {
			return nil, err
		}

		record.ActivityWeekDetail = *detail
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func historyOfDay(db *sql.DB) ([]*models.Activityrecord, error) {
	query := `SELECT 
		DAY(ar.create_date) AS day,
		s.activity AS activity_type,
		COUNT(*) AS total
	FROM 
		activity_record ar
	JOIN 
		service_info s ON ar.service_info_id = s.service_info_id
	WHERE 
		DAY(ar.create_date) = DAY(CURDATE())
	GROUP BY 
		DAY(ar.create_date), s.activity
	ORDER BY 
		day DESC, activity_type;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{}
		detail := &models.ActivityDayDetail{}

		if err := rows.Scan(&detail.Day, &detail.ActivityType, &detail.Total); err != nil {
			return nil, err
		}

		record.ActivityDayDetail = *detail
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
func historyOfSet(db *sql.DB, startDate, endDate string) ([]*models.Activityrecord, error) {
	query := `SELECT 
		DATE(ar.create_date) AS date, 
		s.activity AS activity_type, 
		COUNT(*) AS total
	FROM 
		activity_record ar
	JOIN 
		service_info s ON ar.service_info_id = s.service_info_id
	WHERE 
		ar.create_date BETWEEN ? AND ?
	GROUP BY 
		DATE(ar.create_date), s.activity
	ORDER BY 
		date DESC, activity_type;`

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Activityrecord
	for rows.Next() {
		record := &models.Activityrecord{
			ServiceInfo: models.ServiceInfo{},
		}
		detail := &models.ActivitySetDetail{}

		var date string
		if err := rows.Scan(&date, &record.ServiceInfo.Activity, &detail.Total); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		dayInt, err := strconv.Atoi(date[8:10]) // Extract day from date string
		if err != nil {
			log.Printf("Error converting date to int: %v", err)
			return nil, err
		}
		detail.Date = dayInt
		detail.ActivityType = record.ServiceInfo.Activity

		record.ActivitySetDetail = *detail
		results = append(results, record)

		// Debug log
		log.Printf("Fetched record: Date=%d, ActivityType=%s, Total=%d", detail.Date, detail.ActivityType, detail.Total)
	}

	// Ensure a proper return at the end
	return results, nil
}

