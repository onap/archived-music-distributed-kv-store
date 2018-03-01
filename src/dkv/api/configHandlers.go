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

func HandleConfigUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000) // 2k bytes?
	file, handler, err := r.FormFile("configFile")
	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, "Error in uploaded file.")
		return
	}
	defer file.Close()

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
		return
	}

	token := r.Form.Get("token")
	subdomain := r.Form.Get("subdomain")

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not present in Form data.")
		return
	}

	var filename = ""
	if subdomain != "" {
		filename += token + "/" + subdomain + "/" + handler.Filename
	} else {
		filename += token + "/" + handler.Filename
	}

	f, err := Directory.CreateFile(MOUNTPATH + filename)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
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
		GenerateResponse(w, r, http.StatusBadRequest, "Empty body.")
		return
	}

	err = ValidateLoadConfigBody(body)

	if err != nil {
		GenerateResponse(w, r, http.StatusBadRequest, string(err.Error()))
		return
	}

	err = KeyValues.ConfigReader(body.Token, body.Subdomain, body.Filename)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
		return
	}

	err = KeyValues.WriteKVsToConsul(body.Token, body.Subdomain)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Configuration read and Key Values loaded to Consul.")
	}
}

func HandleDefaultConfigLoad(w http.ResponseWriter, r *http.Request) {
	err := KeyValues.ConfigReader("default", "", "")
	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
		return
	}
	err = KeyValues.WriteKVsToConsul("default", "")
	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Default Configuration read and default Key Values loaded to Consul.")
	}
}

func HandleConfigGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	filename := vars["filename"]
	subdomain := vars["subdomain"]

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not passed.")
		return
	}

	if filename == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "filename not passed.")
		return
	}

	Directory.FetchFile(w, r, token, subdomain, filename)
}

func HandleConfigDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	filename := vars["filename"]
	subdomain := vars["subdomain"]

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not passed.")
		return
	}

	if filename == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "filename not passed.")
		return
	}

	err := Directory.RemoveFile(token, subdomain, filename)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Deletion of config is successful.")
	}
}
