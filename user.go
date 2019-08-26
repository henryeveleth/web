package main

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	Id        int            `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Age       int            `json:"age"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt mysql.NullTime `json:"-"`
}

func GetUser(id int) (User, error) {
	dbConn := DatabaseConnection()
	defer dbConn.Close()

	user := User{}

	/*
	   We are querying the database twice here. Ideally, we wouldn't do this.
	*/
	if exists, _ := UserExists(id); exists {
		query, err := dbConn.Query("SELECT * FROM users WHERE id=? AND deleted_at IS NULL", id)

		if err != nil {
			panic(err.Error())
		}

		for query.Next() {
			var id, age int
			var firstName, lastName string
			var createdAt, updatedAt time.Time
			var deletedAt mysql.NullTime

			err = query.Scan(&id, &firstName, &lastName, &age, &createdAt, &updatedAt, &deletedAt)

			if err != nil {
				panic(err.Error())
			}

			user.Id = id
			user.FirstName = firstName
			user.LastName = lastName
			user.Age = age
			user.CreatedAt = createdAt
			user.UpdatedAt = updatedAt
			user.DeletedAt = deletedAt
		}

		return user, nil
	} else {
		message := fmt.Sprintf("User %d not found.", id)
		respError := ResponseError{Message: message}
		return user, &respError
	}
}

func GetAllUsers() ([]User, error) {
	dbConn := DatabaseConnection()
	defer dbConn.Close()

	query, err := dbConn.Query("SELECT * FROM users WHERE deleted_at IS NULL")

	if err != nil {
		panic(err.Error())
	}

	users := []User{}

	for query.Next() {
		user := User{}

		var id, age int
		var firstName, lastName string
		var createdAt, updatedAt time.Time
		var deletedAt mysql.NullTime

		err = query.Scan(&id, &firstName, &lastName, &age, &createdAt, &updatedAt, &deletedAt)

		if err != nil {
			panic(err.Error())
		}

		user.Id = id
		user.FirstName = firstName
		user.LastName = lastName
		user.Age = age
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt

		users = append(users, user)
	}

	return users, nil
}

func (user *User) persist() error {
	dbConn := DatabaseConnection()
	defer dbConn.Close()

	id := user.Id
	firstName := user.FirstName
	lastName := user.LastName
	age := user.Age
	deletedAt := user.DeletedAt

	if exists, _ := UserExists(id); !exists {
		query, err := dbConn.Prepare("INSERT INTO users (first_name, last_name, age, created_at, updated_at) VALUES(?, ?, ?, now(), now())")

		if err != nil {
			panic(err.Error())
		}

		query.Exec(firstName, lastName, age)
	} else {
		query, err := dbConn.Prepare("UPDATE users SET first_name = ?, last_name = ?, age = ?, updated_at = now(), deleted_at = ? WHERE id = ?")

		if err != nil {
			panic(err.Error())
		}

		_, err = query.Exec(firstName, lastName, age, deletedAt, id)
		if err != nil {
			panic(err.Error())
		}
	}

	return nil
}

func UserExists(id int) (bool, error) {
	dbConn := DatabaseConnection()
	defer dbConn.Close()

	query, err := dbConn.Query("SELECT COUNT(*) FROM users WHERE id=?", id)

	if err != nil {
		panic(err.Error())
	}

	var count int
	for query.Next() {
		err = query.Scan(&count)

		if err != nil {
			panic(err.Error())
		}
	}

	return (count == 1), nil
}
