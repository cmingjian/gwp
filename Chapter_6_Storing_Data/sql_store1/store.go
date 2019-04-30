package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id      int64
	Content string
	Author  string
}

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("mysql", "cmj:123456a@tcp(127.0.0.1:3306)/gwp")
	if err != nil {
		panic(err)
	}
}

// get all posts
func Posts(limit int) (posts []Post, err error) {
	rows, err := Db.Query("select id, content, author from posts limit ?", limit)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Get a single post
func GetPost(id int64) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, content, author from posts where id = ?", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	statement := "insert into posts (content, author) values (?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(post.Content, post.Author)
	if err != nil {
		panic(err)
	}
	post.Id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

// Update a post
func (post *Post) Update() (err error) {
	_, err = Db.Exec("update posts set content = ?, author = ? where id = ?", post.Content, post.Author, post.Id)
	if err != nil {
		panic(err)
	}
	return
}

// Delete a post
func (post *Post) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = ?", post.Id)
	if err != nil {
		panic(err)
	}
	return
}

// Delete all posts
func DeleteAll() (err error) {
	_, err = Db.Exec("delete from posts")
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	DeleteAll()
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	// Create a post
	fmt.Println(post) // {0 Hello World! Sau Sheong}
	post.Create()
	fmt.Println(post) // {1 Hello World! Sau Sheong}

	// Get one post
	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost) // {1 Hello World! Sau Sheong}

	// Update the post
	readPost.Content = "Bonjour Monde!"
	readPost.Author = "Pierre"
	readPost.Update()

	// Get all posts
	posts, _ := Posts(10)
	fmt.Println(posts) // [{1 Bonjour Monde! Pierre}]

	// Delete the post
	readPost.Delete()

	// Get all posts
	posts, _ = Posts(10)
	fmt.Println(posts) // []

	// Delete all posts
  	DeleteAll()
}
