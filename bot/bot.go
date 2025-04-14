package bot

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"qwbot/database"
	"qwbot/model/task"
	"qwbot/utils"
)

type bot struct {
	webhook string
	msgChan chan string
}

func NewBot(webhook string) *bot {
	return &bot{
		webhook: webhook,
		msgChan: make(chan string, 5),
	}
}

func (b *bot) SendMessage(resp string) error {
	data := fmt.Sprintf(`{"msgtype":"markdown","markdown":{"content":"%s"}}`, resp)
	req, err := http.NewRequest("POST", b.webhook, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (b *bot) hourTask() {
	db := database.GetMysqlDb()
	var tasks []task.Task
	curTime := time.Now().AddDate(0, 0, 3).Add(time.Minute)
	nextTime := curTime.Add(time.Hour)
	// 查询并删除任务
    err := db.Where("time >= ? AND time <= ?", curTime, nextTime).Find(&tasks).Error
	fmt.Println(tasks)
    if err != nil {
        fmt.Println("Error querying tasks:", err)
        return
    }
    
    if len(tasks) == 0 {
        fmt.Println("No tasks found")
        return
    }
    
    // Then delete the tasks separately
    err = db.Where("time >= ? AND time <= ?", curTime, nextTime).Delete(&task.Task{}).Error
    if err != nil {
        fmt.Println("Error deleting tasks:", err)
    }
    
    for _, task := range tasks {
        b.msgChan <- fmt.Sprintf("Task: %s\nTime: %s", task.Content, task.Time)
    }
}

func (b *bot) Run() {
	go func() {
		utils.HourTask(b.hourTask)
	}()

	for msg := range b.msgChan {
		err := b.SendMessage(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
