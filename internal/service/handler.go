package service

import (
	"encoding/json"
  "fmt"
  "github.com/sirupsen/logrus"
  "strings"

  "github.com/gorilla/mux"
	"github.com/ricmalta/urlshortner/internal/store"

	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

var (
	httpInternalServerError = HTTPError{Error: "Internal Server Error", Status: http.StatusInternalServerError}
	httpBadRequest          = HTTPError{Error: "Bad Request", Status: http.StatusBadRequest}
)

type AddURLRequest struct {
	URL string `json:"url"`
}

type AddURLResponse struct {
	TinyURL string `json:"tiny_url"`
}

type HTTPError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type Handler struct {
	urlStore *store.Store
  logger *logrus.Logger
  baseURL string
  Router *mux.Router
}

func NewHandler(urlStore *store.Store, logger *logrus.Logger,serviceBaseURL string) *Handler {
	handler := &Handler{
		urlStore: urlStore,
		logger: logger,
		baseURL: serviceBaseURL,
    Router: mux.NewRouter(),
	}
  handler.Router.HandleFunc("/", handler.AddURL).Methods("POST")
  handler.Router.HandleFunc("/{tinyURL}", handler.GetURL).Methods("GET")

  handler.Router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
  handler.Router.MethodNotAllowedHandler = http.HandlerFunc(handler.NotFoundHandler)

	return handler
}

func (s *Handler) AddURL(w http.ResponseWriter, r *http.Request) {
	var reqPayload AddURLRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPayload); err != nil {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
		s.SerializeJSON(w, payload.Status, payload)
		s.logger.Errorf("error decoding url to encode %v", err)
		return
	}

	shortKey, err := s.urlStore.Add(reqPayload.URL)
	if err != nil {
    switch err.(type) {
    case store.ErrorInvalidInputURL:
      payload := HTTPError{
        Error:  err.Error(),
        Status: http.StatusBadRequest,
      }
      s.SerializeJSON(w, payload.Status, payload)
      s.logger.Error(err)
      return
    default:
      payload := HTTPError{
        Error:  err.Error(),
        Status: http.StatusInternalServerError,
      }
      s.logger.Error(err)
      s.SerializeJSON(w, payload.Status, payload)
      return
    }
	}

	respPayload := AddURLResponse{
		TinyURL: fmt.Sprintf("%s/%s", s.baseURL, shortKey),
	}

	s.SerializeJSON(w, http.StatusCreated, respPayload)
	return
}

func (s *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
  parts := strings.Split(r.URL.Path, "/")
  var tinyURL string
  if len(parts) > 0 {
    tinyURL = parts[1]
  }
	if tinyURL == "" {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
    s.logger.Error("parsing tiny url")
		s.SerializeJSON(w, payload.Status, payload)
		return
	}

  originalURL, err := s.urlStore.Get(tinyURL)
	if err != nil {
		switch err.(type) {
		case store.ErrorNotStoredShortURL:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusNotFound,
			}
      s.logger.Warnf("tiny url %s not found", tinyURL)
			s.SerializeJSON(w, payload.Status, payload)
			return
		default:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
      s.logger.Error(err)
			s.SerializeJSON(w, payload.Status, payload)
			return
		}
	}
	http.Redirect(w, r, originalURL,301)
	return
}

func (s *Handler) SerializeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(status)

	if data != nil {
		encoder := json.NewEncoder(w)

		if err := encoder.Encode(data); err != nil {
			encoder.Encode(HTTPError{
				Error:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			s.logger.Error(err)
		}
	}
}

func (s *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	payload := HTTPError{
		Error:  "Not Found",
		Status: http.StatusNotFound,
	}
	s.SerializeJSON(w, payload.Status, payload)
	return
}
