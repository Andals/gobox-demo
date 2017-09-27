package svc

import (
	"andals/gobox/mysql"
	"andals/gobox/redis"
	"andals/golog"

	"reflect"
)

const (
	DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS = 60 * 30
)

var demoEntityType reflect.Type = reflect.TypeOf(DemoEntity{})
var demoColNames []string = reflectColNames(demoEntityType)

type DemoEntity struct {
	SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}

type DemoSvc struct {
	*sqlRedisBindSvc
}

func NewDemoSvc(elogger golog.ILogger, mclient *mysql.Client, redisKeyPrefix string, rclient *redis.Client) *DemoSvc {
	bs := newBaseSvc(elogger)
	sbs := newSqlBaseSvc(bs, mclient, "demo")

	return &DemoSvc{
		&sqlRedisBindSvc{
			baseSvc:      bs,
			sqlBaseSvc:   sbs,
			redisBaseSvc: newRedisBaseSvc(bs, rclient),

			redisKeyPrefix: redisKeyPrefix,
		},
	}
}

func (this *DemoSvc) Insert(entities ...*DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return this.sqlRedisBindSvc.Insert(this.entityName, demoColNames, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, is...)
}

func (this *DemoSvc) GetById(id int64) (*DemoEntity, error) {
	entity := new(DemoEntity)

	find, err := this.sqlRedisBindSvc.GetById(this.entityName, id, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, entity)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (this *DemoSvc) DeleteById(id int64) (bool, error) {
	return this.sqlRedisBindSvc.DeleteById(this.entityName, id)
}

func (this *DemoSvc) UpdateById(id int64, newEntity *DemoEntity, updateFields map[string]bool) (bool, error) {
	return this.sqlRedisBindSvc.UpdateById(this.entityName, id, newEntity, updateFields, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS)
}

func (this *DemoSvc) ListByIds(ids ...int64) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := this.sqlBaseSvc.ListByIds(this.entityName, ids, "id desc", demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (this *DemoSvc) SimpleQueryAnd(sqp *SqlQueryParams) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := this.sqlBaseSvc.SimpleQueryAnd(this.entityName, sqp, demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
