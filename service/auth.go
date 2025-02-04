package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nirun/pkg/models"
)

// RequestBody โครงสร้าง JSON-RPC Request
type RequestBody struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// Params โครงสร้างพารามิเตอร์ที่ส่งไปยัง API
type Params struct {
	Service string        `json:"service"`
	Method  string        `json:"method"`
	Args    []interface{} `json:"args"`
}

// PostRequestByID ฟังก์ชันสำหรับส่ง HTTP POST Request ตาม ID ที่กำหนด
func PostRequestPatientByID(cardID string) (*models.PatientInfo, error) {
	// ตั้งค่าพารามิเตอร์ JSON-RPC
	params := Params{
		Service: "object",
		Method:  "execute_kw",
		Args: []interface{}{
			"nirun-community", // ฐานข้อมูล
			12,                // User ID
			"0809697302",      // รหัสผ่าน
			"ni.patient",      // โมเดล
			"search_read",     // เมธอด
			[]interface{}{
				[][]interface{}{{"identification_id", "=", cardID}},  // ใช้ค่า ID ที่รับเข้ามา
				[]string{"name", "title", "identification_id", "gender", "birthdate", "age", "phone", "mobile"}, // ฟิลด์ที่ต้องการ
			},
			map[string]interface{}{"limit": 80, "offset": 0, "order": "birthdate"},
		},
	}

	// กำหนดโครงสร้าง JSON-RPC
	requestBody := RequestBody{
		JSONRPC: "2.0",
		Method:  "call",
		Params:  params,
		ID:      12345, // รหัสอ้างอิงการร้องขอ
	}

	// แปลง JSON-RPC Request เป็น JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	// ส่ง HTTP POST Request
	url := "https://community.app.nirun.life/jsonrpc"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// กำหนด Header
	req.Header.Set("Content-Type", "application/json")

	// ส่ง Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// แปลง JSON Response เป็นโครงสร้างข้อมูล
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	// ตรวจสอบว่ามีข้อมูลผู้ป่วยหรือไม่
	if res, ok := result["result"].([]interface{}); ok && len(res) > 0 {
		patientData := res[0].(map[string]interface{})
		patient := models.PatientInfo{
			Name: patientData["name"].(string),
			// CardID:      patientData["identification_id"].(string),
			Sex: getGender(patientData, "gender"),
			// DateOfBirth: patientData["birthdate"].(string),
			Age:         fmt.Sprintf("%v ปี", patientData["age"]),
			PhoneNumber: patientData["phone"].(string),
			// RightToTreatmentInfo: models.RightToTreatmentInfo{
			// 	Right_to_treatment: getString(patientData, "right_to_treatment"),
			// },
		}

		// สร้าง Flex Message
		return &patient, nil
	}

	// ถ้าไม่พบข้อมูล ส่งข้อความแจ้งเตือน
	// return linebot.NewFlexMessage("ไม่พบข้อมูล", &linebot.BubbleContainer{
	// 	Body: &linebot.BoxComponent{
	// 		Type:   linebot.FlexComponentTypeBox,
	// 		Layout: linebot.FlexBoxLayoutTypeVertical,
	// 		Contents: []linebot.FlexComponent{
	// 			&linebot.TextComponent{
	// 				Type:  linebot.FlexComponentTypeText,
	// 				Text:  "ไม่พบข้อมูลของหมายเลขบัตรประชาชนนี้",
	// 				Size:  linebot.FlexTextSizeTypeMd,
	// 				Color: "#FF0000",
	// 				Wrap:  true,
	// 			},
	// 		},
	// 	},
	// })
	return nil, err
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "" // คืนค่าว่างถ้าไม่มีข้อมูล
}

// getGender ดึงค่า gender รองรับทั้ง bool และ string
func getGender(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case bool: // ถ้าเป็น boolean
			if v {
				return "ชาย"
			}
			return "ไม่ระบุ" // ถ้า false แสดงว่าไม่มีข้อมูล
		case string: // ถ้าเป็น string
			return v
		}
	}
	return "ไม่ระบุ" // คืนค่าเริ่มต้นหากไม่มีข้อมูล
}

// PostRequestService ฟังก์ชันส่ง JSON-RPC Request
func PostActivitiesByCategory(categoryID int) ([]string, error) {
    //ตั้งค่าพารามิเตอร์ JSON-RPC
    params := Params{
        Service: "object",
        Method:  "execute_kw",
        Args: []interface{}{
            "nirun-community",
            12,               
            "0809697302",     
            "ni.service",    
            "search_read",   
            []interface{}{
                [][]interface{}{
                    {"category_id", "=", categoryID},
                },
            },
            map[string]interface{}{"fields": []string{"name"}}, //map
        },
    }

    //สร้างโครงสร้าง JSON-RPC Request
    requestBody := RequestBody{
        JSONRPC: "2.0",
        Method:  "call",
        Params:  params,
        ID:      12345,
    }

    //แปลงเป็น JSON
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("❌ Error encoding JSON: %v", err)
    }

    //ส่ง HTTP POST Request
    url := "https://community.app.nirun.life/jsonrpc"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("Error creating request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    //ส่ง Request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("Error making request: %v", err)
    }
    defer resp.Body.Close()

    //อ่าน Response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Error reading response: %v", err)
    }

    //พิมพ์ JSON Response เพื่อตรวจสอบข้อมูลที่ได้รับ
    log.Printf("API Response for category_id %d: %s", categoryID, body)

    //แปลง JSON Response เป็น `map`
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("Error decoding JSON response: %v", err)
    }

    //ตรวจสอบและดึงข้อมูล
    var activities []string
    if res, ok := result["result"].([]interface{}); ok && len(res) > 0 {
        for _, item := range res {
            data := item.(map[string]interface{})
            if name, exists := data["name"].(string); exists {
                activities = append(activities, name)
            }
        }
    } else {
        log.Printf("No activities found for category ID %d", categoryID)
        return nil, fmt.Errorf("No activities found for category ID %d", categoryID)
    }

    return activities, nil
}


