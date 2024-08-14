package configs

import "github.com/spf13/viper"

type Conf struct {
	IPLimitMaxReq     int    `mapstructure:"IP_LIMIT_MAX_REQUEST"`
	IPBlockTimeSec    int    `mapstructure:"IP_BLOCK_TIME_SECONDS"`
	WebServerPort     string `mapstructure:"WEBSERVER_PORT"`
	DbAddress         string `mapstructure:"DB_ADDRESS"`
	DbPort            string `mapstructure:"DB_PORT"`
}

func Load(path string) (*Conf, error) {
	var conf *Conf
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
