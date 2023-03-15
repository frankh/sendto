CREATE TABLE IF NOT EXISTS Messages (
	ID         TEXT    PRIMARY KEY,
	CipherText TEXT    NOT NULL,
	ExpiresAt  INTEGER NOT NULL DEFAULT (strftime('%s', 'now') + 2592000),
	FromUser   TEXT    NOT NULL,
	ToUser	   TEXT    NOT NULL
);
