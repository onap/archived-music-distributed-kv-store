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

type CreateRegisterServiceBody struct {
	Domain string `json:"domain"`
}

type CreateServiceSubdomainBody struct {
	Subdomain string `json:"subdomain"`
}

func ValidateCreateRegisterServiceBody(body CreateRegisterServiceBody) error {
	if body.Domain == "" {
		return errors.New("Domain not set. Please set domain in POST.")
	}
	return nil
}

func HandleServiceCreate(w http.ResponseWriter, r *http.Request) {
	var body CreateRegisterServiceBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		GenerateResponse(w, r, http.StatusBadRequest, "Empty body.")
		return
	}

	err = ValidateCreateRegisterServiceBody(body)

	if err != nil {
		GenerateResponse(w, r, http.StatusBadRequest, string(err.Error()))
		return
	}

	token, err := DirectoryOperation.CreateService(body)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Registration Successful. Token: "+token)
	}
}

func HandleServiceGet(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func HandleServiceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	err := DirectoryOperation.RemoveService(token)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Deletion of service is successful.")
	}
}

func HandleServiceSubdomainCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	var body CreateServiceSubdomainBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		GenerateResponse(w, r, http.StatusBadRequest, "Empty body.")
		return
	}

	if body.Subdomain == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Subdomain not found in POST.")
		return
	}

	err = DirectoryOperation.CreateServiceSubdomain(token, body.Subdomain)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Subdomain creation success with token: "+token)
	}

}

func HandleServiceSubdomainGet(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func HandleServiceSubdomainDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	subdomain := vars["subdomain"]

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not passed.")
		return
	}

	err := DirectoryOperation.RemoveServiceSubdomain(token, subdomain)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Deletion of service is successful.")
	}
}
