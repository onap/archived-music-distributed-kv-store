package api

import (
	"errors"
	"os"
)

func Initialise() error {
	if os.Getenv("CONSUL_IP") == "" {
		return errors.New("CONSUL_IP environment variable not set.")
	}

	Consul = &ConsulStruct{}
	KeyValues = &KeyValuesStruct{kvs: make(map[string]string)}

	err := Consul.InitializeConsulClient()
	if err != nil {
		return err
	}

	err = Consul.CheckConsulHealth()
	if err != nil {
		return err
	}

	return nil
}
