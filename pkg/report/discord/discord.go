package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/atomicptr/pity-patrol/pkg/config"
	"github.com/atomicptr/pity-patrol/pkg/constants"
	"github.com/atomicptr/pity-patrol/pkg/report"
)

const colorSuccess = 4431943
const colorFailure = 15022389

func Send(reporter *config.Reporter, cfg *config.Config, account *config.Account, r *report.Report) error {
	// when we already claimed reward send nothing
	if !r.WasClaimed {
		return nil
	}

	fields := []map[string]any{
		{
			"name":  "Game",
			"value": account.GameName(),
		},
	}

	if account.Identifier != "" {
		fields = append(fields, map[string]any{
			"name":  "Account",
			"value": account.Identifier,
		})
	}

	for _, f := range r.CustomFields {
		fields = append(fields, map[string]any{
			"name":  f.Key,
			"value": f.Value,
		})
	}

	var thumbnail map[string]string = nil

	// add reward infos if available
	if r.Reward != nil {
		fields = append(fields, map[string]any{
			"name":  "Reward",
			"value": fmt.Sprintf("%dx %s", r.Reward.Count, r.Reward.Name),
		})

		if r.Reward.Image != "" {
			thumbnail = map[string]string{
				"url": r.Reward.Image,
			}
		}
	}

	payload := map[string]any{
		"content": "",
		"embeds": []map[string]any{
			{
				"author": map[string]string{
					"name": "atomicptr/pity-patrol",
					"url":  "https://github.com/atomicptr/pity-patrol",
				},
				"color":       colorSuccess,
				"description": "Successfully claimed reward",
				"fields":      fields,
				"thumbnail":   thumbnail,
				"timestamp":   time.Now().Format(time.RFC3339),
			},
		},
	}

	return sendToDiscord(payload, reporter, cfg)
}

func SendError(reporter *config.Reporter, cfg *config.Config, account *config.Account, message string) error {
	fields := []map[string]any{
		{
			"name":  "Game",
			"value": account.GameName(),
		},
	}

	if account.Identifier != "" {
		fields = append(fields, map[string]any{
			"name":  "Account",
			"value": account.Identifier,
		})
	}

	payload := map[string]any{
		"content": "",
		"embeds": []map[string]any{
			{
				"author": map[string]string{
					"name": "atomicptr/pity-patrol",
					"url":  "https://github.com/atomicptr/pity-patrol",
				},
				"color":       colorFailure,
				"description": message,
				"fields":      fields,
				"timestamp":   time.Now().Format(time.RFC3339),
			},
		},
	}

	return sendToDiscord(payload, reporter, cfg)
}

func sendToDiscord(payload map[string]any, reporter *config.Reporter, cfg *config.Config) error {
	client := http.Client{Timeout: constants.DefaultTimeoutSecs}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reporter.WebhookUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if cfg.DebugMode {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Discord Response: %s\n", body)
	}

	return nil
}
