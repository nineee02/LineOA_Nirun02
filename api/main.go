package main

import (
	"nirun/pkg/hook"
	"nirun/pkg/linebot"

	"github.com/gin-gonic/gin"
)

// func generateQRCode() {
// 	url := "https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=2006878417&redirect_uri=http%3A%2F%2Fcommunity.app.nirun.life%2Fauth_oauth%2Fsignin&state=random_string&scope=profile+openid+email"
// 	err := qrcode.WriteFile(url, qrcode.Medium, 256, "qrcode_line_login.png")
// 	if err != nil {
// 		log.Fatalf("Failed to generate QR Code: %v", err)
// 	}
// 	// log.Println("QR Code generated successfully.")
// }

func main() {
	// สร้าง QR Code สำหรับ URL
	// generateQRCode()

	// รันเซิร์ฟเวอร์ของ LINE Bot
	// (โค้ดที่เหลือของคุณสำหรับการตั้งค่า Gin และเซิร์ฟเวอร์)
	linebot.InitLineBot()

	router := gin.Default()

	// router.GET("/login", event.LineLoginHandler)
	// router.GET("/auth_oauth/signin", event.LineLoginCallback)
	// router.GET("/callback", event.LineLoginCallback)
	router.POST("/webhook", hook.HandleLineWebhook)

	router.Run(":8080")
}
