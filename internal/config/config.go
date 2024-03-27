package config

import "os"

var (
    ServerAddr        = getEnv("SERVER_ADDRESS", ":50051")
    RedisAddr         = getEnv("REDIS_ADDRESS", ":6379")
    RedisPwd          = getEnv("REDIS_PASSWORD", "")
    KeyPresence       = getEnv("KEY_PRESENCE", "presence")
    FieldPrefixClient = getEnv("REDIS_KEY_PREFIX_CLIENT", "client:")
    FieldPrefixServer = getEnv("REDIS_KEY_PREFIX_SERVER", "server:")
)

func getEnv(k, d string) string {
    val := os.Getenv(k)
    if val == "" {
        return d
    }
    return val
}
