package svc

import (
	"andals/gobox/redis"

	"encoding/json"
	"reflect"
)

const (
	ENTITY_REDIS_HASH_FIELD_TAG = "redis"
)

type redisBaseSvc struct {
	*baseSvc

	rclient *redis.Client
}

func newRedisBaseSvc(bs *baseSvc, rclient *redis.Client) *redisBaseSvc {
	return &redisBaseSvc{
		baseSvc: bs,

		rclient: rclient,
	}
}

func (this *redisBaseSvc) saveJsonDataToRedis(key string, v interface{}, expireSeconds int64) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		this.elogger.Warning([]byte("json_encode " + key + " error: " + err.Error()))
		return err
	}

	_, err = this.rclient.Do("set", key, string(jsonBytes), "ex", expireSeconds)
	if err != nil {
		this.elogger.Warning([]byte("set " + key + " to redis error: " + err.Error()))
		return err
	}

	return nil
}

func (this *redisBaseSvc) getJsonDataFromRedis(key string, v interface{}) (bool, error) {
	reply, err := this.rclient.Do("get", key)
	if err != nil {
		this.elogger.Warning([]byte("get " + key + " from redis error: " + err.Error()))
		return false, err
	}

	if reply == nil {
		return false, nil
	}

	jsonBytes, err := reply.Bytes()
	if err != nil {
		this.elogger.Warning([]byte("reply " + key + " from redis error: " + err.Error()))
		return false, err
	}

	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		this.elogger.Warning([]byte("json_decode " + key + " from redis error: " + err.Error()))
		return false, err
	}

	return true, nil
}

func (this *redisBaseSvc) saveHashEntityToRedis(key string, entityPtr interface{}, expireSeconds int64) error {
	eargs := this.reflectSaveHashEntityArgs(reflect.ValueOf(entityPtr).Elem())
	args := make([]interface{}, len(eargs)+1)
	args[0] = key
	for i, arg := range eargs {
		args[i+1] = arg
	}

	this.rclient.Send("hmset", args...)
	if expireSeconds > 0 {
		this.rclient.Send("expire", key, expireSeconds)
	}
	_, err := this.rclient.ExecPipelining()
	if err != nil {
		this.elogger.Warning([]byte("hmset " + key + " to redis error: " + err.Error()))
		return err
	}

	return nil
}

func (this *redisBaseSvc) getHashEntityFromRedis(key string, entityPtr interface{}) (bool, error) {
	reply, err := this.rclient.Do("hgetall", key)
	if err != nil {
		this.elogger.Warning([]byte("hgetall " + key + " from redis error: " + err.Error()))
		return false, err
	}

	if reply.ArrReplyIsNil() {
		return false, nil
	}

	err = reply.Struct(entityPtr)
	if err != nil {
		this.elogger.Warning([]byte("reply to struct " + key + " from redis error: " + err.Error()))
		return false, err
	}

	return true, nil
}

func (this *redisBaseSvc) reflectSaveHashEntityArgs(rev reflect.Value) []interface{} {
	var args []interface{}
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			args = this.reflectSaveHashEntityArgs(revf)
			continue
		}

		retf := ret.Field(i)
		fn, ok := retf.Tag.Lookup(ENTITY_REDIS_HASH_FIELD_TAG)
		if ok {
			args = append(args, fn, revf.Interface())
		}
	}

	return args
}
