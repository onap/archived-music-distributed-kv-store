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

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/loadconfigs", HandlePOST).Methods("POST")
	router.HandleFunc("/getconfig/{key}", HandleGET).Methods("GET")
	router.HandleFunc("/deleteconfig/{key}", HandleDELETE).Methods("DELETE")
	router.HandleFunc("/getconfigs", HandleGETS).Methods("GET")
	return router
}

func TestHandleGETS(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsul{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("GET", "/getconfigs", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleGETS_err(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsulErr{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("GET", "/getconfigs", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleGET(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsul{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("GET", "/getconfig/key1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleGET_err(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsulErr{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("GET", "/getconfig/key1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandlePOST(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &POSTBodyStruct{
		Type: &TypeStruct{
			FilePath: "default",
		},
	}

	b, _ := json.Marshal(body)

	// json Marshal converts struct to json in Bytes. But bytes doesn't have
	// io reader needed. So the byte is passed to NewBuffer.
	request, _ := http.NewRequest("POST", "/loadconfigs", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandlePOST_no_body(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &POSTBodyStruct{}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/loadconfigs", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandlePOST_no_filepath(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsul{}
	KeyValues = &FakeKeyValues{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &POSTBodyStruct{
		Type: &TypeStruct{},
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/loadconfigs", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandlePOST_ConsulError(t *testing.T) {
	oldConsul := Consul
	oldKeyValues := KeyValues

	Consul = &FakeConsulErr{}
	KeyValues = &FakeKeyValuesErr{}

	defer func() {
		Consul = oldConsul
		KeyValues = oldKeyValues
	}()

	body := &POSTBodyStruct{
		Type: &TypeStruct{
			FilePath: "default",
		},
	}

	b, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/loadconfigs", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code, "500 response is expected")
}

func TestHandleDELETE(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsul{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("DELETE", "/deleteconfig/key1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleDELETE_err(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsulErr{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("DELETE", "/deleteconfig/key1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}
