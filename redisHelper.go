package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
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

func (rs redisSession) HGET(key string, values string, out interface{}) error {
	reply, err := redis.Values(rs.conn.Do("HMGET", key, "value"))
	if err != nil {
		return err
	}

	var jsonData string
	if _, err := redis.Scan(reply, &jsonData); err != nil {
		return err
	}
	return json.NewDecoder(strings.NewReader(jsonData)).Decode(&out)
}

func (rs redisSession) HSET(key string, values string, data interface{}) error {
	jsonData := &bytes.Buffer{}
	if err := json.NewEncoder(jsonData).Encode(data); err != nil {
		return err
	}
	_, err := rs.conn.Do("HSET", key, "value", jsonData)
	return err
}

func (rs *redisSession) ExpiresAt(key string, t time.Time) error {
	_, err := rs.conn.Do("EXPIREAT", key, t.Unix())
	return err
}
