CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(30) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT students_grade_check
        CHECK (grade >= 0 AND grade <= 100)
);

CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key
ON students (LOWER(nim));