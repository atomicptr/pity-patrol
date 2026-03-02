package endfield

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/atomicptr/pity-patrol/pkgs/config"
	"github.com/atomicptr/pity-patrol/pkgs/constants"
)

const (
	BaseURL    = "https://zonai.skport.com"
	ClaimURL   = "/web/v1/game/endfield/attendance"
	RefreshURL = "/web/v1/auth/refresh"
	Platform   = "3"
	VName      = "1.0.0"
)

func Claim(cfg *config.Config, account *config.Account) (bool, error) {
	ua := cfg.UserAgent
	if ua == "" {
		ua = constants.UserAgent
	}

	client := http.Client{Timeout: constants.DefaultTimeoutSecs}

	token, err := refreshToken(&client, account.Credentials, ua)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := generateSign(ClaimURL, "", timestamp, token)

	req, err := http.NewRequest("POST", BaseURL+ClaimURL, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("User-Agent", ua)
	req.Header.Set("cred", account.Credentials)
	req.Header.Set("sk-game-role", account.SkGameRole)
	req.Header.Set("platform", Platform)
	req.Header.Set("vName", VName)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("sign", sign)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sk-language", "en")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	switch result.Code {
	case 0:
		return true, nil
	case 10001:
		return false, nil
	default:
		return false, fmt.Errorf("api error: %s (code %d)", result.Message, result.Code)
	}
}

func refreshToken(client *http.Client, credentials, ua string) (string, error) {
	req, _ := http.NewRequest("GET", BaseURL+RefreshURL, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("cred", credentials)
	req.Header.Set("platform", Platform)
	req.Header.Set("vName", VName)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Code int64 `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Code != 0 {
		return "", fmt.Errorf("token refresh failed: %s", res.Message)
	}

	return res.Data.Token, nil
}

func generateSign(path, body, timestamp, token string) string {
	headerJSON := fmt.Sprintf(`{"platform":"%s","timestamp":"%s","dId":"","vName":"%s"}`,
		Platform, timestamp, VName)

	dataToSign := path + body + timestamp + headerJSON

	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(dataToSign))
	hmacHex := hex.EncodeToString(mac.Sum(nil))

	hasher := md5.New()
	hasher.Write([]byte(hmacHex))
	return hex.EncodeToString(hasher.Sum(nil))
}
