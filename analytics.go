package main

import (
	"os"
	"path/filepath"

	"github.com/posthog/posthog-go"
)

var analyticsClient posthog.Client

// InitAnalytics initializes the PostHog client if analytics are enabled and an API key is configured.
// It reads POSTHOG_API_KEY and POSTHOG_ANALYTICS_ENABLED from ~/.jobdash/.env.
func InitAnalytics(cfg Config) {
	if !cfg.AnalyticsEnabled {
		return
	}
	home, _ := os.UserHomeDir()
	env := loadEnv(filepath.Join(home, ".jobdash", ".env"))
	apiKey := env["POSTHOG_API_KEY"]
	if apiKey == "" {
		return
	}
	client, err := posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint: "https://us.i.posthog.com",
	})
	if err != nil {
		return
	}
	analyticsClient = client
}

// CloseAnalytics flushes and closes the PostHog client. Safe to call even if analytics
// were never initialized (nil client is a no-op).
func CloseAnalytics() {
	if analyticsClient != nil {
		analyticsClient.Close()
		analyticsClient = nil
	}
}

// TrackEvent sends an analytics event to PostHog. If the client is nil (analytics disabled
// or not configured), this is a no-op. All events use an anonymous distinct ID.
// NEVER include: job titles, company names, descriptions, resume text, keywords,
// salary data, location values, API keys, or file paths.
func TrackEvent(name string, properties map[string]interface{}) {
	if analyticsClient == nil {
		return
	}
	if properties == nil {
		properties = map[string]interface{}{}
	}
	analyticsClient.Enqueue(posthog.Capture{
		DistinctId: "jobdash-anonymous",
		Event:      name,
		Properties: properties,
	})
}
