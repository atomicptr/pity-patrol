package endfield

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/atomicptr/pity-patrol/pkg/config"
	"github.com/atomicptr/pity-patrol/pkg/constants"
	"github.com/atomicptr/pity-patrol/pkg/report"
)

const (
	baseURL    = "https://zonai.skport.com"
	claimURL   = "/web/v1/game/endfield/attendance"
	refreshUrl = "/web/v1/auth/refresh"
	platform   = "3"
	vName      = "1.0.0"
)

type attendanceResource struct {
	Id    string `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
}

type attendanceResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Calendar []struct {
			AwardId   string `json:"awardId"`
			Available bool   `json:"available"`
			Done      bool   `json:"done"`
		} `json:"calendar"`
		ResourceInfoMap map[string]attendanceResource
	} `json:"data"`
}

func (res *attendanceResponse) findReward() *attendanceResource {
	var lastId string

	for _, entry := range res.Data.Calendar {
		if !entry.Done {
			break
		}

		lastId = entry.AwardId
	}

	if lastId == "" {
		return nil
	}

	resource, ok := res.Data.ResourceInfoMap[lastId]
	if !ok {
		return nil
	}

	return &resource
}

func Claim(cfg *config.Config, account *config.Account) (*report.Report, error) {
	ua := cfg.UserAgent
	if ua == "" {
		ua = constants.UserAgent
	}

	client := http.Client{Timeout: constants.DefaultTimeoutSecs}

	token, err := refreshToken(&client, account.Credentials, ua)

	resp, err := requestAttendance(&client, "POST", ua, token, cfg, account)
	if err != nil {
		return nil, err
	}

	switch resp.Code {
	case 0:
		rep := &report.Report{
			WasClaimed: true,
		}

		result, err := requestAttendance(&client, "GET", ua, token, cfg, account)
		if err != nil {
			// while this is an error we already claimed so just report without rewards
			log.Printf("error while trying to check rewards: %s", err)
			return rep, nil
		}

		if resource := result.findReward(); resource != nil {
			rep.Reward = &report.Reward{
				Name:  resource.Name,
				Count: resource.Count,
				Image: resource.Icon,
			}
		}

		return rep, nil
	case 10001:
		return &report.Report{
			WasClaimed: false,
		}, nil
	default:
		return nil, fmt.Errorf("api error: %s (code %d)", resp.Message, resp.Code)
	}
}

func refreshToken(client *http.Client, credentials, ua string) (string, error) {
	req, _ := http.NewRequest("GET", baseURL+refreshUrl, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("cred", credentials)
	req.Header.Set("platform", platform)
	req.Header.Set("vName", vName)

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
		platform, timestamp, vName)

	dataToSign := path + body + timestamp + headerJSON

	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(dataToSign))
	hmacHex := hex.EncodeToString(mac.Sum(nil))

	hasher := md5.New()
	hasher.Write([]byte(hmacHex))
	return hex.EncodeToString(hasher.Sum(nil))
}

func requestAttendance(client *http.Client, method, ua, token string, cfg *config.Config, account *config.Account) (*attendanceResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := generateSign(claimURL, "", timestamp, token)

	req, err := http.NewRequest(method, baseURL+claimURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", ua)
	req.Header.Set("cred", account.Credentials)
	req.Header.Set("sk-game-role", account.SkGameRole)
	req.Header.Set("platform", platform)
	req.Header.Set("vName", vName)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("sign", sign)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sk-language", "en")

	if cfg.DebugMode {
		log.Printf("[DEBUG] POST %s", baseURL+claimURL)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result attendanceResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
