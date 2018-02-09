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
	"github.com/gorilla/mux"
	"net/http"
)

var getkvs = GetKVsFromConsul

func HandlePOST(w http.ResponseWriter, r *http.Request) {

	var body LoadStruct

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

	err = KVStruct.ReadConfigs(body)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KVStruct.WriteKVsToConsul()

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

	value, err := GetKVFromConsul(key)

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

	values, err := getkvs()

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
