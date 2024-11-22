-- Users table to store user information
CREATE TABLE IF NOT EXISTS users (
                                     user_id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique identifier for each user
                                     name TEXT NOT NULL,                        -- Name of the user
                                     student_number TEXT NOT NULL UNIQUE,       -- Unique student number for each user
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp for record creation
);

-- Index for faster lookups by student number
CREATE INDEX IF NOT EXISTS idx_student_number ON users (student_number);

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
