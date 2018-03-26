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
	"errors"
	uuid "github.com/hashicorp/go-uuid"
	"net/http"
	"os"
)

type DirectoryOperationer interface {
	// Service Operations.
	CreateService(CreateRegisterServiceBody) (string, error)
	RemoveService(string) error
	CreateServiceSubdomain(string, string) error
	RemoveServiceSubdomain(string, string) error
	// Directory Operations.
	CreateDirectory(string) error
	RemoveDirectory(string) error
	RemoveSubDirectory(string, string) error
	RemoveFile(string, string, string) error
	FindService(string) (string, bool, error)
	FetchFile(http.ResponseWriter, *http.Request, string, string, string)
	CreateFile(string) (*os.File, error)
}

type DirectoryStruct struct {
	directory string
}

/*
TODO(sshank): This needs to be made completely stateless or else each container running this API will not have
same data in token_service_map.json
*/
const (
	JSONPATH = "api/token_service_map.json"
)

var MOUNTPATH = ""

func (d *DirectoryStruct) CreateService(body CreateRegisterServiceBody) (string, error) {

	// Having same name is prohibited?
	found, err := FindServiceInJSON(JSONPATH, body.Domain)
	if err != nil {
		return "", err
	}
	if found {
		return "", errors.New("Service already found. Check name.")
	}

	token, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	err = d.CreateDirectory(token)
	if err != nil {
		return "", err
	}

	err = WriteJSON(JSONPATH, token, body.Domain)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (d *DirectoryStruct) CreateServiceSubdomain(token string, subdomain string) error {
	foundToken, err := FindTokenInJSON(JSONPATH, token)
	if err != nil {
		return err
	}
	if foundToken == false {
		return errors.New("Token not found. Please check token or if service is created.")
	}
	err = d.CreateSubDirectory(token, subdomain)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) RemoveService(token string) error {
	err := DeleteInJSON(JSONPATH, token)
	if err != nil {
		return err
	}
	err = d.RemoveDirectory(token)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) FindService(token string) (string, bool, error) {
	service, found, err := GetServicebyToken(JSONPATH, token)
	if err != nil {
		return "", false, err
	}
	return service, found, nil
}

func (d *DirectoryStruct) RemoveServiceSubdomain(token string, subdomain string) error {
	foundToken, err := FindTokenInJSON(JSONPATH, token)
	if err != nil {
		return err
	}
	if foundToken == false {
		return errors.New("Token not found. Please check token or if service is created.")
	}
	err = d.RemoveSubDirectory(token, subdomain)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) CreateDirectory(token string) error {
	// Permissions inside mount point?
	err := os.Mkdir(MOUNTPATH+token, os.FileMode(0770))
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) CreateSubDirectory(token string, subdomain string) error {
	err := os.Mkdir(MOUNTPATH+token+"/"+subdomain, os.FileMode(0770))
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) RemoveDirectory(token string) error {
	err := os.RemoveAll(MOUNTPATH + token)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) RemoveSubDirectory(token string, subdomain string) error {
	err := os.RemoveAll(MOUNTPATH + token + "/" + subdomain)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) RemoveFile(token string, subdomain string, filename string) error {
	var filepath = ""
	if subdomain != "" {
		filepath += MOUNTPATH + token + "/" + subdomain + "/" + filename
	} else {
		filepath += MOUNTPATH + token + "/" + filename
	}
	// If error, it seems to show the mounthpath back to the client. This is not good
	// error return practise. It shoudn't return the exact file path on the system.
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}

func (d *DirectoryStruct) FetchFile(
	w http.ResponseWriter, r *http.Request, token string, subdomain string, filename string) {

	var filepath = ""
	if subdomain != "" {
		filepath += MOUNTPATH + token + "/" + subdomain + "/" + filename
	} else {
		filepath += MOUNTPATH + token + "/" + filename
	}

	http.ServeFile(w, r, filepath)
}

func (d *DirectoryStruct) CreateFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		return nil, err
	}
	return f, nil
}
