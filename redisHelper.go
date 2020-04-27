package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var DB redisHelper

type redisHelper struct {
	pool *redis.Pool
}

func (rh redisHelper) Session() *redisSession {
	return &redisSession{
		conn: rh.pool.Get(),
	}
}

func (rh redisHelper) Init(host string, port int) error {
	DB.pool = &redis.Pool{
		MaxActive:   100,
		MaxIdle:     80,
		IdleTimeout: time.Duration(8) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
		},
	}
	return nil
}

//db session

type redisSession struct {
	conn redis.Conn
}

func (rs redisSession) Close() {
	rs.conn.Close()
}

func (rs redisSession) HGETALL(key string, out interface{}) error {
	values, err := redis.Values(rs.conn.Do("HGETALL", key))
	if err != nil {
		return err
	}
	err = redis.ScanStruct(values, out)
	return err
}

func (rs redisSession) HMSET(key string, data interface{}) error {
	_, err := rs.conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(data)...)
	return err
}

func (rs redisSession) INC(key string) error {
	_, err := rs.conn.Do("INCR", key)
	return err
}

func (rs redisSession) SetADD(key string, value string) error {
	_, err := rs.conn.Do("SADD", key, value)
	return err
}

func (rs *redisSession) ExpiresAt(key string, t time.Time) error {
	_, err := rs.conn.Do("EXPIREAT", key, t.Unix())
	return err
}
