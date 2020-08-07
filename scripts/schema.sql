-- Initialize the database.

CREATE TABLE user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL DEFAULT '123456',
  realname TEXT NOT NULL,
  dept_id INTEGER NOT NULL,
  role BOOLEAN NOT NULL DEFAULT 0,
  permission TEXT NOT NULL DEFAULT '',
  FOREIGN KEY (dept_id) REFERENCES department (id)
);

CREATE TABLE department (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  dept_name TEXT UNIQUE NOT NULL
);

CREATE TABLE record (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  dept_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  date DATE NOT NULL,
  type BOOLEAN NOT NULL,
  duration INTEGER NOT NULL,
  describe TEXT NOT NULL DEFAULT '',
  status INTEGER NOT NULL DEFAULT 0,
  created TIMESTAMP NOT NULL DEFAULT (datetime('now', 'localtime')),
  createdby TEXT,
  verifiedby TEXT,
  FOREIGN KEY (user_id) REFERENCES user (id),
  FOREIGN KEY (dept_id) REFERENCES department (id)
);

CREATE VIEW employee AS
  SELECT u.id, username, realname, dept_id, dept_name, role, permission
  FROM user u
  JOIN department d ON d.id = dept_id
  ORDER BY dept_name, realname;

CREATE VIEW statistics AS
  SELECT strftime('%Y-%m', date) period, r.dept_id, dept_name, user_id, realname,
  sum(CASE WHEN r.type = 1 THEN duration ELSE 0 END) overtime,
  sum(CASE WHEN r.type = 0 THEN 0 - duration ELSE 0 END) leave,
  sum(duration) summary
  FROM record r
  JOIN employee e ON e.id = user_id
  WHERE status = 1
  GROUP BY period, r.dept_id, user_id
  ORDER BY period DESC, dept_name, realname;

INSERT INTO user (id, username, realname, dept_id, role)
  VALUES (0, 'root', 'root', 0, 1);
