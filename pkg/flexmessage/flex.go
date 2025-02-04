package flexmessage

import (
	"log"
	"nirun/pkg/models"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func getCurrentTime() string {
	format := "02-01-2006 03:04 PM"
	return time.Now().Format(format)
}
func FormatConfirmCheckin(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
	// ตรวจสอบว่า worktimeRecord ไม่เป็น nil
	if worktimeRecord == nil {
		return nil
	}

	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR, Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ลงเวลางาน",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "เช็คอิน",
								Text:  "เช็คอิน",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
					},
				},
			},
		},
	}

	// ส่งกลับ Flex Message
	return linebot.NewFlexMessage("ยืนยันการเช็คอิน/เช็คเอ้าท์", container)
}
func FormatConfirmCheckout(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
	// ตรวจสอบว่า worktimeRecord ไม่เป็น nil
	if worktimeRecord == nil {
		return nil
	}

	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR, Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ลงเวลางาน",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "เช็คเอ้าท์",
								Text:  "เช็คเอ้าท์",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
					},
				},
			},
		},
	}

	// ส่งกลับ Flex Message
	return linebot.NewFlexMessage("ยืนยันการเช็คอิน/เช็คเอ้าท์", container)
}

//แนวนอน
// func FormatConfirmationWorktime(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
// 	// ตรวจสอบว่า worktimeRecord ไม่เป็น nil
// 	if worktimeRecord == nil {
// 		return nil
// 	}

// 	// สร้าง BubbleContainer สำหรับ Flex Message
// 	container := &linebot.BubbleContainer{
// 		Type:      linebot.FlexContainerTypeBubble,
// 		Direction: linebot.FlexBubbleDirectionTypeLTR,
// 		Header: &linebot.BoxComponent{
// 			Type:   linebot.FlexComponentTypeBox,
// 			Layout: linebot.FlexBoxLayoutTypeVertical,
// 			Contents: []linebot.FlexComponent{
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   "ลงเวลางาน",
// 					Weight: linebot.FlexTextWeightTypeBold,
// 					Size:   linebot.FlexTextSizeTypeLg,
// 					Color:  "#FFFFFF",
// 					Align:  linebot.FlexComponentAlignTypeCenter,
// 				},
// 			},
// 			BackgroundColor: "#00bcd4",
// 		},
// 		Body: &linebot.BoxComponent{
// 			Type:    linebot.FlexComponentTypeBox,
// 			Layout:  linebot.FlexBoxLayoutTypeVertical,
// 			Spacing: linebot.FlexComponentSpacingTypeMd,
// 			Contents: []linebot.FlexComponent{
// 				// ข้อความแนะนำ
// 				&linebot.TextComponent{
// 					Type:    linebot.FlexComponentTypeText,
// 					Text:    "กรุณาเลือกดำเนินการ",
// 					Weight:  linebot.FlexTextWeightTypeRegular,
// 					Size:    linebot.FlexTextSizeTypeMd,
// 					Color:   "#212121",
// 					Align:   linebot.FlexComponentAlignTypeStart,
// 					Gravity: linebot.FlexComponentGravityTypeCenter,
// 				},
// 				// เส้นแบ่ง
// 				&linebot.SeparatorComponent{
// 					Color:  "#58BDCF",
// 					Margin: linebot.FlexComponentMarginTypeMd,
// 				},
// 				// ปุ่ม Check-in และ Check-out ในแนวนอน
// 				&linebot.BoxComponent{
// 					Type:    linebot.FlexComponentTypeBox,
// 					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
// 					Spacing: linebot.FlexComponentSpacingTypeMd,
// 					Margin:  linebot.FlexComponentMarginTypeLg,
// 					Contents: []linebot.FlexComponent{
// 						&linebot.ButtonComponent{
// 							Type: linebot.FlexComponentTypeButton,
// 							Action: &linebot.MessageAction{
// 								Label: "เช็คอิน",
// 								Text:  "เช็คอิน",
// 							},
// 							Style: linebot.FlexButtonStyleTypePrimary,
// 							Color: "#00bcd4",
// 							Gravity: linebot.FlexComponentGravityTypeCenter,
// 							Flex:    linebot.IntPtr(1),
// 						},
// 						&linebot.ButtonComponent{
// 							Type: linebot.FlexComponentTypeButton,
// 							Action: &linebot.MessageAction{
// 								Label: "เช็คเอ้าท์",
// 								Text:  "เช็คเอ้าท์",
// 							},
// 							Style: linebot.FlexButtonStyleTypeSecondary,
// 							Gravity: linebot.FlexComponentGravityTypeCenter,
// 							Flex:    linebot.IntPtr(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// ส่งกลับ Flex Message
// 	return linebot.NewFlexMessage("ยืนยันการเช็คอิน/เช็คเอ้าท์", container)
// }

func FormatworktimeCheckin(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
	if worktimeRecord == nil {
		return nil
	}

	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ยินดีต้อนรับ",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				// ชื่อ
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   worktimeRecord.UserInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},

				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เช็คอินที่: " + getCurrentTime(),
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeMd,
				},
			},
		},
	}

	// สร้าง Flex Message
	return linebot.NewFlexMessage("ลงเวลาเข้างาน", container)
}
func FormatworktimeCheckout(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
	if worktimeRecord == nil {
		return nil
	}

	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ลาก่อน",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				// ชื่อ
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   worktimeRecord.UserInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},

				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เช็คอินที่: " + getCurrentTime(),
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeMd,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   worktimeRecord.Period,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeMd,
				},
			},
		},
	}

	// สร้าง Flex Message
	return linebot.NewFlexMessage("ลงเวลาออกงาน", container)
}
func intPtr(i int) *int {
	return &i
}
func FormatPatientInfo(patient *models.PatientInfo) *linebot.FlexMessage {
	if patient == nil || patient.Name == "" {
		log.Println("Error: ไม่พบข้อมูลผู้ป่วย")
		return linebot.NewFlexMessage("ไม่พบข้อมูล", &linebot.BubbleContainer{
			Body: &linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:  linebot.FlexComponentTypeText,
						Text:  "ไม่พบข้อมูลผู้ป่วย กรุณาลองใหม่",
						Size:  linebot.FlexTextSizeTypeMd,
						Color: "#FF0000",
						Wrap:  true,
					},
				},
			},
		})
	}

	// ตรวจสอบสิทธิการรักษา
	rightToTreatment := "ไม่มีข้อมูล"
	if (models.RightToTreatmentInfo{}) != patient.RightToTreatmentInfo {
		rightToTreatment = getSafeString(&patient.RightToTreatmentInfo.Right_to_treatment, "ไม่มีข้อมูล")
	}

	// สร้าง BubbleContainer
	container := &linebot.BubbleContainer{
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   getSafeString(&patient.Name, "ไม่มีข้อมูล"),
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#000000",
					Align:  linebot.FlexComponentAlignTypeStart,
				},
				createTextRow("อายุ", getSafeString(&patient.Age, "ไม่ระบุ")),
				createTextRow("เพศ", formatGender(getSafeString(&patient.Sex, "ไม่ระบุ"))),
				createTextRow("หมู่เลือด", getSafeString(&patient.Blood, "ไม่ระบุ")),
				createTextRow("ADL", getSafeString(&patient.ADL, "ไม่มีข้อมูล")),
				createTextRow("เบอร์", getSafeString(&patient.PhoneNumber, "ไม่มีข้อมูล")),
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				createTextRow("สิทธิการรักษา", rightToTreatment),
			},
			BackgroundColor: "#f3fcfd",
		},
	}

	log.Printf("Flex Message Data: %+v", patient)

	// สร้าง FlexMessage
	return linebot.NewFlexMessage("ข้อมูลผู้สูงอายุ", container)
}

// ฟังก์ชันป้องกัน nil
func getSafeString(value *string, defaultValue string) string {
	if value != nil && *value != "" {
		return *value
	}
	return defaultValue
}

// func FormatPatientInfo(patient *models.PatientInfo) *linebot.FlexMessage {
// 	// สร้าง BubbleContainer
// 	container := &linebot.BubbleContainer{
// 		// 	Type: linebot.FlexContainerTypeBubble,
// 		// 	Size: linebot.FlexBubbleSizeTypeMega,
// 		// 	Hero: &linebot.ImageComponent{
// 		// 		Type:        linebot.FlexComponentTypeImage,
// 		// 		// URL:         imageUrl, // ใช้ URL ของภาพที่ได้รับจาก MinIO
// 		// 		Size:        linebot.FlexImageSizeTypeFull,
// 		// 		AspectRatio: "20:13",
// 		// 		AspectMode:  "cover",
// 		// 	},
// 		// Header: &linebot.BoxComponent{
// 		// 	Type:   linebot.FlexComponentTypeBox,
// 		// 	Layout: linebot.FlexBoxLayoutTypeVertical,
// 		// 	Contents: []linebot.FlexComponent{
// 		// 		&linebot.TextComponent{
// 		// 			Type:   linebot.FlexComponentTypeText,
// 		// 			Text:   patient.PatientInfo.Name,
// 		// 			Weight: linebot.FlexTextWeightTypeBold,
// 		// 			Size:   linebot.FlexTextSizeTypeLg,
// 		// 			Color:  "#FFFFFF",
// 		// 			Align:  linebot.FlexComponentAlignTypeStart,
// 		// 		},
// 		// 		// บรรทัดที่สอง: ข้อมูลเลขประจำตัวประชาชน
// 		// 		&linebot.TextComponent{
// 		// 			Type:   linebot.FlexComponentTypeText,
// 		// 			Text:   "เลขประจำตัวประชาชน: " + patient.PatientInfo.CardID,
// 		// 			Size:   linebot.FlexTextSizeTypeSm,
// 		// 			Color:  "#F8F8F8",
// 		// 			Margin: linebot.FlexComponentMarginTypeXs,
// 		// 			Align:  linebot.FlexComponentAlignTypeStart,
// 		// 		},
// 		// 	},
// 		// 	BackgroundColor: "#08BED7",
// 		// },
// 		Body: &linebot.BoxComponent{
// 			Type:   linebot.FlexComponentTypeBox,
// 			Layout: linebot.FlexBoxLayoutTypeVertical,
// 			Contents: []linebot.FlexComponent{
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   patient.Name,
// 					Weight: linebot.FlexTextWeightTypeBold,
// 					Size:   linebot.FlexTextSizeTypeLg,
// 					Color:  "#000000",
// 					Align:  linebot.FlexComponentAlignTypeStart,
// 				},
// 				// บรรทัดที่สอง: ข้อมูลเลขประจำตัวประชาชน
// 				// &linebot.TextComponent{
// 				// 	Type:   linebot.FlexComponentTypeText,
// 				// 	Text:   "เลขประจำตัวประชาชน: " + patient.PatientInfo.CardID,
// 				// 	Size:   linebot.FlexTextSizeTypeSm,
// 				// 	Color:  "#555555",
// 				// 	Margin: linebot.FlexComponentMarginTypeXs,
// 				// 	Align:  linebot.FlexComponentAlignTypeStart,
// 				// },
// 				// ข้อมูลผู้ป่วย
// 				createTextRow("อายุ", patient.Age+" ปี"),
// 				createTextRow("เพศ", formatGender(patient.Sex)),
// 				createTextRow("หมู่เลือด", patient.Blood),
// 				createTextRow("ADL", patient.ADL),
// 				createTextRow("เบอร์", patient.PhoneNumber),
// 				&linebot.SeparatorComponent{
// 					Type:   linebot.FlexComponentTypeSeparator,
// 					Color:  "#58BDCF",
// 					Margin: linebot.FlexComponentMarginTypeXl,
// 				},

// 				// ข้อมูลสิทธิการรักษา
// 				&linebot.BoxComponent{
// 					Type:   linebot.FlexComponentTypeBox,
// 					Layout: linebot.FlexBoxLayoutTypeVertical,
// 					Contents: []linebot.FlexComponent{
// 						&linebot.TextComponent{
// 							Type:   linebot.FlexComponentTypeText,
// 							Text:   "สิทธิการรักษา:",
// 							Weight: linebot.FlexTextWeightTypeBold,
// 							Size:   linebot.FlexTextSizeTypeMd,
// 							Color:  "#000000",
// 							Margin: linebot.FlexComponentMarginTypeMd,
// 						},
// 						&linebot.TextComponent{
// 							Type:   linebot.FlexComponentTypeText,
// 							Text:   patient.RightToTreatmentInfo.Right_to_treatment,
// 							Size:   linebot.FlexTextSizeTypeMd,
// 							Color:  "#212121",
// 							Align:  linebot.FlexComponentAlignTypeStart,
// 							Wrap:   true, // เปิดตัดข้อความอัตโนมัติ
// 							Margin: linebot.FlexComponentMarginTypeSm,
// 						},
// 						&linebot.SpacerComponent{},
// 					},
// 				},
// 			},
// 			BackgroundColor: "#f3fcfd", // สีพื้นหลังสำหรับ Body

// 		},
// 	}

// 	// สร้าง FlexMessage
// 	return linebot.NewFlexMessage("ข้อมูลผู้สูงอายุ", container)
// }

// &linebot.TextComponent{
// 	Type:   linebot.FlexComponentTypeText,
// 	Text:   patient.PatientInfo.Name,
// 	Weight: linebot.FlexTextWeightTypeBold,
// 	Size:   linebot.FlexTextSizeTypeLg,
// 	Color:  "#FFFFFF", // สีขาวสำหรับข้อความ
// 	Align:  linebot.FlexComponentAlignTypeStart,
// },
// // บรรทัดที่สอง: ข้อมูลเลขประจำตัวประชาชน
// &linebot.TextComponent{
// 	Type:   linebot.FlexComponentTypeText,
// 	Text:   "เลขประจำตัวประชาชน: " + patient.PatientInfo.CardID,
// 	Size:   linebot.FlexTextSizeTypeSm,
// 	Color:  "#FFFFFF", // สีขาวสำหรับข้อความ
// 	Margin: linebot.FlexComponentMarginTypeXs,
// 	Align:  linebot.FlexComponentAlignTypeStart,
// },

// ฟังก์ชันแปลงเพศ
func formatGender(sex string) string {
	if sex == "Male" {
		return "ชาย"
	} else if sex == "Female" {
		return "หญิง"
	}
	return "ไม่ระบุ"
}

func FormatServiceSelection() *linebot.FlexMessage {
	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		// Direction: linebot.FlexBubbleDirectionTypeLTR, Header: &linebot.BoxComponent{
		// 	Type:   linebot.FlexComponentTypeBox,
		// 	Layout: linebot.FlexBoxLayoutTypeVertical,
		// 	Contents: []linebot.FlexComponent{
		// 		&linebot.TextComponent{
		// 			Type:   linebot.FlexComponentTypeText,
		// 			Text:   "ลงเวลางาน",
		// 			Weight: linebot.FlexTextWeightTypeBold,
		// 			Size:   linebot.FlexTextSizeTypeLg,
		// 			Color:  "#FFFFFF",
		// 			Align:  linebot.FlexComponentAlignTypeCenter,
		// 		},
		// 	},
		// 	BackgroundColor: "#00bcd4",
		// },
		Body: &linebot.BoxComponent{

			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,

			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "เลือกบริการ:",
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeMd,
							Color:  "#00bcd4",
							Margin: linebot.FlexComponentMarginTypeMd,
							Align:  linebot.FlexComponentAlignTypeStart,
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Color:  "#58BDCF",
							Margin: linebot.FlexComponentMarginTypeXl,
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "บันทึกกิจกรรม",
								Text:  "บันทึกกิจกรรม",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "รายงานปัญหา",
								Text:  "รายงานปัญหา",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
					},
				},
			},
		},
	}

	// ส่งกลับ Flex Message
	return linebot.NewFlexMessage("บันทึกกิจกรรม/รายงานปัญหา", container)
}

// มิติกิจกรรม
func FormatActivityCategories() *linebot.FlexMessage {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		// Direction: linebot.FlexBubbleDirectionTypeLTR, Header: &linebot.BoxComponent{
		// 	Type:   linebot.FlexComponentTypeBox,
		// 	Layout: linebot.FlexBoxLayoutTypeVertical,
		// 	Contents: []linebot.FlexComponent{
		// 		&linebot.TextComponent{
		// 			Type:   linebot.FlexComponentTypeText,
		// 			Text:   "ลงเวลางาน",
		// 			Weight: linebot.FlexTextWeightTypeBold,
		// 			Size:   linebot.FlexTextSizeTypeLg,
		// 			Color:  "#FFFFFF",
		// 			Align:  linebot.FlexComponentAlignTypeCenter,
		// 		},
		// 	},
		// 	BackgroundColor: "#00bcd4",
		// },
		Body: &linebot.BoxComponent{

			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,

			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "เลือกมิติของกิจกรรม:",
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeMd,
							Color:  "#00bcd4",
							Margin: linebot.FlexComponentMarginTypeMd,
							Align:  linebot.FlexComponentAlignTypeStart,
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Color:  "#58BDCF",
							Margin: linebot.FlexComponentMarginTypeXl,
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติเทคโนโลยี",
								Text:  "มิติเทคโนโลยี",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติสังคม",
								Text:  "มิติสังคม",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติสุขภาพ",
								Text:  "มิติสุขภาพ",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติเศรษฐกิจ",
								Text:  "มิติเศรษฐกิจ",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติสิ่งแวดล้อม",
								Text:  "มิติสิ่งแวดล้อม",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "มิติอื่นๆ",
								Text:  "มิติอื่นๆ",
							},
							Style: linebot.FlexButtonStyleTypePrimary,
							Color: "#00bcd4",
						},
					},
				},
			},
		},
	}

	// ส่งกลับ Flex Message
	return linebot.NewFlexMessage("เลือกมิติกิจกรรม", container)
}

// แสดงรายการกิจกรรม
func FormatActivities(activities []string) *linebot.FlexMessage {
	var activityButtons []linebot.FlexComponent

	// ✅ ฟังก์ชันตัดข้อความให้ยาวไม่เกิน 40 ตัวอักษร
	truncateLabel := func(text string) string {
		if len(text) > 40 {
			return text[:37] + "..." // ตัดให้เหลือ 37 ตัว แล้วเติม "..."
		}
		return text
	}

	// ✅ สร้างปุ่ม "🔙 ย้อนกลับ" เป็นปุ่มแรก
	backButton := &linebot.ButtonComponent{
		Type:   linebot.FlexComponentTypeButton,
		Style:  linebot.FlexButtonStyleTypeSecondary,
		Height: linebot.FlexButtonHeightTypeMd,
		Color:  "#FF9800", // 🔸 ปรับสีให้ดูเด่นขึ้น
		Action: &linebot.MessageAction{
			Label: "🔙 ย้อนกลับ",
			Text:  "ย้อนกลับ",
		},
	}
	activityButtons = append(activityButtons, backButton) // ใส่ปุ่มย้อนกลับก่อนกิจกรรม

	// ✅ สร้างปุ่มกิจกรรมจากรายการที่รับมา
	for _, activity := range activities {
		trimmedActivity := strings.TrimSpace(activity) // ตัดช่องว่างด้านหน้า-หลัง
		if trimmedActivity == "" {
			continue // ข้ามกิจกรรมที่เป็นค่าว่าง
		}

		button := &linebot.ButtonComponent{
			Type:   linebot.FlexComponentTypeButton,
			Style:  linebot.FlexButtonStyleTypePrimary,
			Height: linebot.FlexButtonHeightTypeMd,
			Action: &linebot.MessageAction{
				Label: truncateLabel(trimmedActivity), // ✅ ตรวจสอบความยาว
				Text:  trimmedActivity,
			},
		}
		activityButtons = append(activityButtons, button)
	}

	// ✅ สร้างโครงสร้าง Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Size:      linebot.FlexBubbleSizeTypeMega,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลือกกิจกรรม",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXl,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: append([]linebot.FlexComponent{
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  "กรุณาเลือกกิจกรรมที่ต้องการบันทึก:",
					Size:  linebot.FlexTextSizeTypeMd,
					Align: linebot.FlexComponentAlignTypeStart,
					Wrap:  true,
				},
			}, activityButtons...), // ✅ ใส่ปุ่มกิจกรรม + ปุ่มย้อนกลับ
		},
	}

	return linebot.NewFlexMessage("เลือกกิจกรรม", container)
}


// ฟังก์ชันสร้างแถวข้อความ
func createTextRow(label string, value string) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   label + ":",
				Weight: linebot.FlexTextWeightTypeBold,
				Size:   linebot.FlexTextSizeTypeMd,
				Color:  "#100F0FFF",
				Flex:   intPtr(1), //จัดระเบียบข้อความ
			},
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   value,
				Weight: linebot.FlexTextWeightTypeRegular,
				Size:   linebot.FlexTextSizeTypeMd,
				Color:  "#333333",
				Flex:   intPtr(2),
			},
		},
		Margin: linebot.FlexComponentMarginTypeMd,
	}
}

func getPeriodOrDefault(period string) string {
	if period == "" {
		return "ไม่ระบุ"
	}
	return period
}

func createActivityButtons(activities []struct {
	Label string
	Text  string
	Color string
}) *linebot.BoxComponent {
	rows := []*linebot.BoxComponent{}

	// สร้างแถวของปุ่ม (2 ปุ่มต่อแถว)
	for i := 0; i < len(activities); i += 2 {
		row := &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeHorizontal,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Margin:  linebot.FlexComponentMarginTypeMd,
		}
		for j := 0; j < 2 && i+j < len(activities); j++ {
			activity := activities[i+j]
			row.Contents = append(row.Contents, &linebot.ButtonComponent{
				Type: linebot.FlexComponentTypeButton,
				Action: &linebot.MessageAction{
					Label: activity.Label,               // แสดงข้อความในปุ่ม
					Text:  normalizeText(activity.Text), // ใช้ normalizeText เพื่อลบช่องว่าง
				},
				Style:  linebot.FlexButtonStyleTypeSecondary,
				Color:  activity.Color,
				Height: linebot.FlexButtonHeightTypeMd,
			})
		}
		rows = append(rows, row)
	}

	// รวมทุกแถวเป็นแนวตั้ง
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: func() []linebot.FlexComponent {
			components := []linebot.FlexComponent{}
			for _, row := range rows {
				components = append(components, row)
			}
			return components
		}(),
	}
}

func FormatServiceInfo(activity []models.Activityrecord) *linebot.FlexMessage {
	// รายการกิจกรรมที่สามารถเลือกเพิ่มได้
	activities := []struct {
		Label string
		Text  string
		Color string
	}{
		{"แช่เท้า", "แช่เท้า", "#4DD0E1"},
		{"นวด / ประคบ", "นวด/ประคบ", "#4DD0E1"},

		{"ฝังเข็ม", "ฝังเข็ม", "#4DD0E1"},
		{"คาราโอเกะ", "คาราโอเกะ", "#4DD0E1"},

		{"ครอบแก้ว", "ครอบแก้ว", "#4DD0E1"},
		{"ทำอาหาร", "ทำอาหาร", "#4DD0E1"},

		{"นั่งสมาธิ", "นั่งสมาธิ", "#4DD0E1"},
		{"เล่าสู่กันฟัง", "เล่าสู่กันฟัง", "#4DD0E1"},

		{"ซุโดกุ", "ซุโดกุ", "#4DD0E1"},
		{"จับคู่ภาพ", "จับคู่ภาพ", "#4DD0E1"},
	}

	// สร้าง BubbleContainer
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		// Hero: &linebot.ImageComponent{
		// 	Type:            linebot.FlexComponentTypeImage,
		// 	URL:             "https://www.okmd.or.th/upload/iblock/82c/aging_850x446-fb.png",
		// 	Size:            linebot.FlexImageSizeTypeFull,
		// 	AspectRatio:     "1.51:1",
		// 	AspectMode:      "cover",
		// 	BackgroundColor: "#FFFFFFFF",
		// },
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   activity[0].PatientInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Align:  linebot.FlexComponentAlignTypeCenter,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Wrap:   true,
					Margin: linebot.FlexComponentMarginTypeXs,
				},

				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   "เลขประจำตัวประชาชน: \n" + activity[0].PatientInfo.CardID,
				// 	Weight: linebot.FlexTextWeightTypeRegular,
				// 	Size:   linebot.FlexTextSizeTypeSm,
				// 	Color:  "#FFFFFF",
				// 	Align:  linebot.FlexComponentAlignTypeStart,
				// 	Margin: linebot.FlexComponentMarginTypeXs,
				// },
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   ".",
				// 	Weight: linebot.FlexTextWeightTypeRegular,
				// 	Size:   linebot.FlexTextSizeTypeSm,
				// 	Color:  "#212121",
				// 	Align:  linebot.FlexComponentAlignTypeStart,
				// 	Margin: linebot.FlexComponentMarginTypeLg,
				// },
				// ข้อความเลือกกิจกรรม
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลือกกิจกรรมที่คุณต้องการเพิ่ม:",
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					// Margin: linebot.FlexComponentMarginTypeXxl,
				},
				// ปุ่มกิจกรรม
				createActivityButtons(activities),
			},
			BackgroundColor: "#FFFFFF",
		},
	}

	// สร้าง FlexMessage
	return linebot.NewFlexMessage("รายการกิจกรรม", container)
}

// ฟังก์ชันลบช่องว่างออกจากข้อความ
func normalizeText(input string) string {
	return strings.ReplaceAll(input, " ", "")
}

// flex ปุ่มเริ่มกิจกรรม
func FormatStartActivity(activity string) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เริ่มบันทึกกิจกรรม",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:    linebot.FlexComponentTypeText,
							Text:    `กรุณากดปุ่ม "เริ่มกิจกรรม" เพื่อเริ่มบันทึกเวลา`,
							Size:    linebot.FlexTextSizeTypeSm,
							Align:   linebot.FlexComponentAlignTypeStart,
							Gravity: linebot.FlexComponentGravityTypeCenter,
							// Margin:  linebot.FlexComponentMarginTypeSm,
							Wrap: true,
						},
					},
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Margin:  linebot.FlexComponentMarginTypeXl,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeVertical,
							Spacing: linebot.FlexComponentSpacingTypeMd,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type: linebot.FlexComponentTypeButton,
									Action: &linebot.MessageAction{
										Label: "เริ่มกิจกรรม",
										Text:  "เริ่มกิจกรรม",
									},
									Style: linebot.FlexButtonStyleTypePrimary,
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Margin: linebot.FlexComponentMarginTypeLg,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type: linebot.FlexComponentTypeButton,
									Action: &linebot.MessageAction{
										Label: "ยกเลิก",
										Text:  "ยกเลิก",
									},
									Style: linebot.FlexButtonStyleTypeSecondary,
								},
							},
						},
					},
				},
			},
		},
		Styles: &linebot.BubbleStyle{
			Header: &linebot.BlockStyle{
				BackgroundColor: "#58BDCF",
			},
		},
	}
	return container
}

func FormatactivityRecordStarttime(activityRecord *models.Activityrecord) *linebot.FlexMessage {
	if activityRecord == nil {
		return nil
	}

	// ใช้เวลาเริ่มต้นปัจจุบัน
	// startTime := time.Now().Format("15:04 น.")

	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Size:      linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "กิจกรรมเริ่มต้นแล้ว!",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "กิจกรรม: " + activityRecord.ActivityOther,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Align:  linebot.FlexComponentAlignTypeCenter,
					Color:  "#212121",
					Wrap:   true,
				},
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Margin: linebot.FlexComponentMarginTypeMd,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เวลาเริ่มต้น: " + getCurrentTime(),
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeMd,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  "กรุณากดปุ่ม \"เสร็จสิ้น\" เมื่อทำกิจกรรมเสร็จ",
					Size:  linebot.FlexTextSizeTypeSm,
					Align: linebot.FlexComponentAlignTypeCenter,
					Wrap:  true,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Action: &linebot.MessageAction{Label: "เสร็จสิ้น", Text: "เสร็จสิ้น"},
					Style:  linebot.FlexButtonStyleTypePrimary,
				},
			},
		},
	}

	return linebot.NewFlexMessage("กิจกรรมเริ่มต้นแล้ว", container)
}

// func FormatactivityRecordEndtime(endtime []models.Activityrecord) *linebot.FlexMessage {
// 	if endtime == nil {
// 		return nil
// 	}
// 	// สร้าง BubbleContainer สำหรับข้อความ Flex Message
// 	container := &linebot.BubbleContainer{
// 		Type:      linebot.FlexContainerTypeBubble,
// 		Size:      linebot.FlexBubbleSizeTypeMega,
// 		Direction: linebot.FlexBubbleDirectionTypeLTR,
// 		Header: &linebot.BoxComponent{
// 			Type:   linebot.FlexComponentTypeBox,
// 			Layout: linebot.FlexBoxLayoutTypeVertical,
// 			Contents: []linebot.FlexComponent{
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   fmt.Sprintf("บันทึกกิจกรรม: %s", endtime[0].ServiceInfo.Activity),
// 					Weight: linebot.FlexTextWeightTypeBold,
// 					Size:   linebot.FlexTextSizeTypeLg,
// 					Color:  "#FFFFFF",
// 					Align:  linebot.FlexComponentAlignTypeCenter,
// 				},
// 			},

// 			BackgroundColor: "#00bcd4",
// 		},
// 		Body: &linebot.BoxComponent{
// 			Type:   linebot.FlexComponentTypeBox,
// 			Layout: linebot.FlexBoxLayoutTypeVertical,
// 			Contents: []linebot.FlexComponent{
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   endtime[0].PatientInfo.Name,
// 					Weight: linebot.FlexTextWeightTypeBold,
// 					Size:   linebot.FlexTextSizeTypeMd,
// 					Color:  "#212121",
// 					Align:  linebot.FlexComponentAlignTypeCenter,
// 					Margin: linebot.FlexComponentMarginTypeXs,
// 				},
// 				// เลขประจำตัวประชาชน
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   "เลขประจำตัวประชาชน: " + endtime[0].PatientInfo.CardID,
// 					Weight: linebot.FlexTextWeightTypeRegular,
// 					Size:   linebot.FlexTextSizeTypeSm,
// 					Color:  "#212121",
// 					Align:  linebot.FlexComponentAlignTypeCenter,
// 					Margin: linebot.FlexComponentMarginTypeNone,
// 				},
// 				&linebot.SeparatorComponent{
// 					Type:   linebot.FlexComponentTypeSeparator,
// 					Color:  "#58BDCF",
// 					Margin: linebot.FlexComponentMarginTypeXl,
// 				},
// 				// กิจกรรมที่ทำ
// 				// &linebot.BoxComponent{
// 				// 	Type:   linebot.FlexComponentTypeBox,
// 				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
// 				// 	Contents: []linebot.FlexComponent{
// 				// 		&linebot.TextComponent{
// 				// 			Type:   linebot.FlexComponentTypeText,
// 				// 			Text:   endtime[0].ServiceInfo.Activity,
// 				// 			Weight: linebot.FlexTextWeightTypeBold,
// 				// 			Size:   linebot.FlexTextSizeTypeMd,
// 				// 			Color:  "#212121",
// 				// 			Align:  linebot.FlexComponentAlignTypeCenter,
// 				// 		},
// 				// 	},
// 				// },
// 				// วันที่และเวลา
// 				// &linebot.TextComponent{
// 				// 	Type:   linebot.FlexComponentTypeText,
// 				// 	Text:   "เสร็จสิ้นเมื่อ",
// 				// 	Weight: linebot.FlexTextWeightTypeRegular,
// 				// 	Size:   linebot.FlexTextSizeTypeMd,
// 				// 	Color:  "#212121",
// 				// 	Align:  linebot.FlexComponentAlignTypeCenter,
// 				// 	Margin: linebot.FlexComponentMarginTypeLg,
// 				// },
// 				&linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   "วันที่: " + getCurrentTime(),
// 					Weight: linebot.FlexTextWeightTypeRegular,
// 					Size:   linebot.FlexTextSizeTypeSm,
// 					Color:  "#212121",
// 					Align:  linebot.FlexComponentAlignTypeStart,
// 					Margin: linebot.FlexComponentMarginTypeXl,
// 				}, &linebot.TextComponent{
// 					Type:   linebot.FlexComponentTypeText,
// 					Text:   "ระยะเวลา: " + endtime[0].Period,
// 					Weight: linebot.FlexTextWeightTypeRegular,
// 					Size:   linebot.FlexTextSizeTypeSm,
// 					Color:  "#212121",
// 					Align:  linebot.FlexComponentAlignTypeStart,
// 					Margin: linebot.FlexComponentMarginTypeXs,
// 				},
// 				// createTextRow("วันที่ ", getCurrentTime()),
// 				// createTextRow("ระยะเวลา ", endtime[0].Period),
// 				// &linebot.TextComponent{
// 				// 	Type:   linebot.FlexComponentTypeText,
// 				// 	Text:   "วันที่ " + getCurrentTime(),
// 				// 	Weight: linebot.FlexTextWeightTypeRegular,
// 				// 	Size:   linebot.FlexTextSizeTypeMd,
// 				// 	Color:  "#212121",
// 				// 	Align:  linebot.FlexComponentAlignTypeCenter,
// 				// 	Margin: linebot.FlexComponentMarginTypeNone, // ไม่ต้องการระยะห่าง
// 				// },
// 				&linebot.BoxComponent{
// 					Type:   linebot.FlexComponentTypeBox,
// 					Layout: linebot.FlexBoxLayoutTypeVertical,
// 					Contents: []linebot.FlexComponent{
// 						// &linebot.TextComponent{
// 						// 	Type:   linebot.FlexComponentTypeText,
// 						// 	Text:   "ระยะเวลาในการทำกิจกรรม",
// 						// 	Weight: linebot.FlexTextWeightTypeRegular,
// 						// 	Size:   linebot.FlexTextSizeTypeSm,
// 						// 	Color:  "#212121",
// 						// 	Align:  linebot.FlexComponentAlignTypeCenter,
// 						// 	Margin: linebot.FlexComponentMarginTypeMd,
// 						// },
// 						// &linebot.TextComponent{
// 						// 	Type:   linebot.FlexComponentTypeText,
// 						// 	Text:   endtime[0].Period,
// 						// 	Weight: linebot.FlexTextWeightTypeRegular,
// 						// 	Size:   linebot.FlexTextSizeTypeMd,
// 						// 	Color:  "#212121",
// 						// 	Align:  linebot.FlexComponentAlignTypeCenter,
// 						// 	Margin: linebot.FlexComponentMarginTypeXs,
// 						// },
// 						&linebot.SeparatorComponent{
// 							Type:   linebot.FlexComponentTypeSeparator,
// 							Color:  "#58BDCF",
// 							Margin: linebot.FlexComponentMarginTypeXl,
// 						},
// 						&linebot.TextComponent{
// 							Type:   linebot.FlexComponentTypeText,
// 							Text:   "!!!สำเร็จ!!!",
// 							Weight: linebot.FlexTextWeightTypeBold,
// 							Size:   linebot.FlexTextSizeTypeMd,
// 							Color:  "#212121",
// 							Align:  linebot.FlexComponentAlignTypeCenter,
// 							Margin: linebot.FlexComponentMarginTypeSm,
// 						},
// 					},
// 				},
// 				// // สถานที่
// 				// &linebot.BoxComponent{
// 				// 	Type:   linebot.FlexComponentTypeBox,
// 				// 	Layout: linebot.FlexBoxLayoutTypeHorizontal,
// 				// 	Margin: linebot.FlexComponentMarginTypeLg,
// 				// 	Contents: []linebot.FlexComponent{
// 				// 		&linebot.TextComponent{
// 				// 			Type:   linebot.FlexComponentTypeText,
// 				// 			Text:   "สถานที่:",
// 				// 			Align:  linebot.FlexComponentAlignTypeEnd,
// 				// 		},
// 				// 		&linebot.TextComponent{
// 				// 			Type:   linebot.FlexComponentTypeText,
// 				// 			Text:   "บ้านผู้สูงอายุ",
// 				// 			Align:  linebot.FlexComponentAlignTypeStart,
// 				// 		},
// 				// 	},
// 				// },
// 				// // รูปภาพ
// 				// &linebot.BoxComponent{
// 				// 	Type:   linebot.FlexComponentTypeBox,
// 				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
// 				// 	Contents: []linebot.FlexComponent{
// 				// 		&linebot.ImageComponent{
// 				// 			Type:  linebot.FlexComponentTypeImage,
// 				// 			URL:   "https://www.yanheenursinghome.com/wp-content/uploads/2023/07/ART_0196.jpg",
// 				// 			Align: linebot.FlexComponentAlignTypeCenter,
// 				// 			Size:  linebot.FlexImageSizeTypeFull,
// 				// 		},
// 				// 	},
// 				// },
// 			},
// 		},
// 	}

// 	// สร้างและส่ง Flex Message
// 	return linebot.NewFlexMessage("บันทึกกิจกรรม", container)
// }

func FormatConfirmationCheckIn(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
	if worktimeRecord == nil {
		return nil
	}
	// สร้าง BubbleContainer สำหรับ Flex Message
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ยืนยันการเช็คอิน",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:    linebot.FlexComponentTypeText,
							Text:    worktimeRecord.EmployeeInfo.Name,
							Weight:  linebot.FlexTextWeightTypeBold,
							Size:    linebot.FlexTextSizeTypeMd,
							Color:   "#212121",
							Align:   linebot.FlexComponentAlignTypeCenter,
							Gravity: linebot.FlexComponentGravityTypeTop,
						},
						&linebot.SpacerComponent{}, // เพิ่มช่องว่างระหว่างคอมโพเนนต์
					},
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "รหัสพนักงาน:" + worktimeRecord.EmployeeInfo.EmployeeCode,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeNone,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Margin: linebot.FlexComponentMarginTypeXxl,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Spacing: linebot.FlexComponentSpacingTypeMd,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type: linebot.FlexComponentTypeButton,
									Action: &linebot.MessageAction{
										Label: "เช็คอิน",
										Text:  "เช็คอิน",
									},
									Margin: linebot.FlexComponentMarginTypeXs,
									Height: linebot.FlexButtonHeightTypeMd,
									Style:  linebot.FlexButtonStyleTypePrimary,
								},
							},
						},
					},
				},
				&linebot.SpacerComponent{}, // เพิ่มช่องว่างเพิ่มเติม
			},
		},
	}

	return linebot.NewFlexMessage("ยืนยันการเช็คอิน", container)
}
