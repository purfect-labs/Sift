package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var version = "dev"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	svc := &JobService{}
	if err := svc.Init(); err != nil {
		log.Fatal("DB init failed:", err)
	}

	// Initialize analytics (reads POSTHOG_API_KEY + POSTHOG_ANALYTICS_ENABLED from .env)
	InitAnalytics(svc.config)
	defer CloseAnalytics()

	app := application.New(application.Options{
		Name:        "JobDash",
		Description: "AI-powered job search & application tracker",
		Services: []application.Service{
			application.NewService(svc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Register scrape log event
	application.RegisterEvent[string]("scrape_log")

	// Give service access to app for event emission
	svc.app = app

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "JobDash",
		Width:  1200,
		Height: 800,
		MinWidth:  900,
		MinHeight: 600,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(10, 12, 20),
		URL:              "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
