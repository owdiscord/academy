-- migrate:up
CREATE TABLE IF NOT EXISTS waves (
  id INT AUTO_INCREMENT PRIMARY KEY,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  -- One of 'interviews', 'helper', 'historic',
  -- set to 'interviews' when we start, only showing the interview questions, then
  -- set to 'helper' when we want helpers to be managing things, then 'historic'
  -- when the wave ends and trianees are promoted.
  state VARCHAR(32) NOT NULL DEFAULT 'interviews',
  begin_at DATETIME NOT NULL DEFAULT NOW(),
  close_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS staff (
  id INT AUTO_INCREMENT PRIMARY KEY,
  snowflake VARCHAR(22) NOT NULL,
  username VARCHAR(128) NOT NULL,
  display_name VARCHAR(512) NOT NULL,
  thread_participation_count INT NOT NULL DEFAULT 0,
  message_count INT NOT NULL DEFAULT 0,
  thread_count INT NOT NULL DEFAULT 0,
  case_count INT NOT NULL DEFAULT 0,
  wave_id INT NOT NULL REFERENCES waves(id),
  -- One of 'trainee', 'moderator', 'helper', or 'admin'
  role VARCHAR(64) NOT NULL DEFAULT 'trainee'
);

CREATE TABLE IF NOT EXISTS issues (
  id INT AUTO_INCREMENT PRIMARY KEY,
  wave_id INT REFERENCES waves(id),
  created_by VARCHAR(22) REFERENCES staff(snowflake),
  trainee_id VARCHAR(22) REFERENCES staff(snowflake),
  thread_id VARCHAR(36) NULL DEFAULT NULL,
  message_id VARCHAR(36) NULL DEFAULT NULL,
  -- One of 'pending', 'handled', 'archived', or 'deleted'
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  reason TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS interview_questions (
  id INT PRIMARY KEY,
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NULL DEFAULT NOW(),
  text VARCHAR(512) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id INT AUTO_INCREMENT,
  token VARCHAR(256) NOT NULL UNIQUE,
  user_id INT REFERENCES staff(id),
  wave_id INT REFERENCES waves(id),
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  KEY idx_expires_at (expires_at),
  KEY idx_user_id (user_id)
);

CREATE TABLE IF NOT EXISTS stats_per_date (
  id INT AUTO_INCREMENT,
  date TIMESTAMP NOT NULL DEFAULT now(),
  staff_id INT REFERENCES staff(id),
  wave_id INT REFERENCES waves(id),
  public_messages INT DEFAULT 0,
  private_messages INT DEFAULT 0,
  cases INT DEFAULT 0,
  thread_chat INT DEFAULT 0,
  thread_replies INT DEFAULT 0,
  thread_closures INT DEFAULT 0,
  snippets_used INT DEFAULT 0,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS collection_log (
  id INT AUTO_INCREMENT,
  threads_imported INT DEFAULT 0,
  cases_imported INT DEFAULT 0,
  messages_imported INT DEFAULT 0,
  run_at TIMESTAMP DEFAULT NOW(),

  PRIMARY KEY (id)
);

-- migrate:down
DROP TABLE IF EXISTS collection_log;
DROP TABLE IF EXISTS stats_per_date;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS interview_questions;
DROP TABLE IF EXISTS waves;
