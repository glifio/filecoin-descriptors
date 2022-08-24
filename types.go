package main

type NetworkCode = string
type Network struct {
	Code NetworkCode
	Url  string
}

type ActorName = string
type ActorCode = string
type ActorCodeMap = map[ActorName]ActorCode

type PropName = string
type DataTypes = map[PropName]DataType
type DataType struct {
	Name       string
	Type       string
	Key        *DataType  `json:",omitempty"` // For map type
	Contains   *DataType  `json:",omitempty"` // For map / array type
	Children   DataTypes  `json:",omitempty"` // For object type
	Params     []DataType `json:",omitempty"` // For function type
	Returns    []DataType `json:",omitempty"` // For function type
	IsVariadic bool       `json:",omitempty"` // For function type
}

const (
	DataTypeBool      = "bool"
	DataTypeNumber    = "number"
	DataTypeString    = "string"
	DataTypeMap       = "map"
	DataTypeArray     = "array"
	DataTypeObject    = "object"
	DataTypeFunction  = "function"
	DataTypeInterface = "interface"
)

type ActorMethods = map[string]ActorMethod
type ActorMethod = struct {
	Name   string
	Param  DataType
	Return DataType
}

type ActorDescriptors = map[ActorName]ActorDescriptor
type ActorDescriptor struct {
	State   DataTypes
	Methods ActorMethods
}
