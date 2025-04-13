package utils

import (
	"time"
)

func ScheduleTask(f func(), hour, minute int) {
    now := time.Now()
    year, month, day := now.Date()
    targetTime := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
    
    // If target time has passed today, schedule for tomorrow
    if targetTime.Before(now) {
        targetTime = targetTime.AddDate(0, 0, 1)
    }
    
    duration := time.Until(targetTime)
    time.AfterFunc(duration, f)
    // No blocking <-ctx.Done()
}

// 每个整点执行一次
func HourTask(f func()) {
    now := time.Now()
	nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	initialDelay := nextHour.Sub(now)
	
	// 先安排第一次执行
	time.AfterFunc(initialDelay, func() {
		f()
		
		// 设置每小时执行一次的ticker
		ticker := time.NewTicker(time.Hour)
		go func() {
			for range ticker.C {
                f()
			}
		}()
	})
}