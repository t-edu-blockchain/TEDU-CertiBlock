package educert

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Certificate struct {
	ID          string `json:"ID"`
	StudentName string `json:"StudentName"`
	University  string `json:"University"`
	IssueDate   string `json:"IssueDate"`
	DegreeType  string `json:"DegreeType"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	certificates := []Certificate{
		{ID: "cert1", StudentName: "Nguyen Van A", University: "ABC University", IssueDate: "2024-02-01", DegreeType: "Bachelor"},
		{ID: "cert2", StudentName: "Tran Thi B", University: "XYZ University", IssueDate: "2023-06-15", DegreeType: "Master"},
	}

	for _, cert := range certificates {
		certJSON, err := json.Marshal(cert)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(cert.ID, certJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CreateCertificate(ctx contractapi.TransactionContextInterface, id string, studentName string, university string, issueDate string, degreeType string) error {
	exists, err := s.CertificateExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the certificate %s already exists", id)
	}

	cert := Certificate{
		ID:          id,
		StudentName: studentName,
		University:  university,
		IssueDate:   issueDate,
		DegreeType:  degreeType,
	}
	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, certJSON)
}

func (s *SmartContract) ReadCertificate(ctx contractapi.TransactionContextInterface, id string) (*Certificate, error) {
	certJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if certJSON == nil {
		return nil, fmt.Errorf("the certificate %s does not exist", id)
	}

	var cert Certificate
	err = json.Unmarshal(certJSON, &cert)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (s *SmartContract) UpdateCertificate(ctx contractapi.TransactionContextInterface, id string, studentName string, university string, issueDate string, degreeType string) error {
	exists, err := s.CertificateExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the certificate %s does not exist", id)
	}

	cert := Certificate{
		ID:          id,
		StudentName: studentName,
		University:  university,
		IssueDate:   issueDate,
		DegreeType:  degreeType,
	}
	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, certJSON)
}

func (s *SmartContract) DeleteCertificate(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.CertificateExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the certificate %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) CertificateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	certJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return certJSON != nil, nil
}

func (s *SmartContract) GetAllCertificates(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var certificates []*Certificate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var cert Certificate
		err = json.Unmarshal(queryResponse.Value, &cert)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, &cert)
	}

	return certificates, nil
}
