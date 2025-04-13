package task

type Task struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Time    string `json:"time" gorm:"not null"`
	Content string `json:"content" gorm:"not null"`
}
