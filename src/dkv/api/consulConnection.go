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
	consulapi "github.com/hashicorp/consul/api"
	"os"
)

// Interface to have all signature methods.
type ConsulRequester interface {
	InitializeConsulClient() error
	CheckConsulHealth() error
	RequestPUT(string, string) error
	RequestGET(string) (string, error)
	RequestGETS() ([]string, error)
	RequestDELETE(string) error
}

type ConsulStruct struct {
	consulClient *consulapi.Client
}

/*
This var is an interface used to initialize ConsulStruct when the who API is brought up. This is done this way so
that a fake Consul can be created which satisfies the interface and we can use that fake Consul in unit testing.
*/
var Consul ConsulRequester

/*
The following functions seems like they are not used. But since they are following the ConsulRequest interface,
they can be visible to any Struct which is initiated using the ConsulRequest. This is done for this project in
the initialise.go file where we are creating a ConsulStruct and assigning it to Consul var which is declared
above.
*/
func (c *ConsulStruct) InitializeConsulClient() error {
	config := consulapi.DefaultConfig()
	config.Address = os.Getenv("CONSUL_IP") + ":8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}
	c.consulClient = client

	return nil
}

func (c *ConsulStruct) CheckConsulHealth() error {
	kv := c.consulClient.KV()
	_, _, err := kv.Get("test", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsulStruct) RequestPUT(key string, value string) error {

	kv := c.consulClient.KV()

	p := &consulapi.KVPair{Key: key, Value: []byte(value)}

	_, err := kv.Put(p, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *ConsulStruct) RequestGET(key string) (string, error) {

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

func (c *ConsulStruct) RequestDELETE(key string) error {
	kv := c.consulClient.KV()

	_, err := kv.Delete(key, nil)

	if err != nil {
		return err
	}

	return nil
}
