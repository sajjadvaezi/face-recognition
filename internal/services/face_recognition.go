package services

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal/clients"
	"log/slog"
)

type FlaskResponse struct {
	Hash  string `json:"hash"`
	Error string `json:"error"`
}

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

func AddFace(userNumber string) error {

	fc := clients.NewFlaskClient("http://127.0.0.1:5000")
	slog.Info("calling flask endpoint")
	hash, err := fc.RegisterFace()
	fmt.Println("face hash")
	_, err = db.AddFaceWithUserNumber(userNumber, hash)
	if err != nil {

		fmt.Println("add to db error")

		return err
	}
	return nil
}

func RecognizeFaceWithImage(studentNumber, image string) (string, error) {
	fc := clients.NewFlaskClient("http://127.0.0.1:5000")
	slog.Info("calling flask endpoint")

	// Call the Flask service
	resp, err := fc.UploadImage(image)
	if err != nil {
		slog.Error("Error uploading image to Flask:", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	if resp.Hash == "" {
		return "", fmt.Errorf("flask service did not return a valid hash")
	}

	// Find the user by face hash in the database
	user, err := db.FindUserByFaceHash(resp.Hash)
	if err != nil {
		slog.Error("Error finding user in database:", slog.String("hash", resp.Hash), slog.String("error", err.Error()))
		return "", fmt.Errorf("user not found for hash: %s", resp.Hash)
	}

	// Verify the student number
	if user.UserNumber != studentNumber {
		return "", fmt.Errorf("student number mismatch: expected %s, got %s", user.UserNumber, studentNumber)
	}

	// Return the recognized user's name
	return user.Name, nil
}

func AddFaceWithImage(studentNumber, image string) error {
	fc := clients.NewFlaskClient("http://127.0.0.1:5000")
	slog.Info("calling flask endpoint")

	// sending the face through flask and register it
	resp, err := fc.RegisterImage(image)
	if err != nil {
		if resp.Error != "" {
			fmt.Println("couldn't upload image. error: ", resp.Error)
			return err
		}
		fmt.Println("couldn't upload image.")
		return err
	}

	// adding the face to the user in database
	_, err = db.AddFaceWithUserNumber(studentNumber, resp.Hash)
	if err != nil {
		fmt.Println("add to db error")

		return err
	}
	return nil

}
