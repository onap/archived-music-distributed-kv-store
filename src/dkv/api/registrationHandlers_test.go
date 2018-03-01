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

func RouterRegister() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/register", HandleServiceCreate).Methods("POST")
	router.HandleFunc("/v1/register/{token}", HandleServiceGet).Methods("GET")
	router.HandleFunc("/v1/register/{token}", HandleServiceDelete).Methods("DELETE")
	return router
}

func RouterRegisterSubdomain() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/register/{token}/subdomain", HandleServiceSubdomainCreate).Methods("POST")
	router.HandleFunc("/v1/register/{token}/subdomain/{subdomain}", HandleServiceSubdomainDelete).Methods("DELETE")
	return router
}

func TestHandleServiceCreate(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateRegisterServiceBody{
		Domain: "test",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleServiceCreate_default(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateRegisterServiceBody{
		Domain: "default",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}
func TestHandleServiceCreate_no_body(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateRegisterServiceBody{}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleServiceCreate_no_domain(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateRegisterServiceBody{
		Domain: "",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleServiceGet(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}
	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("GET", "/v1/register/token1", nil)
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleServiceGet_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectoryErr{}
	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("GET", "/v1/register/token1", nil)
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}
func TestHandleServiceDelete(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/register/token1", nil)
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleServiceDelete_default(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/register/default", nil)
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 406, response.Code, "406 response is expected")
}

func TestHandleServiceDelete_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectoryErr{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/register/token1", nil)
	response := httptest.NewRecorder()
	RouterRegister().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleServiceSuddomainCreate(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateServiceSubdomainBody{
		Subdomain: "test",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register/token1/subdomain", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleServiceSuddomainCreate_no_body(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateServiceSubdomainBody{}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register/token1/subdomain", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleServiceSuddomainCreate_no_subdomain(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	body := &CreateServiceSubdomainBody{
		Subdomain: "",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register/token1/subdomain", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleServiceSuddomainCreate_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectoryErr{}

	defer func() { Directory = oldDirectory }()

	body := &CreateServiceSubdomainBody{
		Subdomain: "test",
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/v1/register/token1/subdomain", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleServiceSuddomainDelete(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectory{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/register/token1/subdomain/subdomain1", nil)
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleServiceSuddomainDelete_err(t *testing.T) {
	oldDirectory := Directory
	Directory = &FakeDirectoryErr{}

	defer func() { Directory = oldDirectory }()

	request, _ := http.NewRequest("DELETE", "/v1/register/token1/subdomain/subdomain1", nil)
	response := httptest.NewRecorder()
	RouterRegisterSubdomain().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}
