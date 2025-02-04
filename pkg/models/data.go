package models

import "time"

type PatientInfo struct {
	PatientInfo_ID       int                  `json:"patient_info_id"`
	CardID               string               `json:"card_id"`
	Name                 string               `json:"name_surname"`
	PhoneNumber          string               `json:"phone_number"`
	Email                string               `json:"email"`
	Address              string               `json:"address"`
	DateOfBirth          string               `json:"date_of_birth"`
	Age                  string               `json:"age"`
	Sex                  string               `json:"sex"`
	Blood                string               `json:"blood"`
	ADL                  string               `json:"ADL"`
	CreateDate           string               `json:"create_date"`
	UpdateDate           string               `json:"update_date"`
	Religion             Religion             `json:"religion_info"`
	CountryInfo          CountryInfo          `json:"country_info"`
	RightToTreatmentInfo RightToTreatmentInfo `json:"reght_to_treatment_info"`
	Create_by            string               `json:"create_by"`
	Update_by            string               `json:"update_by"`
}

type CountryInfo struct {
	CountryInfo_ID int    `json:"country_info_id"`
	Country        string `json:"country"`
	CreateDate     string `json:"create_date"`
	UpdateDate     string `json:"update_date"`
	Create_by      string `json:"create_by"`
	Update_by      string `json:"update_by"`
}

type Religion struct {
	ReligionInfo_ID int    `json:"religion_info_id"`
	Religion        string `json:"religion"`
	CreateDate      string `json:"create_date"`
	UpdateDate      string `json:"update_date"`
	Create_by       string `json:"create_by"`
	Update_by       string `json:"update_by"`
}

type RightToTreatmentInfo struct {
	RightToTreatmentInfo_ID int    `json:"right_to_treatment_info_id"`
	Right_to_treatment      string `json:"right_to_treatment"`
	CreateDate              string `json:"create_date"`
	UpdateDate              string `json:"update_date"`
	Create_by               string `json:"create_by"`
	Update_by               string `json:"update_by"`
}

type Activityrecord struct {
	ActivityRecord_ID         int                       `json:"activity_record_id"`
	StartTime                 time.Time                 `json:"start_time"`
	EndTime                   time.Time                 `json:"end_time"`
	Period                    string                    `json:"period"`
	Evidence_activity         []byte                    `json:"evidence_activity"`
	Evidence_time             []byte                    `json:"evidence_time"`
	CreateDate                string                    `json:"create_date"`
	UpdateDate                string                    `json:"update_date"`
	UserInfo                  User_info                 `json:"user_info"`
	Create_by                 string                    `json:"create_by"`
	Update_by                 string                    `json:"update_by"`
	PatientInfo               PatientInfo               `json:"patient_info"`
	EmployeeInfo              EmployeeInfo              `json:"employee_info"`
	ActivityTechnologyInfo    ActivityTechnologyInfo    `json:"activity_technology_info"`
	ActivitySocialInfo        ActivitySocialInfo        `json:"activity_social_info"`
	ActivityHealthInfo        ActivityHealthInfo        `json:"activity_health_info"`
	ActivityEconomicInfo      ActivityEconomicInfo      `json:"activity_economic_info"`
	ActivityEnvironmentalInfo ActivityEnvironmentalInfo `json:"activity_environmental_info"`
	ActivityOther             string                    `json:"activity_other"`
}

type ActivityTechnologyInfo struct {
	ActivityTechnologyInfo_ID int    `json:"activity_technology_id"`
	ActivityTechnology        string `json:"activity"`
	ServiceType               string `json:"service_type"`
	CreateDate                string `json:"create_date"`
	UpdateDate                string `json:"update_date"`
	Create_by                 string `json:"create_by"`
	Update_by                 string `json:"update_by"`
}
type ActivitySocialInfo struct {
	ActivitySocialInfo_ID int    `json:"activity_social_id"`
	ActivitySocial        string `json:"activity"`
	ServiceType           string `json:"service_type"`
	CreateDate            string `json:"create_date"`
	UpdateDate            string `json:"update_date"`
	Create_by             string `json:"create_by"`
	Update_by             string `json:"update_by"`
}
type ActivityHealthInfo struct {
	ActivityHealthInfo_ID int    `json:"activity_health_id"`
	ActivityHealth        string `json:"activity"`
	ServiceType           string `json:"service_type"`
	CreateDate            string `json:"create_date"`
	UpdateDate            string `json:"update_date"`
	Create_by             string `json:"create_by"`
	Update_by             string `json:"update_by"`
}
type ActivityEconomicInfo struct {
	ActivityEconomicInfo_ID int    `json:"activity_economic_id"`
	ActivityEconomic        string `json:"activity"`
	ServiceType             string `json:"service_type"`
	CreateDate              string `json:"create_date"`
	UpdateDate              string `json:"update_date"`
	Create_by               string `json:"create_by"`
	Update_by               string `json:"update_by"`
}
type ActivityEnvironmentalInfo struct {
	ActivityEnvironmentalInfo_ID int    `json:"activity_environmental_id"`
	ActivityEnvironmental        string `json:"activity"`
	ServiceType                  string `json:"service_type"`
	CreateDate                   string `json:"create_date"`
	UpdateDate                   string `json:"update_date"`
	Create_by                    string `json:"create_by"`
	Update_by                    string `json:"update_by"`
}
type ProblemRecord struct {
	ProblemRecord_ID         int                      `json:"problem_record_id"`
	ProblemTechnologyInfo    ProblemTechnologyInfo    `json:"problem_technology_id"`
	ProblemSocialInfo        ProblemSocialInfo        `json:"problem_social_id"`
	ProblemHealthInfo        ProblemHealthInfo        `json:"problem_health_id"`
	ProblemEconomicInfo      ProblemEconomicInfo      `json:"problem_economic_id"`
	ProblemEnvironmentalInfo ProblemEnvironmentalInfo `json:"problem_environmental_id"`
	ProblemOther             string                   `json:"problem_other"`
	CreateDate               string                   `json:"create_date"`
	UpdateDate               string                   `json:"update_date"`
	Create_by                string                   `json:"create_by"`
	Update_by                string                   `json:"update_by"`
	UserInfo                 User_info                `json:"user_info"`
	PatientInfo              PatientInfo              `json:"patient_info"`
}
type ProblemTechnologyInfo struct {
	ProblemTechnologyInfo_ID int    `json:"problem_technology_id"`
	ProblemTechnology        string `json:"problem"`
	ServiceType              string `json:"service_type"`
	CreateDate               string `json:"create_date"`
	UpdateDate               string `json:"update_date"`
	Create_by                string `json:"create_by"`
	Update_by                string `json:"update_by"`
}
type ProblemSocialInfo struct {
	ProblemSocialInfo_ID int    `json:"problem_social_id"`
	ProblemSocial        string `json:"problem"`
	ServiceType          string `json:"service_type"`
	CreateDate           string `json:"create_date"`
	UpdateDate           string `json:"update_date"`
	Create_by            string `json:"create_by"`
	Update_by            string `json:"update_by"`
}
type ProblemHealthInfo struct {
	ProblemHealthInfo_ID int    `json:"problem_health_id"`
	ProblemSocial        string `json:"problem"`
	ServiceType          string `json:"service_type"`
	CreateDate           string `json:"create_date"`
	UpdateDate           string `json:"update_date"`
	Create_by            string `json:"create_by"`
	Update_by            string `json:"update_by"`
}
type ProblemEconomicInfo struct {
	ProblemEconomicInfo_ID int    `json:"problem_economic_id"`
	ProblemEconomic        string `json:"problem"`
	ServiceType            string `json:"service_type"`
	CreateDate             string `json:"create_date"`
	UpdateDate             string `json:"update_date"`
	Create_by              string `json:"create_by"`
	Update_by              string `json:"update_by"`
}
type ProblemEnvironmentalInfo struct {
	ProblemEnvironmentalInfo_ID int    `json:"problem_environmental_id"`
	ProblemEnvironmental        string `json:"problem"`
	ServiceType                 string `json:"service_type"`
	CreateDate                  string `json:"create_date"`
	UpdateDate                  string `json:"update_date"`
	Create_by                   string `json:"create_by"`
	Update_by                   string `json:"update_by"`
}
type EmployeeInfo struct {
	EmployeeInfo_ID int             `json:"employee_info_id"`
	EmployeeCode    string          `json:"employee_code"`
	Sex             string          `json:"sex"`
	Name            string          `json:"name_surname"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	CreateDate      string          `json:"create_date"`
	UpdateDate      string          `json:"update_date"`
	DepartmentInfo  DepartmentInfo  `json:"department_inf"`
	JobPositionInfo JobPositionInfo `json:"job_position_info"`
	Create_by       string          `json:"create_by"`
	Update_by       string          `json:"update_by"`
}
type DepartmentInfo struct {
	DepartmentInfo_id int    `json:"department_info_id"`
	Department        string `json:"department"`
	CreateDate        string `json:"create_date"`
	UpdateDate        string `json:"update_date"`
	Create_by         string `json:"create_by"`
	Update_by         string `json:"update_by"`
}
type JobPositionInfo struct {
	JobPositionInfo_id int    `json:"job_position_info_id"`
	JobPosition        string `json:"job_position"`
	CreateDate         string `json:"create_date"`
	UpdateDate         string `json:"update_date"`
	Create_by          string `json:"create_by"`
	Update_by          string `json:"update_by"`
}

type WorktimeRecord struct {
	WorktimeRecord_ID int           `json:"worktime_record_id"`
	CheckIn           time.Time     `json:"check_in"`
	CheckOut          time.Time     `json:"check_out"`
	Period            string        `json:"period"`
	Location          string        `json:"location"`
	CreateDate        string        `json:"create_date"`
	UpdateDate        string        `json:"update_date"`
	EmployeeInfo      *EmployeeInfo `json:"employee_info"`
	UserInfo          *User_info    `json:"user_info"`
	Create_by         string        `json:"create_by"`
	Update_by         string        `json:"update_by"`
}
type User_info struct {
	UserInfo_ID  int          `json:"user_info_id"`
	Line_user_id string       `json:"line_user_id"`
	Sex          string       `json:"sex"`
	Name         string       `json:"user_name"`
	Email        string       `json:"email"`
	PhoneNumber  string       `json:"phone_number"`
	CreateDate   string       `json:"create_date"`
	UpdateDate   string       `json:"update_date"`
	EmployeeInfo EmployeeInfo `json:"employee_info"`
	Create_by    string       `json:"create_by"`
	Update_by    string       `json:"update_by"`
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
