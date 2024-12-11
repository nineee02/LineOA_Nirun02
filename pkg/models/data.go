package models

import "time"

type PatientInfo struct {
	Activityrecord   Activityrecord `json:activity_record`
	CardID           string         `json:"card_id"`
	PatientID        string         `json:"patient_id"`
	Name             string         `json:"username"`
	PhoneNumber      string         `json:"phone_number"`
	Email            string         `json:"email"`
	Address          string         `json:"address"`
	Country          string         `json:"country"`
	Religion         string         `json:"religion"`
	Sex              string         `json:"sex"`
	Blood            string         `json:"blood"`
	DateOfBirth      string         `json:"date_of_birth"`
	Age              int            `json:"age"`
	RightToTreatment string         `json:"right_to_treatment"`
}

type ServiceInfo struct {
	Activity    string `json:"activity"`
	ServiceType string `json:"service_type"`
	Location    string `json:"location"`
}

type Activityrecord struct {
	Usage_id  int      `json:"Usage_id"`
	StartTime string   `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Period    string   `json:"period"`
	Activity  string   `json:"activity"`
	CardID    string   `json:"card_id"`
}

type Employee struct {
	EmployeeID           string `json:"employee_id"`
	Name_Employee        string `json:"username"`
	Deployment           string `json:"deployment"`
	JobPosition          string `json:"job_position"`
	PhoneNumber_Employee string `json:"phone_number"`
	Email_Employee       string `json:"email"`
}

type WorktimeRecord struct {
	Worktime_UsagId int    `json:"worktime_usage_id"`
	EmployeeID      string `json:"employee_id"`
	CheckIn         string `json:"check_in"`
	CheckOutn       string `json:"check_out"`
	Period          string `json:"period"`
}
