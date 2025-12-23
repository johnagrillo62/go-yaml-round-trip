package nested

type Car struct {
	Maker string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Cars []Car  `json:"cars"`
}

type Company struct {
	Name  string `json:"name"`
	CEO   User   `json:"ceo"`
	Fleet []Car  `json:"fleet"`
}


