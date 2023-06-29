/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// Define structs to be used by chaincode

type User struct {
	UserID    string `json:"userId,required"`
	Name      string `json:"name"`
	Password  string `json:"password,required"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	PaymentID string `json:"paymentID"`
	Timestamp string `json:"timeStamp"`
}
type Candidate struct {
	CandidateID string          `json:"candidateId"`
	Name        string          `json:"name"`
	ElectionID  string          `json:"electionId"`
	Faculty     string          `json:"faculty"`
	Major       string          `json:"major"`
	ClassOf     string          `json:"classOf"`
	Description string          `json:"Description"`
	Jargon      string          `json:"jargon"`
	Photo       string          `json:"photo"`
	Election    []ElectionCount `json:"elections"`
}

type Committee struct {
	UserId    string `json:"userId,required"`
	Name      string `json:"name"`
	Password  string `json:"password,required"`
	Email     string `jsong:"email"`
	Timestamp string `json:"timeStamp"`
}

type Election struct {
	ElectionID  string `json:"electionId"`
	Name        string `json:"name"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"EndDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updateAt"`
	Description string `json:"description"`
}

type ElectionCount struct {
	ElectionID string `json:"electionId"`
	Votes      int    `json:"votes"`
}

type ElectionResult struct {
	Candidate []Candidate `json:"candidates"`
	Winner    Candidate   `json:"winner"`
}

type Voter struct {
	UserId     string `json:"userId,required"`
	Name       string `json:"name"`
	Password   string `json:"password,required"`
	Email      string `json:"email"`
	ElectionID string `json:"electionId"`
	Voted      bool   `json:"voted"`
	Timestamp  string `json:"timeStamp"`
}
type Patient struct {
	PatientID       string `json:"patientID"`
	FirstName       string `json:"fName"`
	LastName        string `json:"lName"`
	DOB             string `json:"dob"`
	Gender          string `json:"gender"`
	Mobile          string `json:"mobile"`
	EmergencyNumber string `json:"emergency_phone"`
	Address         string `json:"address"`
}

type SearchKey struct {
	Key string `json:"searchString"`
}

type PatientPvtData struct {
	PatientID         string `json:"patientID"`
	PastIllness       string `json:"pastIllness"`
	Surgeries         string `json:"surgeries"`
	Medications       string `json:"medications"`
	Allergies         string `json:"allergies"`
	SubstanceAbuse    string `json:"substanceAbuse"`
	DiagnosticResults string `json:"diagnosticResults"`
	TreatmentPlan     string `json:"treatmentPlan"`
	InsuranceInfo     string `json:"insuranceInfo"`
	FileCIDs          string `json:"fileLocations"`
}
