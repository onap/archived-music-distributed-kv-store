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
	// "errors"
	"github.com/magiconair/properties"
	"io/ioutil"
	"log"
	// "path"
	// "runtime"
	"os"
	"sync"
)

type KeyValuesInterface interface {
	WriteKVsToConsul(string, string) error
	ConfigReader(string, string, string) error
	ReadMultiplePropertiesRecursive(string) error
	ReadMultipleProperties(string) error
	ReadProperty(string)
}

type KeyValuesStruct struct {
	sync.RWMutex
	kvs map[string]string
}

var KeyValues KeyValuesInterface

func (kvStruct *KeyValuesStruct) WriteKVsToConsul(token string, subdomain string) error {
	var prefix = ""
	if subdomain != "" {
		prefix += token + "/" + subdomain
	} else {
		prefix += token + "/"
	}
	for key, value := range kvStruct.kvs {
		key = prefix + key
		err := Consul.RequestPUT(key, value)
		if err != nil {
			return err
		}
		log.Println("[INFO] Key: ", key, "| Value: ", value)
	}
	log.Println("[INFO] Wrote KVs to Consul.")
	return nil
}

func (kvStruct *KeyValuesStruct) ConfigReader(token string, subdomain string, filename string) error {
	defer kvStruct.Unlock()

	kvStruct.Lock()
	var filepath = MOUNTPATH

	if filename != "" && subdomain != "" {
		// Specific file in specific domain.
		filepath += token + "/" + subdomain + "/" + filename
		kvStruct.ReadProperty(filepath)
		return nil
	}

	if filename != "" && subdomain == "" {
		// Specific file in Token
		println(filepath)
		filepath += token + "/" + filename
		kvStruct.ReadProperty(filepath)
		return nil
	}

	if filename == "" && subdomain != "" {
		// All files in specific domain
		filepath += token + "/" + subdomain
		kvStruct.ReadMultipleProperties(filepath)
		return nil
	}

	filepath += token
	kvStruct.ReadMultiplePropertiesRecursive(filepath)
	return nil
}

func (kvStruct *KeyValuesStruct) ReadMultiplePropertiesRecursive(path string) error {
	// Go inside each sub directory and run ReadMultipleProperties inside.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fi, _ := os.Stat(path + "/" + f.Name())
		if fi.Mode().IsDir() {
			kvStruct.ReadMultipleProperties(path + "/" + f.Name())
		} else {
			kvStruct.ReadProperty(path + "/" + f.Name())
		}
	}
	return nil
}

func (kvStruct *KeyValuesStruct) ReadMultipleProperties(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		kvStruct.ReadProperty(path + f.Name())
	}

	return nil
}

func (kvStruct *KeyValuesStruct) ReadProperty(path string) {
	p := properties.MustLoadFile(path, properties.UTF8)
	for _, key := range p.Keys() {
		kvStruct.kvs[key] = p.MustGet(key)
	}
}
