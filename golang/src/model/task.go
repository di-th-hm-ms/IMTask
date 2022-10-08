package model

type Task struct {
	Id			int64
	Title		string
	IsAchieved	bool
	UserId		string
	// TODO parent Task
	// PageId 		int64
	// createdAt 	string
}
