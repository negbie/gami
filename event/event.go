//Package event decoder
//This Build Type of Event received
package event

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/negbie/gami"
)

// eventTrap used internal for trap events and cast
var eventTrap = make(map[string]interface{})

//New build a new event Type if not return the AMIEvent
func New(event *gami.AMIEvent) interface{} {
	if intf, ok := eventTrap[event.ID]; ok {
		return build(event, &intf)
	}
	//return *event
	return nil
}

func build(event *gami.AMIEvent, klass *interface{}) interface{} {

	typ := reflect.TypeOf(*klass)
	value := reflect.ValueOf(*klass)
	ret := reflect.New(typ).Elem()
	//fmt.Printf("EVENT=%+v\n", *event)
	for ix := 0; ix < value.NumField(); ix++ {
		field := ret.Field(ix)
		tfield := typ.Field(ix)

		if tfield.Name == "Privilege" {
			field.Set(reflect.ValueOf(event.Privilege))
			continue
		}
		AMIField := tfield.Tag.Get("AMI")
		switch field.Kind() {
		case reflect.String:
			field.SetString(event.Params[AMIField])
		case reflect.Int64, reflect.Int:
			vint, _ := strconv.Atoi(event.Params[AMIField])
			field.SetInt(int64(vint))
		case reflect.Float64:
			if vfloat, err := strconv.ParseFloat(event.Params[AMIField], 64); err == nil {
				field.SetFloat(vfloat)
			}
		default:
			fmt.Print(ix, tfield.Tag.Get("AMI"), ":", field, "\n")
		}

	}
	// fmt.Printf("INTERFACE=%+v\n", ret.Interface())
	return ret.Interface()
}
