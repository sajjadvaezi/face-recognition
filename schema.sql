-- Users table to store both student and teacher information
CREATE TABLE IF NOT EXISTS users (
                                     user_id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique identifier for each user
                                     name TEXT NOT NULL,                        -- Name of the user
                                     user_number TEXT NOT NULL UNIQUE,          -- Unique identifier number for both students and teachers
                                     role TEXT NOT NULL CHECK(role IN ('student', 'teacher')),  -- Role of the user
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp for record creation
);

-- Index for faster lookups by user number
CREATE INDEX IF NOT EXISTS idx_user_number ON users (user_number);

-- User faces table to store face hashes associated with users
CREATE TABLE IF NOT EXISTS user_faces (
                                          face_id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique identifier for each face entry
                                          user_id INTEGER NOT NULL,                  -- Foreign key referencing the users table
                                          face_hash TEXT NOT NULL UNIQUE,            -- Unique face hash for each face entry
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for record creation
                                          FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- Cascade delete
);

-- Index for faster lookups by user ID (useful for joins)
CREATE INDEX IF NOT EXISTS idx_user_id ON user_faces (user_id);

-- Index for faster lookups by face hash
CREATE INDEX IF NOT EXISTS idx_face_hash ON user_faces (face_hash);

-- (Optional) Composite index for combined lookups on user ID and face hash
CREATE INDEX IF NOT EXISTS idx_user_id_face_hash ON user_faces (user_id, face_hash);

-- Classes table to link teachers with classes (teachers are now referenced by user_id)
CREATE TABLE IF NOT EXISTS classes (
                                       class_id INTEGER PRIMARY KEY AUTOINCREMENT,
                                       classname TEXT NOT NULL UNIQUE,  -- e.g., 'XYZ'
                                       teacher_id INTEGER NOT NULL,
                                       FOREIGN KEY (teacher_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Attendance table to track student presence in classes
CREATE TABLE IF NOT EXISTS attendance (
                                          attendance_id INTEGER PRIMARY KEY AUTOINCREMENT,
                                          class_id INTEGER NOT NULL,
                                          student_id INTEGER NOT NULL,
                                          date DATE NOT NULL,
                                          present BOOLEAN NOT NULL,  -- or use INTEGER where 1 is present, 0 is absent
                                          FOREIGN KEY (class_id) REFERENCES classes(class_id) ON DELETE CASCADE,
                                          FOREIGN KEY (student_id) REFERENCES users(user_id) ON DELETE CASCADE,
                                          UNIQUE (class_id, student_id, date)
);

-- Indexes for the attendance table for better query performance
CREATE INDEX IF NOT EXISTS idx_attendance_class ON attendance (class_id);
CREATE INDEX IF NOT EXISTS idx_attendance_user ON attendance (student_id);
CREATE INDEX IF NOT EXISTS idx_attendance_date ON attendance (date);