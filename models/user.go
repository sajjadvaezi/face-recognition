package models

import "time"

type FaceHash string

// FaceRecognition defines the methods for registering and recognizing a face.
type FaceRecognition interface {
	RegisterFace(face Face) error               // Registers a new face
	RecognizeFace(hash FaceHash) (*User, error) // Recognizes a face by its hash and returns the associated user
}

// User represents a user with a unique ID, name, user number (which can be for both teachers and students), role, and creation time.
type User struct {
	UserID     int
	Name       string
	UserNumber string
	Role       string
	CreatedAt  time.Time
}

type Face interface {
	GetFaceHash() string
}
