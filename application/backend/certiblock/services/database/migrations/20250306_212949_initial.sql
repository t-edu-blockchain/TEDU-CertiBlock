CREATE TABLE IF NOT EXISTS countries (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    public_key TEXT NOT NULL,
    encrypted_partial_personal_info_1 LONGTEXT NOT NULL,
    encrypted_partial_personal_info_2 LONGTEXT NOT NULL,

    KEY idx_student_pubkey (public_key(16))
);

CREATE TABLE IF NOT EXISTS certificates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    uuid TEXT NOT NULL,
    university_encrypted_ks_1 LONGTEXT NOT NULL,
    university_encrypted_ks_2 LONGTEXT NOT NULL,
    student_encrypted_ks_1 LONGTEXT NOT NULL,
    student_encrypted_ks_2 LONGTEXT NOT NULL,
    ks_encrypted_file LONGTEXT NOT NULL,

    KEY idx_cert_uuid (uuid(16))
);

CREATE TABLE IF NOT EXISTS qrs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cert_uuid TEXT NOT NULL,
    qr_encrypted_ks_1 LONGTEXT NOT NULL,
    qr_encrypted_ks_2 LONGTEXT NOT NULL
);
