package services

import (
	"github.com/sajjadvaezi/face-recognition/db"
)

func CreateUser(name, studentNumber string) (int64, error) {
	userID, err := db.AddUser(name, studentNumber)
	if err != nil {

		return -1, err
	}

	return userID, nil
}
