package main

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type Car struct {
	Maker string `json:"maker"`
	Year  int    `json:"year"`
}


