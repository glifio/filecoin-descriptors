package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

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
	 * Preparation
	 */

	// Create output directory if not exists
	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

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
	 * Actor descriptors
	 */

	var actorDescriptors = ActorDescriptors{}
	for name, reflectableActor := range reflectableActors {

		// State reflection
		stateType := reflect.TypeOf(reflectableActor.State)
		stateDataType := GetDataType(stateType)
		if stateDataType.Type != DataTypeObject {
			log.Fatalf("%s actor state is not an object", name)
		}

		// Methods reflection
		var actorMethods = ActorMethods{}
		for key, method := range reflectableActor.Methods {
			var actorMethod ActorMethod
			methodType := reflect.TypeOf(method)
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

			// Store method in map
			actorMethods[key] = actorMethod
		}

		// Set actor descriptor
		actorDescriptors[name] = ActorDescriptor{
			State:   stateDataType.Children,
			Methods: actorMethods,
		}
	}

	// Write actor descriptors to JSON file
	if err := writeJsonFile(actorDescriptors, "actor-descriptors"); err != nil {
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
