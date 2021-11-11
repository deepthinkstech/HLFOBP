// SalesTransactions.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=================================================================
//Define Structure which will be used to implement all shim methods
//==================================================================
type SimpleChainCode struct {
}

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
	//Adding additional counter for PTR Track and Trace App requirement
	Count int `json:"count"`
}

type PurchaseOrder struct {
	ObjectType            string    `json:"objectType"`
	OrderNumber           string    `json:"orderNumber"`
	Item                  string    `json:"item"`
	DocumentType          string    `json:"documentType"`
	SoldToLegalEntityID   string    `json:"soldToLegalEntityID"`
	SoldToLegalEntity     string    `json:"soldToLegalEntity"`
	ChangeOrderDesc       string    `json:"changeOrderDesc"`
	PaymentTerms          string    `json:"paymentTerms"`
	Action                string    `json:"change"`
	LineNum               int       `json:"lineNum`
	CurrencyCode          string    `json:"currencyCode"`
	Price                 float64   `json:"price"`
	ChangeReason          string    `json:"changeReason"`
	LineAction            string    `json:"lineAction"`
	RefNumber             string    `json:"refNumber"`
	Event                 string    `json:"event"`
	EventCode             string    `json:"eventCode"`
	OrderHeaderID         string    `json:"orderHeaderID"`
	ProcurementBUID       string    `json:"procurementBUID"`
	ProcurementBU         string    `json:"procurementBU"`
	BillToBUID            string    `json:"billToBUID"`
	BillToBU              string    `json:"billToBU"`
	Buyer                 string    `json:"buyer"`
	SupplierName          string    `json:"supplierName"`
	ShipToLocation        string    `json:"shipToLocation"`
	SupplierSite          string    `json:"supplierSite"`
	Quantity              int       `json:"quantity"`
	OrderedAmount         float64   `json:"orderedAmount"`
	TotalAmount           float64   `json:"totalAmount"`
	RequestedDeliveryDate time.Time `json:"requestedDeliveryDate"`
	SerialNo              string    `json:"serialNo"`
	ActualDeliveryDate    time.Time `json:"actualDeliveryDate"`
	QAResult              string    `json:"qaResult"`
	ItemDesc              string    `json:"itemDesc"`
	Notification          string    `json:"notification"`
	Owner                 string    `json:"owner"`
	Custody               string    `json:"custody"`
	//Adding additional counter for PTR Track and Trace App requirement
	Count int `json:"count"`
}

type OrderInfo struct {
	OrderKey string
	OrderVal string
}

//========================
//Initialize the chaincode
//========================
func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//=======================
//Invocation of chaincode
//=======================
func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createOrder" {
		return t.createOrder(stub, args)
	} else if function == "queryOrder" {
		return t.queryOrder(stub, args)
	} else if function == "queryByParentOrder" {
		return t.queryByParentOrder(stub, args)
	} else if function == "queryTrxHistory" {
		return t.queryTrxHistory(stub, args)
	} else if function == "updateOrder" {
		return t.updateOrder(stub, args)
	} else if function == "queryTrxHistoryByParentOrder" {
		return t.queryTrxHistoryByParentOrder(stub, args)
	} else if function == "queryAllChildOrders" {
		return t.queryAllChildOrders(stub, args)
	} else if function == "queryTrxHistoryV2" {
		return t.queryTrxHistoryV2(stub, args)
	} else if function == "getSOWOLink" {
		return t.getSOWOLink(stub, args)
	} else if function == "queryOrderByLatestOrder" {
		return t.queryOrderByLatestOrder(stub, args)
	} else if function == "getSalesOrderCount" {
		return t.getSalesOrderCount(stub)
	} else if function == "getTrxCount" {
		return t.getTrxCount(stub)
	} else if function == "queryLatestStateByRef" {
		return t.queryLatestStateByRef(stub, args)
	} else {
		return shim.Error("Not a valid function " + function)
	}
}

//===============================================================================
//Declaration of the main function which inturn will call Init and Invoke methods
//===============================================================================
func main() {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Cannot start simple chaincode %s", err)
	}
}

//============================================
//Create Order - Chaincode to create new order
//============================================
func (t *SimpleChainCode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//Code added as work-around for OIC sending NULL instead of '' - START
	var j int
	for j = 0; j < len(args); j++ {
		if args[j] == "null" {
			args[j] = ""
		}
	}
	//Code added as work-around for OIC sending NULL instead of '' - END

	if len(args) != 33 {
		return shim.Error("Error 1 Incorrect number of arguments, expecting 33")
	} else if len(args[0]) == 0 {
		return shim.Error("Error 2 Order ID cannot be null")
	}

	//Map the inputs received to variables
	orderID := args[0]
	item := args[1]
	itemDesc := args[2]
	customer := args[3]
	manufacturer := args[4]
	shipper := args[5]
	supplier := args[6]
	quantity, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
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
	if err != nil {
		netAmount = 0
	}
	unitPrice, err := strconv.ParseFloat(args[20], 64)
	if err != nil {
		unitPrice = 0
	}
	charges, err := strconv.ParseFloat(args[21], 64)
	if err != nil {
		charges = 0
	}
	discount, err := strconv.ParseFloat(args[22], 64)
	if err != nil {
		discount = 0
	}
	tax, err := strconv.ParseFloat(args[23], 64)
	if err != nil {
		tax = 0
	}
	owner := ""
	custody := ""
	currentLoc := args[24]
	countryOfOrigin := args[25]
	destination := args[26]
	maxVib, err := strconv.ParseFloat(args[27], 64)
	if err != nil {
		maxVib = 0
	}
	temperature, err := strconv.ParseFloat(args[28], 64)
	if err != nil {
		temperature = 0
	}
	//userLvl := args[23]
	//Use notification for repetitive exception notification
	notification := args[29]
	crossCountry := ""
	serialNum := args[30]
	lotNum := args[31]
	//Use attr1 for exception counter
	attr1 := strconv.Itoa(0)
	//Use attr2 for event code
	attr2 := args[32]
	//Use attr3 for exception notification
	attr3 := ""
	//Use attr4 for certification checks
	attr4 := ""
	//Use attr5 for terms and conditions check
	attr5 := ""
	attr6 := "a"
	invalidTrx := "N"
	count := 0

	//Convert expected and actual delivery dates to timestamp

	actualDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", actDelDate)
	if err != nil {
		actualDeliveryDate, err = time.Parse("2006-01-02T15:04:05.000Z", "0000-00-00T00:00:00.000Z")

	}

	expectedDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", expDelDate)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}

	//Check if order already exists
	orderBytes, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error(err.Error())
	} else if orderBytes != nil {
		return shim.Error("Error 5 Order already exists " + orderID)
	}
	//Check if order quantity has positive value
	if quantity <= 0 {
		return shim.Error("Error 6 Item quantity cannot be lesser or equal to zero")
	}
	//Check if item is valid
	if len(item) == 0 {
		return shim.Error("Error 7 Item cannot be null")
	}
	//Check if manufacturer and customer information is available
	if len(customer) == 0 {
		return shim.Error("Error 8 Customer cannot be null")
	}
	//Check if Work Order ID is available for event Work Order created
	if event == "Work Order created" && len(workOrder) == 0 {
		return shim.Error("Error 9 Work Order ID cannot be null when event is " + event)
	}
	//Check if Purchase Order ID is available for event Purchase Order created
	if event == "Purchase Order created" && len(purchaseOrder) == 0 {
		return shim.Error("Error 10 Purchase Order ID cannot be null when event is " + event)
	}
	//Check if Country of Origin & Destination are available
	if len(countryOfOrigin) == 0 {
		return shim.Error("Error 11 Country of Origin cannot be null")
	}
	if len(destination) == 0 {
		return shim.Error("Error 12 Destination country cannot be null")
	}

	//Assign values to fields determined by Blockchain based on Item ID
	//Assign whether the shipping involves multi-country travel
	if countryOfOrigin == destination {
		crossCountry = "No"
	} else {
		crossCountry = "Yes"
	}
	//Assign Owner
	//Back-to-Back and Drop Shipment
	if (customer != "Vision Operations" || customer != "RedCube Mfg") && supplier == "" {
		owner = manufacturer
	} else if (customer != "Vision Operations" || customer != "RedCube Mfg") && supplier != "" {
		owner = supplier
	}

	//Contract Manufacturing and Raw Material Procurement
	if (customer != "Vision Operations" || customer != "RedCube Mfg") && manufacturer != "" {
		owner = manufacturer
	} else if (customer != "Vision Operations" || customer != "RedCube Mfg") && manufacturer == "" {
		if event == "Purchase Order Release" {
			owner = customer
		} else {
			if supplier != "" {
				owner = supplier
			} else {
				owner = manufacturer
			}

		}
	}

	//Assign Custody
	//Back-to-Back and Drop Shipment
	if (customer != "Vision Operations" || customer != "RedCube Mfg") && supplier == "" {
		if event == "Order Received" {
			custody = ""
		} else {
			custody = manufacturer
		}
	} else if (customer != "Vision Operations" || customer != "RedCube Mfg") && supplier != "" {
		if event == "Order Received" {
			custody = ""
		} else {
			custody = supplier
		}
	}

	//Contract Manufacturing and Raw Material Procurement
	if (customer != "Vision Operations" || customer != "RedCube Mfg") && manufacturer != "" {
		if event == "Order Received" {
			custody = ""
		} else {
			custody = manufacturer
		}
	} else if (customer != "Vision Operations" || customer != "RedCube Mfg") && manufacturer == "" {
		if event != "Order Received" {
			custody = supplier
		} else {
			custody = ""
		}
	}
	//Check if sales order quantities match the purchase order quantities and sales order is created only if the quantites match
	//Invoke purchaseordertransactions chaincode and get the quantity details
	poNumber := reference
	if len(poNumber) != 0 {
		pochannel := "orderprocessing"
		pochaincode := "purchaseordertransactions"
		pofunction := "checkPOandSOQuantity"
		quantityStr := strconv.Itoa(quantity)

		poargs := util.ToChaincodeArgs(pofunction, poNumber, quantityStr)
		poResp := stub.InvokeChaincode(pochaincode, poargs, pochannel)
		if poResp.Status == shim.OK {
			//return shim.Error("Error 13 Incorrect Purchase Order Number " + poNumber)

			respBytes := poResp.Payload
			var respString bytes.Buffer
			var j int
			for j = 0; j < len(respBytes); j++ {
				respString.WriteByte(respBytes[j])
			}
			fmt.Println(respString.String())

			if respString.String() == "Purchase Order and Sales Order quantities do not match, cannot create Sales Order" {
				return shim.Success(respBytes)
			}
			//notification = respString.String()
		}
	}

	//Assign whether certification is obtained or not
	certification = "Not Obtained"

	//Assign Current Location
	currentLoc = countryOfOrigin

	/*
		//PTR Track and Trace App - Increment the count of orders
		//Check if the order belongs to PTR organizations
		var orderCode string
		var count int
		if customer == "Get Well Hospital" || customer == "MedSupply Corp" {
			if customer == "Get Well Hospital" {
				orderCode = "ManufacturerSO"
			} else {
				orderCode = "DistributorSO"
			}

			orderchannel := "orderprocessing"
			orderchaincode := "latestorders"
			orderfunction := "queryOrder"
			orderargs := util.ToChaincodeArgs(orderfunction, orderCode)
			orderResp := stub.InvokeChaincode(orderchaincode, orderargs, orderchannel)
			if orderResp.Status != shim.OK {
				count = 1
			} else {
				respBytes := orderResp.Payload
				respObj := &OrderInfo{}
				err := json.Unmarshal(respBytes, respObj)
				if err != nil {
					return shim.Error("Error 14 Order Info cannot be retrieved " + err.Error())
				}

				existingOrderID := respObj.OrderVal
				existingOrderBytes, err := stub.GetState(existingOrderID)
				if err != nil {
					return shim.Error("Error 15 " + err.Error())
				}
				existingOrderObj := &order{}
				err = json.Unmarshal(existingOrderBytes, existingOrderObj)
				if err != nil {
					return shim.Error("Error 16 " + err.Error())
				}
				count = existingOrderObj.Count + 1
			}
		} else {
			count = 0
		}

		//PTR Track and Trace App - Increment the count of orders
	*/

	//Create an order object
	objectType := "sales order"
	orderObj := &order{objectType, orderID, item, itemDesc, customer, manufacturer, shipper, supplier, quantity, event, expectedDeliveryDate, actualDeliveryDate, exception, documentType, attachment, workOrder, invoice, purchaseOrder, certification, reference, netAmount, unitPrice, charges, discount, tax, owner, custody, currentLoc, countryOfOrigin, destination, maxVib, temperature, notification, crossCountry, serialNum, lotNum, attr1, attr2, attr3, attr4, attr5, attr6, invalidTrx, count}

	//Convert the order object to JSON object
	orderBytes, err = json.Marshal(orderObj)
	if err != nil {

		return shim.Error("Error 17 " + err.Error())
	}

	//Create Write Set
	err = stub.PutState(orderID, orderBytes)
	if err != nil {
		return shim.Error("Error 18 " + err.Error())
	}

	/*
		//PTR Track and Trace App - Increment the count of SO transactions
		//Check if the order belongs to PTR organizations - if yes, create/update the ledger where count of trx are maintained
		if customer == "Get Well Hospital" || customer == "MedSupply Corp" {
			key := "SOCount"
			var curTrxCount, trxCount int
			countBytes, err := stub.GetState(key)
			if err != nil {
				trxCount = 1
			}
			if countBytes == nil {
				trxCount = 1
			} else {
				curTrxCountStr := string(countBytes)
				curTrxCount, err = strconv.Atoi(curTrxCountStr)
				if err != nil {
					return shim.Error("Error 19 " + err.Error())
				}
				trxCount = curTrxCount + 1
			}
			trxCountStr := strconv.Itoa(trxCount)
			err = stub.PutState(key, []byte(trxCountStr))
			if err != nil {
				return shim.Error("Error 20 " + err.Error())
			}

		}
		//PTR Track and Trace App - Increment the count of SO transactions
	*/

	//Create index based on Parent Order ID and orderID
	if len(reference) != 0 {
		indexName := "orderIndex"
		compositeKey, err := stub.CreateCompositeKey(indexName, []string{reference, orderID})
		if err != nil {
			return shim.Error("Error 21 " + err.Error())
		}
		value := []byte{0X00}
		err = stub.PutState(compositeKey, value)
		if err != nil {
			return shim.Error("Error 22 " + err.Error())
		}
	}

	/*Workaround for the composite key issue
	indexName := "orderIndex"
	compKeyArr := []string{indexName, reference}
	compKey := strings.Join(compKeyArr, "|")
	err = stub.PutState(compKey, []byte(orderID))
	if err != nil {
		return shim.Error("Error 21 " + err.Error())
	}
	*/

	//Create index based on Work Order No and orderID
	if len(workOrder) != 0 {
		iName := "woLink"
		woLinkKey, err := stub.CreateCompositeKey(iName, []string{workOrder, orderID})
		if err != nil {
			return shim.Error("Error 23 " + err.Error())
		}
		woLinkValue := []byte{0X00}
		err = stub.PutState(woLinkKey, woLinkValue)
		if err != nil {
			return shim.Error("Error 24 " + err.Error())
		}
	}

	/*Workaround for the composite key issue
	if len(workOrder) != 0 {
		iName := "woLink"
		compKeyArr := []string{iName, workOrder}
		compKey := strings.Join(compKeyArr, "|")
		err = stub.PutState(compKey, []byte(orderID))
		if err != nil {
			return shim.Error("Error 22 " + err.Error())
		}
	}
	*/
	//Set event in chaincode
	err = stub.SetEvent(event, orderBytes)
	if err != nil {
		return shim.Error("Error 23 " + err.Error())
	}
	return shim.Success(nil)
}

//=======================================================================
//Query Order -Get current state of order from state DB based on Order ID
//=======================================================================
func (t *SimpleChainCode) queryOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//Check if the number of arguments is 1
	/*
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments, expecting 1")
		}
	*/
	orderID := args[0]

	//Get the current state
	orderBytes, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if orderBytes == nil {
		return shim.Error("Error 2 Invalid Order ID " + orderID)
	}
	return shim.Success(orderBytes)
}

//=======================================================================
//Query by Parent Order - Get the list of orders based on parent order ID
//=======================================================================
func (t *SimpleChainCode) queryByParentOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//Check if the number of arguments is 1
	/*
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments, expecting 1")
		}
	*/
	parentOrderID := args[0]
	indexName := "orderIndex"
	//Get the list of orders which have the corresponding parent Order ID
	orderIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{parentOrderID})
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if orderIterator == nil {
		return shim.Error("Error 2 Invalid Order Reference ID " + parentOrderID)
	}
	defer orderIterator.Close()
	//Write individual transactions into a buffer which will be sent as output
	var i int
	var buffer bytes.Buffer
	isRecordWritten := false
	buffer.WriteString("[")
	for i = 0; orderIterator.HasNext(); i++ {
		response, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		returnIndex, returnKeys, err := stub.SplitCompositeKey(response.Key)
		if err != nil {
			return shim.Error("Error 4 " + err.Error())
		}
		returnParentOrderID := returnKeys[0]
		returnOrderID := returnKeys[1]
		fmt.Println("Returned Index Name " + returnIndex)
		if isRecordWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Parent Order ID\":\"")
		buffer.WriteString(returnParentOrderID)
		buffer.WriteString("\",")
		buffer.WriteString("\"Order ID\":\"")
		buffer.WriteString(returnOrderID)
		buffer.WriteString("\"}")
		isRecordWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=======================================================================================
//Query Transaction History V2 - Get the complete transaction history for a particular Order
//=======================================================================================

func (t *SimpleChainCode) queryTrxHistoryV2(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	orderID := args[0]
	//Get the transaction history and write into a buffer
	var buffer bytes.Buffer
	var i int
	isRecordWritten := false
	orderIterator, err := stub.GetHistoryForKey(orderID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	defer orderIterator.Close()

	buffer.WriteString("[")
	for i = 0; orderIterator.HasNext(); i++ {
		if isRecordWritten == true {
			buffer.WriteString(",")
		}
		response, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		buffer.WriteString("{\"Order Type\":\"Sales Order\",")
		buffer.WriteString("\"Order ID\":\"")
		buffer.WriteString(orderID)
		buffer.WriteString("\",")
		buffer.WriteString("\"Transaction ID\":\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\",")
		buffer.WriteString("\"Value\":")
		if response.IsDelete == true {
			buffer.WriteString("null")
		} else {
			//buffer.WriteString("{\"")
			buffer.WriteString(string(response.Value))
			//buffer.WriteString("\"}")
		}
		buffer.WriteString(",")
		buffer.WriteString("\"TimeStamp\":\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")
		/*
			buffer.WriteString("\"IsDelete\":\"")
			buffer.WriteString(strconv.FormatBool(response.IsDelete))
		*/
		buffer.WriteString("}")
		isRecordWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=======================================================================================
//Query Transaction History - Get the complete transaction history for a particular Order
//=======================================================================================
func (t *SimpleChainCode) queryTrxHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//Check if correct number of arguments are sent
	/*
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments, expecting 1")
		}
	*/
	orderID := args[0]
	//Get the transaction history and write into a buffer
	var buffer bytes.Buffer
	var i int
	isRecordWritten := false
	orderIterator, err := stub.GetHistoryForKey(orderID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	defer orderIterator.Close()

	buffer.WriteString("[")
	for i = 0; orderIterator.HasNext(); i++ {
		if isRecordWritten == true {
			buffer.WriteString(",")
		}
		response, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		buffer.WriteString("{\"Transaction ID\":\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\",")
		buffer.WriteString("\"Value\":")
		if response.IsDelete == true {
			buffer.WriteString("null")
		} else {
			//buffer.WriteString("{\"")
			buffer.WriteString(string(response.Value))
			//buffer.WriteString("\"}")
		}
		buffer.WriteString(",")
		buffer.WriteString("\"TimeStamp\":\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\",")
		buffer.WriteString("\"IsDelete\":\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"}")
		isRecordWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=====================================================
//Update Order - Update Order Information into state DB
//=====================================================
func (t *SimpleChainCode) updateOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//Check the number of arguments
	if len(args) != 33 {
		return shim.Error("Error 1 Incorrect number of arguments, expecting 33")
	}
	//Code added as work-around for OIC sending NULL instead of '' - START
	var j int
	for j = 0; j < len(args); j++ {
		if args[j] == "null" {
			args[j] = ""
		}
	}
	//Code added as work-around for OIC sending NULL instead of '' - END
	//Map the inputs received to variables
	orderID := args[0]
	item := args[1]
	itemDesc := args[2]
	customer := args[3]
	manufacturer := args[4]
	shipper := args[5]
	supplier := args[6]
	quantity, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error("Error 2 " + err.Error())
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
	if err != nil {
		netAmount = 0
	}
	unitPrice, err := strconv.ParseFloat(args[20], 64)
	if err != nil {
		unitPrice = 0
	}
	charges, err := strconv.ParseFloat(args[21], 64)
	if err != nil {
		charges = 0
	}
	discount, err := strconv.ParseFloat(args[22], 64)
	if err != nil {
		discount = 0
	}
	tax, err := strconv.ParseFloat(args[23], 64)
	if err != nil {
		tax = 0
	}
	owner := ""
	custody := ""
	currentLoc := args[24]
	countryOfOrigin := args[25]
	destination := args[26]
	maxVib, err := strconv.ParseFloat(args[27], 64)
	if err != nil {
		maxVib = 0
	}
	temperature, err := strconv.ParseFloat(args[28], 64)
	if err != nil {
		temperature = 0
	}
	//userLvl := args[23]
	//Use notification for repetitive exception notification
	notification := args[29]
	crossCountry := ""
	serialNum := args[30]
	lotNum := args[31]
	//Use attr1 for exception counter
	attr1 := strconv.Itoa(0)
	//Use attr2 for event code
	attr2 := args[32]
	//Use attr3 for exception notification
	attr3 := ""
	//Use attr4 for certification checks
	attr4 := ""
	//Use attr5 for terms and conditions check
	attr5 := ""
	attr6 := "a"
	invalidTrx := "N"

	taxStr := args[23]
	count := 0

	//Convert expected and actual delivery dates to timestamp
	actualDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", actDelDate)
	if err != nil {
		actualDeliveryDate, err = time.Parse("2006-01-02T15:04:05.000Z", "0000-00-00T00:00:00.000Z")

	}

	expectedDeliveryDate, err := time.Parse("2006-01-02T15:04:05.000Z", expDelDate)
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
	}

	//Check if order ID exists in the state DB
	orderBytes, err := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	} else if orderBytes == nil {
		return shim.Error("Error 5 Invalid Order ID " + orderID)
	}
	//Convert JSON object to Order object
	orderObject := order{}
	err = json.Unmarshal(orderBytes, &orderObject)
	if err != nil {
		return shim.Error("Error 6 " + err.Error())
	}
	count = orderObject.Count
	//Check if transaction can be committed to blockchain or not
	if orderObject.InvalidTrx == "Y" {
		if orderObject.Attribute6 == "e" {
			return shim.Error("Shipment cannot proceed because the country of import is not in the list of compliant countries")
		} else if orderObject.Attribute6 == "b" {
			if event != "RoHs Compliance Certificate" {
				return shim.Error("Shipment cannot proceed because RoHs compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "c" {
			if event != "Conflict Minerals Compliance" {
				return shim.Error("Shipment cannot proceed because Conflict Minerals compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "d" {
			if event != "Final burn-in and Test Certificate" {
				return shim.Error("Shipment cannot proceed because Final burn-in and Test Certificate compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "bc" {
			if event == "RoHs Compliance Certificate" {
				invalidTrx = "Y"
				attr6 = "c"
				attr3 = "Shipment cannot proceed because Conflict Minerals compliance documents are missing. Please upload the relevant documents"
			} else if event == "Conflict Minerals Compliance" {
				invalidTrx = "Y"
				attr6 = "b"
				attr3 = "Shipment cannot proceed because RoHs compliance documents are missing. Please upload the relevant documents"
			} else {
				return shim.Error("Shipment cannot proceed because RoHs and Conflict Minerals compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "bd" {
			if event == "RoHs Compliance Certificate" {
				invalidTrx = "Y"
				attr6 = "d"
				attr3 = "Shipment cannot proceed because Final burn-in and Test Certificate compliance documents are missing. Please upload the relevant documents"
			} else if event == "Final burn-in and Test Certificate" {
				invalidTrx = "Y"
				attr6 = "b"
				attr3 = "Shipment cannot proceed because RoHs compliance documents are missing. Please upload the relevant documents"
			} else {
				return shim.Error("Shipment cannot proceed because RoHs and Final burn-in and Test compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "cd" {
			if event == "Conflict Minerals Compliance" {
				invalidTrx = "Y"
				attr6 = "d"
				attr3 = "Shipment cannot proceed because Final burn-in and Test compliance documents are missing. Please upload the relevant documents"
			} else if event == "Conflict Minerals Compliance" {
				invalidTrx = "Y"
				attr6 = "c"
				attr3 = "Shipment cannot proceed because Conflict Minerals compliance documents are missing. Please upload the relevant documents"
			} else {
				return shim.Error("Shipment cannot proceed because Conflict Minerals and Final burn-in and Test compliance documents are missing. Please upload the relevant documents")
			}
		} else if orderObject.Attribute6 == "bcd" {
			if event == "RoHs Compliance Certificate" {
				invalidTrx = "Y"
				attr6 = "cd"
				attr3 = "Shipment cannot proceed because Conflict Minerals and Final burn-in and Test compliance documents are missing. Please upload the relevant documents"
			} else if event == "Conflict Minerals Compliance" {
				invalidTrx = "Y"
				attr6 = "bd"
				attr3 = "Shipment cannot proceed because RoHs and Final burn-in and Test compliance documents are missing. Please upload the relevant documents"
			} else if event == "Final burn-in and Test Certificate" {
				invalidTrx = "Y"
				attr6 = "bc"
				attr3 = "Shipment cannot proceed because RoHs and Conflict Minerals compliance documents are missing. Please upload the relevant documents"
			} else {
				return shim.Error("Shipment cannot proceed because compliance documents are missing. Please upload the relevant documents")
			}
		}

	}

	if len(event) == 0 {
		return shim.Error("Event name cannot be null")
	}
	//Check if manufacturer and customer information is available
	if len(customer) == 0 {
		return shim.Error("Customer cannot be null")
	}
	//Check if shipper information is available	when event is Shipment Executed
	if event == "Shipment Executed" && len(shipper) == 0 {
		return shim.Error("Shipper cannot be null when event is " + event)
	}
	//Check if manufacturer information has not changed
	if manufacturer != orderObject.Manufacturer {
		return shim.Error("Manufacturer information cannot be tampered with")
	}
	//Check if customer information has not changed
	if customer != orderObject.Customer {
		return shim.Error("Customer information cannot be tampered with")
	}
	//Check if item information has not changed
	if item != orderObject.Item {
		return shim.Error("Item information cannot be tampered with")
	}
	//Check if item quantity has not changed

	if quantity <= 0 {
		return shim.Error("Item quantity cannot be lesser or equal to zero")
	}
	/*
		//Check if shipper information has not changed
		if event != "Shipment Executed" && shipper != orderObject.Shipper {
			return shim.Error("Shipper information cannot be tampered with")
		}
	*/

	//Check if the country of origin, destination has not changed
	if countryOfOrigin != orderObject.CountryOfOrigin {
		return shim.Error("Country of Origin information has been tampered with")
	}

	if destination != orderObject.Destination {
		return shim.Error("Destination of shipment information has been tampered with")
	}

	//Check if necessary documentation is available
	if (event == "RoHs Compliance Certificate" || event == "Conflict Minerals Compliance" || event == "Final burn-in and Test Certificate") && len(attachment) == 0 {
		return shim.Error("Document is mandatory for event " + event)
	}
	//Check if Invoice ID is available for event Invoice Generated
	if event == "Invoice Generated" && len(invoice) == 0 {
		return shim.Error("Invoice ID cannot be null when event is " + event)
	}
	//Check if Shipper information is available for event Invoice Generated
	if event == "Invoice Generated" && len(shipper) == 0 {
		return shim.Error("Shipper cannot be null when Invoice has been generated")
	}

	//Assign Reference number
	if len(reference) == 0 {
		reference = orderObject.Reference
	}
	//Assign Ownership
	//Back-to-Back and Drop Shipment + Contract Manufacturing and Raw Material Procurement
	if (customer != "Vision Operations" || customer != "RedCube Mfg") && supplier == "" {
		if event == "Equipment Installation – In Progress" {
			owner = manufacturer
		} else if event == "Customer Accepted" || event == "Purchase Order Receipt" || event == "Payment – Approved OK to Pay" || event == "Customer Acceptance" {
			owner = customer
		} else {
			owner = orderObject.Owner
		}
	} else if event == "Customer Accepted" || event == "Purchase Order Receipt" || event == "Payment – Approved OK to Pay" || event == "Customer Acceptance" {
		owner = customer
	} else {
		owner = orderObject.Owner
	}

	//Custody Transfers
	if event == "Export Compliance Documentation" && attr2 == "200" {
		custody = customer
	} else if event == "Order Shipped" {
		custody = shipper
	} else {
		if orderObject.Custody == "" {
			custody = orderObject.Owner
		} else {
			custody = orderObject.Custody
		}
	}
	//Update current location
	if event == "Export Compliance Documentation" && attr2 == "200" {
		currentLoc = destination
	} else {
		currentLoc = orderObject.CurrentLocation
	}
	//Assign Cross Country
	crossCountry = orderObject.CrossCountryTransport

	attr1 = orderObject.Attribute1
	if len(notification) == 0 {
		notification = orderObject.Notification
	}
	//attr3 = ""

	//Check terms and conditions are met based on document availability
	//Check if necessary certificates are available at the time of shipment execution
	if event == "Shipment Executed" {

		//Check for RoHs, Conflict Minerals & Final burn-in and Test Certificate compliance
		tcIterator, err := stub.GetHistoryForKey(orderID)
		if err != nil {
			return shim.Error("Error 7 " + err.Error())
		}
		defer tcIterator.Close()
		var i int
		for i = 0; tcIterator.HasNext(); i++ {
			tcRange, err := tcIterator.Next()
			if err != nil {
				return shim.Error("Error 8 " + err.Error())
			}
			tcBytes := tcRange.Value
			tcObj := order{}
			err = json.Unmarshal(tcBytes, &tcObj)
			if err != nil {
				return shim.Error("Error 9 " + err.Error())
			}
			if tcObj.Event == "RoHs Compliance Certificate" && len(tcObj.Attachment) != 0 { //Added on 17 Aug 2019 to address scenario of random upload of documents
				if attr5 == "c" {
					attr5 = "b"
				} else if attr5 == "f" {
					attr5 = "g"
				} else if attr5 == "e" {
					attr5 = "d"
				} else {
					attr5 = "a"
				}
			} else if tcObj.Event == "Conflict Minerals Compliance" && len(tcObj.Attachment) != 0 {
				if attr5 == "a" {
					attr5 = "b"
				} else if attr5 == "f" {
					attr5 = "e"
				} else if attr5 == "g" {
					attr5 = "d"
				} else {
					attr5 = "c"
				}
			} else if tcObj.Event == "Final burn-in and Test Certificate" && len(tcObj.Attachment) != 0 {
				if attr5 == "a" {
					attr5 = "g"
				} else if attr5 == "b" {
					attr5 = "d"
				} else if attr5 == "c" {
					attr5 = "e"
				} else {
					attr5 = "f"
				}
			}
		}
		if attr5 == "a" {
			attr5 = "{\"RoHs Compliance Certificate\":\"Yes\",\"Conflict Minerals Compliance\":\"No\",\"Final burn-in and Test Certificate\":\"No\""
			attr6 = "cd"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else if attr5 == "b" {
			attr5 = "{\"RoHs Compliance Certificate\":\"Yes\",\"Conflict Minerals Compliance\":\"Yes\",\"Final burn-in and Test Certificate\":\"No\""
			attr6 = "d"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else if attr5 == "c" {
			attr5 = "{\"RoHs Compliance Certificate\":\"No\",\"Conflict Minerals Compliance\":\"Yes\",\"Final burn-in and Test Certificate\":\"No\""
			attr6 = "bd"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else if attr5 == "d" {
			attr5 = "{\"RoHs Compliance Certificate\":\"Yes\",\"Conflict Minerals Compliance\":\"Yes\",\"Final burn-in and Test Certificate\":\"Yes\""
			attr6 = "a"
			invalidTrx = "N"
		} else if attr5 == "e" {
			attr5 = "{\"RoHs Compliance Certificate\":\"No\",\"Conflict Minerals Compliance\":\"Yes\",\"Final burn-in and Test Certificate\":\"Yes\""
			attr6 = "b"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else if attr5 == "f" {
			attr5 = "{\"RoHs Compliance Certificate\":\"No\",\"Conflict Minerals Compliance\":\"No\",\"Final burn-in and Test Certificate\":\"Yes\""
			attr6 = "bc"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else if attr5 == "g" {
			attr5 = "{\"RoHs Compliance Certificate\":\"Yes\",\"Conflict Minerals Compliance\":\"No\",\"Final burn-in and Test Certificate\":\"Yes\""
			attr6 = "c"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		} else {
			attr5 = "{\"RoHs Compliance Certificate\":\"No\",\"Conflict Minerals Compliance\":\"No\",\"Final burn-in and Test Certificate\":\"No\""
			attr6 = "bcd"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			attr3 = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		}

		//Check for Country of Origin Compliance
		chaincodeName := "coocompliance"
		channelName := "orderprocessing"
		f := "createCOORecord"
		countStr := strconv.Itoa(orderObject.Count)
		inputArgs := util.ToChaincodeArgs(f, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22], args[23], orderObject.Owner, orderObject.Custody, args[24], args[25], args[26], args[27], args[28], args[29], orderObject.CrossCountryTransport, args[30], args[31], orderObject.Attribute1, args[32], orderObject.Attribute3, orderObject.Attribute4, orderObject.Attribute5, orderObject.Attribute6, orderObject.InvalidTrx, countStr)
		response := stub.InvokeChaincode(chaincodeName, inputArgs, channelName)
		if response.Status != shim.OK {
			errStatus := "Failed to query chaincode. Got error: " + response.Message
			return shim.Error(errStatus)
		}
		respBytes := response.Payload
		var respString bytes.Buffer
		var j int
		for j = 0; j < len(respBytes); j++ {
			respString.WriteByte(respBytes[j])
		}
		fmt.Println(respString.String())
		if respString.String() != "The Shipment is Country of Origin Compliant" {
			COOConcate := []string{attr5, ",\"Country Of Origin Compliance\":\"No\"}"}
			attr5 = strings.Join(COOConcate, "")
			attr6 = "e"
			invalidTrx = "Y"
		} else {
			COOConcate := []string{attr5, ",\"Country Of Origin Compliance\":\"Yes\"}"}
			attr5 = strings.Join(COOConcate, "")
		}
	}

	//Check if export compliance documentation is available after shipment reaches dest
	if event == "Equipment Installation – In Progress" {
		expCompBytes, err := stub.GetState(orderID)
		if err != nil {
			return shim.Error("Error 10 " + err.Error())
		}
		expCompObj := order{}
		err = json.Unmarshal(expCompBytes, &expCompObj)
		if err != nil {
			return shim.Error("Error 11 " + err.Error())
		}
		if expCompObj.Event == "Export Compliance Documentation" && len(expCompObj.Attachment) != 0 {
			attr5 = "{\"Export Compliance Documentation\":\"Yes\"}"
			invalidTrx = "N"
		} else {
			attr5 = "{\"Export Compliance Documentation\":\"No\"}"
			invalidTrx = "Y"
			notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
		}
	}
	/*
		//Check if export compliance documentation is available before invoicing
		if event == "Invoice Generated" {
			expCompBytes, err := stub.GetState(orderID)
			if err != nil {

				return shim.Error(err.Error())
			}
			expCompObj := order{}
			err = json.Unmarshal(expCompBytes, &expCompObj)
			if err != nil {

				return shim.Error(err.Error())
			}
			if expCompObj.Event == "Export Compliance Documentation" && len(expCompObj.Attachment) != 0 {
				attr5 = "{\"Export Compliance Documentation\":\"Yes\"}"
				invalidTrx = "N"
			} else {
				attr5 = "{\"Export Compliance Documentation\":\"No\"}"
				invalidTrx = "Y"
				notification = event + "  is not valid due to missing compliance documents and has been flagged as Invalid Transaction in Blockchain"
			}
		}
	*/
	//Checks for compliance for attr5

	//a - RoHs complaint
	//b - RoHs and Conflict Minerals compliant
	//c - Only Conflict Minerals compliant
	//d - RoHs, Conflict Minerals and Final burn-in and Test compliant
	//e - Conflict Minerals and Final burn-in and Test compliant
	//f - Final burn-in and Test compliant
	//g - Rohs and Final burn-in and Test compliant

	//Added for V2 for attr6
	//a - Default value
	//b - RoHs compliance document not available
	//c - Conflict Minerals document not available
	//d - Final burn-in and test compliance document not available
	//e - Country of Origin not compliant

	if event == "RoHs Compliance Certificate Verification" {
		//Check if RoHs Compliance Certificate event has occured
		compIterator, err := stub.GetHistoryForKey(orderID)
		if err != nil {
			return shim.Error("Error 12 " + err.Error())
		}
		defer compIterator.Close()
		var i int
		for i = 0; compIterator.HasNext(); i++ {
			compresponse, err := compIterator.Next()
			if err != nil {
				return shim.Error("Error 13 " + err.Error())
			}
			compOrderObj := order{}
			compOrderBytes := compresponse.Value
			err = json.Unmarshal(compOrderBytes, &compOrderObj)
			if err != nil {
				return shim.Error("Error 14 " + err.Error())
			}
			if compOrderObj.Event == "RoHs Compliance Certificate" && len(compOrderObj.Attachment) != 0 {
				event = "RoHs Compliance Certificate Verified"
				documentType = compOrderObj.DocumentType
				certification = "Obtained"
				attachment = compOrderObj.Attachment
				notification = "Shipment is RoHs compliant"
				attr2 = strconv.Itoa(240)
				attr3 = "Shipment is RoHs compliant"
				attr4 = "a"
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 15 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 16 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 17 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}
				break

			} else {
				event = "RoHs Compliance Certificate not verified"
				notification = "RoHs Compliance Certificate missing, hold payment"
				attr2 = strconv.Itoa(240)
				attr3 = "RoHs Compliance Certificate missing, hold payment"
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 18 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 19 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 20 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}

			}

		}
	}
	if event == "Conflict Minerals Compliance Verification" {
		//Check if Conflict Minerals Compliance event has occured
		compIterator, err := stub.GetHistoryForKey(orderID)
		if err != nil {
			return shim.Error("Error 21 " + err.Error())
		}
		defer compIterator.Close()
		var i int
		for i = 0; compIterator.HasNext(); i++ {
			compresponse, err := compIterator.Next()
			if err != nil {
				return shim.Error("Error 22 " + err.Error())
			}
			compOrderObj := order{}
			compOrderBytes := compresponse.Value
			err = json.Unmarshal(compOrderBytes, &compOrderObj)
			if err != nil {
				return shim.Error("Error 23 " + err.Error())
			}
			if compOrderObj.Event == "Conflict Minerals Compliance" && len(compOrderObj.Attachment) != 0 {
				event = "Conflict Minerals Compliance Certificate Verified"
				documentType = compOrderObj.DocumentType
				certification = "Obtained"
				attachment = compOrderObj.Attachment
				//Check if shipment is RoHs complaint as well
				rohsBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 24 " + err.Error())
				}
				rohsObj := order{}
				err = json.Unmarshal(rohsBytes, &rohsObj)
				if rohsObj.Attribute4 == "a" {
					notification = "Shipment is RoHs and Conflict Minerals compliant"
					attr4 = "b"
				} else {
					notification = "Shipment is Conflict Minerals compliant, but RoHs compliance document missing. Hold Payment"

					attr4 = "c"
				}
				attr2 = strconv.Itoa(250)
				attr3 = "Shipment is Conflict Minerals compliant"
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 25 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 26 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 27 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}
				break
			} else {
				event = "Conflict Minerals Compliance Certificate not verified"
				notification = "Conflict Minerals Compliance Certificate missing, hold payment"
				attr2 = strconv.Itoa(250)
				attr3 = "Conflict Minerals Compliance Certificate missing, hold payment"
				rohsBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 28 " + err.Error())
				}
				rohsObj := order{}
				err = json.Unmarshal(rohsBytes, &rohsObj)
				if rohsObj.Attribute4 == "a" {
					attr4 = "a"
					notification = "Shipment is RoHs compliant, but Conflict Minerals Compliance Certificate missing, hold payment"
				}
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 29 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 30 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 31 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}

			}

		}
	}

	if event == "Final burn-in and Test Certificate Verification" {
		//Check if Final burn-in and Test Certificate event has occured
		compIterator, err := stub.GetHistoryForKey(orderID)
		if err != nil {
			return shim.Error("Error 32 " + err.Error())
		}
		defer compIterator.Close()
		var i int
		for i = 0; compIterator.HasNext(); i++ {
			compresponse, err := compIterator.Next()
			if err != nil {
				return shim.Error("Error 33 " + err.Error())
			}
			compOrderObj := order{}
			compOrderBytes := compresponse.Value
			err = json.Unmarshal(compOrderBytes, &compOrderObj)
			if err != nil {
				return shim.Error("Error 34 " + err.Error())
			}
			if compOrderObj.Event == "Final burn-in and Test Certificate" && len(compOrderObj.Attachment) != 0 {
				event = "Final burn-in and Test Certificate Verified"
				documentType = compOrderObj.DocumentType
				certification = "Obtained"
				attachment = compOrderObj.Attachment
				//Check if shipment is RoHs and Conflict Minerals complaint as well
				rohsBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 35 " + err.Error())
				}
				rohsObj := order{}
				err = json.Unmarshal(rohsBytes, &rohsObj)
				if rohsObj.Attribute4 == "a" {
					notification = "Shipment is RoHs and Final burn-in and Test compliant but Conflict Minerals compliance document missing. Hold Payment"
					attr4 = "g"
				} else if rohsObj.Attribute4 == "b" {
					notification = "Shipment is RoHs, Conflict Minerals and Final burn-in and Test compliant"
					attr4 = "d"
				} else if rohsObj.Attribute4 == "c" {
					notification = "Shipment is Conflict Minerals and Final burn-in and Test compliant, but RoHs compliance document missing. Hold Payment"
					attr4 = "e"
				} else {
					notification = "Shipment is Final burn-in and Test compliant, but RoHs and Conflict Minerals compliance document missing. Hold Payment"

					attr4 = "f"
				}
				attr2 = strconv.Itoa(260)
				attr3 = "Shipment is Final burn-in and Test compliant"
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 36 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 37 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 38 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}
				break
			} else {
				event = "Final burn-in and Test Certificate not verified"
				confObj := order{}
				confBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 39 " + err.Error())
				}
				err = json.Unmarshal(confBytes, &confObj)
				if confObj.Attribute4 == "a" {
					notification = "Shipment is RoHs compliant, but Conflict Minerals & Final burn-in and Test Certificate missing, hold payment"
				} else if confObj.Attribute4 == "b" {
					notification = "Shipment is RoHs & Conflict Minerals compliant but Final burn-in and Test Certificate missing, hold payment"
				} else if confObj.Attribute4 == "c" {
					notification = "Shipment is Conflict Minerals compliant but RoHs and Final burn-in and Test Certificate missing, hold payment"
				} else {
					notification = "Final burn-in and Test Certificate missing, hold payment"
				}
				attr2 = strconv.Itoa(260)
				attr3 = "Final burn-in and Test Certificate missing, hold payment"
				//Check if high vibration event has occured for the order
				inspectOrderBytes, err := stub.GetState(orderID)
				if err != nil {
					return shim.Error("Error 40 " + err.Error())
				}
				inspectOrderObj := order{}
				err = json.Unmarshal(inspectOrderBytes, &inspectOrderObj)
				if err != nil {
					return shim.Error("Error 41 " + err.Error())
				}
				excepCount, err := strconv.Atoi(inspectOrderObj.Attribute1)
				if err != nil {
					return shim.Error("Error 42 " + err.Error())
				}
				if excepCount == 1 {
					otherConcate := []string{notification, "Vibration detected once, Post Inspection Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount == 2 {
					otherConcate := []string{notification, "Vibration detected twice, Payment on Hold Required"}
					notification = strings.Join(otherConcate, "~")
				} else if excepCount >= 3 {
					otherConcate := []string{notification, "Vibration detected thrice, Replacement subassembly required"}
					notification = strings.Join(otherConcate, "~")
				}

			}

		}
	}

	if event == "Order replacement for Control System" {
		conc := []string{"Vibration detected thrice, New Order", taxStr, "created for Control System replacement"}
		attr3 = strings.Join(conc, " ")
		notification = strings.Join(conc, " ")
	}

	//Get default information from Shipping Transactions chaincode
	if event == "Export Compliance Documentation" && attr2 == "200" {
		chaincodeName := "shippingtransactions"
		channelName := "shipping"
		f := "queryOrder"
		inputArgs := util.ToChaincodeArgs(f, args[0])
		response := stub.InvokeChaincode(chaincodeName, inputArgs, channelName)
		if response.Status != shim.OK {
			errStatus := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
			return shim.Error(errStatus)
		}
		respBytes := response.Payload
		respObj := order{}
		err := json.Unmarshal(respBytes, &respObj)
		if err != nil {
			return shim.Error("Error 43 " + err.Error())
		}
		notification = respObj.Notification
		crossCountry = respObj.CrossCountryTransport
		attr1 = respObj.Attribute1
		attr3 = respObj.Attribute3
		attr4 = respObj.Attribute4
		attr5 = respObj.Attribute5
		attr6 = respObj.Attribute6
	}

	//Update an order object
	objectType := "sales order"
	orderObj := &order{objectType, orderID, item, itemDesc, customer, manufacturer, shipper, supplier, quantity, event, expectedDeliveryDate, actualDeliveryDate, exception, documentType, attachment, workOrder, invoice, purchaseOrder, certification, reference, netAmount, unitPrice, charges, discount, tax, owner, custody, currentLoc, countryOfOrigin, destination, maxVib, temperature, notification, crossCountry, serialNum, lotNum, attr1, attr2, attr3, attr4, attr5, attr6, invalidTrx, count}

	//Convert the order object to JSON object
	orderBytes, err = json.Marshal(orderObj)
	if err != nil {
		return shim.Error("Error 44 " + err.Error())
	}

	//Create Write Set
	err = stub.PutState(orderID, orderBytes)
	if err != nil {
		return shim.Error("Error 45 " + err.Error())
	}
	err = stub.SetEvent(event, orderBytes)
	if err != nil {
		return shim.Error("Error 46 " + err.Error())
	}
	//PTR Track and Trace App - Increment the count of SO transactions
	//Check if the order belongs to PTR organizations - if yes, create/update the ledger where count of trx are maintained
	if customer == "Get Well Hospital" || customer == "MedSupply Corp" {
		key := "SOCount"
		var curTrxCount, trxCount int
		countBytes, err := stub.GetState(key)
		if err != nil {
			trxCount = 1
		}
		if countBytes == nil {
			trxCount = 1
		} else {
			curTrxCountStr := string(countBytes)
			curTrxCount, err = strconv.Atoi(curTrxCountStr)
			if err != nil {
				return shim.Error("Error 47 " + err.Error())
			}
			trxCount = curTrxCount + 1
		}
		trxCountStr := strconv.Itoa(trxCount)
		err = stub.PutState(key, []byte(trxCountStr))
		if err != nil {
			return shim.Error("Error 48 " + err.Error())
		}

	}
	//PTR Track and Trace App - Increment the count of SO transactions

	return shim.Success(nil)

}

//=================================================================================================================
//queryTrxHistoryByParentOrder - Function to retrieve the sales order history by finding the SO based on ref number
//=================================================================================================================
func (t *SimpleChainCode) queryTrxHistoryByParentOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	refNo := args[0]
	var buffer bytes.Buffer
	index := "orderIndex"
	orderIterator, err := stub.GetStateByPartialCompositeKey(index, []string{refNo})
	if err != nil {
		return shim.Error("Error 1" + err.Error())
	}
	if orderIterator == nil {
		buffer.WriteString("[]")
		return shim.Success(buffer.Bytes())
	}
	defer orderIterator.Close()
	isRecWritten := false
	var i int

	buffer.WriteString("[")
	for i = 0; orderIterator.HasNext(); i++ {
		orderResp, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		compKey := orderResp.Key
		compInd, compVal, err := stub.SplitCompositeKey(compKey)
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		fmt.Println(compInd)
		childOrder := compVal[1]
		childOrderTrxIterator, err := stub.GetHistoryForKey(childOrder)
		if err != nil {
			return shim.Error("Error 4 " + err.Error())
		}

		var j int
		for j = 0; childOrderTrxIterator.HasNext(); j++ {
			orderResp, err := childOrderTrxIterator.Next()
			if err != nil {
				return shim.Error("Error 5 " + err.Error())
			}
			if isRecWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Order Type\":\"Sales Order\",")
			buffer.WriteString("\"Order ID\":\"")
			buffer.WriteString(childOrder)
			buffer.WriteString("\",")
			buffer.WriteString("\"Transaction ID\":\"")
			buffer.WriteString(orderResp.TxId)
			buffer.WriteString("\",")
			buffer.WriteString("\"Value\":")
			buffer.Write(orderResp.Value)
			buffer.WriteString(",")
			buffer.WriteString("\"TimeStamp\":\"")
			buffer.WriteString(time.Unix(orderResp.Timestamp.Seconds, int64(orderResp.Timestamp.Nanos)).String())
			buffer.WriteString("\"")
			buffer.WriteString("}")
			isRecWritten = true
		}
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())

}

//======================================================================================
//queryAllChildOrders - Return all the child sales orders for the parent purchase orders
//======================================================================================

func (t *SimpleChainCode) queryAllChildOrders(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	refNo := args[0]
	index := "orderIndex"
	orderIterator, err := stub.GetStateByPartialCompositeKey(index, []string{refNo})
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	if orderIterator == nil {
		return shim.Error("No reference orders found")
	}
	defer orderIterator.Close()
	isRecWritten := false
	var i int
	var buffer bytes.Buffer
	for i = 0; orderIterator.HasNext(); i++ {
		orderResp, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		compKey := orderResp.Key
		compIndex, compRefID, err := stub.SplitCompositeKey(compKey)
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		fmt.Println(compIndex)
		compID := compRefID[1]
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(compID)
		isRecWritten = true
	}

	return shim.Success(buffer.Bytes())
}

//getSOWOLink - Function to get the Sales Order number based on the number which is a common identifier between sales order and work order
func (t *SimpleChainCode) getSOWOLink(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	woLinkID := args[0]
	indexName := "woLink"
	woLinkIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{woLinkID})
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	if woLinkIterator == nil {
		return shim.Error("No Sales Orders with this reference ID exists in the ledger")
	}
	defer woLinkIterator.Close()
	var i int
	var buffer bytes.Buffer

	for i = 0; woLinkIterator.HasNext(); i++ {
		woLinkBytes, err := woLinkIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		woLinkKey := woLinkBytes.Key
		woLinkIndex, woLinkkeyArray, err := stub.SplitCompositeKey(woLinkKey)
		fmt.Println(woLinkIndex)
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		woLinkSO := woLinkkeyArray[1]
		buffer.WriteString(woLinkSO)
	}
	return shim.Success(buffer.Bytes())
}

//==========================================================
//queryOrderByLatestOrder - Query Sales Order Latest OrderID
func (t *SimpleChainCode) queryOrderByLatestOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]
	if key == "" {
		return shim.Error("Argument cannot be null")
	}
	orderchannel := "orderprocessing"
	orderchaincode := "latestorders"
	orderfunction := "queryOrder"
	orderargs := util.ToChaincodeArgs(orderfunction, key)
	orderResp := stub.InvokeChaincode(orderchaincode, orderargs, orderchannel)
	if orderResp.Status != shim.OK {
		return shim.Error("Order Info does not exist in the system")
	}
	respBytes := orderResp.Payload
	respObj := &OrderInfo{}
	err := json.Unmarshal(respBytes, respObj)
	if err != nil {
		return shim.Error("Order Info cannot be retrieved " + err.Error())
	}
	orderNumber := respObj.OrderVal
	orderBytes, err := stub.GetState(orderNumber)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	if orderBytes == nil {
		return shim.Error("Order Info does not exist")
	}
	return shim.Success(orderBytes)
}

//=========================================================================================================
//getSalesOrderCount - Function to return the number of sales orders created to the PTR Track and Trace App
//=========================================================================================================

func (t *SimpleChainCode) getSalesOrderCount(stub shim.ChaincodeStubInterface) pb.Response {
	orderCodeDSO := "DistributorSO"
	orderCodeMSO := "ManufacturerSO"
	var countStr string
	var countDSO, countMSO int

	//Get the latest order count for Distributor SO
	orderchannel := "orderprocessing"
	orderchaincode := "latestorders"
	orderfunction := "queryOrder"
	orderargs := util.ToChaincodeArgs(orderfunction, orderCodeDSO)
	orderResp := stub.InvokeChaincode(orderchaincode, orderargs, orderchannel)
	if orderResp.Status != shim.OK {
		countStr = strconv.Itoa(0)
	} else {
		respBytes := orderResp.Payload
		respObj := &OrderInfo{}
		err := json.Unmarshal(respBytes, respObj)
		if err != nil {
			return shim.Error("Order Info cannot be retrieved " + err.Error())
		}
		orderNumber := respObj.OrderVal
		orderBytes, err := stub.GetState(orderNumber)
		if err != nil {
			return shim.Error("Error 1 " + err.Error())
		}
		if orderBytes == nil {
			return shim.Error("Order Info does not exist")
		}
		orderObj := &order{}
		err = json.Unmarshal(orderBytes, orderObj)
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		countDSO = orderObj.Count

		//Get the latest order count for Manufacturer SO
		orderchannel := "orderprocessing"
		orderchaincode := "latestorders"
		orderfunction := "queryOrder"
		orderargs := util.ToChaincodeArgs(orderfunction, orderCodeMSO)
		orderResp := stub.InvokeChaincode(orderchaincode, orderargs, orderchannel)
		if orderResp.Status != shim.OK {
			countStr = strconv.Itoa(countDSO)
		} else {
			respBytes := orderResp.Payload
			respObj := &OrderInfo{}
			err := json.Unmarshal(respBytes, respObj)
			if err != nil {
				return shim.Error("Order Info cannot be retrieved " + err.Error())
			}
			orderNumber := respObj.OrderVal
			orderBytes, err := stub.GetState(orderNumber)
			if err != nil {
				return shim.Error("Error 3 " + err.Error())
			}
			if orderBytes == nil {
				return shim.Error("Order Info does not exist")
			}
			orderObj := &order{}
			err = json.Unmarshal(orderBytes, orderObj)
			if err != nil {
				return shim.Error("Error 4 " + err.Error())
			}
			countMSO = orderObj.Count
			if countDSO > countMSO {
				countStr = strconv.Itoa(countDSO)
			} else {
				countStr = strconv.Itoa(countMSO)
			}
		}
	}
	return shim.Success([]byte(countStr))
}

//=======================================================================================
//getTrxCount - Function to return the number of transactions for PTR Track and Trace app
//=======================================================================================

func (t *SimpleChainCode) getTrxCount(stub shim.ChaincodeStubInterface) pb.Response {
	key := "SOCount"
	countBytes, err := stub.GetState(key)
	if err != nil {
		count := 0
		countStr := strconv.Itoa(count)
		countBytes = []byte(countStr)
	}
	if countBytes == nil {
		count := 0
		countStr := strconv.Itoa(count)
		countBytes = []byte(countStr)
	}
	return shim.Success(countBytes)
}

//===============================================================================================
//queryLatestStateByRef - Function to get the latest state of the order based on reference number
//===============================================================================================

func (t *SimpleChainCode) queryLatestStateByRef(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	parentOrderID := args[0]
	indexName := "orderIndex"
	//Get the list of orders which have the corresponding parent Order ID
	orderIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{parentOrderID})
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	} else if orderIterator == nil {
		return shim.Error("Invalid Order Reference ID " + parentOrderID)
	}
	defer orderIterator.Close()
	//Write individual transactions into a buffer which will be sent as output
	var i int
	var returnOrderID string
	for i = 0; orderIterator.HasNext(); i++ {
		response, err := orderIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		returnIndex, returnKeys, err := stub.SplitCompositeKey(response.Key)
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		fmt.Println(returnIndex)
		//returnParentOrderID := returnKeys[0]
		returnOrderID = returnKeys[1]
	}
	orderBytes, err := stub.GetState(returnOrderID)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}
	return shim.Success(orderBytes)

}
