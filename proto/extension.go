package presence

import (
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/timestamppb"
    "time"
)

func (x *ServerPresence) MarshalBinary() ([]byte, error) {
    return proto.Marshal(x)
}

func (x *ClientPresence) MarshalBinary() ([]byte, error) {
    return proto.Marshal(x)
}

func GetPlayerState(updatedAt *timestamppb.Timestamp) PlayerState {
    switch {
    case updatedAt.AsTime().After(time.Now().Add(-2 * time.Minute)):
        return PlayerState_Online
    case updatedAt.AsTime().After(time.Now().Add(-10 * time.Minute)):
        return PlayerState_Idle
    default:
        return PlayerState_Offline
    }
}
