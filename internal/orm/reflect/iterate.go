package reflect

import "reflect"

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}


func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	resKeys := make([]any, 0, val.Len())
	resValue := make([]any, 0, val.Len())

	keys := val.MapKeys()

	for _, key := range keys {
		v := val.MapIndex(key)
		resKeys = append(resKeys, key.Interface())
		resValue = append(resValue, v.Interface())
	}

	return resKeys, resValue, nil
}



func IterateMapRange(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	resKeys := make([]any, 0, val.Len())
	resValue := make([]any, 0, val.Len())

	itr := val.MapRange()
	for itr.Next() {
		resKeys = append(resKeys, itr.Key().Interface())
		resValue = append(resValue, itr.Value().Interface())
	}
	return resKeys, resValue, nil
}