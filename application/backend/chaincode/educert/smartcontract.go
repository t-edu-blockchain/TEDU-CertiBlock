package educert

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}


// Asset describes basic details of what makes up a simple asset
type Certificate struct {
	CertHash            string `json:"certHash"`
	UniversitySignature string `json:"universitySignature"`
	StudentSignature    string `json:"studentSignature"`
	DateOfIssuing       string `json:"dateOfIssuing"`
	CertNumber          string `json:"certNumber"`
	CertUUID            string `json:"certUUID"`
	UniversityPK        string `json:"universityPK"`
	StudentPK           string `json:"studentPK"`
	DataType            string `json:"dataType"`
}

type University struct {
	Name        string `json:"name"`
	PublicKey   string `json:"publicKey"`
	Location    string `json:"location"`
	Description string `json:"description"`
	DataType    string `json:"dataType"`
}

type Schema struct {
	CertificateType string   `json:"certificateType"`
	ID              string   `json:"id"`
	Ordering        []string `json:"ordering"`
	DataType        string   `json:"dataType"`
}

func NewCertificate(certHash string, universitySignature string, studentSignature string, dateOfIssuing string, certNumber string, certUUID string, universityPK string, studentPK string) *Certificate {
	return &Certificate{
		CertHash:            certHash,
		UniversitySignature: universitySignature,
		StudentSignature:    studentSignature,
		DateOfIssuing:       dateOfIssuing,
		CertNumber:          certNumber,
		CertUUID:            certUUID,
		UniversityPK:        universityPK,
		StudentPK:           studentPK,
		DataType:            "certificate",
	}
}

func NewUniversity(name string, publicKey string, location string, description string) *University {
	return &University{
		Name:        name,
		PublicKey:   publicKey,
		Location:    location,
		Description: description,
		DataType:    "university",
	}
}

func NewSchema(certificateType string, id string, ordering []string) *Schema {
	return &Schema{
		CertificateType: certificateType,
		ID:              id,
		Ordering:        ordering,
		DataType:        "schema",
	}
}

// CreateAsset issues a new asset to the world state with given details.

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) (*Schema, error) {

	fmt.Println("--------------initLedger called--------------")

	certificateSchema := NewSchema("Bachelor", "v1", []string{"universityName", "major", "departmentName", "cgpa"})

	certSchemaJSON, err := json.Marshal(certificateSchema)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("schema_"+ certificateSchema.ID, certSchemaJSON)
	if err != nil {
		return nil, err
	}

	return certificateSchema, err
}

func (s *SmartContract) issueCertificate(ctx contractapi.TransactionContextInterface, certHash string, universitySignature string, studentSignature string, dateOfIssuing string, certUUID string, universityPK string, studentPK string) (*Certificate, error) {

	fmt.Println("--------------issueCertificate called--------------")

	certificate := NewCertificate(certHash, universitySignature, studentSignature, dateOfIssuing, "", certUUID, universityPK, studentPK)

	certJSON, err := json.Marshal(certificate)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("cert_"+certificate.CertUUID, certJSON)
	if err != nil {
		return nil, err
	}

	return certificate, err
}

func (s *SmartContract) registerUniversity(ctx contractapi.TransactionContextInterface, name string, publicKey string, location string, description string) (*University, error) {

	fmt.Println("--------------registerUniversity called--------------")

	university := NewUniversity(name, publicKey, location, description)

	universityJSON, err := json.Marshal(university)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState("uni_"+university.Name, universityJSON)

	if err != nil {
		return nil, err
	}
	
	return university, err
}


func (s *SmartContract) queryUniversityProfileByName(ctx contractapi.TransactionContextInterface, name string) (*University, error) {
	
	fmt.Println("--------------queryUniversityProfileByName called--------------")

	universityProfileJSON, err := ctx.GetStub().GetState("uni_" + name)
	if err != nil {
		return nil, err
	}
	if universityProfileJSON == nil {
		return nil, fmt.Errorf("the university %s does not exist", name)
	}

	var universityProfile University
	err = json.Unmarshal(universityProfileJSON, &universityProfile)
	if err != nil {
		return nil, err
	}

	return &universityProfile, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) queryCertificateSchema(ctx contractapi.TransactionContextInterface, schemaVersion string) (*Schema, error) {
	
	fmt.Println("--------------queryCertificateSchema called--------------")

	certificateSchemaJSON, err := ctx.GetStub().GetState("schema_" + schemaVersion)
	if err != nil {
		return nil, err
	}
	if certificateSchemaJSON == nil {
		return nil, fmt.Errorf("the schema %s does not exist", schemaVersion)
	}

	var certificateSchema Schema
	err = json.Unmarshal(certificateSchemaJSON, &certificateSchema)
	if err != nil {
		return nil, err
	}

	return &certificateSchema, nil
	
}

func (s *SmartContract) queryCertificateByUUID(ctx contractapi.TransactionContextInterface, UUID string) (*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUUID called--------------")

	certificateJSON, err := ctx.GetStub().GetState("cert_" + UUID)
	if err != nil {
		return nil, err
	}
	if certificateJSON == nil {
		return nil, fmt.Errorf("the certificate %s does not exist", UUID)
	}

	var certificate Certificate
	err = json.Unmarshal(certificateJSON, &certificate)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Certificate, error) {
	var certificates []*Certificate
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var certificate Certificate
		err = json.Unmarshal(queryResult.Value, &certificate)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}


func queryWithQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryWithQueryString called--------------")
	fmt.Println("queryString: " + queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (s *SmartContract) getAllCertificateByStudent(ctx contractapi.TransactionContextInterface, studentPK string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUniversity called--------------")

	queryString := fmt.Sprintf(`{"selector":{"dataType":"certificate","studentPK":"%s"}}`, studentPK)

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}

func (s *SmartContract) getAllCertificateByUniversity(ctx contractapi.TransactionContextInterface, universityPK string) ([]*Certificate, error) {
	
	fmt.Println("--------------queryCertificateByUniversity called--------------")

	queryString := fmt.Sprintf(`{"selector":{"dataType":"certificate","universityPK":"%s"}}`, universityPK)

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}

func (s *SmartContract) queryAll(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
	
	fmt.Println("--------------queryAll called--------------")

	queryString := `{"selector":{"dataType":"certificate"}}`

	queryResults, err := queryWithQueryString(ctx, queryString)

	return queryResults, err
}
