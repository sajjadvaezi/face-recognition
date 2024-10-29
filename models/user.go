package user

type FaceHash string

// FaceRecognition defines the methods for registering and recognizing a face.
type FaceRecognition interface {
	RegisterFace(face Face) error               // Registers a new face
	RecognizeFace(hash FaceHash) (*User, error) // Recognizes a face by its hash and returns the associated user
}

// User represents a user with a unique ID, name, and associated face data.
type User struct {
	UserID int
	Name   string
	Face   Face
}

// Face represents a face with a unique hash.
type Face struct {
	Hash FaceHash
}
