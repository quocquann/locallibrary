package models

type Book struct {
	Id     int    `json:"id"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Image  string `json:"image"`
	Author Author `json: "author"`
	Genre  string `json:"genre"`
}

type BookBaseInfo struct {
	Title string
	Image string
	Url   string
}
