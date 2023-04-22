package db

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

/*

CREATE TABLE "db_monitor"."tb_test" (
  "id" int8 NOT NULL GENERATED ALWAYS AS IDENTITY (
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
),
  "int_test" int4,
  "string_test" varchar COLLATE "pg_catalog"."default",
  "bool_test" bool,
  "float_test" float4,
  "time_test" timestamp(0),
  CONSTRAINT "tb_test_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "db_monitor"."tb_test"
  OWNER TO "postgres";

*/

var dbpool *DBPool

func TestSqlWrap(t *testing.T) {
	dbpool, err := DBPoolInit(&DBConfig{
		Username:     "postgres",
		Password:     "gosuncn20",
		Host:         "127.0.0.1",
		Port:         5432,
		Dbname:       "db_monitor",
		MaxOpenConns: 100,
		MaxIdleConns: 20,
		DbType:       PostgreSQL,
		Timeout:      "10s",
	})
	if err != nil {
		t.Error(err)
	}

	sql := fmt.Sprintf("select int_test, string_test, bool_test, float_test, time_test from db_monitor.tb_test limit 1 offset 0 ")
	rows, err := dbpool.GetDB().Raw(sql).Rows()
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	rowsWrap, err := RowsWrapScan(rows)

	var sql2 string
	for _, v := range rowsWrap {
		sql := "insert into db_monitor.tb_test("
		valueSql := "values("
		for i, col := range v.Columns {
			sql += col.Name
			switch col.ScanType.Kind() {
			case reflect.String:
				if !col.Valid {
					valueSql += "NULL"
				} else {
					valueSql += fmt.Sprintf("'%s'", col.StrValue)
				}
			case reflect.Struct:
				t, _ := time.Parse(time.RFC3339, col.StrValue)
				str := t.In(time.Local).Format("2006-01-02 15:04:05")
				if !col.Valid {
					valueSql += "NULL"
				} else {
					valueSql += fmt.Sprintf("'%s'", str)
				}
			default:
				if !col.Valid {
					valueSql += "NULL"
				} else {
					valueSql += fmt.Sprintf("%v", col.StrValue)
				}
			}
			if i == len(v.Columns)-1 {
				sql += ")"
				valueSql += ")"
			} else {
				sql += ","
				valueSql += ","
			}
		}
		sql2 = sql + " " + valueSql
	}

	t.Log(sql2)
	dbpool.GetDB().Exec(sql2)
}
