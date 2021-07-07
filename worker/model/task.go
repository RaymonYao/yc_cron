package model

type Task struct {
	TaskId          int   `json:"task_id" gorm:"PRIMARY_KEY"`
	PrevExecuteTime int64 `json:"prev_execute_time"`
}

func SaveTask(task *Task) (err error) {
	db := mdb
	err = db.Model(task).Updates(task).Error
	return
}
