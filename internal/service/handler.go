package service

import (
	"encoding/json"
  "fmt"
  "os"

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

type serviceHandler struct {
	urlStore *store.Store
}

func NewServiceHandler(urlStore *store.Store) *mux.Router {
	handler := &serviceHandler{
		urlStore: urlStore,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.addURL).Methods("POST")
	router.HandleFunc("/{tinyURL}", handler.getURL).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router
}

func (s *serviceHandler) addURL(w http.ResponseWriter, r *http.Request) {
	var reqPayload AddURLRequest
	fmt.Printf("body %v", r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(reqPayload); err != nil {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
		serializeJSON(w, payload.Status, payload)
		fmt.Fprintf(os.Stderr, "error decoding url to encode %v\n", err)
		return
	}

	shortKey, err := s.urlStore.Add(reqPayload.URL)
	if err != nil {
		payload := HTTPError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
		serializeJSON(w, payload.Status, payload)
		return
	}

	respPayload := AddURLResponse{
		TinyURL: shortKey,
	}

	serializeJSON(w, http.StatusCreated, respPayload)
	return
}

func (s *serviceHandler) getURL(w http.ResponseWriter, r *http.Request) {
	tinyURL, ok := mux.Vars(r)["tinyURL"]
	if !ok {
		payload := HTTPError{
			Error:  "Bad Request",
			Status: http.StatusBadRequest,
		}
		serializeJSON(w, payload.Status, payload)
		return
	}
	url, err := s.urlStore.Get(tinyURL)
	if err != nil {

		switch err.(type) {
		case store.ErrorNotStoredShortURL:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusNotFound,
			}
			serializeJSON(w, payload.Status, payload)
			return
		default:
			payload := HTTPError{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			serializeJSON(w, payload.Status, payload)
			return

		}
	}
	http.Redirect(w, r, url, 301)
	return
}

func serializeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(status)

	if data != nil {
		encoder := json.NewEncoder(w)

		if err := encoder.Encode(data); err != nil {
			encoder.Encode(HTTPError{
				Error:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
		}
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	payload := HTTPError{
		Error:  "Not Found",
		Status: http.StatusNotFound,
	}
	serializeJSON(w, payload.Status, payload)
	return
}
