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

const (
	// 	clientID     = "2006767645"
	// 	clientSecret = "68fd27f357fe6cc1c6ea782f1cb9819c"
	redirectURI  = "https://beb6-110-164-198-113.ngrok-free.app/callback"
	clientID     = "2006878417"
	clientSecret = "505aefd4da7ba032ab614c30c550164b"
	state        = "random_string"
	scope        = "profile openid email"
)

// LineLoginHandler สร้าง URL สำหรับ Line Login และ Redirect ผู้ใช้
func LineLoginHandler(c *gin.Context) {
	escapedRedirectURI := url.QueryEscape(redirectURI) //ใช้ redirectURI ที่อัปเดตแล้ว

	lineLoginURL := fmt.Sprintf(
		"https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s&prompt=consent",
		clientID, escapedRedirectURI, state, url.QueryEscape(scope),
	)

	log.Println("Line Login URL:", lineLoginURL)
	c.Redirect(http.StatusFound, lineLoginURL)
}

// LineLoginCallback
func LineLoginCallback(c *gin.Context) {
	code := c.Query("code")
	log.Printf("Received code: %s", code) //Debug log
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	token, err := exchangeToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	profile, err := getProfile(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	log.Printf("LINE User ID: %s", profile.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"line_user_id": profile.UserID,
		"display_name": profile.DisplayName,
	})
	//ไปหน้าเพิ่มเพื่อน
	addFriendURL := "https://line.me/R/ti/p//@321hkfbg"
	c.Redirect(http.StatusFound, addFriendURL)
}

// ฟังก์ชันแลกเปลี่ยน Token
func exchangeToken(code string) (*models.LineTokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	// log.Printf("Exchanging code: %s", code)

	resp, err := http.PostForm("https://api.line.me/oauth2/v2.1/token", data)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// ✅ อ่านข้อมูลจาก API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Debug ดู JSON ที่ได้จาก LINE API
	// log.Printf("LINE Token API Response: %s", string(body))

	// ใช้ json.Unmarshal แทนการใช้ json.NewDecoder เพื่อหลีกเลี่ยง EOF
	var tokenResponse models.LineTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil, err
	}

	// log.Printf("Token exchange success: %+v", tokenResponse)
	return &tokenResponse, nil
}

// ฟังก์ชันดึงข้อมูลโปรไฟล์ผู้ใช้
func getProfile(accessToken string) (*models.LineProfile, error) {
	// log.Printf("Fetching profile with Access Token: %s", accessToken)

	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling LINE Profile API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// อ่านข้อมูลจาก API (ครั้งเดียว)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Debug Response
	// log.Printf("LINE Profile API Response: %s", string(body))

	// ตรวจสอบว่า HTTP Status Code เป็น 200 OK หรือไม่
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	//  ใช้ json.Unmarshal แทน json.NewDecoder เพื่อหลีกเลี่ยง EOF
	var profile models.LineProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Printf("❌ Error decoding profile response: %v", err)
		return nil, err
	}

	log.Printf("Successfully received profile: %+v", profile)
	return &profile, nil
}
