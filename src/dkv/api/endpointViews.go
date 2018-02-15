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
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type ResponseStringStruct struct {
	Response string `json:"response"`
}

type ResponseGETStruct struct {
	Response map[string]string `json:"response"`
}

type ResponseGETSStruct struct {
	Response []string `json:"response"`
}

type POSTBodyStruct struct {
	Type *TypeStruct `json:"type"`
}

type TypeStruct struct {
	FilePath string `json:"file_path"`
}

func ValidateBody(body POSTBodyStruct) error {
	if body.Type == nil {
		return errors.New("Type not set. Recheck POST data.")
	} else if body.Type.FilePath == "" {
		return errors.New("file_path not set")
	} else {
		return nil
	}
}

func HandlePOST(w http.ResponseWriter, r *http.Request) {

	var body POSTBodyStruct

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		req := ResponseStringStruct{Response: "Empty body."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&req)
		return
	}

	err = ValidateBody(body)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KeyValues.ReadConfigs(body)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KeyValues.WriteKVsToConsul()

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
	} else {
		req := ResponseStringStruct{Response: "Configuration read and default Key Values loaded to Consul"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&req)
	}
}

func HandleGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Consul.RequestGET(key)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
	} else {
		req := ResponseGETStruct{Response: map[string]string{key: value}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}
}

func HandleGETS(w http.ResponseWriter, r *http.Request) {

	values, err := Consul.RequestGETS()

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
	} else {
		req := ResponseGETSStruct{Response: values}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	}
}
