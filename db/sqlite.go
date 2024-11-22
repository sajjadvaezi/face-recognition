package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
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
func AddUser(name, studentNumber string) (int64, error) {
	query := `INSERT INTO users (name, student_number,created_at) VALUES (?, ? ,?)`
	result, err := db.Exec(query, name, studentNumber, time.Now())
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

func AddFaceWithStudentNumber(studentNumber, faceHash string) (int64, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Find user within the transaction
	queryFindUser := `SELECT user_id FROM users WHERE student_number = ?`
	var userID int64
	err = tx.QueryRow(queryFindUser, studentNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no user found with the given student number")
		}
		return 0, fmt.Errorf("failed to find user: %w", err)
	}

	// Insert face data within the transaction
	queryInsertFace := `INSERT INTO user_faces (user_id, face_hash, created_at) VALUES (?, ?, ?)`
	result, err := tx.Exec(queryInsertFace, userID, faceHash, time.Now())
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to insert face data: %w", err)
	}

	// Return the last inserted ID
	faceID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return faceID, nil
}

// FindUserByFaceHash looks up a user by their face hash
func FindUserByFaceHash(faceHash string) (*models.User, error) {
	query := `SELECT users.user_id, users.name,users.student_number, users.created_at
	          FROM users
	          JOIN user_faces ON users.user_id = user_faces.user_id
	          WHERE user_faces.face_hash = ?`

	row := db.QueryRow(query, faceHash)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.StudentNumber, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found with the given face hash")
	} else if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func FindUserByStudentNumber(studentNumber string) (*models.User, error) {
	query := `SELECT user_id, name, student_number, created_at
	          FROM users
	          WHERE student_number = ?`

	row := db.QueryRow(query, studentNumber)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.StudentNumber, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found with the given student number")
	} else if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}
