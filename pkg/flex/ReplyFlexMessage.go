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
					Color:  "#F8F8FF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#58BDCF",
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
					Color:  "#000000",
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

func FormatactivityRecordStarttime() *linebot.FlexMessage {
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
					Text:   "เริ่มกิจกรรม",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#F8F8FF",
					Align:  linebot.FlexComponentAlignTypeCenter,
				},
			},
			BackgroundColor: "#58BDCF",
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
							Type:   linebot.FlexComponentTypeText,
							Text:   "ซุโดกุ",
							Weight: linebot.FlexTextWeightTypeBold,
							Size:   linebot.FlexTextSizeTypeLg,
							Align:  linebot.FlexComponentAlignTypeCenter,
						},
					},
				},
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
								&linebot.TextComponent{
									Type:   linebot.FlexComponentTypeText,
									Text:   "เริ่มบันทึกกิจกรรมที่:",
									Weight: linebot.FlexTextWeightTypeRegular,
									Size:   linebot.FlexTextSizeTypeSm,
									Color:  "#212121",
									Align:  linebot.FlexComponentAlignTypeCenter,
									Margin: linebot.FlexComponentMarginTypeSm,
								},
								&linebot.TextComponent{
									Type:   linebot.FlexComponentTypeText,
									Text:   getCurrentTime(),
									Weight: linebot.FlexTextWeightTypeRegular,
									Size:   linebot.FlexTextSizeTypeSm,
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
	quickReply := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("ยกเลิก", "ยกเลิก")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("เสร็จสิ้น", "เสร็จสิ้น")),
	)

	// สร้าง Flex Message พร้อม Quick Reply
	flexMessage := linebot.NewFlexMessage("เริ่มบันทึกกิจกรรม", container)
	flexMessage = flexMessage.WithQuickReplies(quickReply).(*linebot.FlexMessage)

	return flexMessage

	// สร้าง Flex Message
	// return linebot.NewFlexMessage("เริ่มบันทึกกิจกรรม", container)
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

		{"ฝังเข็ม", "ฝังเข็ม", "#FFA726"},
		{"คาราโอเกะ", "คาราโอเกะ", "#FFA726"},

		{"ครอบแก้ว", "ครอบแก้ว", "#4DD0E1"},
		{"ทำอาหาร", "ทำอาหาร", "#4DD0E1"},

		{"นั่งสมาธิ", "นั่งสมาธิ", "#FFA726"},
		{"เล่าสู่กันฟัง", "เล่าสู่กันฟัง", "#FFA726"},

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
					Size:   linebot.FlexTextSizeTypeLg,
					Color:  "#000000",
					Wrap:   true,
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลขประจำตัวประชาชน: \n" + activity[0].PatientInfo.CardID,
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#000000",
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

				// ข้อความเลือกกิจกรรม
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "เลือกกิจกรรมที่คุณต้องการเพิ่ม:",
					Weight: linebot.FlexTextWeightTypeRegular,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#000011",
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
