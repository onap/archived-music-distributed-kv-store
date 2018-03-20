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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func RouterConsul() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/getconfig/{key}", HandleGET).Methods("GET")
	router.HandleFunc("/v1/deleteconfig/{key}", HandleDELETE).Methods("DELETE")
	router.HandleFunc("/v1/getconfigs", HandleGETS).Methods("GET")
	return router
}

func TestHandleGETS(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsul{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("GET", "/v1/getconfigs", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleGETS_err(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsulErr{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("GET", "/v1/getconfigs", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleGET(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsul{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("GET", "/v1/getconfig/key1", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleGET_err(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsulErr{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("GET", "/v1/getconfig/key1", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}

func TestHandleDELETE(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsul{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("DELETE", "/v1/deleteconfig/key1", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestHandleDELETE_err(t *testing.T) {
	oldDataStore := Datastore
	Datastore = &FakeConsulErr{}
	defer func() { Datastore = oldDataStore }()

	request, _ := http.NewRequest("DELETE", "/v1/deleteconfig/key1", nil)
	response := httptest.NewRecorder()
	RouterConsul().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "400 response is expected")
}
