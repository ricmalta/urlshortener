package main_test

import (
  "github.com/ricmalta/urlshortner/internal/config"
  "github.com/ricmalta/urlshortner/internal/service"
  "github.com/ricmalta/urlshortner/internal/store"
  "net/http"
  "net/http/httptest"
  "testing"
  "errors"

  "github.com/stretchr/testify/assert"
)

func setup() (*httptest.Server, *http.Client, error) {
  cfg, err := config.LoadConfig("./internal/config")
  if err != nil {
    return nil, nil, err
  }

  serviceStore, err := store.NewStore(cfg)
  if err != nil {
    return nil, nil, err
  }

  handler := service.NewServiceHandler(serviceStore)

  server := httptest.NewServer(handler)
  client := http.Client{}

  return server, &client, nil
}

func TestCreateTinyURL(t *testing.T) {
  err := errors.New("ww")
  assert.Equal(t, 123, 124)
  assert.Error(t, err)
}



/*for i := range []int{0, 1,2,3,4,5,6,7,8,9} {
go func(i int) {
shortKey, err := storeInstance.Add("http://www.example.com/from_the_instance_add")
if err != nil {
fmt.Println("err", err)
}
fmt.Println("shortKey", i, shortKey)
}(i)
}*/
