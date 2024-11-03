package services

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
)

func RecognizeFace(faceHash string) (name string, err error) {
	if faceHash == "" {

		return "", fmt.Errorf("empty face hash")
	}
	user, err := db.FindUserByFaceHash(faceHash)
	if err != nil {

		return "", err
	}

	return user.Name, nil
}

func AddFace(userId int64, hash string) error {
	if hash == "" {

		return fmt.Errorf("empty face hash")
	}
	_, err := db.AddFace(int(userId), hash)
	if err != nil {
		fmt.Println("add to db error error ")

		return err
	}
	return nil
}
