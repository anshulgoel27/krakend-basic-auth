package auth

import (
	"encoding/json"
	"errors"

	"github.com/luraproject/lura/v2/config"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github_com/anshulgoel27/krakend-basic-auth"

// Credentials contains the pair user:pass
type Credentials struct {
	User string
	Pass string
}

var ErrNoConfig = errors.New("no config defined for the module")

func ParseConfig(cfg config.ExtraConfig) (Credentials, error) {
	res := Credentials{}
	e, ok := cfg[Namespace]
	if !ok {
		return res, ErrNoConfig
	}
	b, err := json.Marshal(e)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
