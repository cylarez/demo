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
    "sort"
    "strings"
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
    err = s.rdb.HSet(ctx, config.KeyPresence, config.FieldPrefixClient+req.PlayerId, &update)
    if err != nil {
        core.Log.Println(err)
    }
    return
}

func (s *Controller) ServerUpdate(ctx context.Context, req *pb.ServerUpdateRequest) (resp *pb.ServerUpdateResponse, err error) {
    if len(req.PlayerIds) == 0 || req.Map == "" || req.ServerId == "" {
        err = status.Error(codes.InvalidArgument, "Invalid request arguments")
        return
    }
    updates := map[string]interface{}{}
    for _, playerId := range req.PlayerIds {
        updates[config.FieldPrefixServer+playerId] = &pb.ServerPresence{
            ServerId:  req.ServerId,
            Map:       req.Map,
            UpdatedAt: timestamppb.Now(),
        }
    }
    err = s.rdb.HSet(ctx, config.KeyPresence, updates)
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
    var fields []string
    for _, playerId := range req.GetPlayerIds() {
        fields = append(fields, config.FieldPrefixServer+playerId, config.FieldPrefixClient+playerId)
    }
    // Get Server Presence values
    results, err := s.rdb.HMGet(ctx, config.KeyPresence, fields...)
    if err != nil {
        core.Log.Println(err)
        return
    }
    // Sort to ensure Server presence comes first
    var keys []string
    for key := range results {
        keys = append(keys, key)
    }
    sort.SliceStable(keys, func(i, j int) bool {
        return keys[i] > keys[j]
    })

    resp = &pb.ListPlayerResponse{Players: make(map[string]*pb.Player)}
    for _, field := range keys {
        item := results[field]
        parts := strings.Split(field, ":")
        prefix := parts[0]
        playerId := parts[1]
        switch prefix + ":" {
        case config.FieldPrefixServer:
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
                ServerId:  p.ServerId,
                Map:       p.Map,
                UpdatedAt: p.UpdatedAt,
            }
            continue
        case config.FieldPrefixClient:
            var (
                presence *pb.Player
                hasItem  bool
            )
            clientPresence := new(pb.ClientPresence)
            err = proto.Unmarshal([]byte(item.(string)), clientPresence)
            if err != nil {
                core.Log.Println(err)
                continue
            }
            if presence, hasItem = resp.Players[playerId]; !hasItem {
                // We could choose to return player with empty Server Map value
                core.Log.Printf("WARNING: Missing Server presence for %s", playerId)
                continue
            }
            state := pb.GetPlayerState(clientPresence.UpdatedAt)
            if clientPresence.UpdatedAt.AsTime().After(presence.UpdatedAt.AsTime()) {
                // Client update is more recent
                presence.UpdatedAt = clientPresence.UpdatedAt
                presence.State = state
            }
            if presence.State == pb.PlayerState_Offline {
                return
            }
            presence.Action = clientPresence.Action
        }
    }

    return
}
