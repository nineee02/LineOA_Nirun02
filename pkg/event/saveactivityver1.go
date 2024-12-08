package event

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/line/line-bot-sdk-go/linebot"
// )

// func createQuickReplyActivities() linebot.QuickReplyItems {
// 	// รายการกิจกรรมทั้งหมด
// 	activities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// 	}

// 	quickReplyItems := linebot.QuickReplyItems{}
// 	for _, activity := range activities {
// 		quickReplyItems.Items = append(quickReplyItems.Items,
// 			linebot.NewQuickReplyButton("", linebot.NewMessageAction(activity, activity)),
// 		)
// 	}
// 	return quickReplyItems
// }

// func SaveActivity(db *sql.DB, activity string) error {
// 	if !validateActivity(activity) {
// 		return fmt.Errorf("กิจกรรม '%s' ไม่ตรงกับค่าที่อนุญาตในฐานข้อมูล", activity)
// 	}

// 	query := `INSERT INTO service_info (activity) VALUES (?)`
// 	_, err := db.Exec(query, activity)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถบันทึกกิจกรรม %s ได้: %v", activity, err)
// 	}
// 	return nil
// }

// // validateActivity ตรวจสอบว่าค่าที่ส่งมาตรงกับฐานข้อมูล หรือไม่
// func validateActivity(activity string) bool {
// 	allowedActivities := []string{
// 		"แช่เท้า", "นวด/ประคบ", "ฝังเข็ม", "คาราโอเกะ", "ครอบแก้ว",
// 		"ทำอาหาร", "นั่งสมาธิ", "เล่าสู่กัน", "ซุโดกุ", "จับคู่ภาพ",
// 	}
// 	for _, allowed := range allowedActivities {
// 		if activity == allowed {
// 			return true
// 		}
// 	}
// 	return false
// }
