# Demo

## Storage

Stored in Redis

`ServerUpdate` will update the players in one HSET command
```
HSET presence server:playerId1 {data} server:playerId2 {data}
```

`ClientUpdate` will update the players in one HSET command
```
HSET presence client:playerId1 {data} client:playerId2 {data}
```

`ListPlayers` will fetch data with one HMGET using 2 fields per playerId:
```
HMGET presence client:playerId1 client:playerId2 server:playerId1 server:playerId2
```
Then we merge the results from client/server into one object per player.


