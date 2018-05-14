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
	consulapi "github.com/hashicorp/consul/api"
	"os"
)

type ConsulStruct struct {
	consulClient *consulapi.Client
}

func (c *ConsulStruct) InitializeDatastoreClient() error {
	if os.Getenv("DATASTORE_IP") == "" {
		return errors.New("DATASTORE_IP environment variable not set.")
	}
	config := consulapi.DefaultConfig()
	config.Address = os.Getenv("DATASTORE_IP") + ":8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}
	c.consulClient = client

	return nil
}

func (c *ConsulStruct) CheckDatastoreHealth() error {
	kv := c.consulClient.KV()
	_, _, err := kv.Get("test", nil)
	if err != nil {
		return errors.New("[ERROR] Cannot talk to Datastore. Check if it is running/reachable.")
	}
	return nil
}

func (c *ConsulStruct) RequestPUT(prefix string, key string, value string) error {
	key = prefix + key
	kv := c.consulClient.KV()

	p := &consulapi.KVPair{Key: key, Value: []byte(value)}

	_, err := kv.Put(p, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *ConsulStruct) RequestGET(prefix string, key string) (string, error) {
	key = prefix + key

	kv := c.consulClient.KV()

	pair, _, err := kv.Get(key, nil)

	if pair == nil {
		return string("No value found for key."), err
	}
	return string(pair.Value), err

}

func (c *ConsulStruct) RequestGETS() ([]string, error) {

	kv := c.consulClient.KV()

	pairs, _, err := kv.List("", nil)

	if len(pairs) == 0 {
		return []string{"No keys found."}, err
	}

	var res []string

	for _, keypair := range pairs {
		res = append(res, keypair.Key)
	}

	return res, err
}

func (c *ConsulStruct) RequestDELETE(prefix string, key string) error {
	key = prefix + key
	kv := c.consulClient.KV()

	_, err := kv.Delete(key, nil)

	if err != nil {
		return err
	}

	return nil
}
