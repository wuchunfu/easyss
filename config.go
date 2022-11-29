package easyss

import (
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

type Config struct {
	Server           string `json:"server"`
	ServerPort       int    `json:"server_port"`
	LocalPort        int    `json:"local_port"`
	Password         string `json:"password"`
	Method           string `json:"method"` // encryption method
	Timeout          int    `json:"timeout"`
	BindALL          bool   `json:"bind_all"`
	DisableUTLS      bool   `json:"disable_utls"`
	EnableForwardDNS bool   `json:"enable_forward_dns"`
	Tun2socksModel   string `json:"tun2socks_model"`
	ConfigFile       string `json:"-"`
}

func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	config = &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	return
}

func UpdateConfig(old, ne *Config) {
	newVal := reflect.ValueOf(ne).Elem()
	oldVal := reflect.ValueOf(old).Elem()

	for i := 0; i < newVal.NumField(); i++ {
		newField := newVal.Field(i)
		oldField := oldVal.Field(i)

		switch newField.Kind() {
		case reflect.String:
			s := newField.String()
			if s != "" {
				oldField.SetString(s)
			}
		case reflect.Int:
			i := newField.Int()
			if i != 0 {
				oldField.SetInt(i)
			}
		case reflect.Bool:
			b := newField.Bool()
			if b {
				oldField.SetBool(b)
			}
		}
	}

	if old.Method == "" {
		old.Method = "aes-256-gcm"
	}
	if old.Timeout <= 0 || old.Timeout > 60 {
		old.Timeout = 60
	}
}

func ExampleJSONConfig() string {
	example := Config{
		Server:     "example.com",
		ServerPort: 9999,
		LocalPort:  2080,
		Password:   "your-pass",
		Method:     "aes-256-gcm",
		Timeout:    30,
		BindALL:    false,
	}

	b, _ := json.MarshalIndent(example, "", "    ")
	return string(b)
}
