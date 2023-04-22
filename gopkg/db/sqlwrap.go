package db

import (
	"database/sql"
	"reflect"
)

type ColumnWrap struct {
	Name     string
	ScanType reflect.Type
	StrValue string
	Valid    bool // Valid is true if StrValue is not NULL
}

type RowWrap struct {
	Columns []ColumnWrap
}

func RowsWrapScan(r *sql.Rows) ([]*RowWrap, error) {
	rows := make([]*RowWrap, 0)
	for r.Next() {
		rowwrap := new(RowWrap)
		columnTypes, err := r.ColumnTypes()
		if err != nil {
			return nil, err
		}
		dest, err := SliceScan(r)
		if err != nil {
			return nil, err
		}

		for i, v := range dest {
			var col ColumnWrap
			col.Name = columnTypes[i].Name()
			col.ScanType = columnTypes[i].ScanType()

			var obj sql.NullString
			if err := obj.Scan(v); err == nil {
				col.StrValue = obj.String
				col.Valid = obj.Valid
			}
			rowwrap.Columns = append(rowwrap.Columns, col)
		}
		rows = append(rows, rowwrap)
	}
	return rows, nil
}

func SliceScan(r *sql.Rows) ([]interface{}, error) {
	columns, err := r.Columns()
	if err != nil {
		return []interface{}{}, err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)

	if err != nil {
		return values, err
	}

	for i := range columns {
		values[i] = *(values[i].(*interface{}))
	}

	return values, r.Err()
}

func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}
