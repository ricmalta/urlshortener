package store

type ErrorNotStoredShortURL struct {
}

func (e ErrorNotStoredShortURL) Error() string {
	return "short url not found"
}
