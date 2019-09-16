package litmus

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func ReadConfig(key string, opts interface{}) error {
	return viper.UnmarshalKey(key, &opts)
}
