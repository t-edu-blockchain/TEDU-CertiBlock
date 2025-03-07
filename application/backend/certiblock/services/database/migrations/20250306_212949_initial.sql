CREATE TABLE IF NOT EXISTS countries (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS bc_enrollment_certificates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    student_public_key TEXT NOT NULL,
    university_public_key TEXT NOT NULL,
    hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bc_certificates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    h TEXT NOT NULL,
    C_X LONGTEXT NOT NULL,
    C_U LONGTEXT NOT NULL,
    C_S LONGTEXT NOT NULL,
    student_public_key TEXT NOT NULL,
    university_public_key TEXT NOT NULL,
    S0 LONGTEXT NOT NULL,
    S LONGTEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bc_universities (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    public_key TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS bc_students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    public_key TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS qrs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    rk TEXT NOT NULL,
    h TEXT NOT NULL,
    -- student_public_key TEXT NOT NULL
);
