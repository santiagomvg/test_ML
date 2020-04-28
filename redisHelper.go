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
	_ = rs.conn.Close()
}

func (rs redisSession) HGETALL(key string, out interface{}) error {
	values, err := redis.Values(rs.conn.Do("HGETALL", key))
	if err != nil {
		return err
	}
	err = redis.ScanStruct(values, out)
	return err
}

func (rs redisSession) Raw(cmd string, args ...interface{}) (reply interface{}, err error) {
	return rs.conn.Do(cmd, args...)
}

func (rs redisSession) StoreJson(key string, data interface{}) error {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}
	_, err := rs.conn.Do("HSET", key, "json", buf.Bytes())
	return err
}

func (rs redisSession) ReadJson(key string, out interface{}) error {

	reply, err := redis.Values(rs.conn.Do("HMGET", key, "json"))
	if err != nil {
		return err
	}

	var jsonData string
	if _, err := redis.Scan(reply, &jsonData); err != nil {
		return err
	}
	return json.NewDecoder(strings.NewReader(jsonData)).Decode(out)
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
