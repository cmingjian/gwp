package main

import (
	"database/sql"
	"errors"
	"fmt"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id       int64
	Content  string
	Author   string
	Comments []Comment
}

type Comment struct {
	Id      int64
	Content string
	Author  string
	Post    *Post
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

func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("Post not found")
		return
	}
	res, err := Db.Exec("insert into comments (content, author, post_id) values (?, ?, ?)", comment.Content, comment.Author, comment.Post.Id)
	if err != nil {
		panic(err)
	}
	comment.Id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

// Get a single post
func GetPost(id int64) (post Post, err error) {
	post = Post{}
	post.Comments = []Comment{}
	err = Db.QueryRow("select id, content, author from posts where id = ?", id).Scan(&post.Id, &post.Content, &post.Author)

	rows, err := Db.Query("select id, content, author from comments where post_id = ?", id)
	if err != nil {
		return
	}
	for rows.Next() {
		comment := Comment{Post: &post}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	rows.Close()
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	res, err := Db.Exec("insert into posts (content, author) values (?, ?)", post.Content, post.Author)
	if err != nil {
		panic(err)
	}
	post.Id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	post.Create()

	// Add a comment
	comment := Comment{Content: "Good post!", Author: "Joe", Post: &post}
	comment.Create()
	readPost, _ := GetPost(post.Id)

	fmt.Println(readPost)                  // {1 Hello World! Sau Sheong [{1 Good post! Joe 0xc20802a1c0}]}
	fmt.Println(readPost.Comments)         // [{1 Good post! Joe 0xc20802a1c0}]
	fmt.Println(readPost.Comments[0].Post) // &{1 Hello World! Sau Sheong [{1 Good post! Joe 0xc20802a1c0}]}
}
