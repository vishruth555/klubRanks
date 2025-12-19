package db

import (
	"klubRanks/config"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	driver := config.AppConfig.Database.Driver
	dsn := config.AppConfig.Database.DSN

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch driver {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(dsn), gormConfig)

	case "postgres":
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)

	default:
		log.Fatalf("unsupported db driver: %s", driver)
	}

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Optional but recommended: tune connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("database connected via GORM")

}

// func createTables(driver string) error {
// 	stmts := []string{}

// 	switch driver {

// 	/* ---------------- SQLITE ---------------- */
// 	case "sqlite3":
// 		stmts = append(stmts,
// 			`CREATE TABLE IF NOT EXISTS users (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				username TEXT NOT NULL UNIQUE,
// 				password TEXT NOT NULL,
// 				avatar_id TEXT DEFAULT 'default',
// 				created_at DATETIME NOT NULL
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS clubs (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				created_by INTEGER,
// 				is_private BOOLEAN,
// 				name TEXT NOT NULL,
// 				description TEXT,
// 				created_at DATETIME NOT NULL,
// 				FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS members (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				userid INTEGER NOT NULL,
// 				clubid INTEGER NOT NULL,
// 				role TEXT NOT NULL,
// 				joined_at DATETIME NOT NULL,
// 				FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
// 				FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
// 				UNIQUE(userid, clubid)
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS leaderboard (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				userid INTEGER NOT NULL,
// 				clubid INTEGER NOT NULL,
// 				score INTEGER NOT NULL DEFAULT 0,
// 				last_checkedin DATETIME,
// 				current_streak INTEGER NOT NULL DEFAULT 0,
// 				longest_streak INTEGER NOT NULL DEFAULT 0,
// 				FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
// 				FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
// 				UNIQUE(userid, clubid)
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS messages (
// 				messageid INTEGER PRIMARY KEY AUTOINCREMENT,
// 				clubid INTEGER NOT NULL,
// 				userid INTEGER NOT NULL,
// 				timestamp DATETIME NOT NULL,
// 				message TEXT NOT NULL,
// 				FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
// 				FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS streaks (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				clubid INTEGER NOT NULL,
// 				userid INTEGER NOT NULL,
// 				currentstreak INTEGER NOT NULL DEFAULT 0,
// 				longeststreak INTEGER NOT NULL DEFAULT 0,
// 				lastcheckedin DATETIME,
// 				FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
// 				FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS activity_log (
// 				id INTEGER PRIMARY KEY AUTOINCREMENT,
// 				clubid INTEGER NOT NULL,
// 				userid INTEGER NOT NULL,
// 				points INTEGER NOT NULL,
// 				timestamp DATETIME NOT NULL,
// 				FOREIGN KEY (clubid) REFERENCES clubs(id) ON DELETE CASCADE,
// 				FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
// 			);`,
// 		)

// 	/* ---------------- POSTGRES (SUPABASE) ---------------- */
// 	case "postgres":
// 		stmts = append(stmts,
// 			`CREATE TABLE IF NOT EXISTS users (
// 				id BIGSERIAL PRIMARY KEY,
// 				username TEXT NOT NULL UNIQUE,
// 				password TEXT NOT NULL,
// 				avatar_id TEXT DEFAULT 'default',
// 				created_at TIMESTAMP NOT NULL
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS clubs (
// 				id BIGSERIAL PRIMARY KEY,
// 				created_by BIGINT REFERENCES users(id) ON DELETE CASCADE,
// 				is_private BOOLEAN,
// 				name TEXT NOT NULL,
// 				description TEXT,
// 				created_at TIMESTAMP NOT NULL
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS members (
// 				id BIGSERIAL PRIMARY KEY,
// 				userid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
// 				clubid BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
// 				role TEXT NOT NULL,
// 				joined_at TIMESTAMP NOT NULL,
// 				UNIQUE(userid, clubid)
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS leaderboard (
// 				id BIGSERIAL PRIMARY KEY,
// 				userid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
// 				clubid BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
// 				score INTEGER NOT NULL DEFAULT 0,
// 				last_checkedin TIMESTAMP,
// 				current_streak INTEGER NOT NULL DEFAULT 0,
// 				longest_streak INTEGER NOT NULL DEFAULT 0,
// 				UNIQUE(userid, clubid)
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS messages (
// 				messageid BIGSERIAL PRIMARY KEY,
// 				clubid BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
// 				userid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
// 				timestamp TIMESTAMP NOT NULL,
// 				message TEXT NOT NULL
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS streaks (
// 				id BIGSERIAL PRIMARY KEY,
// 				clubid BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
// 				userid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
// 				currentstreak INTEGER NOT NULL DEFAULT 0,
// 				longeststreak INTEGER NOT NULL DEFAULT 0,
// 				lastcheckedin TIMESTAMP
// 			);`,

// 			`CREATE TABLE IF NOT EXISTS activity_log (
// 				id BIGSERIAL PRIMARY KEY,
// 				clubid BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
// 				userid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
// 				points INTEGER NOT NULL,
// 				timestamp TIMESTAMP NOT NULL
// 			);`,
// 		)

// 	default:
// 		return fmt.Errorf("unsupported database driver: %s", driver)
// 	}

// 	for _, stmt := range stmts {
// 		if tx := DB.Exec(stmt); tx.Error != nil {
// 			return tx.Error
// 		}
// 	}

// 	return nil
// }
