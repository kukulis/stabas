package entities

type Participant struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
}
