package models

import (
	"database/sql"
	"fmt"
	"log"
)

type PatientInfo struct {
	ID          int    `json:"id"`
	PatientID   string `json:"patiet_id"`
	Image       []byte `json:"image"`
	Name        string `json:"name_"`
	PhoneNumber string `json:"phone_numbers"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	CardID      string `json:"card_id"`
	Religion    string `json:"religion"`
	Sex         string `json:"sex"`
	Blood       string `json:"blood"`
	DateOfBirth string `json:"date_of_birth"`
	Age         int    `json:"age"`
}

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
