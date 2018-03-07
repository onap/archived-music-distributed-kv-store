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
	"github.com/magiconair/properties"
	"io/ioutil"
	"log"
	"os"
)

type KeyValuesInterface interface {
	WriteKVsToConsul(string, string, map[string]string) error
	ConfigReader(string, string, string) (map[string]string, error)
	ReadMultiplePropertiesRecursive(string, *map[string]string) error
	ReadMultipleProperties(string, *map[string]string) error
	ReadProperty(string, *map[string]string) error
}

type KeyValuesStruct struct{}

var KeyValues KeyValuesInterface

func (kvStruct *KeyValuesStruct) WriteKVsToConsul(token string, subdomain string, kvs map[string]string) error {
	var prefix = ""
	if subdomain != "" {
		prefix += token + "/" + subdomain
	} else {
		prefix += token + "/"
	}
	for key, value := range kvs {
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

func (kvStruct *KeyValuesStruct) ConfigReader(token string, subdomain string, filename string) (map[string]string, error) {
	var filepath = MOUNTPATH
	kvs := make(map[string]string)

	if filename != "" && subdomain != "" {
		// Specific file in specific domain.
		filepath += token + "/" + subdomain + "/" + filename
		err := kvStruct.ReadProperty(filepath, &kvs)
		if err != nil {
			return kvs, err
		}
		return kvs, nil
	}

	if filename != "" && subdomain == "" {
		// Specific file in Token
		filepath += token + "/" + filename
		err := kvStruct.ReadProperty(filepath, &kvs)
		if err != nil {
			return kvs, err
		}
		return kvs, nil
	}

	if filename == "" && subdomain != "" {
		// All files in specific domain
		filepath += token + "/" + subdomain
		err := kvStruct.ReadMultipleProperties(filepath, &kvs)
		if err != nil {
			return kvs, err
		}
	}

	filepath += token
	err := kvStruct.ReadMultiplePropertiesRecursive(filepath, &kvs)
	if err != nil {
		return kvs, err
	}
	return kvs, nil
}

func (kvStruct *KeyValuesStruct) ReadMultiplePropertiesRecursive(path string, kvs *map[string]string) error {
	// Go inside each sub directory and run ReadMultipleProperties inside.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fi, _ := os.Stat(path + "/" + f.Name())
		if fi.Mode().IsDir() {
			kvStruct.ReadMultipleProperties(path+"/"+f.Name(), kvs)
		} else {
			kvStruct.ReadProperty(path+"/"+f.Name(), kvs)
		}
	}
	return nil
}

func (kvStruct *KeyValuesStruct) ReadMultipleProperties(path string, kvs *map[string]string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		kvStruct.ReadProperty(path+f.Name(), kvs)
	}

	return nil
}

func (kvStruct *KeyValuesStruct) ReadProperty(path string, kvs *map[string]string) error {
	_, err := os.Stat(path)
	if err != nil {
		return errors.New("File does not exists.")
	}
	p := properties.MustLoadFile(path, properties.UTF8)
	for _, key := range p.Keys() {
		(*kvs)[key] = p.MustGet(key)
	}
	return nil
}
