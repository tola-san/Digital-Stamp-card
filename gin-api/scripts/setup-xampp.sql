CREATE DATABASE IF NOT EXISTS digital_stamp
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'app'@'localhost' IDENTIFIED BY 'app';
CREATE USER IF NOT EXISTS 'app'@'127.0.0.1' IDENTIFIED BY 'app';

GRANT ALL PRIVILEGES ON digital_stamp.* TO 'app'@'localhost';
GRANT ALL PRIVILEGES ON digital_stamp.* TO 'app'@'127.0.0.1';

FLUSH PRIVILEGES;
