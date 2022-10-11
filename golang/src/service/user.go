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

type UserService struct {
}

func (u UserService) CreateUserTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(20) BINARY NOT NULL,
		email VARCHAR(50) NOT NULL,
		username VARCHAR(50) NOT NULL,
		password VARCHAR(256) NOT NULL,
		createdAt DATETIME NOT NULL,
		unique(id)
	);`)
	if err != nil {
		// TODO recover
		panic(err)
	}

	userReq := model.UserReq{
		Id: "",
		Email: "abcde@gmail.com",
		Username: "abcde",
		Password: "12345",
		CreatedAt: "",
	}
	u.InsertUserIntoDB(&userReq, db)
}

func (UserService) DropUserTable(db *sql.DB) {
	if _, err := db.Exec(`DROP TABLE IF EXISTS user`); err != nil {
		log.Fatal("drop user table error: ", err)
	}
}

func (UserService) GetUsersFromDB(db *sql.DB) []model.User {
	users := make([]model.User, 0)
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err)
	}
	user := model.User{}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil { log.Println(err) }
	fmt.Println(users)
	return users
}

func (UserService) GetUserFromDB(id string, password string, db *sql.DB) *model.User {
	user := model.NewUser()
	err := db.QueryRow("SELECT * FROM user WHERE id = ? AND password = ?", id, password).Scan(
		&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no row. It doesn't matter: %s\n", err)
		} else {
			// Scan was failed
			log.Printf("Scan error: %s\n", err)
			// return nil
		}
		return nil
	}
	return user
}

func (UserService) GetUsernameFromDB(username string,  db *sql.DB) bool {
	type Str struct { str string }
	s := Str{}
	if err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&s.str); err != nil {
		log.Printf("Scan error: %s\n", err)
		// no rows or Scan error
	}
	return s.str != ""
}


func (u UserService) InsertUserIntoDB(userReq *model.UserReq, db *sql.DB) (*model.User, error) {
	if userReq.Email == "" || userReq.Username == "" || userReq.Password == "" {
		return nil, errors.New("400")
	}
	// defer func () { recover() }()
	ins, err := db.Prepare("INSERT INTO user(id, email, username, password, createdAt) VALUES(?,?,?,?,NOW())")
	if err != nil {
		panic(err)
	}

	if u.GetUsernameFromDB(userReq.Username, db) {
		return nil, errors.New("Already used")
	} else {
		fmt.Println(err)
		fmt.Println(userReq.Username)
	}
	id := Insert(ins, userReq, 1)
	if id == "InsertError22" {
		return nil, errors.New("400")
	}
	// TODO routine (+jwt)
	user := u.GetUserFromDB(id, userReq.Password, db)
	fmt.Println("insert 2")
	return user, nil
}

func Insert(ins *sql.Stmt, userReq *model.UserReq, cnt int) string {
	if (cnt >= 3) {
		fmt.Printf("try %d times(insert error)\n", cnt)
		return "InsertError22"
	}
	str, err := model.GenerateRandStr(20)
	if err != nil {
		log.Fatal("Generate Str error:",err)
	}
	// if _, err := ins.Exec(str, userReq.Email, userReq.Username, userReq.Password, userReq.CreatedAt); err != nil {
	if _, err := ins.Exec(str, userReq.Email, userReq.Username, userReq.Password); err != nil {
		cnt++
		fmt.Println(err)
		return Insert(ins, userReq, cnt+1)
	}
	fmt.Println("insert 1")
	return str
}

func (u UserService) UpdateUserOnDB(userReq *model.UserReq, db *sql.DB) error {
	// jwt後でも一応 check
	// TODO error の細分化
	user := u.GetUserFromDB(userReq.Id, userReq.Password, db)
	if user == nil {
		return errors.New("not authenticated")
	}
	// TODO password 一致 bcrypt
	// if user.Password != userReq.Password {
	// 	return errors.New("passwrd isn't matched(update)")
	// }
	// _, err := db.Exec("UPDATE user SET title = ? WHERE id = ? AND userId = ?", userTitle, userId, userId)
	if _, err := db.Exec("UPDATE user SET email = ?, username = ? WHERE id = ? AND password = ?", userReq.Email, userReq.Username, userReq.Id, userReq.Password); err != nil {
		log.Printf("Update error: %s", err)
		return err
	}

	return nil
}

func (UserService) UpdatePassword(userReq *model.UserReq, db *sql.DB) {
	// TODO password reset
}


func (u UserService) DeleteUserFromDB(userReq *model.UserReq, db *sql.DB) error {
	user := u.GetUserFromDB(userReq.Id, userReq.Password, db)
	if user == nil {
		return errors.New("not authenticated")
	}
	// if user.Password != userReq.Password {
	// 	return errors.New("Password isn't matched")
	// }
	// TODO jwt
	_, err := db.Exec("DELETE FROM user WHERE id = ? and password = ?", user.Id, user.Password)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
