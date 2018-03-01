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
	if body.Domain == "default" {
		return errors.New("Domain not allowed. Please set another domain in POST.")
	}
	return nil
}

/*
	TODO(sshank): Add validations to check if tokens/sub-domains/files indeed
	exist in the token_service JSON or in the directory. This is to avoid the service
	returning the file system errors to the user.
*/

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

	token, err := Directory.CreateService(body)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Registration Successful. Token: "+token)
	}
}

func HandleServiceGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not present in path.")
		return
	}

	service, found, err := Directory.FindService(token)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		if found == true {
			GenerateResponse(w, r, http.StatusOK, service)
		} else {
			GenerateResponse(w, r, http.StatusNotFound, "Service for Token:"+token+"not found.")
		}

	}
}

func HandleServiceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	if token == "default" {
		GenerateResponse(w, r, http.StatusNotAcceptable, "Default delete not allowed.")
		return
	}

	err := Directory.RemoveService(token)

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

	err = Directory.CreateServiceSubdomain(token, body.Subdomain)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Subdomain creation success with token: "+token)
	}

}

func HandleServiceSubdomainGet(w http.ResponseWriter, r *http.Request) {
	// TODO(sshank): This will list all subdomain in a service.
}

func HandleServiceSubdomainDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	subdomain := vars["subdomain"]

	if token == "" {
		GenerateResponse(w, r, http.StatusBadRequest, "Token not passed.")
		return
	}

	if token == "default" && subdomain == "" {
		GenerateResponse(w, r, http.StatusNotAcceptable, "Not allowerd.")
		return
	}

	err := Directory.RemoveServiceSubdomain(token, subdomain)

	if err != nil {
		GenerateResponse(w, r, http.StatusInternalServerError, string(err.Error()))
	} else {
		GenerateResponse(w, r, http.StatusOK, "Deletion of service is successful.")
	}
}
