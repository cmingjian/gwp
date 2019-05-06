package main

import (
	"database/sql"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("mysql", "cmj:123456a@tcp(127.0.0.1:3306)/gwp")
	if err != nil {
		panic(err)
	}
}

// Get a single post
func retrieve(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, content, author from posts where id = ?", id).Scan(&post.Id, &post.Content, &post.Author)
	if err != nil {
		panic(err)
	}
	return
}

// Create a new post
func (post *Post) create() (err error) {
	statement := "insert into posts (content, author) values (?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	Res, err := stmt.Exec(post.Content, post.Author)
	if err != nil {
		panic(err)
	}
	post.Id, err = Res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

// Update a post
func (post *Post) update() (err error) {
	_, err = Db.Exec("update posts set content = ?, author = ? where id = ?",
		post.Content, post.Author, post.Id)
	if err != nil {
		panic(err)
	}
	return
}

// Delete a post
func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where id = ?", post.Id)
	if err != nil {
		panic(err)
	}
	return
}
