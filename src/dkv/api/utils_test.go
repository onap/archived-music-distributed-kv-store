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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadJSON(t *testing.T) {
	oldIoutilRead := IoutilRead

	defer func() {
		IoutilRead = oldIoutilRead
	}()

	IoutilRead = func(path string) ([]byte, error) {
		return []byte("test"), nil
	}

	_, err := ReadJSON("path")
	assert.Equal(t, nil, err, "Error should be nil.")
}
func TestReadJSON_err(t *testing.T) {
	oldIoutilRead := IoutilRead

	defer func() {
		IoutilRead = oldIoutilRead
	}()

	_, err := ReadJSON("path")
	assert.NotNil(t, err, "Err should not be nil.")
}

func TestWriteJSON(t *testing.T) {
	oldReadJson := JsonReader
	oldIoutilWrite := IoutilWrite

	defer func() {
		JsonReader = oldReadJson
		IoutilWrite = oldIoutilWrite
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	err := WriteJSON("path", "token3", "service3")
	assert.Equal(t, nil, err, "Error should be nil.")

}

func TestDeleteInJSON(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	err := DeleteInJSON("path", "token2")
	assert.Equal(t, nil, err, "Error should be nil.")
}

func TestDeleteInJSON_single(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		tsm_list = append(tsm_list, o1)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	err := DeleteInJSON("path", "token1")
	assert.Equal(t, nil, err, "Error should be nil.")
}

func TestDeleteInJSON_not_found(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	err := DeleteInJSON("path", "token3")
	assert.NotNil(t, err, "Err should not be nil.")
}

func TestFindTokenInJSON(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	found, _ := FindTokenInJSON("path", "token2")

	assert.True(t, found, "Token should be found in JSON.")
}

func TestFindServiceInJSON(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	found, _ := FindServiceInJSON("path", "service2")

	assert.True(t, found, "Token should be found in JSON.")
}

func TestFindServiceInJSON_not_found(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	found, _ := FindServiceInJSON("path", "service3")

	assert.False(t, found, "Token should not be found in JSON.")
}

func TestGetServicebyToken(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	service, found, _ := GetServicebyToken("path", "token1")

	assert.Equal(t, "service1", service, "Service not found")
	assert.True(t, found, "Token should be found in JSON.")
}

func TestGetServicebyToken_not_found(t *testing.T) {
	oldReadJson := JsonReader
	defer func() {
		JsonReader = oldReadJson
	}()

	JsonReader = func(path string) ([]Token_service_map, error) {
		var tsm_list []Token_service_map
		o1 := Token_service_map{
			Token:   "token1",
			Service: "service1",
		}
		o2 := Token_service_map{
			Token:   "token2",
			Service: "service2",
		}
		tsm_list = append(tsm_list, o1, o2)
		return tsm_list, nil
	}

	IoutilWrite = func(val string, b []byte, f os.FileMode) error {
		return nil
	}

	service, found, _ := GetServicebyToken("path", "token3")

	assert.Equal(t, "", service, "Service is found")
	assert.False(t, found, "Token should be found in JSON.")
}
