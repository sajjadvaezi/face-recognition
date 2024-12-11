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
func AddUser(name, userNumber, role string) (int64, error) {
	query := `INSERT INTO users (name, user_number, role, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, name, userNumber, role, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// AddFace inserts a new face hash for a user into the user_faces table.
func AddFace(userID int, faceHash string) (int64, error) {
	query := `INSERT INTO user_faces (user_id, face_hash, created_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, userID, faceHash, time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func AddFaceWithUserNumber(userNumber, faceHash string) (int64, error) {
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
	queryFindUser := `SELECT user_id FROM users WHERE user_number = ?`
	var userID int64
	err = tx.QueryRow(queryFindUser, userNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no user found with the given user number")
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
	query := `SELECT users.user_id, users.name, users.user_number, users.role, users.created_at
              FROM users
              JOIN user_faces ON users.user_id = user_faces.user_id
              WHERE user_faces.face_hash = ?`

	row := db.QueryRow(query, faceHash)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.UserNumber, &user.Role, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found with the given face hash")
	} else if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func FindUserByUserNumber(userNumber string) (*models.User, error) {
	query := `SELECT user_id, name, user_number, role, created_at
              FROM users
              WHERE user_number = ?`

	row := db.QueryRow(query, userNumber)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.UserNumber, &user.Role, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found with the given user number")
	} else if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func AddClass(classname, userNumber string) error {
	queryInsertClass := `INSERT INTO classes (classname, teacher_id) VALUES (?,?)`
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
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
	queryFindUser := `SELECT user_id FROM users WHERE user_number = ?`
	var userID int64
	err = tx.QueryRow(queryFindUser, userNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no user found with the given user number")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	_, err = tx.Exec(queryInsertClass, classname, userID)
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return err
		}
		return fmt.Errorf("failed to insert class: %w", err)
	}
	return nil

}

func Attendance(studentNumber string, className string) (int64, error) {
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

	queryFindClass := `SELECT class_id FROM classes WHERE classname = ?`
	var classID int64
	err = tx.QueryRow(queryFindClass, className).Scan(&classID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no class found with the given class name")
		}
		return 0, fmt.Errorf("failed to find class: %w", err)
	}

	queryFindUser := `SELECT user_id FROM users WHERE user_number = ?`
	var userID int64
	err = tx.QueryRow(queryFindUser, studentNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no user found with the given user number")
		}
		return 0, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if attendance already exists for today
	queryCheckAttendance := `
        SELECT attendance_id FROM attendance 
        WHERE student_id = ? AND class_id = ? AND DATE(date) = DATE('now')
    `
	var existingAttendanceID int64
	err = tx.QueryRow(queryCheckAttendance, userID, classID).Scan(&existingAttendanceID)
	if err == nil {
		// Attendance already exists for today
		return 0, fmt.Errorf("attendance already recorded for this student in this class today")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Some other database error occurred
		return 0, fmt.Errorf("failed to check existing attendance: %w", err)
	}

	// Insert new attendance record
	queryInsertAttendance := `
        INSERT INTO attendance (student_id, class_id, date, present) 
        VALUES (?, ?, DATE('now'), 1)
    `
	result, err := tx.Exec(queryInsertAttendance, userID, classID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert attendance: %w", err)
	}

	// Get the ID of the newly inserted attendance record
	attendanceID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get attendance record ID: %w", err)
	}

	return attendanceID, nil
}
