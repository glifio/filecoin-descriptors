package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ipfs/go-hamt-ipld"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func MapToInterface(input interface{}, output interface{}) error {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonInput, output)
	if err != nil {
		return err
	}

	return nil
}

func DecodeNodeCBOR(data []byte) (datamodel.Node, error) {
	np := basicnode.Prototype.Any
	nb := np.NewBuilder()
	if err := dagcbor.Decode(nb, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return nb.Build(), nil
}

func DecodeHamtNodeCBOR(data []byte) (*hamt.Node, error) {
	var node hamt.Node
	if err := node.UnmarshalCBOR(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return &node, nil
}

func DecodeHamtKvCBOR(data []byte) (*hamt.KV, error) {
	var kv hamt.KV
	if err := kv.UnmarshalCBOR(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return &kv, nil
}

func PrintNode(node datamodel.Node, name string, indent int) {
	switch node.Kind() {

	case datamodel.Kind_Null:
		fmt.Printf("%s%s: (null)\n", getIndentation(indent), name)

	case datamodel.Kind_Bool:
		val, _ := node.AsBool()
		fmt.Printf("%s%s: (bool) %t\n", getIndentation(indent), name, val)

	case datamodel.Kind_Int:
		val, _ := node.AsInt()
		fmt.Printf("%s%s: (int) %d\n", getIndentation(indent), name, val)

	case datamodel.Kind_Float:
		val, _ := node.AsFloat()
		fmt.Printf("%s%s: (float) %f\n", getIndentation(indent), name, val)

	case datamodel.Kind_String:
		val, _ := node.AsString()
		fmt.Printf("%s%s: (string) %s\n", getIndentation(indent), name, val)

	case datamodel.Kind_Bytes:
		val, _ := node.AsBytes()
		fmt.Printf("%s%s: (bytes) %x\n", getIndentation(indent), name, val)

	case datamodel.Kind_Link:
		val, _ := node.AsLink()
		fmt.Printf("%s%s: (link) %s\n", getIndentation(indent), name, val.String())

	case datamodel.Kind_Map:
		fmt.Printf("%s%s: (map)\n", getIndentation(indent), name)
		iter := node.MapIterator()
		for {
			if iter.Done() {
				break
			}
			fmt.Printf("%s(map entry)\n", getIndentation(indent+1))
			key, val, _ := iter.Next()
			PrintNode(key, "key", indent+2)
			PrintNode(val, "value", indent+2)
		}

	case datamodel.Kind_List:
		fmt.Printf("%s%s: (list)\n", getIndentation(indent), name)
		iter := node.ListIterator()
		for {
			if iter.Done() {
				break
			}
			idx, val, _ := iter.Next()
			PrintNode(val, strconv.FormatInt(idx, 10), indent+1)
		}

	case datamodel.Kind_Invalid:
		panic(fmt.Sprintf("Node '%s' has invalid kind", name))

	default:
		panic(fmt.Sprintf("Node '%s' has unhandled kind: %s", name, node.Kind().String()))
	}
}

func getIndentation(indent int) string {
	indents := make([]string, indent+1)
	return strings.Join(indents, "  ")
}
