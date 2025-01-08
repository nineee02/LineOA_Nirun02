package flex

import (
	"fmt"
	"nirun/pkg/models"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func intPtr(i int) *int {
	return &i
}

func FormatPatientInfo(patient *models.Activityrecord) *linebot.FlexMessage {
	// สร้าง BubbleContainer
	
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ข้อมูลผู้สูงอายุ",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXl,
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
				// ชื่อผู้สูงอายุ
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   patient.PatientInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				// เลขประจำตัวประชาชน
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลขประจำตัวประชาชน: " + patient.PatientInfo.CardID,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeSm,
				},
				// ข้อมูลผู้สูงอายุ
				createTextRow("อายุ", patient.PatientInfo.Age+" ปี"),
				createTextRow("เพศ", formatGender(patient.PatientInfo.Sex)),
				createTextRow("หมู่เลือด", patient.PatientInfo.Blood),
				createTextRow("ADL", patient.PatientInfo.ADL),
				createTextRow("เบอร์", patient.PatientInfo.PhoneNumber),
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				// แสดงสิทธิการรักษา
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "สิทธิการรักษา:",
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeMd,
							Color:  "#000000",
							Margin: linebot.FlexComponentMarginTypeMd,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   patient.PatientInfo.RightToTreatmentInfo.Right_to_treatment,
							Size:   linebot.FlexTextSizeTypeMd,
							Color:  "#212121",
							Align:  linebot.FlexComponentAlignTypeStart,
							Wrap:   true, // เปิดตัดข้อความอัตโนมัติ
							Margin: linebot.FlexComponentMarginTypeSm,
						},
					},
				},
			},
			BackgroundColor: "#EEECEA7A",
		},
	}

	// สร้าง FlexMessage
	return linebot.NewFlexMessage("ข้อมูลผู้สูงอายุ", container)
}

// ฟังก์ชันแปลงเพศ
func formatGender(sex string) string {
	if sex == "Male" {
		return "ชาย"
	} else if sex == "Female" {
		return "หญิง"
	}
	return "ไม่ระบุ"
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

func getCurrentTime() string {
	format := "02-01-2006 03:04:05 PM"
	return time.Now().Format(format)
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
		Hero: &linebot.ImageComponent{
			Type:            linebot.FlexComponentTypeImage,
			URL:             "https://www.okmd.or.th/upload/iblock/82c/aging_850x446-fb.png",
			Size:            linebot.FlexImageSizeTypeFull,
			AspectRatio:     "1.51:1",
			AspectMode:      "cover",
			BackgroundColor: "#FFFFFFFF",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// ชื่อผู้รับบริการ
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   fmt.Sprintf(activity[0].PatientInfo.Name),
					Weight: linebot.FlexTextWeightTypeBold,
					Align:  linebot.FlexComponentAlignTypeCenter,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Wrap:   true,
					Margin: linebot.FlexComponentMarginTypeXs,
				},

				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลขประจำตัวประชาชน: \n" + activity[0].PatientInfo.CardID,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeStart,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   "เลขประจำตัวประชาชน: " + activity[0].PatientInfo.CardID,
				// 	Weight: linebot.FlexTextWeightTypeBold,
				// 	Size:   linebot.FlexTextSizeTypeSm,
				// 	Color:  "#000000",
				// 	Wrap:   true,
				// },
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   ".",
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeStart,
					Margin: linebot.FlexComponentMarginTypeLg,
				},
				// ข้อความเลือกกิจกรรม
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลือกกิจกรรมที่คุณต้องการเพิ่ม:",
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Margin: linebot.FlexComponentMarginTypeXxl,
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

// func CreateTextRow(label string, value string, size string) *linebot.BoxComponent {
// 	return &linebot.BoxComponent{
// 		Type:   linebot.FlexComponentTypeBox,
// 		Layout: linebot.FlexBoxLayoutTypeHorizontal,
// 		Contents: []linebot.FlexComponent{
// 			&linebot.TextComponent{
// 				Type: linebot.FlexComponentTypeText,
// 				Text: label,
// 				Size: linebot.FlexTextSizeTypeSm, // ระบุขนาดของข้อความ
// 				Wrap: true,
// 			},
// 			&linebot.TextComponent{
// 				Type: linebot.FlexComponentTypeText,
// 				Text: value,
// 				Size: linebot.FlexTextSizeTypeSm, // ระบุขนาดของข้อความ
// 				Wrap: true,
// 			},
// 		},
// 	}
// }

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
func FormatStartActivity(activity string) *linebot.BubbleContainer {
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.MessageAction{
						Label: fmt.Sprintf("เริ่มกิจกรรม"),
						Text:  fmt.Sprintf("เริ่มกิจกรรม"),
						// 	Label: fmt.Sprintf("เริ่มกิจกรรม: %s", activity),
						// 	Text:  fmt.Sprintf("เริ่มกิจกรรม: %s", activity),
					},
					Margin: linebot.FlexComponentMarginTypeLg,
					Height: linebot.FlexButtonHeightTypeSm,
					Style:  linebot.FlexButtonStyleTypePrimary,
				},
			},
		},
	}
	return container
}

func FormatactivityRecordStarttime(activityRecord *models.Activityrecord) *linebot.FlexMessage {
	if activityRecord == nil {
		return nil
	}
	// สร้าง BubbleContainer สำหรับ Flex Message
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
					Text:   fmt.Sprintf("เริ่มบันทึกกิจกรรม"),
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   fmt.Sprintf(activityRecord.ServiceInfo.Activity),
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#00bcd4",
		},
		// Header: &linebot.BoxComponent{
		// 	Type:   linebot.FlexComponentTypeBox,
		// 	Layout: linebot.FlexBoxLayoutTypeVertical,
		// 	Contents: []linebot.FlexComponent{
		// 		&linebot.TextComponent{
		// 			Type:   linebot.FlexComponentTypeText,
		// 			Text:   fmt.Sprintf("บันทึกกิจกรรม: %s", endtime[0].ServiceInfo.Activity),
		// 			Weight: linebot.FlexTextWeightTypeBold,
		// 			Size:   linebot.FlexTextSizeTypeLg,
		// 			Color:  "#F8F8FF",
		// 			Align:  linebot.FlexComponentAlignTypeCenter,
		// 		},
		// 	},
		// 	BackgroundColor: "#00bcd4",
		// },
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.TextComponent{
				// 			Type:   linebot.FlexComponentTypeText,
				// 			Text:   "ซุโดกุ",
				// 			Weight: linebot.FlexTextWeightTypeBold,
				// 			Size:   linebot.FlexTextSizeTypeLg,
				// 			Align:  linebot.FlexComponentAlignTypeCenter,
				// 			Margin: linebot.FlexComponentMarginTypeSm,
				// 		},
				// 	},
				// },
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
				// 	Margin: linebot.FlexComponentMarginTypeSm,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.TextComponent{
				// 			Type:   linebot.FlexComponentTypeText,
				// 			Text:   "นางทองสุก สุขศรี",
				// 			Weight: linebot.FlexTextWeightTypeBold,
				// 			Size:   linebot.FlexTextSizeTypeLg,
				// 			Align:  linebot.FlexComponentAlignTypeCenter,
				// 		},
				// 	},
				// },
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Margin: linebot.FlexComponentMarginTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical, // ใช้ Vertical Layout
							Margin: linebot.FlexComponentMarginTypeMd,
							Contents: []linebot.FlexComponent{
								// ข้อความ "เริ่มบันทึกกิจกรรมที่:"
								// &linebot.TextComponent{
								// 	Type:   linebot.FlexComponentTypeText,
								// 	Text:   "เริ่มบันทึกกิจกรรม",
								// 	Weight: linebot.FlexTextWeightTypeRegular,
								// 	Size:   linebot.FlexTextSizeTypeMd,
								// 	Color:  "#212121",
								// 	Align:  linebot.FlexComponentAlignTypeCenter,
								// 	Margin: linebot.FlexComponentMarginTypeSm,
								// },
								&linebot.TextComponent{
									Type:   linebot.FlexComponentTypeText,
									Text:   "เริ่มที่ " + getCurrentTime(),
									Weight: linebot.FlexTextWeightTypeRegular,
									Size:   linebot.FlexTextSizeTypeMd,
									Color:  "#212121",
									Align:  linebot.FlexComponentAlignTypeCenter,
									Margin: linebot.FlexComponentMarginTypeNone, // ไม่ต้องการระยะห่าง
								},
							},
						},
						// &linebot.TextComponent{
						// 	Type:   linebot.FlexComponentTypeText,
						// 	Text:   "เริ่ม:",
						// 	Align:  linebot.FlexComponentAlignTypeEnd,
						// },
						// &linebot.TextComponent{
						// 	Type: linebot.FlexComponentTypeText,
						// 	Text: "25-02-02 10:00:00",
						// 	Flex: linebot.IntPtr(2),
						// },
					},
				},
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
				// 	Margin: linebot.FlexComponentMarginTypeMd,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.TextComponent{
				// 			Type:       linebot.FlexComponentTypeText,
				// 			Text:       "กรุณากดปุ่ม \"เสร็จสิ้น\" เมื่อทำกิจกรรมเสร็จ",
				// 			Margin:     linebot.FlexComponentMarginTypeSm,
				// 			Wrap:       true,
				// 			Decoration: linebot.FlexTextDecorationTypeUnderline,
				// 		},
				// 	},
				// },
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Margin: linebot.FlexComponentMarginTypeXs,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "เสร็จสิ้น",
								Text:  "เสร็จสิ้น",
							},
							Margin: linebot.FlexComponentMarginTypeXs,
							Height: linebot.FlexButtonHeightTypeMd,
							Style:  linebot.FlexButtonStyleTypePrimary,
						},
					},
				},
			},
		},
	}
	// เพิ่ม Quick Reply ใน Flex Message
	// quickReply := linebot.NewQuickReplyItems(
	// 	linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก")),
	// 	linebot.NewQuickReplyButton("", linebot.NewMessageAction("เสร็จสิ้น", "เสร็จสิ้น")),
	// )

	// สร้าง Flex Message พร้อม Quick Reply
	flexMessage := linebot.NewFlexMessage("เริ่มบันทึกกิจกรรม", container)
	// flexMessage = flexMessage.WithQuickReplies(quickReply).(*linebot.FlexMessage)

	return flexMessage

	// สร้าง Flex Message
	// return linebot.NewFlexMessage("เริ่มบันทึกกิจกรรม", container)
}
func FormatactivityRecordEndtime(endtime []models.Activityrecord) *linebot.FlexMessage {
	if endtime == nil {
		return nil
	}
	// สร้าง BubbleContainer สำหรับข้อความ Flex Message
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
					Text:   fmt.Sprintf("บันทึกกิจกรรม: %s", endtime[0].ServiceInfo.Activity),
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
					Text:   endtime[0].PatientInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				// เลขประจำตัวประชาชน
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลขประจำตัวประชาชน: " + endtime[0].PatientInfo.CardID,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeNone,
				},
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				// กิจกรรมที่ทำ
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.TextComponent{
				// 			Type:   linebot.FlexComponentTypeText,
				// 			Text:   endtime[0].ServiceInfo.Activity,
				// 			Weight: linebot.FlexTextWeightTypeBold,
				// 			Size:   linebot.FlexTextSizeTypeMd,
				// 			Color:  "#212121",
				// 			Align:  linebot.FlexComponentAlignTypeCenter,
				// 		},
				// 	},
				// },
				// วันที่และเวลา
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   "เสร็จสิ้นเมื่อ",
				// 	Weight: linebot.FlexTextWeightTypeRegular,
				// 	Size:   linebot.FlexTextSizeTypeMd,
				// 	Color:  "#212121",
				// 	Align:  linebot.FlexComponentAlignTypeCenter,
				// 	Margin: linebot.FlexComponentMarginTypeLg,
				// },
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "วันที่: " + getCurrentTime(),
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeStart,
					Margin: linebot.FlexComponentMarginTypeXl,
				}, &linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ระยะเวลา: " + endtime[0].Period,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeStart,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				// createTextRow("วันที่ ", getCurrentTime()),
				// createTextRow("ระยะเวลา ", endtime[0].Period),
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   "วันที่ " + getCurrentTime(),
				// 	Weight: linebot.FlexTextWeightTypeRegular,
				// 	Size:   linebot.FlexTextSizeTypeMd,
				// 	Color:  "#212121",
				// 	Align:  linebot.FlexComponentAlignTypeCenter,
				// 	Margin: linebot.FlexComponentMarginTypeNone, // ไม่ต้องการระยะห่าง
				// },
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						// &linebot.TextComponent{
						// 	Type:   linebot.FlexComponentTypeText,
						// 	Text:   "ระยะเวลาในการทำกิจกรรม",
						// 	Weight: linebot.FlexTextWeightTypeRegular,
						// 	Size:   linebot.FlexTextSizeTypeSm,
						// 	Color:  "#212121",
						// 	Align:  linebot.FlexComponentAlignTypeCenter,
						// 	Margin: linebot.FlexComponentMarginTypeMd,
						// },
						// &linebot.TextComponent{
						// 	Type:   linebot.FlexComponentTypeText,
						// 	Text:   endtime[0].Period,
						// 	Weight: linebot.FlexTextWeightTypeRegular,
						// 	Size:   linebot.FlexTextSizeTypeMd,
						// 	Color:  "#212121",
						// 	Align:  linebot.FlexComponentAlignTypeCenter,
						// 	Margin: linebot.FlexComponentMarginTypeXs,
						// },
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Color:  "#58BDCF",
							Margin: linebot.FlexComponentMarginTypeXl,
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "!!!สำเร็จ!!!",
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeMd,
							Color:  "#212121",
							Align:  linebot.FlexComponentAlignTypeCenter,
							Margin: linebot.FlexComponentMarginTypeSm,
						},
					},
				},
				// // สถานที่
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeHorizontal,
				// 	Margin: linebot.FlexComponentMarginTypeLg,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.TextComponent{
				// 			Type:   linebot.FlexComponentTypeText,
				// 			Text:   "สถานที่:",
				// 			Align:  linebot.FlexComponentAlignTypeEnd,
				// 		},
				// 		&linebot.TextComponent{
				// 			Type:   linebot.FlexComponentTypeText,
				// 			Text:   "บ้านผู้สูงอายุ",
				// 			Align:  linebot.FlexComponentAlignTypeStart,
				// 		},
				// 	},
				// },
				// // รูปภาพ
				// &linebot.BoxComponent{
				// 	Type:   linebot.FlexComponentTypeBox,
				// 	Layout: linebot.FlexBoxLayoutTypeVertical,
				// 	Contents: []linebot.FlexComponent{
				// 		&linebot.ImageComponent{
				// 			Type:  linebot.FlexComponentTypeImage,
				// 			URL:   "https://www.yanheenursinghome.com/wp-content/uploads/2023/07/ART_0196.jpg",
				// 			Align: linebot.FlexComponentAlignTypeCenter,
				// 			Size:  linebot.FlexImageSizeTypeFull,
				// 		},
				// 	},
				// },
			},
		},
	}

	// สร้างและส่ง Flex Message
	return linebot.NewFlexMessage("บันทึกกิจกรรม", container)
}

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
					Text:   worktimeRecord.EmployeeInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				// รหัสพนักงาน
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "รหัสพนักงาน: " + worktimeRecord.EmployeeInfo.EmployeeCode,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeNone,
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

				// createTextRow("เช็คอินที่ ", getCurrentTime()),

				// เช็คอิน
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   "เช็คอินที่",
				// 	Weight: linebot.FlexTextWeightTypeBold,
				// 	Size:   linebot.FlexTextSizeTypeMd,
				// 	Color:  "#000000",
				// 	Margin: linebot.FlexComponentMarginTypeMd,
				// 	Align:  linebot.FlexComponentAlignTypeCenter,
				// },
				// &linebot.TextComponent{
				// 	Type:  linebot.FlexComponentTypeText,
				// 	Text:  getCurrentTime(),
				// 	Size:  linebot.FlexTextSizeTypeMd,
				// 	Color: "#212121",
				// 	Align: linebot.FlexComponentAlignTypeCenter,
				// 	// Wrap:   true, // เปิดตัดข้อความอัตโนมัติ
				// 	// Margin: linebot.FlexComponentMarginTypeSm,
				// 	Margin: linebot.FlexComponentMarginTypeNone,
				// },
			},
		},
	}

	// สร้าง Flex Message
	return linebot.NewFlexMessage("ลงเวลาเข้างาน", container)
}

func FormatConfirmationCheckOut(worktimeRecord *models.WorktimeRecord) *linebot.FlexMessage {
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
					Text:   "ยืนยันการเช็คเอ้าท์",
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
					Text:   "รหัสพนักงาน: " + worktimeRecord.EmployeeInfo.EmployeeCode,
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
										Label: "เช็คเอ้าท์",
										Text:  "เช็คเอ้าท์",
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

	return linebot.NewFlexMessage("ยืนยันการเช็คเอ้าท์", container)
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
					Text:   "ขอบคุณที่ใช้บริการ",
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
					Text:   worktimeRecord.EmployeeInfo.Name,
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "รหัสพนักงาน: " + worktimeRecord.EmployeeInfo.EmployeeCode,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#212121",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Margin: linebot.FlexComponentMarginTypeNone,
				},
				// // รหัสพนักงาน
				// &linebot.TextComponent{
				// 	Type:   linebot.FlexComponentTypeText,
				// 	Text:   worktimeRecord.EmployeeInfo.EmployeeCode,
				// 	Weight: linebot.FlexTextWeightTypeRegular,
				// 	Size:   linebot.FlexTextSizeTypeSm,
				// 	Color:  "#212121",
				// 	Align:  linebot.FlexComponentAlignTypeCenter,
				// 	Margin: linebot.FlexComponentMarginTypeNone,
				// },
				// เช็คอิน
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#58BDCF",
					Margin: linebot.FlexComponentMarginTypeXl,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เช็คเอ้าท์ที่",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeMd,
					Color:  "#000000",
					Margin: linebot.FlexComponentMarginTypeMd,
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  getCurrentTime(),
					Size:  linebot.FlexTextSizeTypeMd,
					Color: "#212121",
					Align: linebot.FlexComponentAlignTypeCenter,
					// Wrap:   true, // เปิดตัดข้อความอัตโนมัติ
					// Margin: linebot.FlexComponentMarginTypeSm,
					Margin: linebot.FlexComponentMarginTypeNone,
				},
			},
		},
	}

	// สร้าง Flex Message
	return linebot.NewFlexMessage("ลงเวลาออกงาน", container)
}
func FormatHistoryType() *linebot.FlexMessage {
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// เพิ่มข้อความเงาซ้อน (เงาสีเทาเข้ม)
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "กรุณาเลือกประเภทการดูประวัติ",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#4a4a4a", // เงาสีเทาเข้ม
					Align:  linebot.FlexComponentAlignTypeCenter,
					Wrap:   true,
				},
				// ข้อความหลัก (สีขาว)
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "กรุณาเลือกประเภทการดูประวัติ",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#FFFFFF",
					Align:  linebot.FlexComponentAlignTypeCenter,
					Wrap:   true,
				},
			},
			BackgroundColor: "#00bcd4",
			// PaddingAll:      "md",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				// แถวที่ 1 (ปุ่ม "ทั้งหมด" และ "ปีนี้")
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "ทั้งหมด",
								Text:  "ทั้งหมด",
							},
							Color:  "#58bdcf",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "ปีนี้",
								Text:  "ปีนี้",
							},
							Color:  "#fecd96",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
					},
					Margin: linebot.FlexComponentMarginTypeMd,
				},
				// แถวที่ 2 (ปุ่ม "เดือนนี้" และ "สัปดาห์นี้")
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "เดือนนี้",
								Text:  "เดือนนี้",
							},
							Color:  "#fcfe96",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "สัปดาห์นี้",
								Text:  "สัปดาห์นี้",
							},
							Color:  "#58bdcf",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
					},
					Margin: linebot.FlexComponentMarginTypeMd,
				},
				// แถวที่ 3 (ปุ่ม "วันนี้" และ "ระบุช่วงเวลา")
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "วันนี้",
								Text:  "วันนี้",
							},
							Color:  "#fecd96",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
						&linebot.ButtonComponent{
							Type: linebot.FlexComponentTypeButton,
							Action: &linebot.MessageAction{
								Label: "ระบุช่วงเวลา",
								Text:  "ระบุช่วงเวลา",
							},
							Color:  "#fcfe96",
							Style:  linebot.FlexButtonStyleTypeSecondary,
							Height: linebot.FlexButtonHeightTypeSm,
						},
					},
					Margin: linebot.FlexComponentMarginTypeMd,
				},
			},
		},
	}

	// สร้าง Flex Message
	return linebot.NewFlexMessage("ประเภทการดูประวัติ", container)
}
