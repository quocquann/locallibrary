package models

type Book struct {
	Title  string
	Image  string
	Author Author
	Genre  string
}

type Author struct {
	Name string
}

type BookBaseInfo struct {
	Title string
	Image string
	Url   string
}
