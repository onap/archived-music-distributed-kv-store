/*
 * Copyright 2018 Intel Corporation, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RouterConfig() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/config", HandleConfigUpload).Methods("POST")
	router.HandleFunc("/v1/config/{token}/{filename}", HandleConfigGet).Methods("GET")
	router.HandleFunc("/v1/config/{token}/{filename}", HandleConfigDelete).Methods("DELETE")
	router.HandleFunc("/v1/config/load", HandleConfigLoad).Methods("POST")
	router.HandleFunc("/v1/config/load-default", HandleDefaultConfigLoad).Methods("GET")
	return router
}

func TestHandleConfigGet(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("GET", "/v1/config/token1/filename1", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}
func TestHandleConfigDelete(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/config/token1/filename1", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleConfigDelete_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectoryErr{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/config/token1/filename1", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleConfigPOST(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &LoadConfigBody{
		Token:     "test",
		Filename:  "test",
		Subdomain: "test",
	}

	b, _ := json.Marshal(body)

	// json Marshal converts struct to json in Bytes. But bytes doesn't have
	// io reader needed. So the byte is passed to NewBuffer.
	request, _ := http.NewRequest("POST", "/v1/config/load", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleConfigPOST_only_token(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &LoadConfigBody{
		Token:     "test",
		Filename:  "",
		Subdomain: "",
	}

	b, _ := json.Marshal(body)

	// json Marshal converts struct to json in Bytes. But bytes doesn't have
	// io reader needed. So the byte is passed to NewBuffer.
	request, _ := http.NewRequest("POST", "/v1/config/load", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleConfigPOST_no_body(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &LoadConfigBody{}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/config/load", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleConfigPOST_ConsulError(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsulErr{}
	KeyValues = &FakeKeyValuesErr{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &LoadConfigBody{
		Token:     "test",
		Filename:  "test",
		Subdomain: "test",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/config/load", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleConfigUpload_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("POST", "/v1/config", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleDefaultConfigLoad(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	request, _ := http.NewRequest("GET", "/v1/config/load-default", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleDefaultConfigLoad_err(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValuesErr{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	request, _ := http.NewRequest("GET", "/v1/config/load-default", nil)
	response := httptest.NewRecorder()
	RouterConfig().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}
