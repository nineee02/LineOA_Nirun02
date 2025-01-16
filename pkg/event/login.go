package event

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"nirun/pkg/models"

	"github.com/gin-gonic/gin"
)

func LineLoginHandler(c *gin.Context) {
	clientID := "2006767645"
	redirectURI := "https://dc3a-49-237-19-181.ngrok-free.app/callback"
	state := "random_string"
	scope := "profile openid email"

	// สร้าง URL สำหรับ Line Login
	lineLoginURL := fmt.Sprintf(
		"https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		clientID, redirectURI, state, scope,
	)
	log.Println("Line Login URL:", lineLoginURL)

	// Redirect ไปยัง URL ที่สร้างขึ้น
	c.Redirect(http.StatusFound, lineLoginURL)
}

func LineLoginCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// แลกเปลี่ยน code เป็น access token
	token, err := exchangeToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		log.Printf("Error exchanging token: %v", err)
		return
	}

	// ดึงข้อมูลโปรไฟล์ผู้ใช้
	profile, err := getProfile(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		log.Printf("Error getting profile: %v", err)
		return
	}

	log.Printf("User Profile: %+v", profile)

	// บันทึกข้อมูลผู้ใช้ลงในฐานข้อมูล
	// err = saveUserToDatabase(profile.UserID, profile.DisplayName, profile.Email)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
	//     log.Printf("Error saving user to database: %v", err)
	//     return
	// }

	// Redirect ไปยังหน้าการเพิ่มเพื่อนหลังจากบันทึกสำเร็จ
	addFriendURL := "https://line.me/R/ti/p//@392avxhp"
	c.Redirect(http.StatusFound, addFriendURL)
}

const (
	clientID     = "2006767645"
	clientSecret = "68fd27f357fe6cc1c6ea782f1cb9819c"
	redirectURI  = "https://dc3a-49-237-19-181.ngrok-free.app/callback"
	state        = "random_string"
	scope        = "profile openid email"
)

// ฟังก์ชันแลกเปลี่ยน Token
func exchangeToken(code string) (*models.LineTokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://dc3a-49-237-19-181.ngrok-free.app/callback"},
		"client_id":     {"2006767645"},
		"client_secret": {"68fd27f357fe6cc1c6ea782f1cb9819c"},
	}

	resp, err := http.PostForm("https://api.line.me/oauth2/v2.1/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse models.LineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

// ฟังก์ชันดึงข้อมูลโปรไฟล์ผู้ใช้
func getProfile(accessToken string) (*models.LineProfile, error) {
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	var profile models.LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
