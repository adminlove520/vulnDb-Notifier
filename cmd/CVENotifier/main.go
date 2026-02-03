// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adminlove520/vulnDb-Notifier/internal/config"
	"github.com/adminlove520/vulnDb-Notifier/internal/db"
	"github.com/adminlove520/vulnDb-Notifier/internal/errors"
	"github.com/adminlove520/vulnDb-Notifier/internal/rss"
	"github.com/joho/godotenv"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("main: Warning: .env file not found, using environment variables: %v", err)
		// 继续执行，不退出
	}

	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	feed, err := rss.ParseFeed("https://vuldb.com/?rss.recent")
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	databasePath := "CVENotifier.db"
	dbConn, err := db.InitDB(databasePath)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	defer dbConn.Close()

	slackWebhook := os.Getenv("SLACK_WEBHOOK")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK")

	if slackWebhook == "" && discordWebhook == "" {
		log.Fatalf("main: At least one of SLACK_WEBHOOK or DISCORD_WEBHOOK environment variable must be set")
	}

	var matchFound = 0

	// 确定推送模式，默认为 daily
	pushMode := cfg.PushMode
	if pushMode == "" {
		pushMode = "daily"
	}

	log.Printf("Push Mode: %s", pushMode)

	for _, item := range feed.Items {
		if pushMode == "daily" {
			// daily 模式：推送所有 RSS 内容
			matchFound++

			log.Printf("Title: " + item.Title)
			log.Printf("Link: " + item.Link)
			log.Printf("Published Date: " + item.Published)
			log.Printf("Categories: " + strings.Join(item.Categories, ","))
			log.Printf("Description: " + item.Description)

			description := item.Description
			if description == "" {
				description = "No description available."
			}

			err = db.InsertData(dbConn, item.Title, item.Link, item.Published, strings.Join(item.Categories, ","), description, slackWebhook, discordWebhook)
			if err != nil {
				if _, ok := err.(*errors.SlackNotificationError); ok {
					log.Printf("main: Failed to send notification: %v", err)
				} else {
					log.Printf("main: Failed to insert data: %v", err)
				}
			}
		} else {
			// keyword 模式：根据关键词过滤推送
			for _, keyword := range cfg.Keywords {
				if strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword)) {
					matchFound++

					log.Printf("Matched Keyword: " + keyword)
					log.Printf("Title: " + item.Title)
					log.Printf("Link: " + item.Link)
					log.Printf("Published Date: " + item.Published)
					log.Printf("Categories: " + strings.Join(item.Categories, ","))
					log.Printf("Description: " + item.Description)

					description := item.Description
					if description == "" {
						description = "No description available."
					}

					err = db.InsertData(dbConn, item.Title, item.Link, item.Published, strings.Join(item.Categories, ","), description, slackWebhook, discordWebhook)
					if err != nil {
						if _, ok := err.(*errors.SlackNotificationError); ok {
							log.Printf("main: Failed to send notification: %v", err)
						} else {
							log.Printf("main: Failed to insert data: %v", err)
						}
					}
				}
			}
		}
	}

	if matchFound == 0 {
		fmt.Printf("main: Result: No CVE matches found in the vuldb RSS feed\n")
	}
}
