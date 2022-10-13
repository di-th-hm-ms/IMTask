package service

import (
	"fmt"
	// "os"
	// "time"
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
		panic(err)
	}

	// mock
	userReq := model.UserReq{
		Id: "",
		Email: "abcde@gmail.com",
		Username: "abcde",
		Password: "12345abCd",
		CreatedAt: "",
	}
	u.InsertUserIntoDB(&userReq, db)
}

func (UserService) DropUserTable(db *sql.DB) {
	if _, err := db.Exec(`DROP TABLE IF EXISTS user`); err != nil {
		log.Fatal("drop user table error: ", err)
	}
}

// DEBUG
func (UserService) GetUsersFromDB(db *sql.DB) ([]model.User, *model.ServerError) {
	users := make([]model.User, 0)
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		return nil, model.NewServerError(model.Status("500"), model.Msg("Internal server error"))
	}
	user := model.User{}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewServerError(model.Status("500"), model.Msg("Internal server error"))
	}
	return users, nil
}

func (UserService) GetUserFromDB(id string, password string, db *sql.DB) (*model.User, *model.ServerError) {
	user := model.NewUser()
	err := db.QueryRow("SELECT * FROM user WHERE id = ? AND password = ?", id, password).Scan(
		&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no row. It doesn't matter: %s\n", err)
			// return nil, model.NewServerError(model.Status("422"), model.Msg(err.Error()))
			return nil, model.NewServerError(model.Status("401"), model.Msg("not authenticated fir"))
		} else {
			// Scan was failed
			log.Printf("Scan error: %s\n", err)
			// return nil
			return nil, model.NewServerError(model.Status("500"), model.Msg("Internal server error"))
		}
	}
	return user, nil
}

func (UserService) GetRegisteredInfoFromDB(username string,  db *sql.DB) bool {
	type Str struct { str string }
	s := Str{}
	if err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&s.str); err != nil {
		log.Printf("Scan error: %s\n", err)
		// no rows or Scan error
	}
	return s.str != ""
}


func (u UserService) InsertUserIntoDB(userReq *model.UserReq, db *sql.DB) (*model.User, *model.ServerError) {
	if userReq.Email == "" || userReq.Username == "" || userReq.Password == "" {
		return nil, model.NewServerError(model.Msg("Required parameter is empty"))
	}
	if !model.ValidateEmail(userReq.Email) {
		return nil, model.NewServerError(model.Msg("Email has a wrong format"))
	}
	if !model.ValidatePassword(userReq.Password) {
		return nil, model.NewServerError(model.Msg("Password has a wrong format"))
	}
	defer func () { recover() }()
	ins, err := db.Prepare("INSERT INTO user(id, email, username, password, createdAt) VALUES(?,?,?,?,NOW())")
	if err != nil {
		return nil, model.NewServerError(model.Status("500"), model.Msg(err.Error()))
	}

	if u.GetRegisteredInfoFromDB(userReq.Username, db) {
		return nil, model.NewServerError(model.Status("422"), model.Msg("Already used"))
	}
	id, serverErr := Insert(ins, userReq, 1)
	if serverErr != nil {
		return nil, serverErr
	}
	// TODO routine (+jwt)
	user, serverErr := u.GetUserFromDB(id, userReq.Password, db)
	return user, serverErr
}

func Insert(ins *sql.Stmt, userReq *model.UserReq, cnt int) (string, *model.ServerError) {
	if (cnt >= 3) {
		fmt.Printf("try %d times(insert error)\n", cnt)
		return "", model.NewServerError()
	}
	str, err := model.GenerateRandStr(20)
	if err != nil {
		return "", model.NewServerError(model.Status("500"), model.Msg(err.Error()))
	}
	// if _, err := ins.Exec(str, userReq.Email, userReq.Username, userReq.Password, userReq.CreatedAt); err != nil {
	if _, err := ins.Exec(str, userReq.Email, userReq.Username, userReq.Password); err != nil {
		cnt++
		fmt.Println(err)
		return Insert(ins, userReq, cnt+1)
	}
	return str, nil
}

func (u UserService) UpdateUsernameOnDB(userReq *model.UserReq, db *sql.DB) *model.ServerError {
	// TODO jwt

	if _, serverError := u.GetUserFromDB(userReq.Id, userReq.Password, db);serverError != nil {
		return serverError
	}
	// TODO password 一致 bcrypt
	// if user.Password != userReq.Password {
	// 	return errors.New("passwrd isn't matched(update)")
	// }
	// _, err := db.Exec("UPDATE user SET title = ? WHERE id = ? AND userId = ?", userTitle, userId, userId)
	// TODO TEST
	if _, err := db.Exec("UPDATE user SET username = ? WHERE id = ? AND password = ?", userReq.Username, userReq.Id, userReq.Password); err != nil {
		log.Printf("Update error: %s", err)
		return model.NewServerError(model.Status("401"), model.Msg("not authenticated upd sec"))
	}

	return nil
}

func (UserService) ChangePassword(userReq *model.UserReq, db *sql.DB) {
	// TODO password reset
}

func (UserService) ChangeEmail(userReq *model.UserReq, db *sql.DB) {
	// TODO email reset
	// if !model.ValidateEmail(userReq.Email) || user.Email != userReq.Email {
	// 	return model.NewServerError(model.Msg("Email is incorrect"))
	// }
}

func (u UserService) DeleteUserFromDB(userReq *model.UserReq, db *sql.DB) *model.ServerError {
	user, serverError := u.GetUserFromDB(userReq.Id, userReq.Password, db)
	if serverError != nil {
		return model.NewServerError(model.Status("401"), model.Msg("not authenticated del fir"))
	}
	// if user.Password != userReq.Password {
	// 	return errors.New("Password isn't matched")
	// }
	// TODO jwt
	if _, err := db.Exec("DELETE FROM user WHERE id = ? and password = ?", user.Id, user.Password); err != nil {
		log.Printf("delete error: %s\n", err)
		return model.NewServerError(model.Status("401"), model.Msg("not authenticated del sec"))
	}

	return nil
}
