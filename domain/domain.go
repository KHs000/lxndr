package domain

type (
	// Error ..
	Error struct {
		Code  int
		Error string
	}

	// Response ..
	Response struct {
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}
)
