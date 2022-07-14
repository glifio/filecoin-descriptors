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
	var dataType = DataType{}
	dataType.Name = t.Name()

	// Handle special types
	switch t.String() {

	case addressType.String():
		dataType.Type = DataTypeString
		return dataType

	case bigIntType.String():
		dataType.Name = "FilecoinNumber"
		dataType.Type = DataTypeString
		return dataType

	case bitFieldType.String():
		containsType := DataType{Name: "Bit", Type: DataTypeNumber}
		dataType.Type = DataTypeArray
		dataType.Contains = &containsType
		return dataType

	case cidType.String():
		dataType.Type = DataTypeObject
		dataType.Children = DataTypeMap{}
		dataType.Children["/"] = DataType{Name: "CidString", Type: DataTypeString}
		return dataType
	}

	// Handle base types
	switch t.Kind() {

	case reflect.Ptr:
		return GetDataType(t.Elem())

	case reflect.Bool:
		dataType.Type = DataTypeBool
		return dataType

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		dataType.Type = DataTypeNumber
		return dataType

	case reflect.String:
		dataType.Type = DataTypeString
		return dataType

	case reflect.Array, reflect.Slice:
		containsType := GetDataType(t.Elem())
		dataType.Name = "[]" + containsType.Name
		dataType.Type = DataTypeArray
		dataType.Contains = &containsType
		return dataType

	case reflect.Struct:
		dataType.Type = DataTypeObject
		dataType.Children = DataTypeMap{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			dataType.Children[f.Name] = GetDataType(f.Type)
		}
		return dataType
	}

	panic(fmt.Sprintf("Unhandled type: %s", t.String()))
}
