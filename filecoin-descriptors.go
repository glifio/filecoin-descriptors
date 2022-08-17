package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"reflect"
)

var networks = []Network{
	{
		Code: "t",
		Url:  "https://api.calibration.node.glif.io",
	},
	{
		Code: "f",
		Url:  "https://api.node.glif.io",
	},
}

func main() {
	/*
	 * Get actor codes
	 */

	for _, network := range networks {

		// Open Lotus API for network
		var lotus Lotus
		if err := lotus.Open(network); err != nil {
			log.Fatalf("Failed to start Lotus API: %s", err)
		}
		defer lotus.Close()

		// Retrieve actor codes from Lotus
		actorCodeMap, err := lotus.GetActorCodeMap()
		if err != nil {
			log.Fatalf("Failed to get actor codes: %v", err)
		}

		// Print actor codes
		fmt.Printf("Network: %s\n", network.Code)
		for name, code := range actorCodeMap {
			fmt.Printf("Actor: %s, Code: %s\n", name, code)
		}
		fmt.Print("\n")
	}

	/*
	 * Get actor descriptors
	 */

	var actorDescriptorMap = ActorDescriptorMap{}
	for name, reflectableActor := range reflectableActors {

		// State reflection
		stateType := reflect.TypeOf(reflectableActor.State)
		stateDataType := GetDataType(stateType)
		if stateDataType.Type != DataTypeObject {
			log.Fatalf("%s actor state is not an object", name)
		}

		// Methods reflection
		var actorMethodMap = ActorMethodMap{}
		for key, method := range reflectableActor.Methods {
			var actorMethod ActorMethod
			methodType := reflect.TypeOf(method)
			methodDataType := GetDataType(methodType)

			// Get method parameter
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

			// Get method return value
			returnsCount := len(methodDataType.Returns)
			if returnsCount != 1 {
				log.Fatalf("%s actor method %s has %d return values, expected 1", name, method, returnsCount)
			}
			actorMethod.Return = methodDataType.Returns[0]

			// Store method in map
			actorMethodMap[key] = actorMethod
		}

		// Set actor descriptor
		actorDescriptorMap[name] = ActorDescriptor{
			State:   stateDataType.Children,
			Methods: actorMethodMap,
		}
	}

	// Marshal actor descriptors to JSON
	actorDescriptorJson, err := json.MarshalIndent(actorDescriptorMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal actor descriptors to JSON: %v", err)
	}

	// Create output directory
	if err = os.MkdirAll("output", os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create file for actor descriptors
	f, err := os.Create("output/actor-descriptors.json")
	if err != nil {
		log.Fatalf("Failed to create actor descriptors file: %v", err)
	}
	defer f.Close()

	// Write actor descriptors to file
	if _, err = f.Write(actorDescriptorJson); err != nil {
		log.Fatalf("Failed to write actor descriptors to file: %v", err)
	}

	fmt.Println("Successfully generated actor descriptors file")
}
