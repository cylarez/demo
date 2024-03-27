package main

import (
    "context"
    "demo-1/internal/config"
    "demo-1/internal/controller"
    "demo-1/internal/core"
    "github.com/redis/go-redis/v9"
    "net"

    pb "demo-1/proto"
    "google.golang.org/grpc"
)

func main() {
    core.Log.Println("Server is starting")
    lis, err := net.Listen("tcp", config.ServerAddr)
    if err != nil {
        core.Log.Fatalf("failed to listen: %v", err)
    }
    // Redis Client
    rc := redis.NewClient(&redis.Options{
        Addr:     config.RedisAddr,
        Password: config.RedisPwd,
    })
    if err = rc.Ping(context.Background()).Err(); err != nil {
        core.Log.Fatalln(err)
    }
    core.Log.Printf("Redis client connected at %s", config.RedisAddr)

    // Setup Controller
    s := grpc.NewServer()
    c := controller.NewController(core.NewRedis(rc))
    pb.RegisterPresenceServer(s, c)

    // Run Server
    core.Log.Printf("Server listening at %v", lis.Addr())
    if err = s.Serve(lis); err != nil {
        core.Log.Fatalf("failed to serve: %v", err)
    }
}
