package hoyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/atomicptr/pity-patrol/pkgs/config"
	"github.com/atomicptr/pity-patrol/pkgs/constants"
	"github.com/atomicptr/pity-patrol/pkgs/util"
)

const loginBaseUrl = "https://act.hoyolab.com"
const lang = "en-us"

const retCodeSuccess = 0
const retCodeNotLoggedIn = -100
const retCodeAlreadySignedIn = -5003

type gameConfig struct {
	EventBaseUrl string
	ActId        string
	LoginUrl     string
	ExtraHeaders map[string]string
}

var games = map[string]gameConfig{
	"genshin": {
		EventBaseUrl: "https://sg-hk4e-api.mihoyo.com/event/sol",
		ActId:        "e202102251931481",
		LoginUrl:     loginBaseUrl + "/ys/event/signin-sea-v3/index.html?act_id=e202102251931481",
	},
	"starrail": {
		EventBaseUrl: "https://sg-public-api.hoyolab.com/event/luna/os",
		ActId:        "e202303301540311",
		LoginUrl:     loginBaseUrl + "/bbs/event/signin/hkrpg/index.html?act_id=e202303301540311",
	},
	"honkai": {
		EventBaseUrl: "https://sg-public-api.hoyolab.com/event/mani",
		ActId:        "e202110291205111",
		LoginUrl:     loginBaseUrl + "/bbs/event/signin-bh3/index.html?act_id=e202110291205111",
	},
	"themis": {
		EventBaseUrl: "https://sg-public-api.hoyolab.com/event/luna/os",
		ActId:        "e202202281857121",
		LoginUrl:     loginBaseUrl + "/bbs/event/signin/nxx/index.html?act_id=e202202281857121",
	},
	"zzz": {
		EventBaseUrl: "https://sg-public-api.hoyolab.com/event/luna/zzz/os",
		ActId:        "e202406031448091",
		LoginUrl:     loginBaseUrl + "/bbs/event/signin/zzz/index.html?act_id=e202406031448091",
		ExtraHeaders: map[string]string{
			"x-rpc-signgame": "zzz",
		},
	},
}

func Claim(cfg *config.Config, account *config.Account) (bool, error) {
	data, ok := games[account.Type]
	if !ok {
		return false, fmt.Errorf("unknown hoyo game: %s", account.Type)
	}

	client := http.Client{Timeout: constants.DefaultTimeoutSecs}

	info, err := getSignInInfo(&client, cfg, account, &data)
	if err != nil {
		return false, err
	}

	// already signed in
	if info.IsSign {
		return false, nil
	}

	if info.FirstBind {
		return false, fmt.Errorf("account hasn't signed in yet, please sign in manually at least once")
	}

	totalSignInDay := info.TotalSignDay

	awards, err := getAwards(&client, cfg, account, &data)
	if err != nil {
		return false, err
	}

	if cfg.DebugMode {
		log.Printf("[DEBUG] Hoyo: Checking in account for day %s...", info.Today)
	}

	res, err := performCheckIn(&client, cfg, account, &data)
	if err != nil {
		return false, err
	}

	switch res.RetCode {
	case retCodeSuccess:
		// success! Do nothing
	case retCodeNotLoggedIn:
		return false, fmt.Errorf("you're not logged in, please log into %s and replace the cookie in the config with a new value", data.LoginUrl)
	case retCodeAlreadySignedIn:
		return false, nil
	default:
		return false, fmt.Errorf("server returned an error: %s", res.Message)
	}

	log.Println(util.ToPrettyString(res))

	info, err = getSignInInfo(&client, cfg, account, &data)
	if err != nil {
		return false, err
	}

	newTotalSignInDay := info.TotalSignDay

	// sign in wasnt successful
	if !info.IsSign || newTotalSignInDay == totalSignInDay {
		return false, fmt.Errorf("could not automatically check-in, please sign into the website manually: %s", data.LoginUrl)
	}

	if isCaptchaRequired(res) {
		return false, fmt.Errorf("captcha is required, please sign into the website: %s", data.LoginUrl)
	}

	reward := awards[newTotalSignInDay-1]

	if cfg.DebugMode {
		log.Printf("[DEBUG] Total Sign-In Days: %d\n", newTotalSignInDay)
		log.Printf("[DEBUG] Reward: %dx %s\n", reward.Count, reward.Name)
		log.Printf("[DEBUG] Message: %s\n", res.Message)
	}

	return true, nil
}

type captchaData struct {
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
	Success   int    `json:"success"`
}

type hoyoResponse struct {
	RetCode  int             `json:"retcode"`
	Message  string          `json:"message"`
	Data     json.RawMessage `json:"data"`
	GtResult *captchaData    `json:"gt_result,omitempty"`
}

type signInInfo struct {
	Today        string `json:"today"`
	TotalSignDay int    `json:"total_sign_day"`
	IsSign       bool   `json:"is_sign"`
	FirstBind    bool   `json:"first_bind"`
}

func hoyoRequest(client *http.Client, method string, url string, body []byte, cfg *config.Config, account *config.Account, data *gameConfig) (*hoyoResponse, error) {
	if cfg.DebugMode {
		log.Printf("[DEBUG] %s %s\n", method, url)
	}

	var b io.Reader

	if body != nil {
		b = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", loginBaseUrl)
	req.Header.Add("Cookie", account.Cookie)

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	for key, value := range data.ExtraHeaders {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if cfg.DebugMode {
		log.Printf("Response: %s\n", string(blob))
	}

	var hoyoRes hoyoResponse
	err = json.Unmarshal(blob, &hoyoRes)
	if err != nil {
		return nil, err
	}

	return &hoyoRes, nil
}

func getSignInInfo(client *http.Client, cfg *config.Config, account *config.Account, data *gameConfig) (*signInInfo, error) {
	url := fmt.Sprintf("%s/info?lang=%s&act_id=%s", data.EventBaseUrl, lang, data.ActId)

	res, err := hoyoRequest(client, "GET", url, nil, cfg, account, data)
	if err != nil {
		return nil, err
	}

	var info signInInfo
	err = json.Unmarshal(res.Data, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type awardResponse struct {
	Awards []award `json:"awards"`
}

type award struct {
	Name  string `json:"name"`
	Count int    `json:"cnt"`
}

func getAwards(client *http.Client, cfg *config.Config, account *config.Account, data *gameConfig) ([]award, error) {
	url := fmt.Sprintf("%s/home?lang=%s&act_id=%s", data.EventBaseUrl, lang, data.ActId)

	res, err := hoyoRequest(client, "GET", url, nil, cfg, account, data)
	if err != nil {
		return nil, err
	}

	var awardRes awardResponse
	err = json.Unmarshal(res.Data, &awardRes)
	if err != nil {
		return nil, err
	}

	return awardRes.Awards, nil
}

func performCheckIn(client *http.Client, cfg *config.Config, account *config.Account, data *gameConfig) (*hoyoResponse, error) {
	url := fmt.Sprintf("%s/sign?lang=%s", data.EventBaseUrl, lang)

	payload, err := json.Marshal(map[string]string{"act_id": data.ActId})
	if err != nil {
		return nil, err
	}

	return hoyoRequest(client, "POST", url, payload, cfg, account, data)
}

func isCaptchaRequired(resp *hoyoResponse) bool {
	if resp == nil || resp.GtResult == nil {
		return false
	}

	if resp.GtResult.GT == "" && resp.GtResult.Challenge == "" {
		return false
	}

	return true
}
