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
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type UploadConfigBody struct {
	Token     string
	File      multipart.File
	Subdomain string
}

type LoadConfigBody struct {
	Token     string `json:"token"`
	Filename  string `json:"filename"`
	Subdomain string `json:"subdomain"`
}

func ValidateLoadConfigBody(body LoadConfigBody) error {
	if body.Token == "" {
		return errors.New("Token not set. Please set Token in POST.")
	}
	return nil
}

func ValidateUploadConfigBody(body UploadConfigBody) error {
	return nil
}

func HandleConfigUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000) // 2k bytes?
	file, handler, err := r.FormFile("configFile")
	if err != nil {
		req := ResponseStringStruct{Response: "Error in uploaded file."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}
	defer file.Close()

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	token := r.Form.Get("token")
	subdomain := r.Form.Get("subdomain")

	if token == "" {
		req := ResponseStringStruct{Response: "Token not present in Form data."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	var filename = ""
	if subdomain != "" {
		filename += token + "/" + subdomain + "/" + handler.Filename
	} else {
		filename += token + "/" + handler.Filename
	}

	f, err := os.OpenFile(MOUNTPATH+filename, os.O_CREATE|os.O_WRONLY, 0770)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func HandleConfigLoad(w http.ResponseWriter, r *http.Request) {

	var body LoadConfigBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		req := ResponseStringStruct{Response: "Empty body."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&req)
		return
	}

	err = ValidateLoadConfigBody(body)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KeyValues.ConfigReader(body.Token, body.Subdomain, body.Filename)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KeyValues.WriteKVsToConsul(body.Token, body.Subdomain)

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

func HandleConfigGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	filename := vars["filename"]
	subdomain := vars["subdomain"]

	if token == "" {
		req := ResponseStringStruct{Response: "Token not passed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	if filename == "" {
		req := ResponseStringStruct{Response: "filename not passed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	DirectoryOperation.FetchFile(w, r, token, subdomain, filename)
}

func HandleConfigDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	filename := vars["filename"]
	subdomain := vars["subdomain"]

	if token == "" {
		req := ResponseStringStruct{Response: "Token not passed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	if filename == "" {
		req := ResponseStringStruct{Response: "filename not passed."}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(req)
		return
	}

	err := DirectoryOperation.RemoveFile(token, subdomain, filename)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
	} else {
		// TODO(sshank): Remove all keys from Consul before returning Request?
		req := ResponseStringStruct{Response: "Deletion of config is successful."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&req)
	}
}
