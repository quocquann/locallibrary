package models

type Book struct {
	Id     int
	Isbn   string
	Title  string
	Image  string
	Author Author
	Genre  string
}

type BookBaseInfo struct {
	Title string
	Image string
	Url   string
}
