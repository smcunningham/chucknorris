package models

// NamesResponse contains name returned from NameAPI
type NamesResponse struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender,omitempty"`
	Region  string `json:"region,omitempty"`
}

// JokesResponse contains the joke data returned from JokeAPI
type JokesResponse struct {
	TypeStatus string `json:"type"`
	Value      Value  `json:"value"`
}

// Value contains actual joke data
type Value struct {
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}
