package models

type Book struct {
	Id        int
	Isbn      string
	Title     string
	Image     string
	Describe  string
	Author    Author
	Genre     Genre
	Publisher Publisher
}

type BookBaseInfo struct {
	Title string
	Image string
	Url   string
}
