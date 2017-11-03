package idgen

import "github.com/andals/gobox/mysql"

type SqlIdGenter struct {
	client *mysql.Client

	updateSql string
	selectSql string
}

func NewSqlIdGenter(client *mysql.Client) *SqlIdGenter {
	return &SqlIdGenter{
		client: client,

		updateSql: "UPDATE id_gen SET max_id = last_insert_id(max_id + 1) WHERE name = ?",
		selectSql: "SELECT last_insert_id()",
	}
}

func (this *SqlIdGenter) GenId(name string) (int64, error) {
	_, err := this.client.Exec(this.updateSql, name)
	if err != nil {
		return 0, err
	}

	var id int64
	err = this.client.QueryRow(this.selectSql).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
