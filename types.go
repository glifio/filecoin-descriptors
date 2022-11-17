package main

import (
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/iancoleman/orderedmap"
)

type ActorName = string
type ActorCode = string
type PropName = string

type ActorCodeMap = map[ActorName]ActorCode

type NetworkActorCodeMap = map[dtypes.NetworkName]ActorCodeMap

const (
	TypeBool      = "boolean"
	TypeNumber    = "number"
	TypeString    = "string"
	TypeBytes     = "bytes"
	TypeMap       = "map"
	TypeArray     = "array"
	TypeChan      = "channel"
	TypeObject    = "object"
	TypeFunction  = "function"
	TypeInterface = "interface"
)

type DataType struct {
	Type       string
	Name       string
	Key        *DataType   `json:",omitempty"` // For map type
	Contains   *DataType   `json:",omitempty"` // For map / array / channel type
	Children   DataTypeMap `json:",omitempty"` // For object type
	Methods    DataTypeMap `json:",omitempty"` // For interface type
	Params     []DataType  `json:",omitempty"` // For function type
	Returns    []DataType  `json:",omitempty"` // For function type
	IsVariadic bool        `json:",omitempty"` // For function type
	ChanDir    string      `json:",omitempty"` // For channel type
}

type DataTypeMap = *orderedmap.OrderedMap

type ActorMethod struct {
	Name   string
	Param  DataType
	Return DataType
}

type ActorMethodMap = map[uint64]ActorMethod

type ActorDescriptor struct {
	State   DataTypeMap
	Methods ActorMethodMap
}

type ActorDescriptorMap = map[ActorName]ActorDescriptor
