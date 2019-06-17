package database

import(
	"github.com/pkg/errors"
	"github.com/YaleOpenLab/openclimate/utils"
)

type User struct {
	Id int
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
