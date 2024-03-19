CREATE TABLE IF NOT EXISTS Accounts (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    nick VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    authorization_token VARCHAR(128) NOT NULL,
    alive DATE DEFAULT (CURRENT_DATE()),
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS Switches (
    account_id INT UNSIGNED NOT NULL,
    id INT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    subject TEXT NOT NULL,
    run_after INT UNSIGNED NOT NULL,
    recipients TEXT,
    FOREIGN KEY (account_id) REFERENCES Accounts(id) ON DELETE CASCADE,
    PRIMARY KEY (account_id, id)
);
