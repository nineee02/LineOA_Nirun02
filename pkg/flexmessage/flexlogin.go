package flexmessage

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func SendRegisterLink(bot *linebot.Client, replyToken string) {
	// URL สำหรับลงทะเบียน
	registerURL := "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=2006767645&redirect_uri=https://dc3a-49-237-19-181.ngrok-free.app/callback&state=random_string&scope=profile%20openid%20email"

	// สร้าง Flex Message
	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: "vertical", // ใช้ string แทน
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   "text",
					Text:   "ลงทะเบียน",
					Size:   "lg",
					Weight: "bold",
					Color:  "#000000",
				},
				&linebot.TextComponent{
					Type:  "text",
					Text:  "กดปุ่มด้านล่างเพื่อลงทะเบียนเข้าสู่ระบบ",
					Size:  "md",
					Color: "#555555",
					Wrap:  true,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: "vertical", // ใช้ string แทน
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Action: linebot.NewURIAction("ลงทะเบียน", registerURL),
					Style:  "primary",
					Color:  "#00B900", // สีเขียว
				},
			},
		},
	}

	// ส่งข้อความ
	flexMessage := linebot.NewFlexMessage("ลงทะเบียน", flexContainer)
	if _, err := bot.ReplyMessage(replyToken, flexMessage).Do(); err != nil {
		log.Printf("Error sending Flex Message: %v", err)
	}
}
