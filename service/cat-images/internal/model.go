package internal

type Response[T any] struct {
	Data    []T    `json:"data"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
