// FITTrips project FITTrips.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//================================================
//Define structures for implementing the functions
//================================================

type SimpleChaincode struct{}
type revenueShare struct {
	TripID          string  `json:"tripID"`
	FulfilledBy     string  `json:"fulfilledBy"`
	Stakeholder1    string  `json:"stakeholder1"`
	Stakeholder2    string  `json:"stakeholder2"`
	TripRevenue     float64 `json:"tripRevenue"`
	RevenueCurrency string  `json:"revenueCurrency"`
	Stakeholder1Rev float64 `json:"stakeholder1rev"`
	Stakeholder2Rev float64 `json:"stakeholder2rev"`
}
type tripDetails struct {
	ObjectType               string         `json:"objectType"`
	TktID                    string         `json:"tktID"`
	FromLOC                  string         `json:"fromLOC"`
	ToLOC                    string         `json:"toLOC"`
	RiderID                  string         `json:"riderID"`
	Products                 string         `json:"products"`
	Price                    float64        `json:"price"`
	Currency                 string         `json:"currency"`
	Progress                 string         `json:"progress"`
	Status                   string         `json:"status"`
	HasBookedUsingSegment    bool           `json:"hasBookesUsingSegment"`
	HasFulfilledUsingSegment bool           `json:"hasFulfilledUsingSegment"`
	BookedUsingProduct       string         `json:"bookedUsingProduct"`
	Event                    string         `json:"event"`
	SeqNo                    int            `json:"seqNo"`
	Duration                 string         `json:"duration"`
	CreationDate             time.Time      `json:"creationDate"`
	Geography                string         `json:"geography"`
	RevenueSharing           []revenueShare `json:"revenuesharing"`
	RevenueEarned            float64        `json:"revenueEarned"`
}
type KPI struct {
	NoOfOrgs  int     `json:"noOfOrgs"`
	Revenue   float64 `json:"revenue"`
	Cost      float64 `json:"cost"`
	NoOfTrips int     `json:"noOfTrips"`
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
	if function == "buyTkt" {
		return t.createTrip(stub, args)
	} else if function == "updateTrip" {
		return t.updateTrip(stub, args)
	} else if function == "queryTrip" {
		return t.queryTrip(stub, args)
	} else if function == "getKPIData" {
		return t.getKPIData(stub, args)
	} else if function == "queryTripHistory" {
		return t.queryTripHistory(stub, args)
	} else if function == "queryTripsByStakeholder" {
		return t.queryTripsByStakeholder(stub, args)
	} else if function == "queryTripsByRider" {
		return t.queryTripsByRider(stub, args)
	} else {
		return shim.Error("Invalid function name " + function)
	}
	return shim.Success(nil)
}

//============================================
//createTrip - Function to create Trip booking
//============================================
func (t *SimpleChaincode) createTrip(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fromLOC := args[0]
	toLOC := args[1]
	riderID := args[2]
	products := args[3]
	price, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		price = 0
	}
	currency := args[5]
	progress := args[6]
	status := args[7]
	hasBookedUsingSegment, err := strconv.ParseBool(args[8])
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	hasFulfilledUsingSegment, err := strconv.ParseBool(args[9])
	if err != nil {
		return shim.Error("Error 2 " + err.Error())
	}
	bookedUsingProduct := args[10]
	event := args[11]
	tktID := args[12]
	seqNo, err := strconv.Atoi(args[13])
	if err != nil {
		return shim.Error("Error 8 " + err.Error())
	}
	duration := args[14]
	//Determine the new trip ID
	//Last trip ID is stored as separate key in the ledger. If the key exists, increment the value
	if tktID == "eTkt100" || tktID == "eTkt076" {
		tktID = args[12]
	} else {
		tripIDCurrentBytes, err := stub.GetState("LatestTripID")
		if err != nil {
			return shim.Error("Error 3 " + err.Error())
		}
		var tripCounter int
		if tripIDCurrentBytes == nil {
			tripCounter = 001
		} else {
			tripCounterStr := string(tripIDCurrentBytes)
			tripCounter, err = strconv.Atoi(tripCounterStr)
			if err != nil {
				return shim.Error("Error 4 " + err.Error())
			}
			tripCounter = tripCounter + 1
		}

		tktID = "e-Ticket" + strconv.Itoa(tripCounter)
		err = stub.PutState("LatestTripID", []byte(strconv.Itoa(tripCounter)))
		if err != nil {
			return shim.Error("Error 9 " + err.Error())
		}
	}

	//Determine the trip date
	var creationDate time.Time
	if tktID == "eTkt100" {
		creationDate = time.Now().AddDate(0, 0, -2)
	} else if tktID == "eTkt076" {
		creationDate = time.Now().AddDate(0, 0, -1)
	} else {
		creationDate = time.Now()
	}

	//Determine Geography
	var geography string
	if strings.Contains(toLOC, "USA") {
		geography = "US"
	} else if strings.Contains(toLOC, "Amsterdam") || strings.Contains(toLOC, "London") {
		geography = "Europe"
	} else {
		geography = "Worldwide"
	}

	//Assign null value to revenue sharing
	revShare := []revenueShare{}
	//Assign null value to revenue earned
	revenueEarned := 0.00
	objectType := "Trips"
	tripObj := &tripDetails{objectType, tktID, fromLOC, toLOC, riderID, products, price, currency, progress, status, hasBookedUsingSegment, hasFulfilledUsingSegment, bookedUsingProduct, event, seqNo, duration, creationDate, geography, revShare, revenueEarned}
	tripBytes, err := json.Marshal(tripObj)
	if err != nil {

		return shim.Error("Error 5 " + err.Error())
	}

	//Create Write Set
	err = stub.PutState(tktID, tripBytes)
	if err != nil {
		return shim.Error("Error 6 " + err.Error())
	}
	err = stub.SetEvent(event, tripBytes)
	if err != nil {
		return shim.Error("Error 7 " + err.Error())
	}
	//return shim.Success([]byte(tktID))
	return shim.Success(tripBytes)
}

//============================================
//updateTrip - Function to update Trip booking
//============================================
func (t *SimpleChaincode) updateTrip(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//Check if ticket ID exists in the system
	tktID := args[12]
	tripBytes, err := stub.GetState(tktID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	if tripBytes == nil {
		return shim.Error("Trip ID " + tktID + " does not exist in the system")
	}
	fromLOC := args[0]
	toLOC := args[1]
	riderID := args[2]
	products := args[3]
	price, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		price = 0
	}
	currency := args[5]
	progress := args[6]
	status := args[7]
	hasBookedUsingSegment, err := strconv.ParseBool(args[8])
	if err != nil {
		return shim.Error("Error 2 " + err.Error())
	}
	hasFulfilledUsingSegment, err := strconv.ParseBool(args[9])
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
	}
	bookedUsingProduct := args[10]
	event := args[11]
	seqNo, err := strconv.Atoi(args[13])
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}
	duration := args[14]

	//Determine the trip date
	var creationDate time.Time
	if tktID == "eTkt100" {
		creationDate = time.Now().AddDate(0, 0, -2)
	} else if tktID == "eTkt076" {
		creationDate = time.Now().AddDate(0, 0, -1)
	} else {
		creationDate = time.Now()
	}
	//Determine Geography
	var geography string
	if strings.Contains(toLOC, "USA") {
		geography = "US"
	} else if strings.Contains(toLOC, "Amsterdam") || strings.Contains(toLOC, "London") {
		geography = "Europe"
	} else {
		geography = "Worldwide"
	}
	//Calculate revenue if Payment has been made
	var revShare []revenueShare
	var revenueEarned float64
	var segmentPrice float64
	segmentPrice = 0.00
	tripObject := &tripDetails{}
	err = json.Unmarshal(tripBytes, tripObject)
	if err != nil {
		return shim.Error("Error 5 " + err.Error())
	}
	if event == "Payment Completed" {
		tripIterator, err := stub.GetHistoryForKey(tktID)
		if err != nil {
			return shim.Error("Error 8 " + err.Error())
		}
		defer tripIterator.Close()
		//products = ""
		var i int
		isWritten := false
		revShare = tripObject.RevenueSharing
		revenueEarned = tripObject.RevenueEarned
		//If this is one of the 2 default tickets, update the data differently
		if tktID == "eTkt100" || tktID == "eTkt076" {
			tripID := tktID
			fulfilledBy := products
			stakeholder1 := products
			stakeholder2 := products
			tripRevenue := price
			revCurrency := currency
			segmentPrice = tripRevenue
			tripRevenueString := strconv.FormatFloat(tripRevenue, 'f', 2, 64)
			tripRevenue, err = strconv.ParseFloat(tripRevenueString, 64)
			if err != nil {
				return shim.Error("Error 11 " + err.Error())
			}
			revShareObj := revenueShare{tripID, fulfilledBy, stakeholder1, stakeholder2, tripRevenue, revCurrency, tripRevenue, tripRevenue}
			revShare = append(revShare, revShareObj)
		} else {

			for i = 0; tripIterator.HasNext(); i++ {
				tripResponse, err := tripIterator.Next()
				if err != nil {
					return shim.Error("Error 9 " + err.Error())
				}
				tripRespBytes := tripResponse.Value
				tripRespObj := &tripDetails{}
				err = json.Unmarshal(tripRespBytes, tripRespObj)
				if err != nil {
					return shim.Error("Error 10 " + err.Error())
				}

				//Calculate revenue sharing if one of these events are part of trip history
				if tripRespObj.Event == "Segment Completed" || tripRespObj.Event == "Segment Updated" || tripRespObj.Event == "Trip Completed" || tripRespObj.Event == "Offer Redeemed" {
					if tripRespObj.Price == 0 {
						continue
					}
					tripID := tripRespObj.TktID
					fulfilledBy := tripRespObj.Products
					stakeholder1 := tripRespObj.BookedUsingProduct
					stakeholder2 := tripRespObj.Products
					tripRevenue := tripRespObj.Price
					revCurrency := tripRespObj.Currency
					segmentPrice = segmentPrice + tripRevenue
					//Calculate stakeholder revenue share only if booked by and fulfilled by parties are different
					var stakeholder1Rev, stakeholder2Rev float64
					if tripRespObj.BookedUsingProduct == tripRespObj.Products {
						stakeholder1Rev = tripRevenue
						stakeholder2Rev = tripRevenue
					} else {
						if tripRespObj.Event == "Offer Redeemed" {
							stakeholder1Rev = 0.2 * tripRevenue
							stakeholder2Rev = 0.8 * tripRevenue
						} else {
							stakeholder1Rev = 0.1 * tripRevenue
							stakeholder2Rev = 0.9 * tripRevenue
						}
					}
					stakeholder1RevString := strconv.FormatFloat(stakeholder1Rev, 'f', 2, 64)
					stakeholder2RevString := strconv.FormatFloat(stakeholder2Rev, 'f', 2, 64)
					stakeholder1Rev, err = strconv.ParseFloat(stakeholder1RevString, 64)
					if err != nil {
						return shim.Error("Error 11 " + err.Error())
					}
					stakeholder2Rev, err = strconv.ParseFloat(stakeholder2RevString, 64)
					if err != nil {
						return shim.Error("Error 12 " + err.Error())
					}
					tripRevenueString := strconv.FormatFloat(tripRevenue, 'f', 2, 64)
					tripRevenue, err = strconv.ParseFloat(tripRevenueString, 64)
					revShareObj := revenueShare{tripID, fulfilledBy, stakeholder1, stakeholder2, tripRevenue, revCurrency, stakeholder1Rev, stakeholder2Rev}
					revShare = append(revShare, revShareObj)
					revenueEarned = revenueEarned + stakeholder1Rev
					revenueEarnedString := strconv.FormatFloat(revenueEarned, 'f', 2, 64)
					revenueEarned, err = strconv.ParseFloat(revenueEarnedString, 64)
					if err != nil {
						return shim.Error("Error 13 " + err.Error())
					}
					if isWritten == true {
						products = products + ", " + tripRespObj.Products
					} else {
						products = tripRespObj.Products
						isWritten = true
					}
				} else {
					continue
				}

			}
		}
		price = segmentPrice
		event = "Trip Payment based Revenue Sharing"
	} else {

		revShare = tripObject.RevenueSharing
		revenueEarned = tripObject.RevenueEarned
	}
	objectType := "Trips"
	tripObj := &tripDetails{objectType, tktID, fromLOC, toLOC, riderID, products, price, currency, progress, status, hasBookedUsingSegment, hasFulfilledUsingSegment, bookedUsingProduct, event, seqNo, duration, creationDate, geography, revShare, revenueEarned}
	tripNewBytes, err := json.Marshal(tripObj)
	if err != nil {

		return shim.Error("Error 6 " + err.Error())
	}

	//Create Write Set
	err = stub.PutState(tktID, tripNewBytes)
	if err != nil {
		return shim.Error("Error 7 " + err.Error())
	}
	return shim.Success([]byte(tktID))

}

//===========================================================
//queryTrip - Function to query Trip booking based on trip ID
//===========================================================
func (t *SimpleChaincode) queryTrip(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	tripID := args[0]
	tripBytes, err := stub.GetState(tripID)
	if err != nil {
		return shim.Error("Error 4 " + err.Error())
	}
	return shim.Success(tripBytes)
}
func (t *SimpleChaincode) queryTripHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	tripID := args[0]
	//Get the transaction history and write into a buffer
	var buffer bytes.Buffer
	var i int
	isRecordWritten := false
	tripIterator, err := stub.GetHistoryForKey(tripID)
	if err != nil {
		return shim.Error("Error 1 " + err.Error())
	}
	defer tripIterator.Close()

	buffer.WriteString("[")
	for i = 0; tripIterator.HasNext(); i++ {
		if isRecordWritten == true {
			buffer.WriteString(",")
		}
		response, err := tripIterator.Next()
		if err != nil {
			return shim.Error("Error 2 " + err.Error())
		}
		buffer.WriteString("{\"Transaction Type\":\"Trip\",")
		buffer.WriteString("\"Trip ID\":\"")
		buffer.WriteString(tripID)
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

		buffer.WriteString("}")
		isRecordWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=====================================================================================================================
//queryTripsByStakeholder - Function to query the trips where the trips have been either booked or/and fulfilled by the stakeholder
//=====================================================================================================================
func (t *SimpleChaincode) queryTripsByStakeholder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if args[0] == "" {
		return shim.Error("Stakeholder cannot be null")
	}
	stakeholder := args[0]
	//geo := args[1]
	queryString := fmt.Sprintf("SELECT key as tktID, json_extract(valueJson,'$.creationDate') as creationDate, json_extract(valueJson,'$.progress') as progress, json_extract(valueJson, '$.status') as status, json_extract(valueJson, '$.bookedUsingProduct') as bookedUsingProduct, json_extract(valueJson,'$.products') as products, json_extract(valueJson, '$.fromLOC') as fromLOC, json_extract(valueJson,'$.toLOC') as toLOC, json_extract(valueJson,'$.revenueEarned')as revenueEarned,json_extract(valueJson,'$.revenuesharing') as revenuesharing  FROM <STATE> WHERE json_extract(valueJson, '$.objectType', '$.bookedUsingProduct') = '[\"Trips\",\"%s\"]'", stakeholder)
	//	queryString := fmt.Sprintf("SELECT key as tktID,valueJson FROM <state> WHERE json_extract(valueJson, '$.objectType') = '\"Trips\"' AND json_extract(valueJson, '$.bookedUsingProduct') =" stakeholder)
	TripsList, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(TripsList)
}

//=====================================================================
//queryTripsByRider - Function to query the trips for a particular user
//=====================================================================
func (t *SimpleChaincode) queryTripsByRider(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if args[0] == "" {
		return shim.Error("Rider ID cannot be null")
	}
	riderID := args[0]
	queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.objectType', '$.riderID') = '[\"Trips\",\"%s\"]'", riderID)

	TripsList, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(TripsList)
}

//===============================================================
//getQueryResultForQueryString - Common Function for Rich queries
//===============================================================

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	//fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//buffer.WriteString(string(queryResponse.Value))
		buffer.Write(queryResponse.Value)
		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	//fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

//==========================================
//getKPIData - Common Function for all KPIs
//==========================================

func (t *SimpleChaincode) getKPIData(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	stakeholder := args[0]
	//geography := args[1]
	if stakeholder == "" {
		return shim.Error("Stakeholder cannot be null")
	}

	queryString := fmt.Sprintf("SELECT json_extract(valueJson,'$.revenueEarned')as revenueEarned FROM <STATE> WHERE json_extract(valueJson, '$.objectType', '$.bookedUsingProduct') = '[\"Trips\",\"%s\"]'", stakeholder)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("Error 0 " + err.Error())
	}
	defer resultsIterator.Close()
	tripCount := 0
	var revenue, cost, tripRevenue float64
	revenue = 0.00
	//	var buffer bytes.Buffer
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Error 1 " + err.Error())
		}

		respString := string(response.Value)
		respJson := strings.Split(respString, ":")
		tripRev := respJson[1]
		tripRevJson := strings.Split(tripRev, "}")
		tripRevStr := tripRevJson[0]
		/*
			if respString == "" || respString == "null" || response.Value == nil {
				respString = "0"
			}

			buffer.Write(response.Value)
			buffer.WriteString(",")
		*/
		//return shim.Error(tripRevStr)
		tripRevenue, err = strconv.ParseFloat(tripRevStr, 64)
		if err != nil {
			tripRevenue = 0
			//shim.Error("Error 1.5 " + err.Error())
		}

		/*
			respObj := &tripDetails{}
			err = json.Unmarshal(respBytes, respObj)
			if err != nil {
				return shim.Error("Error 2 " + err.Error())
			}
		*/
		//Increment the trip count
		tripCount += 1
		if tripRevenue >= 0 {
			revenue = revenue + tripRevenue
		}

	}
	//Total number of organizations are hardcoded
	numberOfOrgs := 3
	//Determine cost relative to revenue (40% of revenue is cost)
	cost = 0.97 * revenue

	//Create a KPI object
	kpiObj := &KPI{numberOfOrgs, math.Round(revenue), math.Round(cost), tripCount}
	kpiBytes, err := json.Marshal(kpiObj)
	if err != nil {
		return shim.Error("Error 3 " + err.Error())
	}
	//return shim.Success(buffer.Bytes())
	return shim.Success(kpiBytes)
}
