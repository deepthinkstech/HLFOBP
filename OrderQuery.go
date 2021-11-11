// MasterQuery project MasterQuery.go
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
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct{}

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
type WorkOrder struct {
	OrderID              string    `json:"orderid"`
	OrderNumber          string    `json:"ordernumber"`
	Status               string    `json:"status"`
	ReleaseDate          time.Time `json:"releasedate"`
	WorkDefinitionID     string    `json:"workdefinitionid"`
	Organization         string    `json:"organization"`
	ItemID               string    `json:"item"`
	CompletedQuantity    string    `json:"completequantity"`
	SerialControlled     string    `json:"serialcontrolled"`
	ActualCompletionDate time.Time `json:"actualcompletiondate"`
	Quantity             int       `json:"quantity"`
	ReferenceNo          string    `json:"referenceno"`
	SerialNo             string    `json:"serialno"`
	QcResult             string    `json:"qcresult"`
	Event                string    `json:"event"`
	EventCode            string    `json:"eventcode"`
	ItemDesc             string    `json:"itemDescription"`
	Product              string    `json:"product"`
	//ParentOrder          string    `json:"parentorder"`
	Notification string `json:"notification"`
	Owner        string `json:"owner"`
	Custody      string `json:"custody"`
	Count        int    `json:"count"`
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
	funct, args := stub.GetFunctionAndParameters()
	if funct == "queryOrder" {
		if args[0] == "orderprocessing" {
			if args[1] == "SO" {
				return t.querySO(stub, args)
			} else if args[1] == "PO" {
				return t.queryPO(stub, args)
			} else {
				return shim.Error("Incorrect order type " + args[1])
			}
		} else {
			return shim.Error("Incorrect Channel name " + args[0])
		}
	} else if funct == "queryChildOrders" {
		if args[0] == "orderprocessing" {
			if args[1] == "SO" {
				return t.queryChildSO(stub, args)
			} else if args[1] == "PO" {
				return t.queryChildPO(stub, args)
			} else {
				return shim.Error("Incorrect order type " + args[1])
			}
		} else {
			return shim.Error("Incorrect Channel name " + args[0])
		}
	} else if funct == "getLatestOrderStatus" {
		return t.getLatestOrderStatus(stub, args)
	} else {
		return shim.Error("Incorrect function name " + funct)
	}

}

//===============================================
//querySO - Function to query Sales Order details
//===============================================
func (t *SimpleChaincode) querySO(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	queryChannel := args[0]
	//queryPar := args[1]
	var buffer bytes.Buffer

	orderID := args[2]

	//1. Return the transaction history for the order queried

	//Check if the order exists in the world state DB
	orderIterator, err := stub.GetHistoryForKey(orderID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if orderIterator == nil {
		return shim.Error("Order ID " + orderID + "does not exist in the system")
	}
	defer orderIterator.Close()

	//Write the Sales Order Trx History
	//buffer.WriteString("[")
	/*
		var isRecWritten bool
		if len(args) == 4 {
			isRecWritten = true
		} else {
			isRecWritten = false
		}
	*/
	isRecWritten := false
	soChannel := queryChannel
	soFunction := "queryTrxHistoryV2"
	soChaincode := "salestransactions"
	soArgs := util.ToChaincodeArgs(soFunction, orderID)
	soResp := stub.InvokeChaincode(soChaincode, soArgs, soChannel)
	if soResp.Status == shim.OK {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(soResp.Payload)
		isRecWritten = true

	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}
	//2. Return the transaction history for the sales orders with the queried sales order as reference
	//Check for information on Sales Orders with the queried sales order number as reference
	soChannel = queryChannel
	soFunction = "queryTrxHistoryByParentOrder"
	soChaincode = "salestransactions"
	soArgs = util.ToChaincodeArgs(soFunction, orderID)
	soResp = stub.InvokeChaincode(soChaincode, soArgs, soChannel)
	if soResp.Status == shim.OK {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(soResp.Payload)
		isRecWritten = true

	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}

	//3. Return the purchase orders & their dependents transaction history with the queried sales order as reference
	//Check for information on Purchase Orders with the sales order number as reference
	poChannel := queryChannel
	poFunction := "queryAllChildOrders"
	poChaincode := "purchaseordertransactions"
	poArgs := util.ToChaincodeArgs(poFunction, orderID)
	poResp := stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	if poResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := poResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {

					poArray := []string{queryChannel, "PO", childOrders[j], "true"}
					poResponse := t.queryPO(stub, poArray)
					if poResponse.Status == shim.OK {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.Write(poResponse.Payload)
						isRecWritten = true
					} else {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.WriteString("[]")
					}

					//buffer.WriteString(",")
					//buffer.WriteString(childOrders[j])

				}
			} else {
				poArray := []string{queryChannel, "PO", respString.String(), "true"}
				poResponse := t.queryPO(stub, poArray)
				if poResponse.Status == shim.OK {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.Write(poResponse.Payload)
					isRecWritten = true
				} else {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.WriteString("[]")
				}
			}
		} else {
			if isRecWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("[]")
		}

	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}
	//buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//==================================================
//queryPO - Function to query Purchase Order details
//==================================================
func (t *SimpleChaincode) queryPO(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryChannel := args[0]
	//queryPar := args[1]
	var buffer bytes.Buffer

	orderID := args[2]

	//1. Return the transaction history for the order queried

	//Check if the order exists in the world state DB
	orderIterator, err := stub.GetHistoryForKey(orderID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if orderIterator == nil {
		return shim.Error("Order ID " + orderID + "does not exist in the system")
	}
	defer orderIterator.Close()

	//Write the Purchase Order Trx History
	//buffer.WriteString("[")
	/*
		var isRecWritten bool
		if len(args) == 4 {
			isRecWritten = true
		} else {
			isRecWritten = false
		}
	*/
	isRecWritten := false
	poChannel := queryChannel
	poFunction := "queryTrxHistory"
	poChaincode := "purchaseordertransactions"
	poArgs := util.ToChaincodeArgs(poFunction, orderID)
	poResp := stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	if poResp.Status == shim.OK {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(poResp.Payload)
		isRecWritten = true

	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}
	//2. Return the transaction history for the purchase orders with the queried purchase order as reference
	//Check for information on Sales Orders with the queried purchase order number as reference
	poChannel = queryChannel
	poFunction = "queryTrxHistoryByParentOrder"
	poChaincode = "purchaseordertransactions"
	poArgs = util.ToChaincodeArgs(poFunction, orderID)
	poResp = stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	if poResp.Status == shim.OK {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(poResp.Payload)
		isRecWritten = true

	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}

	//3. Return the purchase orders & their dependents transaction history with the queried sales order as reference
	//Check for information on Purchase Orders with the sales order number as reference
	soChannel := queryChannel
	soFunction := "queryAllChildOrders"
	soChaincode := "salestransactions"
	soArgs := util.ToChaincodeArgs(soFunction, orderID)
	soResp := stub.InvokeChaincode(soChaincode, soArgs, soChannel)
	if soResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := soResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {
					soArray := []string{queryChannel, "SO", childOrders[j], "true"}
					soResponse := t.querySO(stub, soArray)
					if soResponse.Status == shim.OK {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.Write(soResponse.Payload)
						isRecWritten = true
					} else {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.WriteString("[]")
					}
				}
			} else {
				soArray := []string{queryChannel, "SO", respString.String(), "true"}
				soResponse := t.querySO(stub, soArray)
				if soResponse.Status == shim.OK {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.Write(soResponse.Payload)
					isRecWritten = true
				} else {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.WriteString("[]")
				}
			}
		} else {
			if isRecWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("[]")
		}
	} else {
		if isRecWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("[]")
	}
	//buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//===============================================
//querySO - Function to query Sales Order details
//===============================================
func (t *SimpleChaincode) queryChildSO(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	queryChannel := args[0]
	//queryPar := args[1]
	var buffer bytes.Buffer

	orderID := args[2]

	//1. Return the transaction history for the order queried

	//Check if the order exists in the world state DB

	//1.Write the Sales Order into array
	//buffer.WriteString("[")
	var isRecWritten bool
	if len(args) == 4 {
		isRecWritten = true
	} else {
		isRecWritten = false
	}
	buffer.WriteString(orderID)
	isRecWritten = true

	//2. Return the child orders for the sales orders with the queried sales order as reference
	//Check for information on Sales Orders with the queried sales order number as reference
	soChannel := queryChannel
	soFunction := "queryAllChildOrders"
	soChaincode := "salestransactions"
	soArgs := util.ToChaincodeArgs(soFunction, orderID)
	soResp := stub.InvokeChaincode(soChaincode, soArgs, soChannel)
	if soResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := soResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {
					if isRecWritten == true {
						buffer.WriteString(",")
					}

					buffer.WriteString(childOrders[j])
					isRecWritten = true
				}
			} else {
				if isRecWritten == true {
					buffer.WriteString(",")
				}

				buffer.WriteString(respString.String())
				isRecWritten = true

			}
		}
	}

	//3. Return the purchase orders & their dependents transaction history with the queried sales order as reference
	//Check for information on Purchase Orders with the sales order number as reference
	poChannel := queryChannel
	poFunction := "queryAllChildOrders"
	poChaincode := "purchaseordertransactions"
	poArgs := util.ToChaincodeArgs(poFunction, orderID)
	poResp := stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	if poResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := poResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {

					poArray := []string{queryChannel, "PO", childOrders[j], "true"}
					poResponse := t.queryChildPO(stub, poArray)
					if poResponse.Status == shim.OK {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.Write(poResponse.Payload)
						isRecWritten = true
					}

					//buffer.WriteString(",")
					//buffer.WriteString(childOrders[j])

				}
			} else {
				poArray := []string{queryChannel, "PO", respString.String(), "true"}
				poResponse := t.queryChildPO(stub, poArray)
				if poResponse.Status == shim.OK {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.Write(poResponse.Payload)
					isRecWritten = true
				}
			}
		}

	} else {
		buffer.WriteString("No child orders")
	}
	//buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//==================================================
//queryPO - Function to query Purchase Order details
//==================================================
func (t *SimpleChaincode) queryChildPO(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	queryChannel := args[0]
	//queryPar := args[1]
	var buffer bytes.Buffer

	orderID := args[2]

	//1. Return the transaction history for the order queried

	//Check if the order exists in the world state DB

	//buffer.WriteString("[")
	var isRecWritten bool
	if len(args) == 4 {
		isRecWritten = true
	} else {
		isRecWritten = false
	}
	buffer.WriteString(orderID)
	isRecWritten = true

	//2. Return the child purchase orders for the purchase orders with the queried purchase order as reference
	//Check for information on Purchase Orders with the queried purchase order number as reference
	poChannel := queryChannel
	poFunction := "queryAllChildOrders"
	poChaincode := "purchaseordertransactions"
	poArgs := util.ToChaincodeArgs(poFunction, orderID)
	poResp := stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	if poResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := poResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {
					if isRecWritten == true {
						buffer.WriteString(",")
					}

					buffer.WriteString(childOrders[j])
					isRecWritten = true
				}
			} else {
				if isRecWritten == true {
					buffer.WriteString(",")
				}

				buffer.WriteString(respString.String())
				isRecWritten = true

			}
		}
	}

	//3. Return the sales orders & their dependents  with the queried purchase order as reference

	soChannel := queryChannel
	soFunction := "queryAllChildOrders"
	soChaincode := "salestransactions"
	soArgs := util.ToChaincodeArgs(soFunction, orderID)
	soResp := stub.InvokeChaincode(soChaincode, soArgs, soChannel)
	if soResp.Status == shim.OK {
		//buffer.Write(poResp.Payload)
		//isRecWritten = true
		respBytes := soResp.Payload
		var respString bytes.Buffer
		var i int
		for i = 0; i < len(respBytes); i++ {
			respString.WriteByte(respBytes[i])
		}
		if len(respString.String()) != 0 {
			if strings.Contains(respString.String(), ",") {
				childOrders := strings.Split(respString.String(), ",")
				var j int
				for j = 0; j < len(childOrders); j++ {
					soArray := []string{queryChannel, "SO", childOrders[j], "true"}
					soResponse := t.queryChildSO(stub, soArray)
					if soResponse.Status == shim.OK {
						if isRecWritten == true {
							buffer.WriteString(",")
						}
						buffer.Write(soResponse.Payload)
						isRecWritten = true
					}
				}
			} else {
				soArray := []string{queryChannel, "SO", respString.String(), "true"}
				soResponse := t.queryChildSO(stub, soArray)
				if soResponse.Status == shim.OK {
					if isRecWritten == true {
						buffer.WriteString(",")
					}
					buffer.Write(soResponse.Payload)
					isRecWritten = true
				}
			}
		}
	}
	//buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=========================================================================================
//getLatestOrderStatus - Function to return the latest order status for the DSCB V3 Chatbot
//=========================================================================================
func (t *SimpleChaincode) getLatestOrderStatus(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	orderID := args[0]
	var soEvent, soEventCode, poEvent, poEventCode, woEvent, woEventCode, shipEvent, shipEventCode, csoEvent, csoEventCode, pShipEvent, pShipEventCode string
	var err error
	//Check if purchase order ledger has entries for the order
	poChannel := "orderprocessing"
	poFunction := "queryOrder"
	poChaincode := "purchaseordertransactions"
	poArgs := util.ToChaincodeArgs(poFunction, orderID)
	poResp := stub.InvokeChaincode(poChaincode, poArgs, poChannel)
	//If the order is not purchase order, then check if the order is Sales Order
	if poResp.Status != shim.OK {
		soChannel := "orderprocessing"
		soFunction := "queryOrder"
		soChaincode := "salestransactions"
		soArgs := util.ToChaincodeArgs(soFunction, orderID)
		soResp := stub.InvokeChaincode(soChaincode, soArgs, soChannel)
		if soResp.Status != shim.OK {
			return shim.Error("Order ID " + orderID + " is invalid. The records for the order does not exist in the system.")
		} else {
			soBytes := soResp.Payload
			soObj := &order{}
			err = json.Unmarshal(soBytes, soObj)
			if err != nil {
				return shim.Error("Error 1 " + err.Error())
			}
			soEventCode = soObj.Attribute2
			soEvent = soObj.Event
			poEventCode = "0"
			poEvent = "No Event"

		}
		//Check if the order is a shipment
		pShipChannel := "shipping"
		pShipFunction := "queryOrder"
		pShipChaincode := "shippingtransactions"
		pShipArgs := util.ToChaincodeArgs(pShipFunction, orderID)
		pShipResp := stub.InvokeChaincode(pShipChaincode, pShipArgs, pShipChannel)
		if pShipResp.Status != shim.OK {
			pShipEventCode = "0"
			pShipEvent = "No Event"
		} else {
			pShipBytes := pShipResp.Payload
			pShipObj := &order{}
			err = json.Unmarshal(pShipBytes, pShipObj)
			if err != nil {
				return shim.Error("Error 2 " + err.Error())
			}
			pShipEventCode = pShipObj.Attribute2
			pShipEvent = pShipObj.Event
			poEventCode = "0"
			poEvent = "No Event"
		}

	} else {
		poBytes := poResp.Payload
		poObj := &PurchaseOrder{}
		err = json.Unmarshal(poBytes, poObj)
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		poEventCode = poObj.EventCode
		poEvent = poObj.Event
		soEventCode = "0"
		soEvent = "No Event"
		pShipEventCode = "0"
		pShipEvent = "No Event"
	}

	//Find the sales orders with the order ID as reference
	csoChannel := "orderprocessing"
	csoFunction := "queryLatestStateByRef"
	csoChaincode := "salestransactions"
	csoArgs := util.ToChaincodeArgs(csoFunction, orderID)
	csoResp := stub.InvokeChaincode(csoChaincode, csoArgs, csoChannel)
	if csoResp.Status != shim.OK {
		csoEventCode = "0"
		csoEvent = "No Event"

	} else if csoResp.Payload != nil {
		csoBytes := csoResp.Payload
		csoObj := &order{}
		err = json.Unmarshal(csoBytes, csoObj)
		if err != nil {
			return shim.Error("Error 5 " + err.Error())
		}
		csoEventCode = csoObj.Attribute2
		csoEvent = csoObj.Event

	} else {
		csoEventCode = "0"
		csoEvent = "No Event"
	}

	//Find the work orders with the order ID as reference
	woChannel := "spmanufacturing"
	woFunction := "queryLatestStateByRef"
	woChaincode := "workordertransactions"
	woArgs := util.ToChaincodeArgs(woFunction, orderID)
	woResp := stub.InvokeChaincode(woChaincode, woArgs, woChannel)
	if woResp.Status != shim.OK {
		woEventCode = "0"
		woEvent = "No Event"

	} else if woResp.Payload == nil {
		woEventCode = "0"
		woEvent = "No Event"
	} else {
		woBytes := woResp.Payload
		woObj := &WorkOrder{}
		err = json.Unmarshal(woBytes, woObj)
		if err != nil {
			return shim.Error("Error 6 " + err.Error())
		}
		woEventCode = woObj.EventCode
		woEvent = woObj.Event

	}
	//Find the shipments with the order ID as reference
	shipChannel := "shipping"
	shipFunction := "queryLatestStateByRef"
	shipChaincode := "shippingtransactions"
	shipArgs := util.ToChaincodeArgs(shipFunction, orderID)
	shipResp := stub.InvokeChaincode(shipChaincode, shipArgs, shipChannel)
	if shipResp.Status != shim.OK {
		shipEventCode = "0"
		shipEvent = "No Event"

	} else if shipResp.Payload == nil {
		shipEventCode = "0"
		shipEvent = "No Event"
	} else {
		shipBytes := shipResp.Payload
		shipObj := &order{}
		err = json.Unmarshal(shipBytes, shipObj)
		if err != nil {
			return shim.Error("Error 7 " + err.Error())
		}
		shipEventCode = shipObj.Attribute2
		shipEvent = shipObj.Event

	}
	//Check which is the latest of the statuses
	soEventInt, err := strconv.Atoi(soEventCode)
	if err != nil {
		return shim.Error("Error 8 " + err.Error())
	}
	poEventInt, err := strconv.Atoi(poEventCode)
	if err != nil {
		return shim.Error("Error 9 " + err.Error())
	}
	csoEventInt, err := strconv.Atoi(csoEventCode)
	if err != nil {
		return shim.Error("Error 10 " + err.Error())
	}
	woEventInt, err := strconv.Atoi(woEventCode)
	if err != nil {
		return shim.Error("Error 11 " + err.Error())
	}
	shipEventInt, err := strconv.Atoi(shipEventCode)
	if err != nil {
		return shim.Error("Error 12 " + err.Error())
	}
	pShipEventInt, err := strconv.Atoi(pShipEventCode)
	if err != nil {
		return shim.Error("Error 13 " + err.Error())
	}

	var event string
	var eventCode string

	//Check if Sales Order event is the latest of all events
	if soEventInt > poEventInt {
		event = soEvent
		eventCode = strconv.Itoa(soEventInt)
	}
	if soEventInt > woEventInt {
		event = soEvent
		eventCode = strconv.Itoa(soEventInt)
	}
	if soEventInt > shipEventInt {
		event = soEvent
		eventCode = strconv.Itoa(soEventInt)
	}
	if soEventInt > csoEventInt {
		event = soEvent
		eventCode = strconv.Itoa(soEventInt)
	}
	if soEventInt > pShipEventInt {
		event = soEvent
		eventCode = strconv.Itoa(soEventInt)
	}

	//Check if Purchase Order event is the latest of all events
	if poEventInt > soEventInt {
		event = poEvent
		eventCode = strconv.Itoa(poEventInt)
	}
	if poEventInt > woEventInt {
		event = poEvent
		eventCode = strconv.Itoa(poEventInt)
	}
	if poEventInt > shipEventInt {
		event = poEvent
		eventCode = strconv.Itoa(poEventInt)
	}
	if poEventInt > csoEventInt {
		event = poEvent
		eventCode = strconv.Itoa(poEventInt)
	}
	if poEventInt > pShipEventInt {
		event = poEvent
		eventCode = strconv.Itoa(poEventInt)
	}

	//Check if the Child Sales Order event is the latest of all events
	if csoEventInt > soEventInt {
		event = csoEvent
		eventCode = strconv.Itoa(csoEventInt)
	}
	if csoEventInt > poEventInt {
		event = csoEvent
		eventCode = strconv.Itoa(csoEventInt)
	}
	if csoEventInt > woEventInt {
		event = csoEvent
		eventCode = strconv.Itoa(csoEventInt)
	}
	if csoEventInt > shipEventInt {
		event = csoEvent
		eventCode = strconv.Itoa(csoEventInt)
	}
	if csoEventInt > pShipEventInt {
		event = csoEvent
		eventCode = strconv.Itoa(csoEventInt)
	}

	//Check if Work Order event is the latest of all events
	if woEventInt > poEventInt {
		event = woEvent
		eventCode = strconv.Itoa(woEventInt)
	}
	if woEventInt > soEventInt {
		event = woEvent
		eventCode = strconv.Itoa(woEventInt)
	}
	if woEventInt > shipEventInt {
		event = woEvent
		eventCode = strconv.Itoa(woEventInt)
	}
	if woEventInt > csoEventInt {
		event = woEvent
		eventCode = strconv.Itoa(woEventInt)
	}
	if woEventInt > pShipEventInt {
		event = woEvent
		eventCode = strconv.Itoa(woEventInt)
	}

	//Check if Shipment event is the latest of all events
	if shipEventInt > poEventInt {
		event = shipEvent
		eventCode = strconv.Itoa(shipEventInt)
	}
	if shipEventInt > woEventInt {
		event = shipEvent
		eventCode = strconv.Itoa(shipEventInt)
	}
	if shipEventInt > soEventInt {
		event = shipEvent
		eventCode = strconv.Itoa(shipEventInt)
	}
	if shipEventInt > csoEventInt {
		event = shipEvent
		eventCode = strconv.Itoa(shipEventInt)
	}
	if shipEventInt > pShipEventInt {
		event = shipEvent
		eventCode = strconv.Itoa(shipEventInt)
	}

	if pShipEventInt > poEventInt {
		event = pShipEvent
		eventCode = strconv.Itoa(pShipEventInt)
	}
	if pShipEventInt > woEventInt {
		event = pShipEvent
		eventCode = strconv.Itoa(pShipEventInt)
	}
	if pShipEventInt > soEventInt {
		event = pShipEvent
		eventCode = strconv.Itoa(pShipEventInt)
	}
	if pShipEventInt > csoEventInt {
		event = pShipEvent
		eventCode = strconv.Itoa(pShipEventInt)
	}
	if pShipEventInt > shipEventInt {
		event = pShipEvent
		eventCode = strconv.Itoa(pShipEventInt)
	}

	//Prepare appropriate response
	var response string
	if event == "Purchase Order Release" || event == "Order Received" || event == "Work Order Created" || event == "Purchase Order and Sales Order quantity matching" || event == "Work Order Created" {
		response = "Your order " + orderID + " has been placed"
	} else if event == "Work Order Complete" || event == "RoHs Compliance Certificate" || event == "Conflict Minerals Compliance" || event == "Final burn-in and Test Certificate" || event == "Country of Origin Certificate" {
		response = "Your order " + orderID + " is getting ready to be shipped along with the requisite compliance certificates"
	} else if event == "Shipment Executed" || event == "Invoice Generated" || event == "Invoice Documentation" {
		response = "Shipment of your order " + orderID + " is just initiated."
	} else if event == "Shipment En-route to Airport by Road – Started" {
		response = "Your order " + orderID + " is en-route to airport by road."
	} else if event == "Shipment En-route by Road - Arrived" {
		response = "Your order " + orderID + " is getting ready for air-freight."
	} else if event == "Shipment En-route by Air - Started" {
		response = "Your order " + orderID + " is en-route by air."
	} else if event == "Shipment En-route by Air - Arrived" {
		response = "Your order " + orderID + " has arrived at destination airport."
	} else if event == "Customs Clearance – Completed" {
		response = "Your order " + orderID + " has cleared customs and is ready for transport to final destination."
	} else if event == "Customs Clearance – Initiated" {
		response = "Your order " + orderID + " has arrived at destination airport and is awaiting customs clearance."
	} else if event == "Shipment En-route by Road - Started" {
		response = "Your order " + orderID + " is near your location and is expected to be delivered shortly."
	} else if event == "Shipment Reached Destination" || event == " Purchase Price Updated" {
		response = "Your order " + orderID + " has reached destination."
	} else if event == "Equipment Installation – In Progress" {
		response = "Your order " + orderID + " has been delivered and installation is in progress."
	} else if event == "Equipment Installation – Completed" {
		response = "The installation of your order " + orderID + " has been completed."
	} else if event == "Customer Accepted" {
		response = "Customer has accepted the order " + orderID + "."
	} else if event == " RoHs Compliance Certificate Verified" || event == "Conflict Minerals Compliance Verified" || event == "Final burn-in and Test Certificate Verified" {
		response = "Customer has accepted the order " + orderID + ", and certificates are being verified."
	} else if event == "Export Compliance Documentation" {
		if eventCode == "70" {
			response = "Shipment of your order " + orderID + " is just initiated."
		} else if eventCode == "110" {
			response = "Your order " + orderID + " is getting ready to be air-lifted along with the requisite documents."
		} else if eventCode == "140" {
			response = "Your order " + orderID + " has arrived at destination airport and documentation is in progress."
		} else if eventCode == "170" {
			response = "Your order " + orderID + " is getting ready to be transported by road with the requisite documents."
		} else {
			response = "Your order" + orderID + "has requisite documents"
		}
	} else if event == "Payment Terms Updated" {
		response = "Your order " + orderID + "   had experienced 2 vibrations."
	} else if strings.ToLower(event) == "order for replacement part approved" {
		response = "Your order " + orderID + "   had experienced 3 vibrations."
	} else {
		response = "Customer has verified all certificates and accepted the order " + orderID + "."
	}

	return shim.Success([]byte(response))
}
