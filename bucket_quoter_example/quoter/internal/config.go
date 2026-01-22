package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CmdGlobal struct {
	Cmd  *cobra.Command
	Opts *ConfYaml
	Log  *TLog
}

type ConfYaml struct {
	LogOptions *TLogOptions
	Overrides  *ConfigOverrides

	RateLimiter RateLimiterSection `yaml:"rateLimiter"`

	Buckets BucketsSection `yaml:"buckets"`
}

type RateLimiterSection struct {
	Socket   string `yaml:"socket"`
	Port     int    `yaml:"port"`
	Certfile string `yaml:"certfile"`
	Keyfile  string `yaml:"keyfile"`
	PidFile  string `yaml:"pidFile"`
}

type BucketsSection struct {
	Buckets map[string]*BucketSettings
}

type BucketSettings struct {
	Inflow   int `yaml:"inflow"`
	Capacity int `yaml:"capacity"`
}

var defaultConf = []byte(`
log:
  # logging format could be "string" or "json"
  format: "string"

  log: "stdout"

  # level of debugging could be "debug", "info", could be
  # overrides by command line switch
  level: "debug"

rateLimiter:
  socket: "/var/run/ratelimiter.socket"
  port: 8443
  certfile: "certs/dns-api.crt"
  keyfile: "certs/dns-api.key"
  # detach process mode: pidfile
  pidfile: "/var/run/ratelimiter.pid"

buckets:
  "897d9f58-6b42-4ca7-8229-2e04056490b7":
    "inflow": 10
    "capacity": 10
`)

func LoadConf(confPath string, Overrides ConfigOverrides) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ratelimiter") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var err error
	if confPath != "" {
		if _, err = os.Stat(confPath); err != nil {
			return conf, err
		}

		var content []byte
		if content, err = ioutil.ReadFile(confPath); err != nil {
			return conf, err
		}
		if err = viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			//fmt.Printf("using config file as: \"%s\"\n", viper.ConfigFileUsed())
			return conf, err
		}
	} else {
		// Search config in home directory with name "hhs"
		viper.AddConfigPath("/etc/ratelimiter")
		viper.AddConfigPath("$HOME/.ratelimiter")
		viper.AddConfigPath(".")
		viper.SetConfigName("ratelimiter")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			//fmt.Printf("using config file as: \"%s\"\n", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	conf.Overrides = &Overrides

	// logging parameters could be overridden by global
	// command line switches
	var LogOptions TLogOptions
	LogOptions.Format = viper.GetString("log.format")

	LogOptions.Level = viper.GetString("log.level")
	if Overrides.Debug {
		LogOptions.Level = "debug"
	}

	//  # "stdout" (used in systemd startup case)
	//  log: "stdout"

	LogOptions.Log = viper.GetString("log.log")
	if len(Overrides.Log) > 0 {
		LogOptions.Log = Overrides.Log
	}

	conf.LogOptions = &LogOptions

	// limiter settings

	conf.RateLimiter.Socket = viper.GetString("rateLimiter.socket")
	conf.RateLimiter.Port = viper.GetInt("rateLimiter.port")
	conf.RateLimiter.Certfile = viper.GetString("rateLimiter.certfile")
	conf.RateLimiter.Keyfile = viper.GetString("rateLimiter.keyfile")
	conf.RateLimiter.PidFile = viper.GetString("rateLimiter.pidFile")

	var buckets = make(map[string]*BucketSettings)
	for key, item := range viper.GetStringMap("buckets") {
		b, ok := item.(map[string]interface{})
		if ok && b != nil {
			buckets[key] = &BucketSettings{
				Inflow:   b["inflow"].(int),
				Capacity: b["capacity"].(int),
			}
		}
	}
	conf.Buckets = BucketsSection{
		Buckets: buckets,
	}

	return conf, nil
}
