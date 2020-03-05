package app

// Config - Application config
type Config struct {
	IsMuni IsMuni
	Parser string
}

//IsMuni - Is muni config
type IsMuni struct {
	URL string `json:"url" yaml:"url"`
	Token string `json:"token" yaml:"token"`
	Course string `json:"course" yaml:"course"`
	FacultyID int `json:"faculty_id" yaml:"faculty_id"`
}




