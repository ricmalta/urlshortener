package store

type ErrorNotStoredShortURL struct {
}

func (e ErrorNotStoredShortURL) Error() string {
	return "short url not found"
}

type ErrorInvalidInputURL struct {
}

func (e ErrorInvalidInputURL) Error() string {
  return "invalid input URL"
}
