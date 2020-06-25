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



#### POST /counter/:key/init

* initiate the counter, gets the data using the query given from the database given
* datasource supported : mysql/postgres

input
```javascript
{
  "query": "select count(*) as count from tasks;",
  "connection-details": {
    "datasource": "mysql",
    "host": "127.0.0.1",
    "username": "test",
    "password": "1234",
    "port": 3306,
    "dbname": "test"
  }
}
```

samples

```shell
curl -XPOST localhost:4118/counter/taskcounter/init -H "Content-Type: application/json" -d '{"query":"select count(*) as count from tasks;", "connection-details": { "datasource": "mysql", "host": "127.0.0.1", "username": "test", "password": "1234", "port": 3306, "dbname": "test"}}'

curl -XPOST localhost:4118/counter/counter1/init -H "Content-Type: application/json" -d '{"query":"select count(*) from counter1;", "connection-details": { "datasource": "postgres", "host": "127.0.0.1", "username": "postgres", "password": "1234", "port": 5432, "dbname": "test"}}'
```
