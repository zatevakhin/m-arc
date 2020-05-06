package database

// UsersTable - Statement to create `users` table
const UsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id       INTEGER      PRIMARY KEY NOT NULL,
	u_hash   CHAR(40)     UNIQUE      NOT NULL,
	login    VARCHAR(16)  UNIQUE      NOT NULL,
	password VARCHAR(128)             NOT NULL,
	group_id INTEGER                  NOT NULL,

	FOREIGN KEY (group_id) REFERENCES groups(id)
);`

// GroupsTable - Statement to create `groups` table
const GroupsTable = `
CREATE TABLE IF NOT EXISTS groups (
	id   INTEGER PRIMARY KEY NOT NULL,
	name VARCHAR(16) UNIQUE NOT NULL
);`

// AuthorizedTable - Statement to create `authorized` table
const AuthorizedTable = `
CREATE TABLE IF NOT EXISTS authorized (
	id      INTEGER     PRIMARY KEY      NOT NULL,
	u_id    INTEGER                      NOT NULL,
	ip      VARCHAR(16)                  NOT NULL,
	cookie  CHAR(40)                     NOT NULL,
	stamp   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (u_id) REFERENCES users(id)
);`

// AuthorsTable - Statement to create `authors` table
const AuthorsTable = `
CREATE TABLE IF NOT EXISTS authors (
	id    INTEGER      PRIMARY KEY NOT NULL,
	value VARCHAR(256) UNIQUE      NOT NULL
);`

// GenresTable - Statement to create `genres` table
const GenresTable = `
CREATE TABLE IF NOT EXISTS genres (
	id    INTEGER      PRIMARY KEY NOT NULL,
	value VARCHAR(256) UNIQUE      NOT NULL
);`

// TranslatorsTable - Statement to create `translators` table
const TranslatorsTable = `
CREATE TABLE IF NOT EXISTS translators (
	id    INTEGER      PRIMARY KEY NOT NULL,
	value VARCHAR(256) UNIQUE      NOT NULL
);`
