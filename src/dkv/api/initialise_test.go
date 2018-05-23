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
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitialise_consul(t *testing.T) {
	oldDatastore_ip := os.Getenv("DATASTORE_IP")
	oldDatastore_type := os.Getenv("DATASTORE")
	oldJsonExists := JsonChecker
	oldJsonCreate := JsonCreate

	os.Setenv("DATASTORE_IP", "localhost")
	os.Setenv("DATASTORE", "consul")

	defer func() {
		os.Setenv("DATASTORE_IP", oldDatastore_ip)
		os.Setenv("DATASTORE", oldDatastore_type)
		JsonCreate = oldJsonCreate
		JsonChecker = oldJsonExists
	}()

	JsonChecker = func(path string) (bool, error) {
		return false, nil
	}

	JsonCreate = func(path string) error {
		return nil
	}

	err := Initialise()
	assert.NotNil(t, err)
}

func TestInitialise_cassandra(t *testing.T) {
	oldDatastore_ip := os.Getenv("DATASTORE_IP")
	oldDatastore_type := os.Getenv("DATASTORE")
	oldMOUNTPATH := os.Getenv("MOUNTPATH")
	oldJsonChecker := JsonChecker

	os.Setenv("DATASTORE_IP", "localhost")
	os.Setenv("DATASTORE", "cassandra")

	defer func() {
		os.Setenv("DATASTORE_IP", oldDatastore_ip)
		os.Setenv("DATASTORE", oldDatastore_type)
		os.Setenv("MOUNTPATH", oldMOUNTPATH)
		JsonChecker = oldJsonChecker
	}()

	JsonChecker = func(path string) (bool, error) {
		return true, nil
	}

	err := Initialise()
	assert.Nil(t, err)
}

func TestInitialise_datastoreUnknown(t *testing.T) {
	datastore := os.Getenv("DATASTORE")
	defer os.Setenv("DATASTORE", datastore)
	os.Setenv("DATASTORE", "test")

	err := Initialise()
	assert.NotNil(t, err)
}

func TestInitialise_datastoreEmpty(t *testing.T) {
	datastore := os.Getenv("DATASTORE")
	defer os.Setenv("DATASTORE", datastore)
	os.Setenv("DATASTORE", "")

	err := Initialise()
	assert.NotNil(t, err)
}

func TestInitialise_noJSON(t *testing.T) {
	oldDatastore_ip := os.Getenv("DATASTORE_IP")
	oldDatastore_type := os.Getenv("DATASTORE")
	oldJsonChecker := JsonChecker

	os.Setenv("DATASTORE_IP", "localhost")
	os.Setenv("DATASTORE", "consul")

	defer func() {
		os.Setenv("DATASTORE_IP", oldDatastore_ip)
		os.Setenv("DATASTORE", oldDatastore_type)
		JsonChecker = oldJsonChecker
	}()

	JsonChecker = func(path string) (bool, error) {
		return false, nil
	}

	err := Initialise()
	assert.NotNil(t, err)
}
