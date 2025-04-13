package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"qwbot/bot"
	"qwbot/database"
	"qwbot/model/task"
)

type SetRequest struct {
	Year    int    `json:"year"`
	Month   int    `json:"month"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Content string `json:"content"`
}

func main() {
	if os.Getenv("GO_END") != "" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}

	webhook := os.Getenv("WEBHOOK")
	if webhook == "" {
		panic("WEBHOOK is not set in .env file")
	}

	database.Init()
	r := gin.Default()
	b := bot.NewBot(webhook)
	go b.Run()

	r.POST("/set", func(c *gin.Context) {
		var sr SetRequest
		if err := c.ShouldBindJSON(&sr); err != nil {
			c.JSON(400, "error parsing request body")
			return
		}

		db := database.GetMysqlDb()
		var t task.Task
		t.Time = time.Date(sr.Year, time.Month(sr.Month), sr.Day, sr.Hour, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
		t.Content = sr.Content
		if err := db.Create(&t).Error; err != nil {
			c.JSON(500, "error creating task")
			return
		}
		c.JSON(200, "success")
	})

	r.Run(":8080")
}
