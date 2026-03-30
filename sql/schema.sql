
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP
);

CREATE TABLE master_status (
  id INT AUTO_INCREMENT PRIMARY KEY,
  status VARCHAR(50) NOT NULL
);


CREATE TABLE tasks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL ,
  description TEXT,
  status_id INT,
  user_id INT,
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  
  FOREIGN KEY(status_id) REFERENCES master_status(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);


