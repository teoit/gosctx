package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
)

var (
	LogLevel = "info"
	AppEnv   = "dev"

	AppDev  = "dev"
	AppProd = "prd"
	AppStg  = "stg"

	RedisUri  = ""
	MaxActive = 0
	MaxIde    = 0

	SessionLifetime = 0
	Location        = "Asia/Ho_Chi_Minh"

	defaultRedisMaxActive  = 0
	defaultRedisMaxIdle    = 10
	defaultSessionLifetime = 120
)

func init() {

	if err := ReadFileEnv(); err != nil {
		panic("read file env")
	}

	flag.StringVar(&LogLevel, "log-level", "info", "logger level for appsLog level: panic | fatal | error | warn | info | debug | trace")
	flag.StringVar(&AppEnv, "app-env", "dev1", "Env for service. Ex: dev | stg | prd")
	flag.IntVar(&SessionLifetime, "session-lifetime", defaultSessionLifetime, "(For auth) Session Lifetime default 120 minutes")
	flag.StringVar(&Location, "location", "Asia/Ho_Chi_Minh", "local time, should be local time")

	flag.StringVar(&RedisUri, "redis-uri", "redis://localhost:6379", "(For go-redis) Redis connection-string. Ex: redis://localhost/0")
	flag.IntVar(&MaxActive, "redis-pool-max-active", defaultRedisMaxActive, "(For go-redis) Override redis pool MaxActive")
	flag.IntVar(&MaxIde, "redis-pool-max-idle", defaultRedisMaxIdle, "(For go-redis) Override redis pool MaxIdle")

	// override value found in MY_RAZ_VALUE with command line flag value -raz-value=foo
	flag.Parse()
}

func ReadFileEnv() error {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			return err
		}
	}
	return nil
}
