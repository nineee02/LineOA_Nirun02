package flexmessage

import (
	"nirun/pkg/models"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func getCurrentTime() string {
	format := "02-01-2006 03:04 PM"
	return time.Now().Format(format)
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
