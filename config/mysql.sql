/* *****************************************************************************
// Setup the preferences
// ****************************************************************************/
SET NAMES utf8 COLLATE 'utf8_unicode_ci';
SET foreign_key_checks = 1;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
SET default_storage_engine = InnoDB;
SET CHARACTER SET utf8;

/* *****************************************************************************
// Remove old database
// ****************************************************************************/
DROP DATABASE IF EXISTS gowebapp;

/* *****************************************************************************
// Create new database
// ****************************************************************************/
CREATE DATABASE gowebapp DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
USE gowebapp;

/* *****************************************************************************
// Create the tables
// ****************************************************************************/
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);

CREATE TABLE user (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    
    status_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    UNIQUE KEY (email),
    CONSTRAINT `f_user_status` FOREIGN KEY (`status_id`) REFERENCES `user_status` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);

INSERT INTO `user_status` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);

CREATE TABLE note (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    content TEXT NOT NULL,
    
    user_id INT(10) UNSIGNED NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    CONSTRAINT `f_note_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);