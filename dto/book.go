package dto

type Book struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Authors     []string          `json:"authors"`
	Published   int               `json:"published"`
	ISBN        string            `json:"isbn,omitempty"`
	Subjects    []string          `json:"subjects"`
	Covers      map[string]string `json:"covers"`
	Description string            `json:"description"`
}
