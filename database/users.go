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

// NewUser creates a new user
func NewUser(name string, pwhash string, email string) (User, error) {
	var err error
	var user User

	if len(pwhash) != 128 {
		return user, errors.New("pwhash not of length 128, quitting")
	}

	user.Name = name
	user.Pwhash = pwhash
	user.Email = email

	user, err = PutUser(user)
	if err != nil {
		return user, errors.Wrap(err, "could not put user into the db while creating user, quitting")
	}

	return user, nil // you can replace this with return PutUser but that doesn't expand the error wrap that we want to have
}

// RetrieveUser retrieves a user given his name
func RetrieveUser(name string, pwhash string) (User, error) {
	var x User
	db, err := OpenDB()
	if err != nil {
		return x, errors.Wrap(err, "could not open db, quitting")
	}
	defer db.Close()

	var id, dbName, email, providedHash string
	err = db.QueryRow("SELECT * FROM users WHERE name = $1 AND pwhash = $2", name, pwhash).Scan(&id, &dbName, &email, &providedHash)
	if err != nil {
		return x, errors.Wrap(err, "could not get user by name")
	}
	x.Id, err = utils.StoICheck(id)
	if err != nil {
		return x, errors.Wrap(err, "could not convert string to integer")
	}
	x.Name = dbName
	x.Email = email
	x.Pwhash = providedHash

	return x, nil
}

// AuthUser returns true if the user's name and pwhashes match
func AuthUser(name string, pwhash string) bool {
	user, err := RetrieveUser(name, pwhash)
	return !(err != nil) && user.Pwhash == pwhash
}

// PutUser creates a new user in the database
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

// Save updates the user struct stored in the database
func (user *User) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "could not open db, quitting")
	}

	defer db.Close()
	sqlTx := `
	UPDATE users
	SET name=$2, email=$3, pwhash=$4
	WHERE id=$1
	RETURNING id
	`

	var returnedIdS string
	err = db.QueryRow(sqlTx, user.Id, user.Name, user.Email, user.Pwhash).Scan(&returnedIdS)
	if err != nil {
		return errors.Wrap(err, "could not insert user into db, quitting")
	}

	returnedId, err := utils.StoICheck(returnedIdS)
	if err != nil {
		return err
	}

	if returnedId != user.Id {
		return errors.New("ids don't match, quitting")
	}
	return nil
}

// RetrieveAllUsers retrieves all users from the database
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

// DeleteUser deletes a given user from the database
func DeleteUser(name string, pwhash string) error {
	user, err := RetrieveUser(name, pwhash)
	if err != nil {
		return errors.Wrap(err, "could not retrieve user from db, quitting")
	}

	if user.Name != name || user.Pwhash != pwhash {
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
	WHERE name = $1 AND pwhash = $2
	RETURNING id
	`
	var id2 string
	err = db.QueryRow(sqlTx, user.Name, user.Pwhash).Scan(&id2)
	if err != nil {
		return errors.Wrap(err, "could not execute sql to delete user from db, quitting")
	}

	if utils.StoI(id2) != user.Id {
		return errors.New("deleted user id and provided user id don't match, quitting")
	}

	return nil
}
