/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	//"time"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var pointIndexStr = "_pointindex"				//name for the key/value that will store a list of all known points
var transactionStr = "_tx"				        //name for the key/value that will store all completed transactions

var testStr = "_testIndex"


type Point struct{
	Id string `json:"id"`					   //the fieldtags are needed to keep case from bouncing around
	Owner string `json:"owner"`
	Amount int `json:"amount"`	
}

type Transaction struct{
	Id string `json:"txID"`					   //Transaction ID from cppe system
	Timestamp string `json:"ex_time"`		   //utc timestamp of creation
	TraderA string  `json:"user_A_ID"`		   //UserA ID
	TraderB string  `json:"user_B_ID"`         //UserB ID
	SellerA string  `json:"seller_A_ID"`	   //UserA's Seller ID
	SellerB string  `json:"seller_B_ID"`       //UserB's Seller ID
	PointA string  `json:"point_A"`            //Points owned by UserA after exchange
	PointB string  `json:"point_B"`            //Points owned by UserB after exchange
	Prev_Transaction_id_A string `json:"prev_transaction_id_A"`
	Prev_Transaction_id_B string `json:"prev_transaction_id_B"`
	//Related []Point `json:"related"`		   //array of points willing to trade away
}


type AllTx struct{
	TXs []Transaction `json:"tx"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error    

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))                                   //making a test var "abc"
    if err != nil {
	        return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(pointIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	/*
	err = stub.PutState(testStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}*/

	var trades AllTx
	jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	err = stub.PutState(transactionStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
        return t.write(stub, args)
    }/* else if function == "init_transaction" {
        return t.init_transaction(stub, args)
    }*/
    /* else if function == "test"{
		return t.test(stub, args)
	}*/
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions

	if function == "read" {													//read a variable
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var fcn, jsonResp string
	var err error

/*	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}*/

	fcn = args[0]
	if fcn == "read" {
		valAsbytes, err := stub.GetState(args[1])									//get the var from chaincode state
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get state for " + args[1] + "\"}"
			return nil, errors.New(jsonResp)
		}		
		return valAsbytes, nil
	}
	return nil, err													//send it onward
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

 
/*
func (t *SimpleChaincode) init_transaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error	
	//	0        1      2     3      4      5       6
	//["000009", "claudio", "alex", "Taobao", "TMall", "500", "600", "4453", "4456", "22/3/2017 2:34:12"] *"blue", "35*



	completed := Transaction{}
	completed.Id = args[0]	
	completed.TraderA = args[1]
	completed.TraderB = args[2]
	completed.SellerA = args[3]
	completed.SellerB = args[4]
	completed.PointA = args[5]
	completed.PointB = args[6]
	completed.Prev_Transaction_id_A = args[7]
	completed.Prev_Transaction_id_B = args[8]
	completed.Timestamp = args[9]
	
	fmt.Println("- start completed trade")
	jsonAsBytes, _ := json.Marshal(completed)
	err = stub.PutState("_debug1", jsonAsBytes)

	//get the completed trade struct
	tradesAsBytes, err := stub.GetState(pointIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get TXs")
	}

	var trades AllTx
	json.Unmarshal(tradesAsBytes, &trades)										//un stringify it aka JSON.parse()
	
	trades.TXs = append(trades.TXs, completed);						//append to completed trades
	fmt.Println("! appended completed to trades")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(pointIndexStr, jsonAsBytes)								//rewrite completed orders
	if err != nil {
		return nil, err
	}
	fmt.Println("- end completed trade ")
	return nil, nil
}*/

/*
func (t *SimpleChaincode) test(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start test fcn")
	fmt.Println(args[0] + " - " + args[1])

	//get the open trade struct
	testAsBytes, err := stub.GetState(testStr)
	if err != nil {
		return nil, errors.New("Failed to get TXs")
	}
	var test []string

	json.Unmarshal(testAsBytes, &test)	
	
	return nil, nil
}*/