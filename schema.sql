CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    student_number TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on student_number for faster lookups
CREATE INDEX IF NOT EXISTS idx_student_number ON users (student_number);

CREATE TABLE IF NOT EXISTS user_faces (
    face_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    face_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Create an index on face_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_face_id ON user_faces (face_id);
