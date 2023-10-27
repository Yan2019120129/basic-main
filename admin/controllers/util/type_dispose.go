package util

import "reflect"

func GetFieldAddr(a interface{}) (values []interface{}) {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		structValue := v.Elem()
		fieldCount := structValue.NumField()
		values = make([]interface{}, fieldCount)
		for i := 0; i < fieldCount; i++ {
			values[i] = structValue.Field(i).Addr().Interface()
		}
	}
	return values
}
