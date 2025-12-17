package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to db")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()

}

func createTables() {

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	avatar_id TEXT,
	created_at DATETIME NOT NULL
	);
	`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("could not create users table.")
	}

	createClubsTable := `
	CREATE TABLE IF NOT EXISTS clubs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_by INTEGER,
	is_private BOOL,
	name TEXT NOT NULL,
	description TEXT,
	created_at DATETIME NOT NULL,
	FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(createClubsTable)
	if err != nil {
		panic("could not create clubs table")
	}

	createMembersTable := `
	CREATE TABLE IF NOT EXISTS members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    clubid INTEGER NOT NULL,
    role TEXT NOT NULL,
    joined_at DATETIME NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
    UNIQUE(userid, clubid)
	);
	`
	_, err = DB.Exec(createMembersTable)
	if err != nil {
		panic("could not create members table")
	}

	createLeaderboardTable := `
	CREATE TABLE IF NOT EXISTS leaderboard (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    clubid INTEGER NOT NULL,
    score INTEGER NOT NULL DEFAULT 0,
    last_checkedin DATETIME,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
    UNIQUE(userid, clubid)
	);
	`
	_, err = DB.Exec(createLeaderboardTable)
	if err != nil {
		panic("could not create leaderboard table")
	}

	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
    messageid INTEGER PRIMARY KEY AUTOINCREMENT,
    clubid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    timestamp DATETIME NOT NULL,
    message TEXT NOT NULL,
    FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(createMessagesTable)
	if err != nil {
		panic("could not create messages table")
	}

	createStreaksTable := `
	CREATE TABLE IF NOT EXISTS streaks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clubid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    currentstreak INTEGER NOT NULL DEFAULT 0,
	longeststreak INTEGER NOT NULL DEFAULT 0,
	lastcheckedin DATETIME,
    FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = DB.Exec(createStreaksTable)
	if err != nil {
		panic("could not create streaks table")
	}




	//TODO remove this
	// createEventsTable := `
	// CREATE TABLE IF NOT EXISTS events (
	// id INTEGER PRIMARY KEY AUTOINCREMENT,
	// name TEXT NOT NULL,
	// description TEXT NOT NULL,
	// location TEXT NOT NULL,
	// createdAt DATETIME NOT NULL,
	// user_id INTEGER,
	// FOREIGN KEY(user_id) REFERENCES users(id)
	// )
	// `

	// _, err = DB.Exec(createEventsTable)

	// if err != nil {
	// 	panic("could not create events table.")
	// }

	// createRegistrationsTable := `
	// CREATE TABLE IF NOT EXISTS registrations (
	// id INTEGER PRIMARY KEY AUTOINCREMENT,
	// event_id INTEGER,
	// user_id INTEGER,
	// FOREIGN KEY(event_id) REFERENCES events(id),
	// FOREIGN KEY(user_id) REFERENCES users(id)
	// )
	// `
	// _, err = DB.Exec(createRegistrationsTable)
	// if err != nil {
	// 	panic("could not create registration table")
	// }

}
