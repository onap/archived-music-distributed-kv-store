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

func TestInitialise_cassandra(t *testing.T) {
	oldDatastore_ip := os.Getenv("DATASTORE_IP")
	oldDatastore_type := os.Getenv("DATASTORE")

	os.Setenv("DATASTORE_IP", "localhost")
	os.Setenv("DATASTORE", "cassandra")

	defer func() {
		os.Setenv("DATASTORE_IP", oldDatastore_ip)
		os.Setenv("DATASTORE", oldDatastore_type)
	}()

	err := Initialise()
	assert.Nil(t, err)
}
func TestInitialise_consulError(t *testing.T) {
	oldDatastore_ip := os.Getenv("DATASTORE_IP")
	oldDatastore_type := os.Getenv("DATASTORE")

	os.Setenv("DATASTORE_IP", "localhost")
	os.Setenv("DATASTORE", "consul")

	defer func() {
		os.Setenv("DATASTORE_IP", oldDatastore_ip)
		os.Setenv("DATASTORE", oldDatastore_type)
	}()

	err := Initialise()
	assert.NotNil(t, err)
}

func TestInitialise_datastoreEmptyError(t *testing.T) {
	datastore := os.Getenv("DATASTORE")
	os.Unsetenv("DATASTORE")
	defer os.Setenv("DATASTORE", datastore)

	err := Initialise()
	assert.NotNil(t, err)
}
