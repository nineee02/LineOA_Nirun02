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
	redirectURI := "https://eab3-49-237-5-107.ngrok-free.app/callback" // เปลี่ยนลิงก์ให้ตรงกับ ngrok ของคุณ
	state := "random_string"
	scope := "profile openid email"

	lineLoginURL := fmt.Sprintf(
		"https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		clientID, redirectURI, state, scope,
	)
	log.Println("lineloginURL", lineLoginURL)

	c.Redirect(http.StatusFound, lineLoginURL)
}

func LineLoginCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	token, err := exchangeToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		log.Printf("Error exchanging token: %v", err)
		return
	}

	profile, err := getProfile(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		log.Printf("Error getting profile: %v", err)
		return
	}
	log.Println("proflie")

	// ตอบกลับด้วยข้อมูลโปรไฟล์
	c.JSON(http.StatusOK, gin.H{
		"userID":      profile.UserID,
		"displayName": profile.DisplayName,
		"pictureUrl":  profile.PictureURL,
		"email":       profile.Email,
	})
}

const (
	clientID     = "2006767645"
	clientSecret = "68fd27f357fe6cc1c6ea782f1cb9819c"
	redirectURI  = "https://4fe3-110-164-198-127.ngrok-free.app/callback"
	state        = "random_string"
	scope        = "profile openid email"
)    

// ฟังก์ชันแลกเปลี่ยน Token
func exchangeToken(code string) (*models.LineTokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://4fe3-110-164-198-127.ngrok-free.app/callback"},
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

