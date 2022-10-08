package service

import (
	"fmt"
	// "os"
	// "time"
	"errors"
	"log"
	"IMTask/golang/src/model"
	"database/sql"
	// "github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

type TaskService struct {}
// var DbEngine *xorm.Engine

// TODO 臨時
func (TaskService) CreateTaskTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS task (
		id INT(10) AUTO_INCREMENT NOT NULL primary key,
		title VARCHAR(50) NOT NULL,
		isAchieved BOOLEAN,
		userId VARCHAR(20)
	);`)
			// unique(userId) // tasktable
			// userId INT(10) )`)
	// _, err := db.Exec("create table if not exists task ( id INT(10), title varchar(50) )")
	if err != nil {
		// TODO recover
		panic(err)
	}
	ins, err := db.Prepare("INSERT INTO task(title, isAchieved, userId) VALUES(?,?,?)")
	if err != nil {
		panic(err)
	}
	ins.Exec("task", true, "abcde")
}

func (TaskService) DropTaskTable(db *sql.DB) {
	if _, err := db.Exec(`DROP TABLE IF EXISTS task`); err != nil {
		log.Fatal("drop table error: ", err)
	}
}
// func (TaskService) CloseConnection(db *sql.DB) {
// 	if err := db.Close(); err != nil {
// 		log.Fatal("Failed to close DB connection:", err)
// 	}
// }

func (TaskService) GetTasksFromDB(db *sql.DB) []model.Task {
	tasks := make([]model.Task, 0)
	// err :=
	// var rows *sql.Rows, err error = db.Query("SELECT * FROM tasks")
	rows, err := db.Query("SELECT * FROM task")
	if err != nil {
		panic(err)
	}
	task := model.Task{}
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Title, &task.IsAchieved, &task.UserId)
		tasks = append(tasks, task)
	}
	err = rows.Close()
	if err != nil { fmt.Println("close error") }
	fmt.Println(tasks)
	return tasks
}

// AfterInserting, update
func GetTaskFromDB(id int64, db *sql.DB) *model.Task {
	task := model.Task{}
	err := db.QueryRow("SELECT * FROM task WHERE id = ?", id).Scan(
		&task.Id, &task.Title, &task.IsAchieved, &task.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no row. It doesn't matter: %s\n", err)
		} else {
			// Scan was failed
			log.Fatal(err)
			// return nil
		}
	}
	return &task
}

// TODO recover
func (TaskService) InsertTaskIntoDB(taskTitle string, isAchieved bool, userId string, db *sql.DB) *model.Task {
	// defer func () { recover() }()
	fmt.Println(taskTitle, isAchieved, userId)
	ins, err := db.Prepare("INSERT INTO task(title, isAchieved, userId) VALUES(?,?,?)")
	if err != nil {
		panic(err)
	}
	result, err := ins.Exec(taskTitle, isAchieved, userId)
	if err != nil {
		panic(err)
	}
	// TODO return val
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := result.RowsAffected()
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	return GetTaskFromDB(lastId, db)
}

// model.Task
// func (TaskService) UpdateTaskOnDB(taskId int64, taskTitle string, isAchieved bool, userId int64, db *sql.DB) bool {
func (TaskService) UpdateTaskOnDB(taskId int64, taskTitle string, isAchieved bool, userId string, db *sql.DB) error {
	task := GetTaskFromDB(taskId, db)
	if task.Id == 0 {
		return errors.New("no task (id)")
	}
	if task.UserId != userId {
		return errors.New("no task is matched(update)")
	}
	// _, err := db.Exec("UPDATE task SET title = ? WHERE id = ? AND userId = ?", taskTitle, taskId, userId)
	_, err := db.Exec("UPDATE task SET title = ?, isAchieved = ? WHERE id = ? AND userId = ?", taskTitle, isAchieved, taskId,  userId)
	if err != nil {
		// panic(err)
		return err
	}
	return nil
	// return result
}


func (TaskService) DeleteTaskFromDB(taskId int64, userId string, db *sql.DB) (*model.Task, error) {
	task := GetTaskFromDB(taskId, db)
	if task.Id == 0 {
		// log.Println(
		return nil, errors.New("targeted task not found")
	}
	if task.UserId != userId {
		return nil, errors.New("not your task")
	}
	_, err := db.Exec("DELETE FROM task WHERE id = ?", taskId)
	if err != nil {
		panic(err)
		// return false
	}
	return task, nil
}

func Rollback(tx *sql.Tx) {
	if err := recover(); err != nil {
		tx.Rollback()
	}
}
