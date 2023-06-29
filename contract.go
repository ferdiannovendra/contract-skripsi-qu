package main

import (
	"encoding/json"
	"fmt"

	// "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// pb "github.com/hyperledger/fabric-protos-go/peer"
)

const candidatePrefix = "Candidate"
const electionPrefix = "Election"

// SmartContract of this fabric sample
type SmartContract struct {
	contractapi.Contract
}

// AddCandidate issues a new candidate to the world state with given details.
func (s *SmartContract) AddCandidate(ctx contractapi.TransactionContextInterface, args string) Response {
	candidate := &Candidate{}

	err := JSONtoObject([]byte(args), candidate)
	fmt.Println(args)

	electionID := candidate.ElectionID
	candidateID := candidate.CandidateID

	compositeKey := fmt.Sprintf("Election%s_Candidate%s", electionID, candidateID)
	// Generate a unique composite key from the candidate details
	// candidateKey, err := ctx.GetStub().CreateCompositeKey(candidatePrefix, []string{candidate.CandidateID})

	isCdtExists, err := s.CandidateExist(ctx, compositeKey)
	if isCdtExists {
		return BuildResponse("DUPLICATE", fmt.Sprintf("Candidate record already exists in the blockchain"), nil)
	}
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to add new Candidate to the blockchain"), nil)
	}

	objCdtBytes, err := ObjecttoJSON(candidate)
	err = ctx.GetStub().PutState(compositeKey, objCdtBytes)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to add new Candidate to the blockchain"), nil)
	}
	return BuildResponse("SUCCESS", fmt.Sprintf("New candidate record has been added to the blockchain successfully."), nil)

}

// candidateExists returns true when candidate with given EmpId exists in world state
func (s *SmartContract) CandidateExist(ctx contractapi.TransactionContextInterface, candidateId string) (bool, error) {
	candidateJSON, err := ctx.GetStub().GetState(candidateId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return candidateJSON != nil, nil
}

// ReadCandidate returns the candidate stored in the world state with given id.
func (s *SmartContract) ReadCandidate(ctx contractapi.TransactionContextInterface, args string) Response {
	candidate := &Candidate{}
	err := JSONtoObject([]byte(args), candidate)

	candidateKey, err := ctx.GetStub().CreateCompositeKey(candidatePrefix, []string{candidate.CandidateID})
	candidateBytes, err := ctx.GetStub().GetState(candidateKey)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
	}
	if candidateBytes == nil {
		return BuildResponse("ERROR", fmt.Sprintf("The candidate %s does not exist", candidate.Name), nil)
	}
	return BuildResponse("SUCCESS", "", candidateBytes)
}

// Search candidate record, checking the string starts with the key...(Remove '^' from regex to search anywhere in the string)
func (s *SmartContract) QueryByPartialKey(ctx contractapi.TransactionContextInterface, args string) Response {
	key := &SearchKey{}
	err := JSONtoObject([]byte(args), key)

	if err != nil {
		fmt.Println("Error when marshall json:", err)
	}

	keyVal := key.Key

	// Get the query iterator for the given key
	queryString := fmt.Sprintf(`{
        "selector": {
            "$or": [
                {"name": {"$regex": "(?i)^%s"}},
                {"jargon": {"$regex": "(?i)^%s"}},
                {"faculty": {"$regex": "(?i)^%s"}},
            ]
        }
    }`, keyVal, keyVal, keyVal)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		fmt.Println("error:", err)
		return BuildResponse("ERROR", fmt.Sprintf("Error occurred when query the database"), nil)
	}
	defer resultsIterator.Close()

	// Iterate over the results and create an array of CandidateRecord objects
	var records []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Error occurred when iterate records"), nil)
		}

		var record Candidate
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Error occurred when Unmarshal object"), nil)
		}
		records = append(records, &record)
	}

	candidateBytes, err := ObjecttoJSON(records)
	return BuildResponse("SUCCESS", "", candidateBytes)
}

// Updatecandidate updates an existing candidate in the world state with provided parameters.
func (s *SmartContract) UpdateCandidate(ctx contractapi.TransactionContextInterface, args string) Response {
	candidate := &Candidate{}
	err := JSONtoObject([]byte(args), candidate)

	candidateKey, err := ctx.GetStub().CreateCompositeKey(candidatePrefix, []string{candidate.CandidateID})

	objEmpBytes, err := ObjecttoJSON(candidate)
	err = ctx.GetStub().PutState(candidateKey, objEmpBytes)
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to update candidate record in the blockchain"), nil)
	}
	return BuildResponse("SUCCESS", fmt.Sprintf("Candidate record has been updated in the blockchain successfully."), nil)
}

// GetAllcandidates returns all candidates found in world state
func (s *SmartContract) GetAllCandidate(ctx contractapi.TransactionContextInterface) Response {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(candidatePrefix, []string{})
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
	}
	defer resultsIterator.Close()

	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}

		var candidate Candidate
		err = json.Unmarshal(queryResponse.Value, &candidate)
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}
		candidates = append(candidates, &candidate)
	}
	candidateBytes, err := ObjecttoJSON(candidates)
	return BuildResponse("SUCCESS", "", candidateBytes)
}
func (s *SmartContract) GetAllElections(ctx contractapi.TransactionContextInterface) Response {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(electionPrefix, []string{})
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
	}
	defer resultsIterator.Close()

	var elections []*Election
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}

		var election Election
		err = json.Unmarshal(queryResponse.Value, &election)
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}
		elections = append(elections, &election)
	}
	electionBytes, err := ObjecttoJSON(elections)
	return BuildResponse("SUCCESS", "", electionBytes)
}

// Deletecandidate deletes a given candidate from the world state.
func (s *SmartContract) DeleteCandidate(ctx contractapi.TransactionContextInterface, args string) Response {
	candidate := &Candidate{}
	err := JSONtoObject([]byte(args), candidate)

	candidateKey, err := ctx.GetStub().CreateCompositeKey(candidatePrefix, []string{candidate.CandidateID})

	candidateBytes, err := ctx.GetStub().GetState(candidateKey)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate.CandidateID data from the blockchain"), nil)
	}
	if candidateBytes == nil {
		return BuildResponse("ERROR", fmt.Sprintf("The candidate.CandidateID %s does not exist", candidate.Name), nil)
	}

	err = ctx.GetStub().DelState(candidateKey)
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to delete candidate.CandidateID record from the blockchain"), nil)
	}
	return BuildResponse("SUCCESS", "candidate record deleted successfully.", nil)

}

func (s *SmartContract) GetCandidatesByElectionId(ctx contractapi.TransactionContextInterface, args string) Response {
	// electionId := args[0]
	// partialCompositeKey, err := ctx.GetStub().CreateCompositeKey(candidatePrefix, []string{args})
	// if err != nil {
	// 	return BuildResponse("ERROR", fmt.Sprintf("Failed to create candidate composite key from the blockchain"), nil)
	// }
	startKey := fmt.Sprintf("Election%s_", args)
	endKey := fmt.Sprintf("Election%s_\uffff", args)

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
	}
	defer resultsIterator.Close()

	var candidates []*Candidate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}

		var candidate Candidate
		// return BuildResponse("Cek", fmt.Sprintf("hasilnya : "+string(queryResponse.Value)), nil)

		err = json.Unmarshal(queryResponse.Value, &candidate)
		if err != nil {
			return BuildResponse("ERROR", fmt.Sprintf("Failed to read candidate data from the blockchain"), nil)
		}
		candidates = append(candidates, &candidate)
	}
	candidateBytes, err := ObjecttoJSON(candidates)

	if err != nil {
		return BuildResponse("Failed", "Failed to get election: {$electionId}", nil)
	}
	return BuildResponse("SUCCESS", "", candidateBytes)
}

// func (s *SmartContract) DeleteCandidate(ctx contractapi.TransactionContextInterface, args string) Response {
// 	electionId := args[0]

// }

// func (c *SmartContract) GetCandidatesById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}
// electionId := args[0]

// 	// get candidates for election by id
// 	// electionid is stored in the candidate object
// 	// so you need to get all candidateId keys and get the candidate object
// 	// that match the electionId
// 	studentIdsAsBytes, err := stub.GetStateByRange("candidate.", "candidate.z")
// 	if err != nil {
// 		return shim.Error("Failed to get candidate: " + electionId)
// 	}
// 	defer studentIdsAsBytes.Close()
// 	// buffer is a JSON array containing QueryResults
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")
// 	bArrayMemberAlreadyWritten := false
// 	for studentIdsAsBytes.HasNext() {
// 		queryResponse, err := studentIdsAsBytes.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}
// 		candidateAsBytes, err := stub.GetState(queryResponse.Key)
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}
// 		candidate := Candidate{}
// 		json.Unmarshal(candidateAsBytes, &candidate)
// 		for _, election := range candidate.Election {
// 			if election.ElectionID == electionId {
// 				if bArrayMemberAlreadyWritten {
// 					buffer.WriteString(",")
// 				}
// 				buffer.WriteString("{\"Key\":")
// 				buffer.WriteString("\"")
// 				buffer.WriteString(queryResponse.Key)
// 				buffer.WriteString("\"")

//					buffer.WriteString(", \"Record\":")
//					// Record is a JSON object, so we write as-is
//					buffer.WriteString(string(candidateAsBytes))
//					buffer.WriteString("}")
//					bArrayMemberAlreadyWritten = true
//				}
//			}
//		}
//		buffer.WriteString("]")
//		return shim.Success(buffer.Bytes())
//	}
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, args []string) Response {
	voterId := args[0]
	candidateId := args[1]
	electionId := args[2]

	voterAsBytes, err := ctx.GetStub().GetState(voterId)
	if err != nil {
		return BuildResponse("Error", fmt.Sprint("Failed to get Voter from World State"+voterId), nil)
	}
	if voterAsBytes == nil {
		return BuildResponse("Error", fmt.Sprint("Voter doesnt exist"+voterId), nil)
	}

	voter := &Voter{}
	err = JSONtoObject(voterAsBytes, voter)

	if voter.Voted {
		return BuildResponse("Error", fmt.Sprintf("Voter has already casted vote"), nil)
	}

	candidateAsBytes, err := ctx.GetStub().GetState(candidateId)
	if err != nil {
		return BuildResponse("Error", fmt.Sprint("Failed to get candidate from World State"+candidateId), nil)
	}
	if candidateAsBytes == nil {
		return BuildResponse("Error", fmt.Sprint("Candidate doesnt exist"+candidateId), nil)
	}

	candidate := &Candidate{}
	err = JSONtoObject(candidateAsBytes, candidate)
	for i := 0; i < len(candidate.Election); i++ {
		if candidate.Election[i].ElectionID == electionId {
			candidate.Election[i].Votes++
			break
		}
	}

	err = ctx.GetStub().PutState(candidateId, candidateAsBytes)
	if err != nil {
		return BuildResponse("Error", fmt.Sprint("Failed to update candidate from World State"+candidateId), nil)
	}

	voter.Voted = true
	voterAsBytes, err = ObjecttoJSON(voter)
	err = ctx.GetStub().PutState(voterId, voterAsBytes)

	return BuildResponse("Success!", fmt.Sprintf("Voter Success to Vote"), nil)
}

func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, args string) Response {
	election := &Election{}

	err := JSONtoObject([]byte(args), election)

	electionKey, err := ctx.GetStub().CreateCompositeKey(electionPrefix, []string{election.ElectionID})

	electionBytes, err := ObjecttoJSON(election)

	err = ctx.GetStub().PutState(electionKey, electionBytes)

	if err != nil {
		return BuildResponse("ERROR", fmt.Sprintf("failed to add new election record to blockchain"), nil)
	}
	return BuildResponse("Success!", fmt.Sprintf("Election record added to the blockchain network sucessfully."), nil)

}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create employee details chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting employee details chaincode: %s", err.Error())
	}
}
