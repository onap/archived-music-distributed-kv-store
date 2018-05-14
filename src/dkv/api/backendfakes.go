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
	"net/http"
)

/*
A ConsulStruct is added inside this so that FakeConsul becomes an implementation of the Consul interface.
If we don't add ConsulStruct inside this, it complains that the FakeConsul Struct doesn't implement all the methods
defined in Consul interface.
*/
// Correct
type FakeConsul struct {
	ConsulStruct
}

func (f *FakeConsul) InitializeDatastoreClient() error {
	return nil
}

func (f *FakeConsul) CheckDatastoreHealth() error {
	return nil
}

func (f *FakeConsul) RequestGETS() ([]string, error) {
	return []string{"key1", "key2"}, nil
}

func (f *FakeConsul) RequestGET(key string, token string) (string, error) {
	return key, nil
}

func (f *FakeConsul) RequestPUT(key string, value string, token string) error {
	return nil
}

func (f *FakeConsul) RequestDELETE(key string, token string) error {
	return nil
}

// Error
type FakeConsulErr struct {
	ConsulStruct
}

func (f *FakeConsulErr) InitializeDatastoreClient() error {
	return errors.New("Internal Server Error")
}

func (f *FakeConsulErr) CheckDatastoreHealth() error {
	return errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestGETS() ([]string, error) {
	return []string{"", ""}, errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestGET(key string, token string) (string, error) {
	return "", errors.New("Internal Server Error")
}

func (f *FakeConsulErr) RequestDELETE(key string, token string) error {
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

func (f *FakeKeyValues) ConfigReader(token string, subdomain string, filename string) (map[string]string, error) {
	kvs := make(map[string]string)
	return kvs, nil
}

func (f *FakeKeyValues) WriteKVsToDatastore(token string, subdomain string, kvs map[string]string) error {
	return nil
}

// Error
type FakeKeyValuesErr struct {
	KeyValuesStruct
}

func (f *FakeKeyValuesErr) ConfigReader(token string, subdomain string, filename string) (map[string]string, error) {
	kvs := make(map[string]string)
	return kvs, errors.New("Internal Server Error")
}

func (f *FakeKeyValuesErr) WriteKVsToDatastore(token string, subdomain string, kvs map[string]string) error {
	return errors.New("Internal Server Error")
}

// Correct
type FakeDirectory struct {
	DirectoryStruct
}

func (f *FakeDirectory) CreateService(CreateRegisterServiceBody) (string, error) {
	return "", nil
}

func (f *FakeDirectory) RemoveService(token string) error {
	return nil
}

func (f *FakeDirectory) CreateServiceSubdomain(token string, subdomain string) error {
	return nil
}

func (f *FakeDirectory) RemoveServiceSubdomain(token string, subdomain string) error {
	return nil
}

func (f *FakeDirectory) FindService(token string) (string, bool, error) {
	return "service1", true, nil
}

func (f *FakeDirectory) FetchFile(
	w http.ResponseWriter, r *http.Request, token string, subdomain string, filename string) {
}

func (f *FakeDirectory) RemoveFile(token string, subdomain string, filename string) error {
	return nil
}

// Error
type FakeDirectoryErr struct {
	DirectoryStruct
}

func (f *FakeDirectoryErr) CreateService(CreateRegisterServiceBody) (string, error) {
	return "", errors.New("Internal Server Error.")
}

func (f *FakeDirectoryErr) RemoveService(token string) error {
	return errors.New("Internal Server Error.")
}

func (f *FakeDirectoryErr) CreateServiceSubdomain(token string, subdomain string) error {
	return errors.New("Internal Server Error.")
}

func (f *FakeDirectoryErr) RemoveServiceSubdomain(token string, subdomain string) error {
	return errors.New("Internal Server Error.")
}

func (f *FakeDirectoryErr) FindService(token string) (string, bool, error) {
	return "", false, errors.New("Internal Server Error.")
}

func (f *FakeDirectoryErr) FetchFile(
	w http.ResponseWriter, r *http.Request, token string, subdomain string, filename string) {

}

func (f *FakeDirectoryErr) RemoveFile(token string, subdomain string, filename string) error {
	return errors.New("Internal Server Error.")
}
