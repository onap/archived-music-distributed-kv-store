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

import "errors"

/*
A ConsulStruct is added inside this so that FakeConsul becomes an implementation of the Consul interface.
If we don't add ConsulStruct inside this, it complains that the FakeConsul Struct doesn't implement all the methods
defined in Consul interface.
*/
// Correct
type FakeConsul struct {
	ConsulStruct
}

func (f *FakeConsul) RequestGETS() ([]string, error) {
	return []string{"key1", "key2"}, nil
}

func (f *FakeConsul) RequestGET(key string) (string, error) {
	return key, nil
}

func (f *FakeConsul) RequestPUT(key string, value string) error {
	return nil
}

func (f *FakeConsul) RequestDELETE(key string) error {
	return nil
}

// Error
type FakeConsulErr struct {
	ConsulStruct
}

func (f *FakeConsulErr) RequestGETS() ([]string, error) {
	return []string{"", ""}, errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestGET(key string) (string, error) {
	return "", errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestPUT(key string, value string) error {
	return errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestDELETE(key string) error {
	return errors.New("Internal Server Error")
}

/*
This is done similar to the fake Consul above to pass FakeKeyValues to the interface and control method's outputs
as required.
*/
//Correct
type FakeKeyValues struct {
	KeyValuesStruct
}

func (f *FakeKeyValues) ReadConfigs(body POSTBodyStruct) error {
	return nil
}

func (f *FakeKeyValues) WriteKVsToConsul() error {
	return nil
}

// Error
type FakeKeyValuesErr struct {
	KeyValuesStruct
}

func (f *FakeKeyValuesErr) ReadConfigs(body POSTBodyStruct) error {
	return errors.New("Internal Server Error")
}

func (f *FakeKeyValuesErr) WriteKVsToConsul() error {
	return errors.New("Internal Server Error")
}
