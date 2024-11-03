package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sajjadvaezi/face-recognition/models"
	"time"
)

var db *sql.DB

func InitSQLite() {
	var err error
	db, err = sql.Open("sqlite3", "./face_recognition.db")
	if err != nil {
		panic(err)
	}

}

// AddUser inserts a new user into the users table and returns the inserted user ID.
func AddUser(name string) (int64, error) {
	query := `INSERT INTO users (name, created_at) VALUES (?, ?)`
	result, err := db.Exec(query, name, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// AddFace inserts a new face hash for a user into the user_faces table.
func AddFace(userID int, faceHash string) (int64, error) {
	query := `INSERT INTO face_data (user_id, face_hash, created_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, userID, faceHash, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindUserByFaceHash looks up a user by their face hash
func FindUserByFaceHash(faceHash string) (*models.User, error) {
	query := `SELECT users.user_id, users.name, users.created_at
	          FROM users
	          JOIN user_faces ON users.user_id = user_faces.user_id
	          WHERE user_faces.face_hash = ?`

	row := db.QueryRow(query, faceHash)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found with the given face hash")
	} else if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}
