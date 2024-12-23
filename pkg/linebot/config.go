// pkg/linebot/config.go
package linebot

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"gopkg.in/yaml.v2"
)

// Config โครงสร้างข้อมูลสำหรับการอ่านค่าจากไฟล์ config.yaml
type Config struct {
	LineBot struct {
		ChannelSecret string `yaml:"LINE_CHANNEL_SECRET"`
		ChannelToken  string `yaml:"LINE_CHANNEL_ACCESS_TOKEN"`
	} `yaml:"line_bot"`
}

// ตัวแปร global สำหรับเก็บ instance ของ LINE Bot
var lineBot *linebot.Client

// LoadConfig ฟังก์ชันอ่านไฟล์ YAML
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// InitLineBot ตั้งค่า LINE Bot โดยใช้ค่า ChannelSecret และ ChannelToken จากไฟล์ config.yaml
func InitLineBot() {
	// โหลดค่า config จากไฟล์ config.yaml
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config file:", err)
	}

	// ใช้ค่า Channel Secret และ Channel Token จาก config พร้อมตั้งค่า Endpoint
	bot, err := linebot.New(
		config.LineBot.ChannelSecret,
		config.LineBot.ChannelToken,
		// linebot.WithEndpointBase("https://api.line.me/v2/bot"), // ระบุ Endpoint โดยตรง
	)
	if err != nil {
		log.Fatal("Error initializing LINE Bot:", err)
	}

	lineBot = bot
	fmt.Println("LINE Bot initialized successfully")
}

// GetLineBot returns the LINE Bot client instance
func GetLineBot() *linebot.Client {
	return lineBot
}
