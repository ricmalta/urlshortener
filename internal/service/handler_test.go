package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ricmalta/urlshortner/internal/store"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

const (
	addValidURL            = "https://www.example.com/test"
	addInvalidURL          = "--https://www.example.com/test"
	getNonexistentShortURL = "/aaaaa"
	baseURL                = "http://tiny.test.com"
	tinyURL                = "http://tiny.test.com/70"
)

var (
	cache          *lru.Cache
	mr             *miniredis.Miniredis
	redisClient    *redis.Client
	logger         *logrus.Logger
	storeInstance  *store.Store
	serviceHandler *Handler
)

func getParams(s string) string {
	res := strings.Split(s, "/")
	if len(res) > 0 {
		return res[len(res)-1]
	}

	return ""
}

func init() {
	cache, err := lru.New(10)
	if err != nil {
		panic(err)
	}
	mr, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	logger, _ = test.NewNullLogger()
	storeInstance, err = store.NewStore(cache, redisClient, logger)
	if err != nil {
		panic(err)
	}
	serviceHandler = NewHandler(storeInstance, logger, baseURL)
}

func TestAddValidURL(t *testing.T) {
	payloadBytes, err := json.Marshal(AddURLRequest{
		URL: addValidURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler.addURL)

	handler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusCreated, respRecorder.Code)
	var response AddURLResponse
	if err := json.Unmarshal(respRecorder.Body.Bytes(), &response); err != nil {
		t.Errorf("expect POST / to return AddURLResponse type")
	}
}

func TestAddInvalidURL(t *testing.T) {
	payloadBytes, err := json.Marshal(AddURLRequest{
		URL: addInvalidURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler.addURL)

	handler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
	var response HTTPError
	if err := json.Unmarshal(respRecorder.Body.Bytes(), &response); err != nil {
		t.Errorf("expect POST / to return HTTPError type")
	}
}

func TestGetInvalidURL(t *testing.T) {
	req, err := http.NewRequest("GET", getNonexistentShortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler.getURL)

	handler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusNotFound, respRecorder.Code)
	var response HTTPError
	if err := json.Unmarshal(respRecorder.Body.Bytes(), &response); err != nil {
		t.Errorf("expect POST / to return HTTPError type")
	}
}

func TestGetValidURL(t *testing.T) {
	payloadBytes, err := json.Marshal(AddURLRequest{
		URL: addValidURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler.addURL)

	handler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusCreated, respRecorder.Code)
	var AddResponse AddURLResponse
	if err := json.Unmarshal(respRecorder.Body.Bytes(), &AddResponse); err != nil {
		t.Errorf("expect POST / to return AddURLResponse type")
	}

	key := getParams(AddResponse.TinyURL)
	if key == "" {
		t.Error("cannot extract created short url key")
	}
	req, err = http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
	if err != nil {
		t.Fatal(err)
	}

	respRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(serviceHandler.getURL)

	handler.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusMovedPermanently, respRecorder.Code)
}

func TestGetURLParam(t *testing.T) {
	url, err := url.Parse(tinyURL)
	if err != nil {
		t.Error(err)
	}
	request := http.Request{URL: url}
	tiny, ok := getTinyURLParam(&request)
	assert.Equal(t, "70", tiny)
	assert.Equal(t, true, ok)
}
