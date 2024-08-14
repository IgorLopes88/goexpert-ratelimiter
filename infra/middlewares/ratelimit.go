package middlewares

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/igorlopes88/goexpert-ratelimiter/configs"
	"github.com/igorlopes88/goexpert-ratelimiter/infra/storage"
	"github.com/spf13/viper"
)

type AccessSettings struct {
	LimitRPS    int
	BlockForSec int
}

type RateLimiter struct {
	IPSettings     *AccessSettings
	TokensSettings *map[string]*AccessSettings
	StorageAdapter storage.Storage
}

func NewRateLimiter(config *RateLimiter) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return config.Handler(handler)
	}
}

func LoadTokens(path string) (map[string]*AccessSettings, error) {
	var listToken map[string]*AccessSettings
	viper.SetConfigType("json")
	viper.AddConfigPath(path)
	viper.SetConfigFile("list-token.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&listToken)
	if err != nil {
		panic(err)
	}
	return listToken, err
}

func Load(conf configs.Conf, storage storage.Storage) (*RateLimiter, error) {
	listTokens, err := LoadTokens(".")
	if err != nil {
		return nil, err
	}
	var settings = &RateLimiter{
		IPSettings: &AccessSettings{
			LimitRPS:    conf.IPLimitMaxReq,
			BlockForSec: conf.IPBlockTimeSec,
		},
		TokensSettings: &listTokens,
		StorageAdapter: storage,
	}
	return settings, nil
}

func (rt *RateLimiter) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var block *time.Time
		var err error
		token := r.Header.Get("API_KEY")
		if token != "" {
			var accessToken *AccessSettings
			listTokens, exists := (*rt.TokensSettings)[token]
			if exists {
				accessToken = listTokens
			} else {
				accessToken = rt.IPSettings
			}
			block, err = rt.validadeAccess(r.Context(), "tk", token, accessToken)
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			block, err = rt.validadeAccess(r.Context(), "ip", host, rt.IPSettings)
		}
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}
		if block != nil {
			w.WriteHeader(429)
			w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame."))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (rt *RateLimiter) validadeAccess(ctx context.Context, key string, value string, settings *AccessSettings) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	block, err := rt.StorageAdapter.Search(ctx, key, value)
	if err != nil {
		return nil, err
	}
	if block == nil {
		success, count, err := rt.StorageAdapter.RegisterAccess(ctx, key, value, settings.LimitRPS)
		if err != nil {
			return nil, err
		}
		if success {
			log.Printf("request number %d from ip/token %v\n", count, value)
		} else {
			block, err = rt.StorageAdapter.Block(ctx, key, value, settings.BlockForSec)
			log.Printf("access denied to ip/token %v due to too many requests per second\n", value)
			if err != nil {
				return nil, err
			}
		}
	}
	if block != nil {
		log.Printf("ip/token %v blocked for %d seconds", value, settings.BlockForSec)
		return block, nil
	}
	return nil, nil
}
