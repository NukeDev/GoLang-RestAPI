package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func dbConn() (db *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	//Getting vars from .env file
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbpassword := os.Getenv("DBPASSWORD")
	dbname := os.Getenv("DBNAME")
	db, erro := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":3306)/"+dbname)
	if erro != nil {
		panic(erro.Error())
	}
	return db
}

//GetAuthByKey DB
func GetAuthByKey(authkey string) AuthKey {
	db := dbConn()
	results, err := db.Query("SELECT * FROM auth WHERE authkey = ?;", authkey)
	if err != nil {
		panic(err.Error())
	}
	var ak AuthKey
	for results.Next() {
		err = results.Scan(&ak.ID, &ak.USERID, &ak.AUTHKEY, &ak.TYPE, &ak.TIMESTAMP)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	return ak
}

//GetAuthByUser DB
func GetAuthByUser(user_id int) AuthKey {
	db := dbConn()
	results, err := db.Query("SELECT * FROM auth WHERE user_id = ?;", user_id)
	if err != nil {
		panic(err.Error())
	}
	var ak AuthKey
	for results.Next() {
		err = results.Scan(&ak.ID, &ak.USERID, &ak.AUTHKEY, &ak.TYPE, &ak.TIMESTAMP)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	return ak
}

//InsertNewUser DB
func InsertNewUser(user User) (string, bool) {
	db := dbConn()

	results, err := db.Exec("INSERT IGNORE INTO users (name, surname, email, password) VALUES (?,?,?,?);", user.NAME, user.SURNAME, user.EMAIL, user.PASSWORD)
	count, err2 := results.RowsAffected()
	defer db.Close()

	if err != nil && err2 != nil {
		return err.Error(), false
	} else if count > 0 {
		return "Success", true
	} else {
		return "Already Exists!", false
	}

}

//InsertNewAuthkey DB
func InsertNewAuthkey(user User, authkey string) (string, bool) {
	db := dbConn()

	results, err := db.Exec("INSERT IGNORE INTO auth (user_id, authkey, type) VALUES (?,?,?);", user.ID, authkey, user.TYPE)
	count, err2 := results.RowsAffected()
	defer db.Close()

	if err != nil && err2 != nil {
		return err.Error(), false
	} else if count > 0 {
		return "Success", true
	} else {
		return "Already Exists!", false
	}

}

func UpdateNewAuthkey(user User, authkey string) (string, bool) {
	db := dbConn()

	log.Println(authkey, user.ID)
	results, err := db.Exec("UPDATE auth SET authkey = ? WHERE user_id = ?;", authkey, user.ID)
	count, err2 := results.RowsAffected()
	defer db.Close()

	if err != nil && err2 != nil {
		return err.Error(), false
	} else if count > 0 {
		return "Success", true
	} else {
		return "Already Exists!", false
	}

}

//GetUserByMail DB
func GetUserByMail(email string) (User, bool) {
	db := dbConn()
	var ak User
	var status bool
	results, err := db.Query("SELECT * FROM users WHERE email = ?;", email)
	if err != nil {
		status = false
	}

	for results.Next() {
		err = results.Scan(&ak.ID, &ak.NAME, &ak.SURNAME, &ak.EMAIL, &ak.PASSWORD, &ak.TYPE, &ak.TIMESTAMP)
		if err != nil {
			status = false
		}
	}
	defer db.Close()
	if ak.EMAIL == email {
		status = true
		return ak, status
	} else {
		return ak, status
	}

}
