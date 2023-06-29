package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const voterPrefix = "Voter"

func (c *SmartContract) CreateVoter(ctx contractapi.TransactionContextInterface, args string) Response {
	voter := &Voter{}
	err := JSONtoObject([]byte(args), voter)

	voterKey, err := ctx.GetStub().CreateCompositeKey(voterPrefix, []string{voter.UserId})

	objVoterBytes, err := ObjecttoJSON(voter)

	err = ctx.GetStub().PutState(voterKey, objVoterBytes)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("failed to add new voter record to blockchain"), nil)
	}
	return BuildResponse("Success!", fmt.Sprintf("New Voter record added to the blockchain network sucessfully."), nil)

}

func (c *SmartContract) GetVoterDetails(ctx contractapi.TransactionContextInterface, args string) Response {

	voterKey, err := ctx.GetStub().CreateCompositeKey(voterPrefix, []string{args})

	voterBytes, err := ctx.GetStub().GetState(voterKey)
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("failed to Read Voter details from the ledger"), nil)
	}

	return BuildResponse("SUCCESS", "", voterBytes)
}
