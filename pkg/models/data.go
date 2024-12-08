package models

type PatientInfo struct {
	ServiceInfo    ServiceInfo    `json:service_info`
	Activityrecord Activityrecord `json:activity_record`
	CardID         string         `json:"card_id"`
	PatientID      string         `json:"patient_id"`
	Image          []byte         `json:"image"`
	Name           string         `json:"username"`
	PhoneNumber    string         `json:"phone_number"`
	Email          string         `json:"email"`
	Address        string         `json:"address"`
	Country        string         `json:"country"`
	Religion       string         `json:"religion"`
	Sex            string         `json:"sex"`
	Blood          string         `json:"blood"`
	DateOfBirth    string         `json:"date_of_birth"`
	Age            int            `json:"age"`
}

type ServiceInfo struct {
	//PatientInfo           PatientInfo `json:"patient_info"`
	CardID                string `json:"card_id"`
	ServiceCode           string `json:"service_code"`
	RightToTreatment      string `json:"right_to_treatment"`
	ServiceType           string `json:"service_type"`
	Activity              string `json:"activity"`
	Location              string `json:"location"`
	StartTime_ServiceInfo string `json:"start_time"`
	EndTime_ServiceInfo   string `json:"end_time"`
	Period_ServiceInfo    string `json:"period"`
	Selected              bool   `json:"selected"`
}

type Employee struct {
	EmployeeID           string `json:"employee_id"`
	Image_Employee       []byte `json:"image"`
	Name_Employee        string `json:"username"`
	Deployment           string `json:"deployment"`
	JobPosition          string `json:"job_position"`
	PhoneNumber_Employee string `json:"phone_number"`
	Email_Employee       string `json:"email"`
	Starttime_Employee   string `json:"start_time"`
	Endtime_Employee     string `json:"end_time"`
	Period_Employee      string `json:"period"`
}

type Activityrecord struct {
	CardID   string `json:"card_id"`
	Activity string `json:"activity"`
}
