// FIT_UserProfiles project FIT_UserProfiles.go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//================================================
//Define structures for implementing the functions
//================================================
type SimpleChaincode struct{}
type UserProfile struct {
	ObjectType      string `json:"objectType"`
	ProfileID       string `json:"profileID"`
	UserPreferences string `json:"userPreferences"`
	Event           string `json:"event"`
}

//=============
//Main function
//=============
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting the chaincode " + err.Error())
	}
}

//=======================
//Initialization function
//=======================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//===================
//Invocation function
//===================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createProfile" {
		return t.createProfile(stub, args)
	} else if function == "queryProfile" {
		return t.queryProfile(stub, args)
	} else if function == "updateProfile" {
		return t.updateProfile(stub, args)
	} else {
		return shim.Error("Invalid function name " + function)
	}
	return shim.Success(nil)
}

//===============================================
//createProfile - Function to create user profile
//===============================================
func (t *SimpleChaincode) createProfile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Check if the profile ID is not null
	if args[0] == "" {
		return shim.Error("Unique profile identifier cannot be null")
	}

	//Assign the arguments to variables
	objectType := "User Profile"
	profileID := args[0]
	userPreferences := args[1]
	event := "User profile created and preferences added"

	//Check if the profile ID already exists
	userBytes, err := stub.GetState(profileID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if userBytes != nil {
		return shim.Error("Profile with ID " + profileID + " already exists in the system")
	}
	//Create the JSON object interface
	profileObj := &UserProfile{objectType, profileID, userPreferences, event}

	//Convert the JSON object to bytes
	profileBytes, err := json.Marshal(profileObj)
	if err != nil {
		return shim.Error("Error 2 " + err.Error())
	}

	//Create a write set with profile ID as key
	err = stub.PutState(profileID, profileBytes)
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
	}

	//Publish the event of profile creation
	err = stub.SetEvent(event, profileBytes)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}

	return shim.Success(profileBytes)
}

//=========================================================================
//queryProfile - Function to query for the latest state of the user profile
//=========================================================================
func (t *SimpleChaincode) queryProfile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Check if the profile ID is not null
	if args[0] == "" {
		return shim.Error("Unique profile identifier cannot be null")
	}
	profileID := args[0]
	//Check if the profile ID exists and get the latest state
	profileBytes, err := stub.GetState(profileID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if profileBytes == nil {
		return shim.Error("Profile with ID " + profileID + " does not exist in the system")
	}

	return shim.Success(profileBytes)
}

//===============================================
//updateProfile - Function to create user profile
//===============================================
func (t *SimpleChaincode) updateProfile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Check if the profile ID is not null
	if args[0] == "" {
		return shim.Error("Unique profile identifier cannot be null")
	}

	//Assign the arguments to variables
	objectType := "User Profile"
	profileID := args[0]
	userPreferences := args[1]
	event := "User Preferences updated"

	//Check if the profile ID already exists
	userBytes, err := stub.GetState(profileID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if userBytes == nil {
		return shim.Error("Profile with ID " + profileID + " does not exist in the system")
	}
	//Create the JSON object interface
	profileObj := &UserProfile{objectType, profileID, userPreferences, event}

	//Convert the JSON object to bytes
	profileBytes, err := json.Marshal(profileObj)
	if err != nil {
		return shim.Error("Error 2 " + err.Error())
	}

	//Create a write set with profile ID as key
	err = stub.PutState(profileID, profileBytes)
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
	}

	//Publish the event of profile update
	err = stub.SetEvent(event, profileBytes)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}

	return shim.Success(profileBytes)
}
