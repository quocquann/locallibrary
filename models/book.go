package models

type Book struct {
	Title  string
	Image  string
	Author string
	Genre  string
}

type BookBaseInfo struct {
	Title string
	Image string
	Url   string
}
