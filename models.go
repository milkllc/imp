package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

type User struct {
	UserId int64
	Handle string
	Status string
	Biography string
	Email string
	PasswordHash string
	JoinedDate mysql.NullTime
}

func (u *User) Fetch(db *sql.DB, handle string) (err error) {
    u.UserId = -1

	stmt, err := db.Prepare("SELECT `UserId`, `Handle`, `Status`, `Biography`, `JoinedDate` FROM `User` WHERE Handle LIKE ? LIMIT 1")
	if err != nil {
	    log.Fatal(err)
	    return
	}
	defer stmt.Close()

	rows, err := stmt.Query(handle)
	if err != nil {
	    log.Fatal(err)
	    return
	}

	if rows.Next() {
	    if err := rows.Scan(&u.UserId, &u.Handle, &u.Status, &u.Biography, &u.JoinedDate); err != nil {
	        log.Fatal(err)
	    }
	}
	if err := rows.Err(); err != nil {
	    log.Fatal(err)
	}
	return
}

func (u *User) Save(db *sql.DB) (err error) {
	if u.UserId > 0 {

	} else {

	}
	// TODO: 
	return
}

// TODO: filters
func QueryUsers(db *sql.DB) (users []*User) {
	rows, err := db.Query("SELECT `UserId`, `Handle`, `Status`, `Biography`, `JoinedDate` FROM `User` WHERE 1")	// TODO: user id
	if err != nil {
	    log.Fatal(err)
	    // TODO: 500
	}

	users = make([]*User, 0)
	for rows.Next() {
	    var user User
	    if err := rows.Scan(&user.UserId, &user.Handle, &user.Status, &user.Biography, &user.JoinedDate); err != nil {
	        log.Fatal(err)
	    }
	    users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
	    log.Fatal(err)
	    // TODO: 500
	}
	return
}