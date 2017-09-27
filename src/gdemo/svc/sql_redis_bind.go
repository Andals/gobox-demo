package svc

import (
	"gdemo/dao"
	"strconv"
)

type sqlRedisBindSvc struct {
	*baseSvc
	*sqlBaseSvc
	*redisBaseSvc

	redisKeyPrefix string
}

func NewSqlRedisBindSvc(bs *baseSvc, ss *sqlBaseSvc, rs *redisBaseSvc, redisKeyPrefix string) *sqlRedisBindSvc {
	return &sqlRedisBindSvc{bs, ss, rs, redisKeyPrefix}
}

func (this *sqlRedisBindSvc) redisKeyForEntity(id int64) string {
	return this.redisKeyPrefix + "_entity_" + this.entityName + "_id_" + strconv.FormatInt(id, 10)
}

func (this *sqlRedisBindSvc) Insert(tableName string, colNames []string, expireSeconds int64, entities ...interface{}) ([]int64, error) {
	ids, err := this.sqlBaseSvc.Insert(tableName, colNames, entities...)
	if err != nil {
		return nil, err
	}

	for i, id := range ids {
		this.saveHashEntityToRedis(this.redisKeyForEntity(id), entities[i], expireSeconds)
	}

	return ids, nil
}

func (this *sqlRedisBindSvc) GetById(tableName string, id, expireSeconds int64, entityPtr interface{}) (bool, error) {
	rk := this.redisKeyForEntity(id)

	find, err := this.getHashEntityFromRedis(rk, entityPtr)
	if err != nil {
		return this.sqlBaseSvc.GetById(tableName, id, entityPtr)
	}
	if find {
		return true, nil
	}

	find, err = this.sqlBaseSvc.GetById(tableName, id, entityPtr)
	if err != nil {
		this.elogger.Error([]byte("getById from mysql error"))
		return false, err
	}
	if !find {
		return false, nil
	}

	this.saveHashEntityToRedis(rk, entityPtr, expireSeconds)

	return true, nil
}

func (this *sqlRedisBindSvc) DeleteById(tableName string, id int64) (bool, error) {
	result := this.dao.DeleteById(tableName, id)
	if result.Err != nil {
		this.elogger.Error([]byte("delete from mysql error: " + result.Err.Error()))
		return false, result.Err
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	rk := this.redisKeyForEntity(id)
	_, err := this.rclient.Do("del", rk)
	if err != nil {
		this.elogger.Warning([]byte("del key " + rk + " from redis failed: " + err.Error()))
	}

	return true, nil
}

func (this *sqlRedisBindSvc) UpdateById(tableName string, id int64, newEntityPtr interface{}, updateFields map[string]bool, expireSeconds int64) (bool, error) {
	setItems, err := this.sqlBaseSvc.UpdateById(tableName, id, newEntityPtr, updateFields)

	if err != nil {
		return false, err
	}
	if setItems == nil {
		return false, nil
	}

	this.updateSqlHashEntity(this.redisKeyForEntity(id), setItems, expireSeconds)

	return true, nil
}

func (this *sqlRedisBindSvc) updateSqlHashEntity(key string, setItems []*dao.SqlColQueryItem, expireSeconds int64) error {
	cnt := len(setItems)*2 + 1
	args := make([]interface{}, cnt)
	args[0] = key

	for si, ai := 0, 1; ai < cnt; si++ {
		args[ai] = setItems[si].Name
		ai++
		args[ai] = setItems[si].Value
		ai++
	}

	this.rclient.Send("hmset", args...)
	if expireSeconds > 0 {
		this.rclient.Send("expire", expireSeconds)
	}
	_, err := this.rclient.ExecPipelining()
	if err != nil {
		this.elogger.Warning([]byte("hmset key " + key + " to redis failed: " + err.Error()))
	}

	return err
}
