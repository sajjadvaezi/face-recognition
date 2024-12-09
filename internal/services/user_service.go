package services

import (
	"github.com/sajjadvaezi/face-recognition/db"
)

func CreateUser(name, userNumber, role string) (int64, error) {
	userID, err := db.AddUser(name, userNumber, role)
	if err != nil {

		return -1, err
	}

	return userID, nil
}
