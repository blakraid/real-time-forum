package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Sql *sql.DB

func CreateDatabase() error {
	var err error
	Sql , err = sql.Open("sqlite3","./data/ForumDatabase.db")
	if err != nil {
		return err
	}
	if err = CreateTabel(); err != nil {
		return err
	}
	return nil
}


func CreateTabel() error{
	// Users table
	_, err := Sql.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL ,
            email TEXT UNIQUE NOT NULL ,
            password TEXT NOT NULL,
			firtsName TEXT NOT NULL,
			lastName TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender TEXT NOT NULL,
			session_token TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Printf("Error creating 'users' table: %v", err)
		return err
	} else {
		log.Println("'users' table created or already exists")
	}

	// Posts table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (session_token) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Printf("Error creating 'posts' table: %v", err)
		return err
	} else {
		log.Println("'posts' table created or already exists")
	}

	// Categories table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL
        );
    `)
	if err != nil {
		log.Printf("Error creating 'categories' table: %v", err)
		return err
	} else {
		log.Println("'categories' table created or already exists")
	}

	// Post categories relationship table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS post_categories (
            post_id INTEGER NOT NULL,
            category_id INTEGER NOT NULL,
            PRIMARY KEY (post_id, category_id),
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
            FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Printf("Error creating 'post_categories' table: %v", err)
		return err
	} else {
		log.Println("'post_categories' table created or already exists")
	}

	// Comments table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Printf("Error creating 'comments' table: %v", err)
		return err
	} else {
		log.Println("'comments' table created or already exists")
	}

	// post_Likes table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS post_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            post_id INTEGER NOT NULL,
            is_like BOOLEAN NOT NULL,	
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			UNIQUE (user_id, post_id)
        );
    `)
	if err != nil {
		log.Printf("Error creating 'post_likes' table: %v", err)
		return err
	} else {
		log.Println("'post_likes' table created or already exists")
	}

	// comment_Likes table
	_, err = Sql.Exec(`
        CREATE TABLE IF NOT EXISTS comment_likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            comment_id INT NOT NULL,
            is_like BOOLEAN NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(comment_id, user_id),
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Printf("Error creating 'comment_likes' table: %v", err)
		return err
	} else {
		log.Println("'comment_likes' table created or already exists")
	}

	// Insert default categories
	_, err = Sql.Exec(`
        INSERT OR IGNORE INTO categories (name) VALUES 
        ('Technology'),
        ('Lifestyle'),
        ('Travel'),
        ('Food'),
		('Other'),
		('Study')
    `)
	if err != nil {
		log.Printf("Error inserting default categories: %v", err)
		return err
	} else {
		log.Println("Default categories inserted or already exist")
	}

	// Create messages table
	_, err = Sql.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender TEXT,
			receiver TEXT,
			text TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Error inserting default messages: %v", err)
		return err
	} else {
		log.Println("Default messages inserted or already exist")
	}

	return nil
}