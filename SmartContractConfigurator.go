// SmartContractConfigurator project SmartContractConfigurator.go
package main

import (
	//"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct{}
type Term struct {
	TermName  string `json:"termName"`
	TermValue string `json:"termValue"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Cannot start chaincode " + err.Error())
	}
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functions, args := stub.GetFunctionAndParameters()
	if functions == "createUpdateTerms" {
		return t.createUpdateTerms(stub, args)
	} else if functions == "queryTerms" {
		return t.queryTerms(stub, args)
	} else {
		return shim.Error("Incorrect function name " + functions)
	}
}

func (t *SimpleChaincode) createUpdateTerms(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	termName := args[0]
	termValue := args[1]
	if len(termName) == 0 {
		return shim.Error("Term Name cannot be null")
	}
	if (termName != "Penalty Percentage") && (termName != "Discount Percentage") && (termName != "Updated Payment Terms") {
		return shim.Error("Term Name is not in the list of values")
	}
	if len(termValue) == 0 {
		return shim.Error("Term Value cannot be null")
	}

	err := stub.PutState(termName, []byte(termValue))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Term " + termName + " updated with value " + termValue))

}

func (t *SimpleChaincode) queryTerms(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	termName := args[0]
	if len(termName) == 0 {
		return shim.Error("Term Name cannot be null")
	}
	termValBytes, err := stub.GetState(termName)
	if err != nil {
		return shim.Error(err.Error())
	}
	if termValBytes == nil {
		return shim.Error("Term " + termName + " does not exist in the system")
	}
	return shim.Success(termValBytes)
}
