package database

import (
	//"log"
	"github.com/YaleOpenLab/openclimate/utils"
	"github.com/pkg/errors"
)

type User struct {
	Id     int
	Name   string
	Email  string
	Pwhash string
}

func RetrieveUser(name string) (User, error) {
	var x User
	db, err := OpenDB()
	if err != nil {
		return x, errors.Wrap(err, "could not open db, quitting")
	}
	defer db.Close()

	var id, dbName, email, pwhash string
	err = db.QueryRow("SELECT * FROM users WHERE name = $1", name).Scan(&id, &dbName, &email, &pwhash)
	if err != nil {
		return x, errors.Wrap(err, "could not get user by name")
	}
	x.Id, err = utils.StoICheck(id)
	if err != nil {
		return x, errors.Wrap(err, "could not convert string to integer")
	}
	x.Name = dbName
	x.Email = email
	x.Pwhash = pwhash

	return x, nil
}

func PutUser(user User) (User, error) {
	db, err := OpenDB()
	if err != nil {
		return user, errors.Wrap(err, "could not open db, quitting")
	}

	defer db.Close()
	sqlTx := `
	INSERT INTO users (name, email, pwhash)
	VALUES($1, $2, $3)
	RETURNING id
	`
	err = db.QueryRow(sqlTx, user.Name, user.Email, user.Pwhash).Scan(&user.Id)
	if err != nil {
		return user, errors.Wrap(err, "could not insert user into db, quitting")
	}
	return user, nil
}

func RetrieveAllUsers() ([]User, error) {
	var users []User
	db, err := OpenDB()
	if err != nil {
		return users, errors.Wrap(err, "could not open db, quitting")
	}

	defer db.Close()

	sqlTx := `
	SELECT * FROM users
	`

	rows, err := db.Query(sqlTx)
	if err != nil {
		return users, errors.Wrap(err, "could not query db for all users, quitting")
	}

	for rows.Next() {
		var user User
		var id, name, email, pwhash string
		if err := rows.Scan(&id, &name, &email, &pwhash); err != nil {
			return users, err
		}
		user.Name = name
		user.Email = email
		user.Pwhash = pwhash
		user.Id, err = utils.StoICheck(id)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func DeleteUser(name string, id int) (error) {
	user, err := RetrieveUser(name)
	if err != nil {
		return errors.Wrap(err, "could not retrieve user from db, quitting")
	}

	if user.Name != name || user.Id != id {
		return errors.Wrap(err, "did not delete, user names don't match")
	}
	// open db and delete the user now
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "could not open db, quitting")
	}

	defer db.Close()
	sqlTx := `
	DELETE FROM users
	WHERE name = $1 AND id = $2
	`
	var id2 string
	err = db.QueryRow(sqlTx, user.Name, user.Id).Scan(&id2)
	if err != nil {
		return errors.Wrap(err, "could not insert user into db, quitting")
	}
	return nil
}
