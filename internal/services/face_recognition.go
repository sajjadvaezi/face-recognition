package services

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal/clients"
	"log/slog"
)

func RecognizeFace() (name string, err error) {

	fc := clients.NewFlaskClient("http://127.0.0.1:5000")
	slog.Info("calling flask endpoint")
	hash, err := fc.RecognizeFace()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	user, err := db.FindUserByFaceHash(hash)
	if err != nil {

		return "", err
	}

	return user.Name, nil
}

func AddFace(studentNumber string) error {

	fc := clients.NewFlaskClient("http://127.0.0.1:5000")
	slog.Info("calling flask endpoint")
	hash, err := fc.RegisterFace()
	fmt.Println("face hash")
	_, err = db.AddFaceWithStudentNumber(studentNumber, hash)
	if err != nil {

		fmt.Println("add to db error")

		return err
	}
	return nil
}
