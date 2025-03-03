// api/main.go
package main

import (
	"net/http"
	"time"
	"fmt"
	"backend/gateway"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/gin-gonic/gin"
)

var contract *client.Contract

func main() {
	// Khởi tạo kết nối Fabric Gateway
	gw, conn, err := educert.NewGateway()
	if err != nil {
		panic(err)
	}
	defer gw.Close()
	defer conn.Close()

	// Kết nối tới channel và chaincode
	network := gw.GetNetwork("mychannel")
	contract = network.GetContract("certicontract")

	// Khởi tạo REST API
	router := gin.Default()
	educert.InitLedger(contract)
	// Các endpoint
	router.GET("/certificates", getAllCertificates)
	router.GET("/certificates/:id", getCertificate)
	router.POST("/certificates", createCertificate)
	router.PUT("/certificates/:id", updateCertificate)
	router.DELETE("/certificates/:id", deleteCertificate)

	// Khởi chạy server
	router.Run(":8080")
}

// Struct để nhận dữ liệu từ request POST/PUT
type CertificateInput struct {
	ID          string `json:"id"`
	StudentName string `json:"studentName"`
	University  string `json:"university"`
	IssueDate   string `json:"issueDate"`
	DegreeType  string `json:"degreeType"`
}

// Các handler cho REST API
func getAllCertificates(c *gin.Context) {
	result, err := educert.GetAllCertificates(contract)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func getCertificate(c *gin.Context) {
	id := c.Param("id")
	result, err := educert.ReadCertificate(contract, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func createCertificate(c *gin.Context) {
	var input CertificateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Tạo ID tự động nếu không cung cấp
	if input.ID == "" {
		input.ID = fmt.Sprintf("cert%d", time.Now().Unix())
	}

	err := educert.CreateCertificate(contract, input.ID, input.StudentName, input.University, input.IssueDate, input.DegreeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Certificate created", "id": input.ID})
}

func updateCertificate(c *gin.Context) {
	id := c.Param("id")
	var input CertificateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := educert.UpdateCertificate(contract, id, input.StudentName, input.University, input.IssueDate, input.DegreeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Certificate updated"})
}

func deleteCertificate(c *gin.Context) {
	id := c.Param("id")
	err := educert.DeleteCertificate(contract, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted"})
}