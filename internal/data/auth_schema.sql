CREATE TABLE IF NOT EXISTS identity_user (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_no VARCHAR(32) NOT NULL,
  phone_cipher VARCHAR(255) NOT NULL DEFAULT '',
  phone_hash CHAR(64) NULL,
  email VARCHAR(128) NULL,
  password_hash VARCHAR(255) NOT NULL,
  role_code VARCHAR(32) NOT NULL DEFAULT 'user',
  account_status VARCHAR(32) NOT NULL DEFAULT 'normal',
  locked_until DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_identity_user_user_no (user_no),
  UNIQUE KEY uk_identity_user_phone_hash (phone_hash),
  UNIQUE KEY uk_identity_user_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS identity_real_name_auth (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  real_name_cipher VARCHAR(255) NOT NULL,
  id_card_cipher VARCHAR(255) NOT NULL,
  id_card_hash CHAR(64) NOT NULL,
  image_urls JSON NULL,
  auth_status VARCHAR(32) NOT NULL DEFAULT 'pending',
  reject_reason VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_identity_real_name_auth_user_id (user_id),
  UNIQUE KEY uk_identity_real_name_auth_id_card_hash (id_card_hash),
  CONSTRAINT fk_identity_real_name_auth_user FOREIGN KEY (user_id) REFERENCES identity_user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS identity_security_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  event_type VARCHAR(64) NOT NULL,
  result VARCHAR(32) NOT NULL,
  ip VARCHAR(64) NOT NULL DEFAULT '',
  device_id VARCHAR(128) NOT NULL DEFAULT '',
  failure_reason VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_identity_security_log_user_created_at (user_id, created_at),
  CONSTRAINT fk_identity_security_log_user FOREIGN KEY (user_id) REFERENCES identity_user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
