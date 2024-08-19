package tools

import (
	"log"
	"reflect"
)

func ExtendedPrint(v interface{}) {
	val := reflect.ValueOf(v)
	//  проверяем, а не передали ли нам указатель на структуру
	switch val.Kind() {
	case reflect.Ptr:
		if val.Elem().Kind() != reflect.Struct {
			log.Printf("Pointer to %v : %v", val.Elem().Type(), val.Elem())
			return
		}
		// если всё-таки это указатель на структуру, дальше будем работать с самой структурой
		val = val.Elem()
	case reflect.Struct: // работаем со структурой
	default:
		log.Printf("%v : %v", val.Type(), val)
		return
	}
	log.Printf("Struct of type %v and number of fields %d:\n", val.Type(), val.NumField())
	for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
		field := val.Field(fieldIndex) // field — тоже Value
		log.Printf("\tField %v: %v - val :%v\n", val.Type().Field(fieldIndex).Name, field.Type(), field)
		// имя поля мы получаем не из значения поля, а из его типа.
	}
}

func PadRight(s *string, pad string, ln int64) {
	if len(*s) > int(ln) {
		return
	}
	for len(*s) < int(ln) {
		*s = *s + pad
	}
}
