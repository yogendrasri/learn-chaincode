package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SampleChaincode struct {
}

/************** Store Waybill Starts ************************/
/**
Expected Input is {WayBillID : "123456", "CreatedDate" : "2016-10-09"}
***
**/
type WayBillIndex struct {
	index []string
}

//custom data models
type WayBill struct {
	WayBillID        string   `json:"wayBillID"`
	CreatedDate      string   `json:"createdDate"`
	LastModifiedDate string   `json:"lastModifiedDate"`
	Status           string   `json:"status"`
	CreatedBy        string   `json:"createdBy"`
	PendingWith      string   `json:"pendingWith"`
	Palettes         []string `json:"palettes"`
}

type MasterWayBill struct {
	WayBillID        string   `json:"wayBillID"`
	CreatedDate      string   `json:"createdDate"`
	LastModifiedDate string   `json:"lastModifiedDate"`
	Status           string   `json:"status"`
	CreatedBy        string   `json:"createdBy"`
	PendingWith      string   `json:"pendingWith"`
	Palettes         []string `json:"palettes"`
}

func CreateWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateWayBill")

	if len(args) < 1 {
		fmt.Println("JSON Data missing")
		return nil, errors.New("JSON Data missing")
	}
	wayBill := parseWayBill(args[0])

	err := stub.PutState(wayBill.WayBillID, []byte(args[0]))
	if err != nil {
		fmt.Println("Could not save way bill to ledger", err)
		return nil, err
	}

	nerr := addWayBillIdToIndex(stub, wayBill.WayBillID)
	if nerr != nil {
		fmt.Println("Could not save way bill to ledger", err)
		return nil, nerr
	}

	fmt.Println("Successfully saved way bill")
	return nil, nil

}

func CreateMasterWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateWayBill")

	if len(args) < 1 {
		fmt.Println("JSON Data missing")
		return nil, errors.New("JSON Data missing")
	}
	wayBill := parseWayBill(args[0])

	err := stub.PutState(wayBill.WayBillID, []byte(args[0]))
	if err != nil {
		fmt.Println("Could not save way bill to ledger", err)
		return nil, err
	}

	/*nerr := addWayBillIdToIndex(stub, wayBill.WayBillID)
	if nerr != nil {
		fmt.Println("Could not save way bill to ledger", err)
		return nil, nerr
	}*/

	fmt.Println("Successfully saved master way bill")
	return nil, nil

}

func addWayBillIdToIndex(stub shim.ChaincodeStubInterface, wayBillId string) error {
	var wayBillIndex WayBillIndex
	indexByte, err := stub.GetState("WAYBILL_INDEX")
	if err != nil {
		json.Unmarshal(indexByte, &wayBillIndex)
		newWayBillIds := wayBillIndex.index
		newWayBillIds = append(newWayBillIds, wayBillId)
		wayBillIndex.index = newWayBillIds

		wayBillIndexString, _ := json.Marshal(wayBillIndex)

		nerr := stub.PutState("WAYBILL_INDEX", []byte(wayBillIndexString))
		if nerr != nil {
			fmt.Println("Could not save way bill to ledger", err)
			return nerr
		}
	} else {
		var tmpIndex []string
		tmpIndex = append(tmpIndex, wayBillId)
		wayBillIndex.index = tmpIndex

		wayBillIndexString, _ := json.Marshal(wayBillIndex)

		nerr := stub.PutState("WAYBILL_INDEX", []byte(wayBillIndexString))
		if nerr != nil {
			fmt.Println("Could not save way bill to ledger", err)
			return nerr
		}
	}

	return nil
}

func parseWayBill(jsondata string) WayBill {
	res := WayBill{}
	json.Unmarshal([]byte(jsondata), &res)
	fmt.Println(res)
	return res
}

func GetWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetWayBill")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing way bill ID")
	}

	var wayBillId = args[0]
	bytes, err := stub.GetState(wayBillId)
	if err != nil {
		fmt.Println("Could not fetch loan way bill with id "+wayBillId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func GetMasterWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetWayBill")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing way bill ID")
	}

	var masterWayBillId = args[0]
	bytes, err := stub.GetState(masterWayBillId)
	if err != nil {
		fmt.Println("Could not fetch loan master way bill with id "+masterWayBillId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

/************** Waybill Ends ************************/

// Init resets all the things
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Inside INIT for test chaincode")
	return nil, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateWayBill" {
		return CreateWayBill(stub, args)
	} else if function == "CreateMasterWayBill" {
		return CreateMasterWayBill(stub, args)
	} else {
		return nil, errors.New("Invalid function name " + function)
	}
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetWayBill" {
		return GetWayBill(stub, args)
	} else if function == "GetMasterWayBill" {
		return GetMasterWayBill(stub, args)
	}

	return nil, nil
}

func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("Could not start SampleChaincode")
	} else {
		fmt.Println("SampleChaincode successfully started")
	}
}
