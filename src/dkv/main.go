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

package main

import (
	"dkv/api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	err := api.Initialise()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	// Sevice Registration
	// Domain CRUD
	router.HandleFunc("/v1/register", api.HandleServiceCreate).Methods("POST")
	router.HandleFunc("/v1/register/{token}", api.HandleServiceGet).Methods("GET")
	router.HandleFunc("/v1/register/{token}", api.HandleServiceDelete).Methods("DELETE")
	// Subdomain CRUD
	router.HandleFunc("/v1/register/{token}/subdomain", api.HandleServiceSubdomainCreate).Methods("POST")
	router.HandleFunc("/v1/register/{token}/subdomain", api.HandleServiceSubdomainGet).Methods("GET")
	router.HandleFunc("/v1/register/{token}/subdomain/{subdomain}", api.HandleServiceSubdomainDelete).Methods("DELETE")
	// Configuration CRUD
	router.HandleFunc("/v1/config", api.HandleConfigUpload).Methods("POST")
	router.HandleFunc("/v1/config", api.HandleConfigGet).Methods("GET")
	router.HandleFunc("/v1/config", api.HandleConfigDelete).Methods("DELETE")
	router.HandleFunc("/v1/loadconfig", api.HandleConfigLoad).Methods("POST")
	// Direct Consul queries.
	router.HandleFunc("/v1/getconfig/{key}", api.HandleGET).Methods("GET")
	router.HandleFunc("/v1/deleteconfig/{key}", api.HandleDELETE).Methods("DELETE")
	router.HandleFunc("/v1/getconfigs", api.HandleGETS).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("[INFO] Started Distributed KV Store server.")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
