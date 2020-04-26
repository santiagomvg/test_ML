package main

var db redisHelper

type redisHelper struct{}

func (rh redisHelper) Session() *redisSession {

}

type redisSession struct{}

func (rs redisSession) Close() {

}

func (rs redisSession) HGET(key string, values string, out interface{}) error {

}

func (rs redisSession) HSET(key string, values string, in interface{}) error {

}
