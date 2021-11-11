// COOCompliance project COOCompliance.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChainCode struct{}
type order struct {
	ObjectType           string    `json:"objectType"`
	SalesOrderID         string    `json:"salesOrderID"`
	Item                 string    `json:"item"`
	ItemDescription      string    `json:"itemDescription"`
	Customer             string    `json:"customer"`
	Manufacturer         string    `json:"manufacturer"`
	Shipper              string    `json:"shipper"`
	Supplier             string    `json:"supplier"`
	Quantity             int       `json:"quantity"`
	Event                string    `json:"event"`
	ExpectedDeliveryDate time.Time `json:"expectedDeliveryDate"`
	ActualDeliveryDate   time.Time `json:"actualDeliveryDate"`
	//Accept                string `json:"accept"`
	Exception       string  `json:"exception"`
	DocumentType    string  `json:"documentType"`
	Attachment      string  `json:"attachment"`
	WorkOrderNumber string  `json:"workOrderNumber"`
	InvoiceNumber   string  `json:"invoiceNumber"`
	PONumber        string  `json:"poNumber"`
	Certification   string  `json:"certification"`
	Reference       string  `json:"reference"`
	NetAmount       float64 `json:"netAmount"`
	UnitPrice       float64 `json:"unitPrice"`
	Charges         float64 `json:"charges"`
	Discount        float64 `json:"discount"`
	Tax             float64 `json:"tax"`
	Owner           string  `json:"owner"`
	Custody         string  `json:"custody"`
	CurrentLocation string  `json:"currentLoc"`
	CountryOfOrigin string  `json:"countryOfOrigin"`
	Destination     string  `json:"destination"`
	MaxVibration    float64 `json:"maxVibration"`
	Temperature     float64 `json:"temperature"`
	//UserLevel             string `json:"userLevel"`
	Notification          string `json:"notification"`
	CrossCountryTransport string `json:"crossCountryTrans`
	SerialNumber          string `json:"serialNumber"`
	LotNumber             string `json:"lotNumber"`
	Attribute1            string `json:"attribute1"`
	Attribute2            string `json:"attribute2"`
	Attribute3            string `json:"attribute3"`
	Attribute4            string `json:"attribute4"`
	Attribute5            string `json:"attribute5"`
	Attribute6            string `json:"attribute6"`
	InvalidTrx            string `json:"invalidTrx"`
	Count                 int    `json:"count"`
}

func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}
func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createCOORecord" {
		return t.createCOORecord(stub, args)
	} else if function == "queryCOORecord" {
		return t.queryCOORecord(stub, args)
	} else {
		return shim.Error("Invalid function name " + function)
	}

}
func main() {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Println("Error starting simple chain code " + err.Error())
	}
}

//========================================
//Verify and create COO Compliance record
//========================================
func (t *SimpleChainCode) createCOORecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderID := args[0]
	item := args[1]
	itemDesc := args[2]
	customer := args[3]
	manufacturer := args[4]
	shipper := args[5]
	supplier := args[6]
	quantity, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error(err.Error())
	}
	event := args[8]
	expDelDate := args[9]
	actDelDate := args[10]
	//accept := args[9]
	exception := args[11]
	documentType := args[12]
	attachment := args[13]
	workOrder := args[14]
	invoice := args[15]
	purchaseOrder := args[16]
	certification := args[17]
	reference := args[18]
	netAmount, err := strconv.ParseFloat(args[19], 64)
	unitPrice, err := strconv.ParseFloat(args[20], 64)
	charges, err := strconv.ParseFloat(args[21], 64)
	discount, err := strconv.ParseFloat(args[22], 64)
	tax, err := strconv.ParseFloat(args[23], 64)
	owner := args[24]
	custody := args[25]
	currentLoc := args[26]
	countryOfOrigin := args[27]
	destination := args[28]
	maxVib, err := strconv.ParseFloat(args[29], 64)
	temperature, err := strconv.ParseFloat(args[30], 64)
	//userLvl := args[23]
	//Use notification for repetitive exception notification
	notification := args[31]
	crossCountry := args[32]
	serialNum := args[33]
	lotNum := args[34]
	//Use attr1 for exception counter
	attr1 := args[35]
	//Use attr2 for event code
	attr2 := args[36]
	//Use attr3 for exception notification
	attr3 := args[37]
	//Use attr4 for certification checks
	attr4 := args[38]
	//Use attr5 for terms and conditions check
	attr5 := args[39]
	attr6 := args[40]
	invalidTrx := args[41]
	count, err := strconv.Atoi(args[42])
	if err != nil {
		return shim.Error(err.Error())
	}

	//Convert expected and actual delivery dates to timestamp
	actualDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", actDelDate)
	if err != nil {
		actualDeliveryDate, err = time.Parse("2006-01-02T15:04:05.000Z", "0000-00-00T00:00:00.000Z")

	}

	expectedDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", expDelDate)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Check if the transaction is Country Of Origin Compliant
	if strings.ToLower(countryOfOrigin) == "cuba" {
		//notification = "The Shipment is not Country of Origin Compliant"
		attr3 = countryOfOrigin + " is not in the approved list of countries for importing of goods"
		attr5 = "{\"Country of Origin Compliant\":\"No\"}"
	} else if strings.ToLower(countryOfOrigin) == "iran" {
		//notification = "The Shipment is not Country of Origin Compliant"
		attr3 = countryOfOrigin + " is not in the approved list of countries for importing of goods"
		attr5 = "{\"Country of Origin Compliant\":\"No\"}"
	} else if strings.ToLower(countryOfOrigin) == "north korea" {
		//notification = "The Shipment is not Country of Origin Compliant"
		attr3 = countryOfOrigin + " is not in the approved list of countries for importing of goods"
		attr5 = "{\"Country of Origin Compliant\":\"No\"}"
	} else if strings.ToLower(countryOfOrigin) == "syria" {
		//notification = "The Shipment is not Country of Origin Compliant"
		attr3 = countryOfOrigin + " is not in the approved list of countries for importing of goods"
		attr5 = "{\"Country of Origin Compliant\":\"No\"}"
	} else {
		//notification = "The Shipment is Country of Origin Compliant"
		attr3 = "The Shipment is Country of Origin Compliant"
		attr5 = "{\"Country of Origin Compliant\":\"Yes\"}"
	}

	//Update event name
	event = "Country of Origin Compliance Verification"
	var response bytes.Buffer
	response.WriteString(attr3)

	//Update an order object
	objectType := "sales order"
	orderObj := &order{objectType, orderID, item, itemDesc, customer, manufacturer, shipper, supplier, quantity, event, expectedDeliveryDate, actualDeliveryDate, exception, documentType, attachment, workOrder, invoice, purchaseOrder, certification, reference, netAmount, unitPrice, charges, discount, tax, owner, custody, currentLoc, countryOfOrigin, destination, maxVib, temperature, notification, crossCountry, serialNum, lotNum, attr1, attr2, attr3, attr4, attr5, attr6, invalidTrx, count}

	//Convert the order object to JSON object
	orderBytes, err := json.Marshal(orderObj)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Create Write Set
	err = stub.PutState(orderID, orderBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent(event, orderBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(response.Bytes())
}

//===========================
//Query COO Compliance record
//===========================
func (t *SimpleChainCode) queryCOORecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//Check if the number of arguments is 1
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	orderID := args[0]

	//Get the current state
	orderBytes, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error(err.Error())
	} else if orderBytes == nil {
		return shim.Error("Invalid Order ID " + orderID)
	}
	return shim.Success(orderBytes)
}
