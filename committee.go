package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const committeePrefix = "Committee"

func (c *SmartContract) CreateCommittee(ctx contractapi.TransactionContextInterface, args string) Response {
	committee := &Committee{}
	err := JSONtoObject([]byte(args), committee)

	committeeKey, err := ctx.GetStub().CreateCompositeKey(committeePrefix, []string{committee.UserId})

	objCommitBytes, err := ObjecttoJSON(committee)

	err = ctx.GetStub().PutState(committeeKey, objCommitBytes)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("failed to add new committee record to blockchain"), nil)
	}
	return BuildResponse("Success!", fmt.Sprintf("New Committee record added to the blockchain network sucessfully."), nil)
}

func (c *SmartContract) GetCommitteeDetails(ctx contractapi.TransactionContextInterface, args string) Response {
	committeKey, err := ctx.GetStub().CreateCompositeKey(committeePrefix, []string{args})

	committeeBytes, err := ctx.GetStub().GetState(committeKey)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("failed to Read Committee details from the ledger"), nil)
	}

	return BuildResponse("SUCCESS", "", committeeBytes)
}
