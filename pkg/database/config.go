package database

import (
	"database/sql"
	"fmt"
	"log"
	"nirun/pkg/models"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v2"
)

// ConnectToDB เชื่อมต่อกับฐานข้อมูลโดยใช้ข้อมูลจาก config.yaml
func ConnectToDB() (*sql.DB, error) {
	// โหลด config จากไฟล์ config.yaml
	var config models.Config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Error reading config file:", err)
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error parsing config file:", err)
		return nil, err
	}

	// รูปแบบการเชื่อมต่อฐานข้อมูล (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Name,
	)

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ฟังก์ชัน LoadConfig เพื่อโหลดค่าการตั้งค่าจาก config.yaml
func LoadConfig() (models.Config, error) {
	var config models.Config

	// อ่านไฟล์ config.yaml
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Println("Error reading config file:", err)
		return config, err
	}

	// แปลงข้อมูลจาก YAML เป็น struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Error parsing config file:", err)
		return config, err
	}

	return config, nil
}

func ConnectToMinio() (*minio.Client, error) {
	// โหลดค่าการตั้งค่าจาก config.yaml
	config, err := LoadConfig()
	if err != nil {
		log.Println("Error loading config for MinIO:", err)
		return nil, err
	}

	// สร้าง MinIO client
	minioClient, err := minio.New(config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.AccessKey, config.Minio.SecretKey, ""),
		Secure: config.Minio.UseSSL,
	})
	if err != nil {
		log.Println("Error connecting to MinIO:", err)
		return nil, err
	}
	log.Printf("MinIO Credentials Loaded: AccessKey=%s, SecretKey=%s", config.Minio.AccessKey, config.Minio.SecretKey)
	log.Println("Successfully connected to MinIO")
	return minioClient, nil
}
