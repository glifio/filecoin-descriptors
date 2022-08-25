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

	"github.com/filecoin-project/lotus/node/modules/dtypes"
)

var apiUrls = []string{
	"https://api.node.glif.io",
	"https://api.calibration.node.glif.io",
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

	var networkActorCodes = map[dtypes.NetworkName]ActorCodes{}

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
		actorCodes, err := lotus.GetActorCodes()
		if err != nil {
			log.Fatalf("Failed to get actor codes: %v", err)
		}

		// Store actor codes in map
		networkActorCodes[networkName] = actorCodes
	}

	// Write actor codes
	if err := writeJsonFile(networkActorCodes, "actor-codes"); err != nil {
		log.Fatalf("Failed to write actor codes to JSON file: %v", err)
	}

	/*
	 * Actor descriptors
	 */

	var actorDescriptors = ActorDescriptors{}
	for name, reflectableActor := range reflectableActors {

		// State reflection
		stateType := reflect.TypeOf(reflectableActor.State)
		stateDataType := GetDataType(stateType)
		if stateDataType.Type != TypeObject {
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
