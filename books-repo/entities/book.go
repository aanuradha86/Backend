package entities

// Book represents a single book
type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"writer"`
	ISBN        string `json:"ISBNCode"`
	Description string `json:"synopsis"`
	Price       string `json:"price"`
}
