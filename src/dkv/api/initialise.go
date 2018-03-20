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
	"os"
)

var (
	Datastore DatastoreConnector
	KeyValues KeyValuesInterface
	Directory DirectoryOperationer
)

func Initialise() error {
	if os.Getenv("DATASTORE") == "" {
		return errors.New("DATASTORE environment variable not set.")
	}
	if os.Getenv("DATASTORE") == "consul" {
		Datastore = &ConsulStruct{}
	} else if os.Getenv("DATASTORE") == "cassandra" {
		Datastore = &CassandraStruct{}
	}
	KeyValues = &KeyValuesStruct{}
	Directory = &DirectoryStruct{directory: ""}

	err := Datastore.InitializeDatastoreClient()
	if err != nil {
		return err
	}

	err = Datastore.CheckDatastoreHealth()
	if err != nil {
		return err
	}

	if os.Getenv("MOUNTPATH") != "" {
		MOUNTPATH = os.Getenv("MOUNTPATH")
	} else {
		MOUNTPATH = "../../mountpath/"
	}

	return nil
}
