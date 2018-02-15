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
	router.HandleFunc("/getconfigs", HandleGETS).Methods("GET")
	return router
}

/*
A ConsulStruct is added inside this so that FakeConsul becomes an implementation of the Consul interface.
If we don't add ConsulStruct inside this, it complains that the FakeConsul Struct doesn't implement all the methods
defined in Consul interface.
*/
type FakeConsul struct {
	ConsulStruct
}

func (f *FakeConsul) requestGETS() ([]string, error) {
	return []string{"key1", "key2"}, nil
}

func (f *FakeConsul) requestGET(key string) (string, error) {
	return key, nil
}

func (f *FakeConsul) requestPUT(key string, value string) error {
	return nil
}

/*
This is done similar to the fake Consul above to pass FakeKeyValues to the interface and control method's outputs
as required.
*/
type FakeKeyValues struct {
	KeyValuesStruct
}

func (f *FakeKeyValues) ReadConfigs(body POSTBodyStruct) error {
	return nil
}

func (f *FakeKeyValues) WriteKVsToConsul() error {
	return nil
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

func TestHandleGET(t *testing.T) {
	oldConsul := Consul
	Consul = &FakeConsul{}
	defer func() { Consul = oldConsul }()

	request, _ := http.NewRequest("GET", "/getconfig/key1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "200 response is expected")
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
