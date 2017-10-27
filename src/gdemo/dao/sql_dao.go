package dao

import (
	"andals/gobox/mysql"

	"database/sql"
)

type SqlExecResult struct {
	Err          error
	LastInsertId int64
	RowsAffected int64
}

type SqlDao struct {
	*mysql.Client
}

func NewSqlDao(client *mysql.Client) *SqlDao {
	return &SqlDao{
		Client: client,
	}
}

func (this *SqlDao) Insert(tableName string, colNames []string, colsValues ...[]interface{}) *SqlExecResult {
	sqb := new(SqlQueryBuilder)
	sqb.Insert(tableName, colNames...).
		Values(colsValues...)

	return convertSqlResultToSqlExecResult(this.Exec(sqb.Query(), sqb.Args()...))
}

func (this *SqlDao) DeleteById(tableName string, id int64) *SqlExecResult {
	sqb := new(SqlQueryBuilder)
	sqb.Delete(tableName).
		WhereConditionAnd(NewSqlColQueryItem("id", SQL_COND_EQUAL, id))

	return convertSqlResultToSqlExecResult(this.Exec(sqb.Query(), sqb.Args()...))
}

func (this *SqlDao) UpdateById(tableName string, id int64, setItems ...*SqlColQueryItem) *SqlExecResult {
	sqb := new(SqlQueryBuilder)
	sqb.Update(tableName).
		Set(setItems...).
		WhereConditionAnd(NewSqlColQueryItem("id", SQL_COND_EQUAL, id))

	return convertSqlResultToSqlExecResult(this.Exec(sqb.Query(), sqb.Args()...))
}

func (this *SqlDao) SelectById(tableName string, what string, id int64) *sql.Row {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(NewSqlColQueryItem("id", SQL_COND_EQUAL, id))

	return this.QueryRow(sqb.Query(), sqb.Args()...)
}

func (this *SqlDao) SelectByIds(tableName string, what string, orderyBy string, ids ...int64) (*sql.Rows, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(NewSqlColQueryItem("id", SQL_COND_IN, ids)).
		OrderBy(orderyBy)

	return this.Query(sqb.Query(), sqb.Args()...)
}

func (this *SqlDao) SimpleQueryAnd(tableName string, what string, orderBy string, offset, cnt int64, condItems ...*SqlColQueryItem) (*sql.Rows, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(condItems...).
		OrderBy(orderBy).
		Limit(offset, cnt)

	return this.Query(sqb.Query(), sqb.Args()...)
}

func (this *SqlDao) SimpleQueryOr(tableName string, what string, orderBy string, offset, cnt int64, condItems ...*SqlColQueryItem) (*sql.Rows, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionOr(condItems...).
		OrderBy(orderBy).
		Limit(offset, cnt)

	return this.Query(sqb.Query(), sqb.Args()...)
}

func (this *SqlDao) SimpleTotalAnd(tableName string, condItems ...*SqlColQueryItem) (int64, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select("count(1)", tableName).
		WhereConditionAnd(condItems...)

	var total int64
	err := this.QueryRow(sqb.Query(), sqb.Args()...).Scan(&total)

	return total, err
}

func (this *SqlDao) SimpleTotalOr(tableName string, condItems ...*SqlColQueryItem) (int64, error) {
	sqb := new(SqlQueryBuilder)
	sqb.Select("count(1)", tableName).
		WhereConditionOr(condItems...)

	var total int64
	err := this.QueryRow(sqb.Query(), sqb.Args()...).Scan(&total)

	return total, err
}

func convertSqlResultToSqlExecResult(sqlResult sql.Result, err error) *SqlExecResult {
	execResult := new(SqlExecResult)
	if err != nil {
		execResult.Err = err
	} else {
		lid, err := sqlResult.LastInsertId()
		if err != nil {
			execResult.Err = err
		} else {
			execResult.LastInsertId = lid
			ra, err := sqlResult.RowsAffected()
			if err != nil {
				execResult.Err = err
			} else {
				execResult.RowsAffected = ra
			}
		}
	}

	return execResult
}
