package models

type Job struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Company  string   `json:"company"`
	Location Location `json:"location"`
}

type Location struct {
	City       string     `json:"city"`
	Department Department `json:"department"`
}

type Department struct {
	IsoCode string `json:"isoCode"`
}
