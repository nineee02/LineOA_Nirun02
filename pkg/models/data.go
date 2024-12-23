package models

import "time"

type PatientInfo struct {
	PatientInfo_ID       int                  `json:"patient_info_id"`
	CardID               string               `json:"card_id"`
	Name                 string               `json:"username"`
	PhoneNumber          string               `json:"phone_number"`
	Email                string               `json:"email"`
	Address              string               `json:"address"`
	DateOfBirth          string               `json:"date_of_birth"`
	Age                  string               `json:"age"`
	Sex                  string               `json:"sex"`
	Blood                string               `json:"blood_"`
	ADL                  string               `json:"ADL"`
	CreateDate           string               `json:"create_date"`
	WriteDate            string               `json:"write_date"`
	Religion             Religion             `json:"religion_info"`
	CountryInfo          CountryInfo          `json:"country_info"`
	RightToTreatmentInfo RightToTreatmentInfo `json:"reght_to_treatment_info"`
}

type CountryInfo struct {
	CountryInfo_ID int    `json:"country_info_id"`
	Country        string `json:"country"`
	CreateDate     string `json:"create_date"`
	WriteDate      string `json:"write_date"`
}

type Religion struct {
	ReligionInfo_ID int    `json:"religion_info_id"`
	Religion        string `json:"religion"`
	CreateDate      string `json:"create_date"`
	WriteDate       string `json:"write_date"`
}

type RightToTreatmentInfo struct {
	RightToTreatmentInfo_ID int    `json:"right_to_treatment_info_id"`
	Right_to_treatment      string `json:"right_to_treatment"`
	CreateDate              string `json:"create_date"`
	WriteDate               string `json:"write_date"`
}

type ServiceInfo struct {
	ServiceInfo_Id int    `json:"service_info_id"`
	Activity       string `json:"activity"`
	ServiceType    string `json:"service_type"`
	CreateDate     string `json:"create_date"`
	WriteDate      string `json:"write_date"`
}

type Activityrecord struct {
	ActivityRecord_ID int          `json:"activity_record_id"`
	StartTime         string       `json:"start_time"`
	EndTime           time.Time    `json:"end_time"`
	Period            string       `json:"period"`
	Evidence_activity []byte       `json:"evidence_activity"`
	Evidence_time     []byte       `json:"evidence_time"`
	Location          string       `json:"location"`
	CreateDate        string       `json:"create_date"`
	WriteDate         string       `json:"write_date"`
	ServiceInfo_Id    int          `json:"service_info_id"`
	PatientInfo_Id    int          `json:"patient_info_id"`
	PatientInfo       PatientInfo  `json:"patient_info"`
	ServiceInfo       ServiceInfo  `json:"service_info"`
	EmployeeInfo      EmployeeInfo `json:"employee_info"`
}

type EmployeeInfo struct {
	EmployeeInfo_ID int             `json:"employee_info_id"`
	EmployeeCode    string          `json:"employee_code"`
	Name            string          `json:"username"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	CreateDate      string          `json:"create_date"`
	WriteDate       string          `json:"write_date"`
	DepartmentInfo  DepartmentInfo  `json:"department_inf"`
	JobPositionInfo JobPositionInfo `json:"job_position_info"`
}
type DepartmentInfo struct {
	DepartmentInfo_id int    `json:"department_info_id"`
	Department        string `json:"department"`
	CreateDate        string `json:"create_date"`
	WriteDate         string `json:"write_date"`
}
type JobPositionInfo struct {
	JobPositionInfo_id int    `json:"job_position_info_id"`
	JobPosition        string `json:"job_position"`
	CreateDate         string `json:"create_date"`
	WriteDate          string `json:"write_date"`
}

type WorktimeRecord struct {
	WorktimeRecord_ID int          `json:"worktime_record_id"`
	CheckIn           string       `json:"check_in"`
	CheckOut          string       `json:"check_out"`
	Period            string       `json:"period"`
	CreateDate        string       `json:"create_date"`
	WriteDate         string       `json:"write_date"`
	EmployeeInfo_ID   int          `json:"employee_info_id"`
	EmployeeInfo      EmployeeInfo `json:"employee_info"`
}
