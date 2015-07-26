/* *****************************************************************************
// Create the tables
// ****************************************************************************/
CREATE TABLE user_status (
    id INTEGER PRIMARY KEY,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted UNSIGNED TINYINT(1) NOT NULL DEFAULT 0
);

CREATE TRIGGER user_status_update
AFTER UPDATE
ON user_status
FOR EACH ROW
BEGIN
UPDATE user_status SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;

CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    
    status_id INTEGER NOT NULL DEFAULT 1,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted UNSIGNED TINYINT(1) NOT NULL DEFAULT 0,
    
    UNIQUE (email),
    CONSTRAINT `f_user_user_status` FOREIGN KEY (`status_id`) REFERENCES `user_status` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER user_update
AFTER UPDATE
ON user
FOR EACH ROW
BEGIN
UPDATE user SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;

INSERT INTO `user_status` (`id`, `status`, `created_at`, `updated_at`, `deleted`) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);