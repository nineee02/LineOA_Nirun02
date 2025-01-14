package main

import (
	"fmt"
	"log"
	"nirun/pkg/database" // ถ้าฟังก์ชัน HandleLineWebhook อยู่ใน pkg/hook
	"nirun/pkg/event"
	"nirun/pkg/hook"
	"nirun/pkg/linebot"

	"github.com/gin-gonic/gin"
)

func main() {
	// เรียกใช้ InitLineBot เพื่อสร้าง instance ของ LINE Bot
	linebot.InitLineBot()

	router := gin.Default()

	router.POST("/webhook", hook.HandleLineWebhook)

	// เส้นทางสำหรับ LINE Login
	router.GET("/login", event.LineLoginHandler)

	// เส้นทางสำหรับ Callback
	router.GET("/callback", event.LineLoginCallback)

	// รันเซิร์ฟเวอร์ที่พอร์ต 8080
	router.Run(":8080")
	// เชื่อมต่อฐานข้อมูล
	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Connected to database successfully!")

	// รันเซิร์ฟเวอร์ที่พอร์ต 8080
	router.Run(":8080")
}
