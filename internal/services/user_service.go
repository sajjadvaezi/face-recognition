package services

import (
	"github.com/sajjadvaezi/face-recognition/db"
)

func CreateUser(name string) (int64, error) {
	userID, err := db.AddUser(name)
	if err != nil {

		return -1, err
	}

	return userID, nil
}

