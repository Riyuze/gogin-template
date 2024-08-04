package helper

import (
	"reflect"
	"strconv"
	"strings"
)

func RepoPGGetColumns(typ reflect.Type) (key string, insert string, update string, field string, order string) {
	key, insert, update, field, order = "", "", "", "", ""
	for i := 0; i < typ.NumField(); i++ {
		objectField := typ.Field(i)
		if tagValue, ok := objectField.Tag.Lookup("db"); ok {
			insert += tagValue + ", "
			field += "obj." + objectField.Name + ", "
			if dbxValue, ok2 := objectField.Tag.Lookup("dbx"); ok2 {
				if strings.Contains(dbxValue, "key") {
					key += tagValue + ", "
				} else {
					update += tagValue + ", "
				}
				if strings.Contains(dbxValue, "sort") {
					order += tagValue + ", "
				}
			} else {
				update += tagValue + ", "
			}
		}
	}
	key = strings.TrimSuffix(key, ", ")
	insert = strings.TrimSuffix(insert, ", ")
	update = strings.TrimSuffix(update, ", ")
	field = strings.TrimSuffix(field, ", ")
	order = strings.TrimSuffix(order, ", ")
	return key, insert, update, field, order
}

func RepoPGGetInputs(typ reflect.Type) (key string, insert string, update string, field string, where string, whereF string) {
	key, insert, update, field, where, whereF = "", "", "", "", "", ""
	index := 0
	for i := 0; i < typ.NumField(); i++ {
		objectField := typ.Field(i)
		if tagValue, ok := objectField.Tag.Lookup("db"); ok {
			index++
			insert += "$" + strconv.Itoa(index) + ", "
			field += "obj." + objectField.Name + ", "
			if dbxValue, ok2 := objectField.Tag.Lookup("dbx"); ok2 {
				if strings.Contains(dbxValue, "key") {
					key += "$" + strconv.Itoa(index) + ", "
					where += tagValue + " = $" + strconv.Itoa(index) + " AND "
				} else {
					update += "$" + strconv.Itoa(index) + ", "
				}
				if strings.Contains(dbxValue, "foreign") {
					whereF += tagValue + " = $" + strconv.Itoa(index) + " AND "
				}
			} else {
				update += "$" + strconv.Itoa(index) + ", "
			}
		}
	}
	key = strings.TrimSuffix(key, ", ")
	insert = strings.TrimSuffix(insert, ", ")
	update = strings.TrimSuffix(update, ", ")
	field = strings.TrimSuffix(field, ", ")
	where = strings.TrimSuffix(where, " AND ")
	whereF = strings.TrimSuffix(whereF, " AND ")
	return key, insert, update, field, where, whereF
}

func RepoPGGetSelect(typ reflect.Type, schema string, tableName string) string {
	_, insertN, _, _, _ := RepoPGGetColumns(typ)
	_, _, _, _, whereI, _ := RepoPGGetInputs(typ)
	query := `SELECT ` + insertN + ` FROM ` + schema + `.` + tableName + ` WHERE ` + whereI
	return query
}

func RepoPGGetSelectForeign(typ reflect.Type, schema string, tableName string) string {
	_, insertN, _, _, _ := RepoPGGetColumns(typ)
	_, _, _, _, _, whereFI := RepoPGGetInputs(typ)
	query := `SELECT ` + insertN + ` FROM ` + schema + `.` + tableName + ` WHERE ` + whereFI
	return query
}

func RepoPGGetInsert(typ reflect.Type, schema string, tableName string) string {
	_, insertN, _, _, _ := RepoPGGetColumns(typ)
	_, insertI, _, _, _, _ := RepoPGGetInputs(typ)
	query := `INSERT INTO ` + schema + `.` + tableName + ` (` + insertN + `) VALUES (` + insertI + `) `
	return query
}

func RepoPGGetUpsert(typ reflect.Type, schema string, tableName string) string {
	keyN, insertN, updateN, _, _ := RepoPGGetColumns(typ)
	_, insertI, updateI, _, _, _ := RepoPGGetInputs(typ)
	query := `INSERT INTO ` + schema + `.` + tableName + ` (` + insertN + `) VALUES (` + insertI + `) ON CONFLICT (` + keyN + `) DO UPDATE SET (` + updateN + `) = (` + updateI + `) `
	return query
}

func RepoPGGetDelete(typ reflect.Type, schema string, tableName string) string {
	_, _, _, _, whereI, _ := RepoPGGetInputs(typ)
	query := `DELETE FROM ` + schema + `.` + tableName + ` WHERE ` + whereI
	return query
}

func RepoPGGetDeleteForeign(typ reflect.Type, schema string, tableName string) string {
	_, _, _, _, _, whereFI := RepoPGGetInputs(typ)
	query := `DELETE FROM ` + schema + `.` + tableName + ` WHERE ` + whereFI
	return query
}

func RepoPGGetTypeArgValue[T any](obj T) []interface{} {
	v := reflect.ValueOf(obj)
	index := 0
	for i := 0; i < v.NumField(); i++ {
		objectField := v.Type().Field(i)
		if _, ok := objectField.Tag.Lookup("db"); ok {
			index++
		}
	}

	values := make([]interface{}, index)

	indexInsert := 0
	for i := 0; i < v.NumField(); i++ {
		objectField := v.Type().Field(i)
		if _, ok := objectField.Tag.Lookup("db"); ok {
			values[indexInsert] = v.Field(i).Interface()
			indexInsert++
		}
	}
	return values
}

func GetLimitAndOffset(pageSize int, page int) (*int, *int) {
	var limit, offset *int

	if pageSize > 0 && page >= 0 {
		limit = &pageSize
		offsetCalc := pageSize * page
		offset = &offsetCalc
	}

	return limit, offset
}
