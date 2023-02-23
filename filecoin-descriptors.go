package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"reflect"

	"github.com/filecoin-project/go-state-types/abi"
)

var apiUrls = []string{
	"https://api.node.glif.io/rpc/v1",
	"https://api.calibration.node.glif.io/rpc/v1",
	"https://api.hyperspace.node.glif.io/rpc/v1",
}

func main() {
	/*
	 * Preparation
	 */

	// Create output directory if not exists
	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	/*
	 * Actor codes
	 */

	var networkActorCodeMap = NetworkActorCodeMap{}

	for _, url := range apiUrls {

		// Open Lotus API for network
		var lotus Lotus
		if err := lotus.Open(url); err != nil {
			log.Fatalf("Failed to start Lotus API: %s", err)
		}
		defer lotus.Close()

		// Retrieve network name from Lotus
		networkName, err := lotus.api.StateNetworkName(context.Background())
		if err != nil {
			log.Fatalf("Failed to get network name: %s", err)
		}

		// Retrieve actor codes from Lotus
		actorCodeMap, err := lotus.GetActorCodeMap()
		if err != nil {
			log.Fatalf("Failed to get actor codes: %v", err)
		}

		// Store actor codes in map
		networkActorCodeMap[networkName] = actorCodeMap
	}

	// Write actor codes
	if err := writeJsonFile(networkActorCodeMap, "actor-codes"); err != nil {
		log.Fatalf("Failed to write actor codes to JSON file: %v", err)
	}

	/*
	 * Actor descriptors
	 */

	var actorDescriptorMap = ActorDescriptorMap{}
	for name, reflectableActor := range reflectableActors {

		// State reflection
		var actorState DataTypeMap = nil
		if stateType := reflect.TypeOf(reflectableActor.State); stateType != nil {
			stateDataType := GetDataType(stateType)
			if stateDataType.Type != TypeObject {
				log.Fatalf("%s actor state is not an object", name)
			}
			actorState = stateDataType.Children
		}

		// Methods reflection
		var actorMethodMap = ActorMethodMap{}

		// Add Send method
		if name != "system" {
			emptyType := reflect.TypeOf((*abi.EmptyValue)(nil))
			emptyDataType := GetDataType(emptyType)
			actorMethodMap[0] = ActorMethod{
				Name:   "Send",
				Param:  emptyDataType,
				Return: emptyDataType,
			}
		}

		// Iterate over actor methods
		for key, method := range reflectableActor.Methods {
			var actorMethod ActorMethod
			methodType := reflect.TypeOf(method)

			if methodType.Name() == "CustomMethod" {
				var customMethod = method.(CustomMethod)
				actorMethod.Name = customMethod.Name
				actorMethod.Param = GetDataType(reflect.TypeOf(customMethod.Param))
				actorMethod.Return = GetDataType(reflect.TypeOf(customMethod.Return))
			} else {
				// Get method DataType
				methodDataType := GetDataType(methodType)

				// Set method name
				fullName := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
				nameParts := strings.Split(fullName, ".")
				actorMethod.Name = nameParts[len(nameParts)-1]

				// Set method parameter
				paramsCount := len(methodDataType.Params)
				if paramsCount != 3 {
					log.Fatalf("%s actor method %s has %d parameters, expected 3", name, method, paramsCount)
				}
				firstParamName := methodDataType.Params[0].Name
				if firstParamName != "Actor" {
					log.Fatalf("%s actor method %s has %s as first parameter, should be Actor", name, method, firstParamName)
				}
				secondParamName := methodDataType.Params[1].Name
				if secondParamName != "Runtime" {
					log.Fatalf("%s actor method %s has %s as first parameter, should be Runtime", name, method, secondParamName)
				}
				actorMethod.Param = methodDataType.Params[2]

				// Set method return value
				returnsCount := len(methodDataType.Returns)
				if returnsCount != 1 {
					log.Fatalf("%s actor method %s has %d return values, expected 1", name, method, returnsCount)
				}
				actorMethod.Return = methodDataType.Returns[0]
			}

			// Store method in map
			actorMethodMap[key] = actorMethod
		}

		// Set actor descriptor
		actorDescriptorMap[name] = ActorDescriptor{
			State:   actorState,
			Methods: actorMethodMap,
		}
	}

	// Write actor descriptors to JSON file
	if err := writeJsonFile(actorDescriptorMap, "actor-descriptors"); err != nil {
		log.Fatalf("Failed to write actor descriptors to JSON file: %v", err)
	}

	/*
	 * Done
	 */

	fmt.Println("Done")
}

func writeJsonFile(data interface{}, filename string) error {

	// Marshal data to JSON
	dataJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Create file
	f, err := os.Create(fmt.Sprintf("output/%s.json", filename))
	if err != nil {
		return err
	}
	defer f.Close()

	// Write file
	if _, err = f.Write(dataJson); err != nil {
		return err
	}

	return nil
}
