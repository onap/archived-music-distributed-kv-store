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
	"path"
	"runtime"
	"sync"
)

type KeyValuesInterface interface {
	WriteKVsToConsul(string) error
	ReadConfigs(POSTBodyStruct) error
	PropertiesFilesToKV(string) error
	ReadMultipleProperties(string) error
	ReadProperty(string)
}

type KeyValuesStruct struct {
	sync.RWMutex
	kvs map[string]string
}

var KeyValues KeyValuesInterface

func (kvStruct *KeyValuesStruct) WriteKVsToConsul(prefix string) error {
	for key, value := range kvStruct.kvs {
		key = prefix + "." + key
		err := Consul.RequestPUT(key, value)
		if err != nil {
			return err
		}
		log.Println("[INFO] Key: ", key, "| Value: ", value)
	}
	log.Println("[INFO] Wrote KVs to Consul.")
	return nil
}

func (kvStruct *KeyValuesStruct) ReadConfigs(body POSTBodyStruct) error {
	defer kvStruct.Unlock()

	kvStruct.Lock()

	err := kvStruct.PropertiesFilesToKV(body.Type.FilePath)
	if err != nil {
		return err
	}
	return nil
}

func (kvStruct *KeyValuesStruct) PropertiesFilesToKV(directory string) error {

	if directory == "default" {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			return errors.New("No caller")
		}

		defaultDir := path.Dir(filename) + "/../configurations/"
		err := kvStruct.ReadMultipleProperties(defaultDir)
		if err != nil {
			return err
		}

		return nil

	} else {
		directory += "/"
		err := kvStruct.ReadMultipleProperties(directory)
		if err != nil {
			return err
		}

		return nil
	}
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
