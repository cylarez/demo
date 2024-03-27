# Demo

## Storage

Stored in Redis using `HSET`

Server updates use the "server:" prefix on each field.

Client updates use the "client:" prefix on each field.

Which allows us to fetch all the players presence in one single Redis call (HMGET)


## APIs

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

## Alternative

Simpler implementation is possible using two different HSET but `ListPlayers` needs to make 2 calls to Redis.

Code can be found in the branch: [`with-2-hset`](https://github.com/cylarez/demo/tree/with-2-hset)


