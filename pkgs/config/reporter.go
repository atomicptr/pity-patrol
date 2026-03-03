package config

import (
	"fmt"
	"slices"
)

type Reporter struct {
	Type string   `toml:"type,omitempty"`
	On   []string `toml:"on,omitempty"`

	// Discord
	WebhookUrl string `toml:"webhook-url"`
}

func (r *Reporter) validate() error {
	switch r.Type {
	case "discord":
		if r.WebhookUrl == "" {
			return fmt.Errorf("discord: no webhook-url specified")
		}
	default:
		return fmt.Errorf("invalid reporter type '%s'", r.Type)
	}

	if len(r.On) == 0 {
		return fmt.Errorf("no `on` event registered, please use either `success`, `failure` or both")
	}

	for _, on := range r.On {
		switch on {
		case "success", "failure":
		default:
			return fmt.Errorf("invalid `on` value, please only use `success` and/or `failure`")
		}
	}

	return nil
}

func (r *Reporter) ReportOnSuccess() bool {
	return slices.Contains(r.On, "success")
}

func (r *Reporter) ReportOnFailure() bool {
	return slices.Contains(r.On, "failure")
}
