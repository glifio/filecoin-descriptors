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
type DataTypeMap = map[PropName]DataType
type DataType struct {
	Name       string
	Type       string
	Contains   *DataType   `json:",omitempty"` // For array type
	Children   DataTypeMap `json:",omitempty"` // For object type
	Params     []DataType  `json:",omitempty"` // For function type
	Returns    []DataType  `json:",omitempty"` // For function type
	IsVariadic bool        `json:",omitempty"` // For function type
}

const (
	DataTypeBool      = "bool"
	DataTypeNumber    = "number"
	DataTypeString    = "string"
	DataTypeArray     = "array"
	DataTypeObject    = "object"
	DataTypeFunction  = "function"
	DataTypeInterface = "interface"
)

type ActorMethodMap = map[string]ActorMethod
type ActorMethod = struct {
	Params DataTypeMap
	Return DataType
}

type ActorDescriptorMap = map[ActorName]ActorDescriptor
type ActorDescriptor struct {
	State   DataTypeMap
	Methods ActorMethodMap
}
