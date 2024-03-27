# Demo

## Storage

Stored in Redis

`ServerUpdate` will update the players in one HSET command
```
HSET presence:server playerId1 {data} playerId2 {data}
```

`ClientUpdate` will update the players in one HSET command
```
HSET presence:client playerId1 {data} playerId2 {data}
```

`ListPlayers` will fetch both HSET and merge data into one map.
```
HMGet presence:server playerId1 playerId2
HMGet presence:client playerId1 playerId2
```

If 2 calls to Redis becomes an issue, it can be made in one single op using a Lua script.
