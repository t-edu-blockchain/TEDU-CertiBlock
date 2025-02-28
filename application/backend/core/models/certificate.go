package models

import (
	"time"
)

type Certificate struct {
	CertificateID  string    `json:"certificate_id"`  // Mã chứng chỉ (UUID hoặc Hash)
	StudentID      string    `json:"student_id"`      // Mã sinh viên / căn cước công dân
	StudentName    string    `json:"student_name"`    // Họ và tên sinh viên
	University     string    `json:"university"`      // Tên trường cấp chứng chỉ
	Degree         string    `json:"degree"`          // Loại bằng cấp (VD: Cử nhân CNTT)
	IssueDate      time.Time `json:"issue_date"`      // Ngày cấp chứng chỉ
	ExpirationDate time.Time `json:"expiration_date"` // Ngày hết hạn (có thể null)
	Hash           string    `json:"hash"`            // Mã băm SHA256 để kiểm tra toàn vẹn dữ liệu
}

func NewCertificate(certID, studentID, studentName, university, degree string, issueDate, expirationDate time.Time) *Certificate {
	cert := &Certificate{
		CertificateID:  certID,
		StudentID:      studentID,
		StudentName:    studentName,
		University:     university,
		Degree:         degree,
		IssueDate:      issueDate,
		ExpirationDate: expirationDate,
	}
	return cert
}
