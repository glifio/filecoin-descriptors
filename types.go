package main

type ActorName = string
type ActorCode = string
type ActorCodes = map[ActorName]ActorCode

type PropName = string
type DataTypes = map[PropName]DataType
type DataType struct {
	Type       string
	Name       string
	Key        *DataType  `json:",omitempty"` // For map type
	Contains   *DataType  `json:",omitempty"` // For map / array / channel type
	Children   DataTypes  `json:",omitempty"` // For object type
	Methods    DataTypes  `json:",omitempty"` // For interface type
	Params     []DataType `json:",omitempty"` // For function type
	Returns    []DataType `json:",omitempty"` // For function type
	IsVariadic bool       `json:",omitempty"` // For function type
	ChanDir    string     `json:",omitempty"` // For channel type
}

const (
	DataTypeBool      = "bool"
	DataTypeNumber    = "number"
	DataTypeString    = "string"
	DataTypeMap       = "map"
	DataTypeArray     = "array"
	DataTypeChan      = "channel"
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
