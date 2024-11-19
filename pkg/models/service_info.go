package models
import(
	// "fmt"
	// "database/sql" 
	// "log"
)
type ServiceInfo struct {
	PatientInfo     PatientInfo `json:"patient_info"`
	ID              int         `json:"id"`
	IntoNumber      string      `json:"into_number"`
	ServiceCode     string      `json:"service_code"`
	RightToTreatment string      `json:"right_to_treatment"`
	ServiceType     string      `json:"service_type"`
	Activity        string      `json:"activity"`
	Location        string      `json:"location"`
	Day             string      `json:"day"`
	StartTime       string      `json:"start_time"`
	EndTime         string      `json:"end_time"`
	Period          string      `json:"period"`
}

type ServiceInfoSummary struct {
    ID       int    `json:"id"`
    Count    int    `json:"count"`
    Activity string `json:"activity"`
}

// func GetServiceInfoByActivity(db *sql.DB, activity string) (*ServiceInfo, error) {
//     const query = `
//         SELECT 
//             p.name_, p.patiet_id, p.age, p.sex, p.blood, p.phone_numbers,
//             p.email, p.address, p.country, p.card_id, p.religion, p.date_of_birth,
//             s.id, s.activity, s.service_code, s.into_number, s.location, 
//             s.start_time, s.end_time, s.period, s.right_to_treatment, s.service_type,
//             s.day
//         FROM patient_info p
//         JOIN service_info s ON p.patiet_id = s.patient_id
//         WHERE s.activity LIKE ?`

//     var service ServiceInfo
//     err := db.QueryRow(query, "%"+activity+"%").Scan(
//         // Patient info fields
//         &service.PatientInfo.Name,
//         &service.PatientInfo.PatientID,
//         &service.PatientInfo.Age,
//         &service.PatientInfo.Sex,
//         &service.PatientInfo.Blood,
//         &service.PatientInfo.PhoneNumber,
//         &service.PatientInfo.Email,
//         &service.PatientInfo.Address,
//         &service.PatientInfo.Country,
//         &service.PatientInfo.CardID,
//         &service.PatientInfo.Religion,
//         &service.PatientInfo.DateOfBirth,
//         // Service info fields
//         &service.ID,
//         &service.Activity,
//         &service.ServiceCode,
//         &service.IntoNumber,
//         &service.Location,
//         &service.StartTime,
//         &service.EndTime,
//         &service.Period,
//         &service.RightToTreatment,
//         &service.ServiceType,
//         &service.Day,
//     )

//     if err == sql.ErrNoRows {
//         return nil, fmt.Errorf("ไม่พบข้อมูลกิจกรรม: %s", activity)
//     }
//     if err != nil {
//         return nil, fmt.Errorf("database error: %v", err)
//     }

//     log.Printf("Retrieved service info: %+v", service)
//     return &service, nil
// }
