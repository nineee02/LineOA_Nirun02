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
	Blood                string               `json:"blood"`
	ADL                  string               `json:"ADL"`
	CreateDate           string               `json:"create_date"`
	WriteDate            string               `json:"write_date"`
	Religion             Religion             `json:"religion_info"`
	CountryInfo          CountryInfo          `json:"country_info"`
	RightToTreatmentInfo RightToTreatmentInfo `json:"reght_to_treatment_info"`
	Create_by            string               `json:"create_by"`
	Write_by             string               `json:"write_by"`
}

type CountryInfo struct {
	CountryInfo_ID int    `json:"country_info_id"`
	Country        string `json:"country"`
	CreateDate     string `json:"create_date"`
	WriteDate      string `json:"write_date"`
	Create_by      string `json:"create_by"`
	Write_by       string `json:"write_by"`
}

type Religion struct {
	ReligionInfo_ID int    `json:"religion_info_id"`
	Religion        string `json:"religion"`
	CreateDate      string `json:"create_date"`
	WriteDate       string `json:"write_date"`
	Create_by       string `json:"create_by"`
	Write_by        string `json:"write_by"`
}

type RightToTreatmentInfo struct {
	RightToTreatmentInfo_ID int    `json:"right_to_treatment_info_id"`
	Right_to_treatment      string `json:"right_to_treatment"`
	CreateDate              string `json:"create_date"`
	WriteDate               string `json:"write_date"`
	Create_by               string `json:"create_by"`
	Write_by                string `json:"write_by"`
}

type ServiceInfo struct {
	ServiceInfo_ID int    `json:"service_info_id"`
	Activity       string `json:"activity"`
	ServiceType    string `json:"service_type"`
	CreateDate     string `json:"create_date"`
	WriteDate      string `json:"write_date"`
	Create_by      string `json:"create_by"`
	Write_by       string `json:"write_by"`
}

type Activityrecord struct {
	ActivityRecord_ID int          `json:"activity_record_id"`
	StartTime         time.Time    `json:"start_time"`
	EndTime           time.Time    `json:"end_time"`
	Period            string       `json:"period"`
	Evidence_activity []byte       `json:"evidence_activity"`
	Evidence_time     []byte       `json:"evidence_time"`
	CreateDate        string       `json:"create_date"`
	WriteDate         string       `json:"write_date"`
	PatientInfo       PatientInfo  `json:"patient_info"`
	ServiceInfo       ServiceInfo  `json:"service_info"`
	EmployeeInfo      EmployeeInfo `json:"employee_info"`
	UserInfo          User_info    `json:"user_info"`
	Create_by         string       `json:"create_by"`
	Write_by          string       `json:"write_by"`

	// Location            string              `json:"location"`
}
type EmployeeInfo struct {
	EmployeeInfo_ID int             `json:"employee_info_id"`
	EmployeeCode    string          `json:"employee_code"`
	Sex             string          `json:"sex"`
	Name            string          `json:"username"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	CreateDate      string          `json:"create_date"`
	WriteDate       string          `json:"write_date"`
	DepartmentInfo  DepartmentInfo  `json:"department_inf"`
	JobPositionInfo JobPositionInfo `json:"job_position_info"`
	Create_by       string          `json:"create_by"`
	Write_by        string          `json:"write_by"`
}
type DepartmentInfo struct {
	DepartmentInfo_id int    `json:"department_info_id"`
	Department        string `json:"department"`
	CreateDate        string `json:"create_date"`
	WriteDate         string `json:"write_date"`
	Create_by         string `json:"create_by"`
	Write_by          string `json:"write_by"`
}
type JobPositionInfo struct {
	JobPositionInfo_id int    `json:"job_position_info_id"`
	JobPosition        string `json:"job_position"`
	CreateDate         string `json:"create_date"`
	WriteDate          string `json:"write_date"`
	Create_by          string `json:"create_by"`
	Write_by           string `json:"write_by"`
}

type WorktimeRecord struct {
	WorktimeRecord_ID int          `json:"worktime_record_id"`
	CheckIn           time.Time    `json:"check_in"`
	CheckOut          time.Time    `json:"check_out"`
	Period            string       `json:"period"`
	CreateDate        string       `json:"create_date"`
	WriteDate         string       `json:"write_date"`
	EmployeeInfo_ID   int          `json:"employee_info_id"`
	EmployeeInfo      EmployeeInfo `json:"employee_info"`
	UserInfo          User_info    `json:"user_info"`
	Create_by         string       `json:"create_by"`
	Write_by          string       `json:"write_by"`
}
type User_info struct {
	UserInfo_ID  int          `json:"user_info_id"`
	Line_user_id string       `json:"line_user_id"`
	Sex          string       `json:"sex"`
	Name         string       `json:"user_name"`
	Email        string       `json:"email"`
	PhoneNumber  string       `json:"phone_number"`
	CreateDate   string       `json:"create_date"`
	WriteDate    string       `json:"write_date"`
	EmployeeInfo EmployeeInfo `json:"employee_info"`
	Create_by    string       `json:"create_by"`
	Write_by     string       `json:"write_by"`
	
}

type LineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type LineProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	Email       string `json:"email"`
}
