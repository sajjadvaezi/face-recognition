package models

type FaceHash string

// FaceRecognition defines the methods for registering and recognizing a face.
type FaceRecognition interface {
	RegisterFace(face Face) error               // Registers a new face
	RecognizeFace(hash FaceHash) (*User, error) // Recognizes a face by its hash and returns the associated models
}

// User represents a models with a unique ID, name, and associated face data.
type User struct {
	UserID        int
	Name          string
	StudentNumber string
	CreatedAt     string
}

type Face interface {
	GetFaceHash() string
}
