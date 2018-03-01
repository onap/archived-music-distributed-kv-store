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
	"io/ioutil"
)

type Token_service_map struct {
	Token   string `json:"token"`
	Service string `json:"service"`
}

func ReadJSON(path string) ([]Token_service_map, error) {
	var tsm_list []Token_service_map
	// raw, err := ioutil.ReadFile("./token_service_map.json")
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return tsm_list, err
	}
	json.Unmarshal(raw, &tsm_list)
	return tsm_list, nil
}

func WriteJSON(path string, token string, service string) error {
	tsm_list, err := ReadJSON(path)
	if err != nil {
		return err
	}
	var tsm Token_service_map
	tsm.Token = token
	tsm.Service = service
	tsm_list = append(tsm_list, tsm)
	raw, err := json.Marshal(tsm_list)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, raw, 0644)
	if err != nil {
		return err
	}
	return nil
}

func DeleteInJSON(path string, token string) error {
	serviceList, err := ReadJSON(path)
	if err != nil {
		return err
	}

	var resultList []Token_service_map
	var foundFlag = false

	// Linear search for the token. If found set found flag. If not, keep
	// copying the different values to resultList.
	for _, service := range serviceList {
		if service.Token == token {
			foundFlag = true
		} else {
			resultList = append(resultList, service)
		}
	}

	if foundFlag == false {
		return errors.New("Service not found. Check if Token is correct or service is registered.")
	} else {
		// This is done to avoid writing 'null' in the json file.
		if len(serviceList) == 1 {
			dummy := Token_service_map{Token: "", Service: ""}
			resultList = append(resultList, dummy)
		}
		raw, err := json.Marshal(resultList)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, raw, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func FindTokenInJSON(path string, token string) (bool, error) {
	serviceList, err := ReadJSON(path)
	if err != nil {
		return false, err
	}

	for _, service := range serviceList {
		if service.Token == token {
			return true, nil
		}
	}
	return false, nil
}

func FindServiceInJSON(path string, serviceName string) (bool, error) {
	serviceList, err := ReadJSON(path)
	if err != nil {
		return false, err
	}

	for _, service := range serviceList {
		if service.Service == serviceName {
			return true, nil
		}
	}
	return false, nil
}

// func GenerateResponse(){
// 	req := ResponseStringStruct{Response: string(err.Error())}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(req)
// 		return
// }
