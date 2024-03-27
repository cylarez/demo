package controller

import (
    "context"
    "demo-1/internal/config"
    "demo-1/internal/core"
    pb "demo-1/proto"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
    pb.UnimplementedPresenceServer
    rdb core.RedisDal
}

func NewController(rdb core.RedisDal) *Controller {
    return &Controller{rdb: rdb}
}

func (s *Controller) ClientUpdate(ctx context.Context, req *pb.ClientUpdateRequest) (resp *pb.ClientUpdateResponse, err error) {
    if req.PlayerId == "" || req.Action == "" {
        err = status.Error(codes.InvalidArgument, "Invalid request arguments")
        return
    }
    update := pb.ClientPresence{
        Action:    req.Action,
        UpdatedAt: timestamppb.Now(),
    }
    err = s.rdb.HSet(ctx, config.RedisKeyClient, req.PlayerId, &update)
    if err != nil {
        core.Log.Println(err)
    }
    return
}

func (s *Controller) ListPlayer(ctx context.Context, req *pb.ListPlayerRequest) (resp *pb.ListPlayerResponse, err error) {
    if len(req.PlayerIds) == 0 {
        err = status.Error(codes.InvalidArgument, "Invalid request arguments")
        return
    }
    // Get Server Presence values
    results, err := s.rdb.HMGet(ctx, config.RedisKeyServer, req.PlayerIds...)
    if err != nil {
        core.Log.Println(err)
        return
    }
    resp = &pb.ListPlayerResponse{Players: make(map[string]*pb.Player)}
    for playerId, item := range results {
        p := new(pb.ServerPresence)
        err = proto.Unmarshal([]byte(item.(string)), p)
        if err != nil {
            core.Log.Println(err)
            continue
        }
        state := pb.GetPlayerState(p.UpdatedAt)
        if state == pb.PlayerState_Offline {
            continue
        }
        resp.Players[playerId] = &pb.Player{
            State:     state,
            Map:       p.Map,
            UpdatedAt: p.UpdatedAt,
        }
    }
    // Get Client Presence values
    results, err = s.rdb.HMGet(ctx, config.RedisKeyClient, req.PlayerIds...)
    if err != nil {
        core.Log.Println(err)
        return
    }
    // merge data between client and server presence values
    for playerId, item := range results {
        var (
            presence *pb.Player
            hasItem  bool
        )
        if presence, hasItem = resp.Players[playerId]; !hasItem {
            // We could choose to return player with empty Server Map value
            core.Log.Printf("WARNING: Missing Server presence for %s", playerId)
            continue
        }
        p := new(pb.ClientPresence)
        err = proto.Unmarshal([]byte(item.(string)), p)
        if err != nil {
            core.Log.Println(err)
            continue
        }
        state := pb.GetPlayerState(p.UpdatedAt)
        if p.UpdatedAt.AsTime().After(presence.UpdatedAt.AsTime()) {
            // Client update is more recent
            presence.UpdatedAt = p.UpdatedAt
            presence.State = state
        }
        if presence.State == pb.PlayerState_Offline {
            return
        }
        presence.Action = p.Action
    }
    return
}

func (s *Controller) ServerUpdate(ctx context.Context, req *pb.ServerUpdateRequest) (resp *pb.ServerUpdateResponse, err error) {
    if len(req.PlayerIds) == 0 || req.Map == "" {
        err = status.Error(codes.InvalidArgument, "Invalid request arguments")
        return
    }
    updates := map[string]interface{}{}
    for _, i := range req.PlayerIds {
        updates[i] = &pb.ServerPresence{
            Map:       req.Map,
            UpdatedAt: timestamppb.Now(),
        }
    }
    err = s.rdb.HSet(ctx, config.RedisKeyServer, updates)
    if err != nil {
        core.Log.Println(err)
    }
    return
}
