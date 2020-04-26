package main

var db redisHelper

type redisHelper struct{}

func (rh redisHelper) Session() *redisSession {
	return nil
}

type redisSession struct{}

func (rs redisSession) Close() {

}

func (rs redisSession) HGET(key string, values string, out interface{}) error {
	return nil
}

func (rs redisSession) HSET(key string, values string, in interface{}) error {
	return nil
}
