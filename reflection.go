package main

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/ipfs/go-cid"
)

// Special data types
var addressType = reflect.TypeOf((*address.Address)(nil)).Elem()
var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
var bitFieldType = reflect.TypeOf((*bitfield.BitField)(nil)).Elem()
var cidType = reflect.TypeOf((*cid.Cid)(nil)).Elem()

func GetDataType(t reflect.Type) DataType {
	var dataType DataType
	dataType.Name = t.Name()

	// Handle special types
	switch t.String() {

	case addressType.String():
		dataType.Type = TypeString
		return dataType

	case bigIntType.String():
		dataType.Name = "FilecoinNumber"
		dataType.Type = TypeString
		return dataType

	case bitFieldType.String():
		containsType := DataType{Name: "Bit", Type: TypeNumber}
		dataType.Type = TypeArray
		dataType.Contains = &containsType
		return dataType

	case cidType.String():
		dataType.Type = TypeObject
		dataType.Children = DataTypeMap{}
		dataType.Children["/"] = DataType{Name: "CidString", Type: TypeString}
		return dataType
	}

	// Handle base types
	switch t.Kind() {

	case reflect.Ptr:
		return GetDataType(t.Elem())

	case reflect.Bool:
		dataType.Type = TypeBool
		return dataType

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		dataType.Type = TypeNumber
		return dataType

	case reflect.String:
		dataType.Type = TypeString
		return dataType

	case reflect.Chan:
		containsType := GetDataType(t.Elem())
		dataType.Type = TypeChan
		dataType.ChanDir = t.ChanDir().String()
		dataType.Contains = &containsType
		return dataType

	case reflect.Map:
		keyType := GetDataType(t.Key())
		containsType := GetDataType(t.Elem())
		dataType.Type = TypeMap
		dataType.Key = &keyType
		dataType.Contains = &containsType
		return dataType

	case reflect.Array, reflect.Slice:
		containsType := GetDataType(t.Elem())
		dataType.Name = "[]" + containsType.Name
		dataType.Type = TypeArray
		dataType.Contains = &containsType
		return dataType

	case reflect.Struct:
		dataType.Type = TypeObject
		dataType.Children = DataTypeMap{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			dataType.Children[f.Name] = GetDataType(f.Type)
		}
		return dataType

	case reflect.Func:
		dataType.Type = TypeFunction
		dataType.IsVariadic = t.IsVariadic()
		for i := 0; i < t.NumIn(); i++ {
			dataType.Params = append(dataType.Params, GetDataType(t.In(i)))
		}
		for i := 0; i < t.NumOut(); i++ {
			dataType.Returns = append(dataType.Returns, GetDataType(t.Out(i)))
		}
		return dataType

	case reflect.Interface:
		dataType.Type = TypeInterface
		dataType.Methods = DataTypeMap{}
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			dataType.Methods[m.Name] = GetDataType(m.Type)
		}
		return dataType
	}

	// Unhandled type
	panic(fmt.Sprintf("Unhandled type with string: %s, name: %s, kind: %s", t.String(), t.Name(), t.Kind().String()))
}
