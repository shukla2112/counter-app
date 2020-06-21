# counter-app
redis based counter app

## APIS
#### GET /counter/:key

```shell
 curl -XGET localhost:4118/counter/mykey
 {"key":"mykey","status":1,"value":"11"}
 ```

#### POST /counter/:key

**data**

operation: string : values can be - `incr/decr`

multiplier: int64 : values can be > 1

```shell
curl -XPOST localhost:4118/counter/mykey -H "Content-Type: application/json" -d '{"operation":"incr"}'
{"key":"mykey","status":1,"value":"12"}



curl -XPOST localhost:4118/counter/mykey -H "Content-Type: application/json" -d '{"operation":"decr"}'
{"key":"mykey","status":1,"value":"11"}


curl -XPOST localhost:4118/counter/mykey -H "Content-Type: application/json" -d '{"operation":"incr", "multiplier" : 3}'
{"key":"mykey","status":1,"value":"14"}

curl -XPOST localhost:4118/counter/mykey -H "Content-Type: application/json" -d '{"operation":"decr", "multiplier" : 4}'
{"key":"mykey","status":1,"value":"10"}
```
