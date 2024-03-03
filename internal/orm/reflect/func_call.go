package reflect

import "reflect"

func IterateFunc(entity any) map[string]FuncInfo {

}
type FuncInfo struct {
	Name 			string
	InputTypes 		[]reflect.Type
	OutputTypes  	[]reflect.Type
	Result          []any
}