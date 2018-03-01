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
	Domain string      `json:"domain"`
	Type   *TypeStruct `json:"type"`
}

func ValidateLoadConfigBody(body LoadConfigBody) error {
	if body.Domain == "" {
		return errors.New("Domain not set. Please set domain in POST.")
	}
	if body.Type == nil {
		return errors.New("Type not set. Recheck POST data.")
	} else if body.Type.FilePath == "" {
		return errors.New("file_path not set")
	} else {
		return nil
	}
}

func ValidateUploadConfigBody(body UploadConfigBody) error {
	return nil
}

func HandleConfigUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2000) // 2k bytes?
	file, handler, err := r.FormFile("configFile")
	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}
	defer file.Close()

	print(r.Body)

	if err != nil {
		println("Helere")
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	token := r.Form.Get("token")
	subdomain := r.Form.Get("subdomain")

	var filename = ""
	if subdomain != "" {
		filename += token + "/" + subdomain + "/" + handler.Filename
	} else {
		filename += token + "/" + handler.Filename
	}

	println(token)

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

	err = KeyValues.ReadConfigs(body)

	if err != nil {
		req := ResponseStringStruct{Response: string(err.Error())}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(req)
		return
	}

	err = KeyValues.WriteKVsToConsul(body.Domain)

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
	// TODO
}

func HandleConfigDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
}
