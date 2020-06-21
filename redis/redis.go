package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// RNewPool : Make redis pool with max 20 idle connections
func RNewPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// RExists : Check if key is there in redis
func RExists(key string, pool *redis.Pool) (exists bool, err error) {
	conn := pool.Get()
	defer conn.Close()
	exists, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return exists, err
	}
	return exists, nil
}

// RSet : set the redis keys
func RSet(key string, setValue string, pool *redis.Pool) (value string, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.String(conn.Do("SET", key, setValue))
	if err != nil {
		return value, err
	}
	return value, nil
}

// RGet : get the redis keys
func RGet(key string, pool *redis.Pool) (value string, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return value, err
	}
	return value, nil
}

// RIncr : incr the redis keys
func RIncr(key string, pool *redis.Pool) (value int64, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return value, err
	}
	return value, nil
}

// RDecr : get the redis keys
func RDecr(key string, pool *redis.Pool) (value int64, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.Int64(conn.Do("DECR", key))
	if err != nil {
		return value, err
	}
	return value, nil
}

// RIncrby : incr the redis keys by x
func RIncrby(key string, incrby int64, pool *redis.Pool) (value int64, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.Int64(conn.Do("INCRBY", key, incrby))
	if err != nil {
		return value, err
	}
	return value, nil
}

// RDecrby : decr the redis keys by y
func RDecrby(key string, decrby int64, pool *redis.Pool) (value int64, err error) {
	conn := pool.Get()
	defer conn.Close()
	// Perform redis command get
	value, err = redis.Int64(conn.Do("DECRBY", key, decrby))
	if err != nil {
		return value, err
	}
	return value, nil
}
