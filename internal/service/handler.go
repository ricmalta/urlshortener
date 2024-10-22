package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

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
	logger   *logrus.Logger
	baseURL  string
	Router   *mux.Router
}

func NewHandler(urlStore *store.Store, logger *logrus.Logger, serviceBaseURL string) *Handler {
	handler := &Handler{
		urlStore: urlStore,
		logger:   logger,
		baseURL:  serviceBaseURL,
		Router:   mux.NewRouter(),
	}
	handler.Router.HandleFunc("/", handler.addURL).Methods("POST")
	handler.Router.HandleFunc("/{tinyURL}", handler.getURL).Methods("GET")

	handler.Router.NotFoundHandler = http.HandlerFunc(handler.notFoundHandler)
	handler.Router.MethodNotAllowedHandler = http.HandlerFunc(handler.notFoundHandler)

	return handler
}

func (s *Handler) addURL(w http.ResponseWriter, r *http.Request) {
	var reqPayload AddURLRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPayload); err != nil {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
		s.serializeJSON(w, payload.Status, payload)
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
			s.serializeJSON(w, payload.Status, payload)
			s.logger.Error(err)
			return
		default:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			s.logger.Error(err)
			s.serializeJSON(w, payload.Status, payload)
			return
		}
	}

	respPayload := AddURLResponse{
		TinyURL: fmt.Sprintf("%s/%s", s.baseURL, shortKey),
	}

	s.serializeJSON(w, http.StatusCreated, respPayload)
	return
}

func (s *Handler) getURL(w http.ResponseWriter, r *http.Request) {
	tinyURL, ok := getTinyURLParam(r)
	if !ok {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
		s.logger.Error("parsing tiny url")
		s.serializeJSON(w, payload.Status, payload)
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
			s.serializeJSON(w, payload.Status, payload)
			return
		default:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			s.logger.Error(err)
			s.serializeJSON(w, payload.Status, payload)
			return
		}
	}
	http.Redirect(w, r, originalURL, 301)
	return
}

func (s *Handler) serializeJSON(w http.ResponseWriter, status int, data interface{}) {
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

func (s *Handler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	payload := HTTPError{
		Error:  "Not Found",
		Status: http.StatusNotFound,
	}
	s.serializeJSON(w, payload.Status, payload)
	return
}

func getTinyURLParam(r *http.Request) (value string, ok bool) {
	ok = false
	if tinyURL, ok := mux.Vars(r)["tinyURL"]; ok {
		return tinyURL, ok
	}
	// workaround to support net/httptest
	if parts := strings.Split(r.URL.Path, "/"); len(parts) > 0 {
		value = parts[1]
		ok = true
	}
	return
}
