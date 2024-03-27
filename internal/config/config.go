package config

import "os"

var (
    ServerAddr     = getEnv("SERVER_ADDRESS", ":50051")
    RedisAddr      = getEnv("REDIS_ADDRESS", ":6379")
    RedisPwd       = getEnv("REDIS_PASSWORD", "")
    RedisKeyClient = getEnv("REDIS_KEY_PRESENCE_CLIENT", "presence:client")
    RedisKeyServer = getEnv("REDIS_KEY_PRESENCE_SERVER", "presence:server")
)

func getEnv(k, d string) string {
    val := os.Getenv(k)
    if val == "" {
        return d
    }
    return val
}
