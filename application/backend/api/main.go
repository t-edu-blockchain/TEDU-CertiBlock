package main

import (
	"fmt"
	"net/http"
	"backend/gateway"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var contract *client.Contract

func main() {
	// Khởi tạo kết nối Fabric Gateway
	gw, conn, err := educert.NewGateway()
	if err != nil {
		panic(fmt.Errorf("failed to initialize gateway: %w", err))
	}
	defer gw.Close()
	defer conn.Close()
	network := gw.GetNetwork("mychannel")
	contract = network.GetContract("certicontract") 

	// Khởi tạo ledger
	// result, err := educert.InitLedger(contract)
	// if err != nil {
	// 	fmt.Printf("Failed to initialize ledger: %v\n", err)
	// 	return
	// }
	// fmt.Println("Ledger initialized:", result)

	// Khởi tạo REST API
	router := gin.Default()

	// Các endpoint
	router.GET("/certificates", getAllCertificates)
	router.GET("/certificates/:uuid", getCertificateByUUID)
	// router.POST("/certificates", issueCertificate)
	router.GET("/certificates/student/:studentPK", getCertificatesByStudent)
	router.GET("/certificates/university/:universityPK", getCertificatesByUniversity)
	// router.POST("/universities", registerUniversity)

	// Khởi chạy server
	err = router.Run(":8080")
	if err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}

// Struct để nhận dữ liệu từ request POST cho chứng chỉ
type CertificateInput struct {
	CertHash          string `json:"certHash"`
	UniversitySig     string `json:"universitySignature"`
	StudentSig        string `json:"studentSignature"`
	DateOfIssuing     string `json:"dateOfIssuing"`
	CertUUID          string `json:"certUUID"`
	UniversityPK      string `json:"universityPK"`
	StudentPK         string `json:"studentPK"`
}

// Struct để nhận dữ liệu từ request POST cho trường đại học
type UniversityInput struct {
	Name        string `json:"name"`
	PublicKey   string `json:"publicKey"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// Các handler cho REST API
func getAllCertificates(c *gin.Context) {
	result, err := educert.QueryAllCertificates(contract)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func getCertificateByUUID(c *gin.Context) {
	certUUID := c.Param("uuid")
	result, err := educert.QueryCertificateByUUID(contract, certUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// func issueCertificate(c *gin.Context) {
// 	var input CertificateInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Tạo UUID tự động nếu không cung cấp
// 	if input.CertUUID == "" {
// 		input.CertUUID = fmt.Sprintf("cert-%d", time.Now().UnixNano())
// 	}

// 	result, err := educert.IssueCertificate(
// 		contract,
// 		input.CertHash,
// 		input.UniversitySig,
// 		input.StudentSig,
// 		input.DateOfIssuing,
// 		input.CertUUID,
// 		input.UniversityPK,
// 		input.StudentPK,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Certificate issued", "certUUID": input.CertUUID, "data": result})
// }

func getCertificatesByStudent(c *gin.Context) {
	studentPK := c.Param("studentPK")
	result, err := educert.GetAllCertificatesByStudent(contract, studentPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func getCertificatesByUniversity(c *gin.Context) {
	universityPK := c.Param("universityPK")
	result, err := educert.GetAllCertificatesByUniversity(contract, universityPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}


//Chuyển qua nodeapp
// func registerUniversity(c *gin.Context) {
// 	var input UniversityInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	result, err := educert.RegisterUniversity(
// 		contract,
// 		input.Name,
// 		input.PublicKey,
// 		input.Location,
// 		input.Description,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "University registered", "name": input.Name, "data": result})
// }