import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Connect, InitLedger, IssueCertificate, RegisterUniversity, GetAll } from '../wailsjs/go/main/App';

function App() {
    const [resultText, setResultText] = useState("Click a button to interact with the blockchain 👇");
    const [certInputs, setCertInputs] = useState({
        certHash: '',
        universitySignature: '',
        studentSignature: '',
        dateOfIssuing: '',
        certUUID: '',
        universityPK: '',
        studentPK: '',
    });
    // const [univInputs, setUnivInputs] = useState({
    //     name: '',
    //     publicKey: '',
    //     location: '',
    //     description: '',
    // });
    const [certificates, setCertificates] = useState([]);

    const handleConnect = () => {
        Connect()
            .then((result) => setResultText(result))
            .catch((err) => setResultText(`Error: ${err}`));
    };

    const handleInitLedger = () => {
        InitLedger()
            .then((result) => setResultText(result))
            .catch((err) => setResultText(`Error: ${err}`));
    };

    const handleGetAll = () => {
        GetAll()
            .then((result) => {
                try {
                    const data = JSON.parse(result);
                    if (Array.isArray(data)) {
                        setCertificates(data);
                        setResultText("Certificates retrieved successfully");
                    } else {
                        setResultText("Invalid data format");
                    }
                } catch (e) {
                    setResultText(`Error parsing JSON: ${e.message}`);
                }
            })
            .catch((err) => setResultText(`Error: ${err}`));
    };

    const updateCertInput = (e) => {
        setCertInputs({ ...certInputs, [e.target.name]: e.target.value });
    };

    const handleIssueCertificate = (e) => {
        e.preventDefault();
        IssueCertificate(
            certInputs.certHash,
            certInputs.universitySignature,
            certInputs.studentSignature,
            certInputs.dateOfIssuing,
            certInputs.certUUID,
            certInputs.universityPK,
            certInputs.studentPK
        )
            .then((result) => {
                setResultText(result);
                setCertInputs({
                    certHash: '',
                    universitySignature: '',
                    studentSignature: '',
                    dateOfIssuing: '',
                    certUUID: '',
                    universityPK: '',
                    studentPK: '',
                });
            })
            .catch((err) => setResultText(`Error: ${err}`));
    };

    // const updateUnivInput = (e) => {
    //     setUnivInputs({ ...univInputs, [e.target.name]: e.target.value });
    // };

    // const handleRegisterUniversity = (e) => {
    //     e.preventDefault();
    //     RegisterUniversity(
    //         univInputs.name,
    //         univInputs.publicKey,
    //         univInputs.location,
    //         univInputs.description
    //     )
    //         .then((result) => {
    //             setResultText(result);
    //             setUnivInputs({
    //                 name: '',
    //                 publicKey: '',
    //                 location: '',
    //                 description: '',
    //             });
    //         })
    //         .catch((err) => setResultText(`Error: ${err}`));
    // };

    // // Hàm cắt ngắn nội dung dài
    const truncateText = (text, maxLength) => {
        if (text && text.length > maxLength) {
            return `${text.substring(0, maxLength)}...`;
        }
        return text || '-';
    };

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <div id="result" className="result">
                {resultText}
                {certificates.length > 0 && (
                    <div className="table-wrapper">
                        <table className="certificate-table">
                            <thead>
                                <tr>
                                    <th>Cert Hash</th>
                                    <th>University Sig</th>
                                    <th>Student Sig</th>
                                    <th>Date Issued</th>
                                    <th>Cert UUID</th>
                                    <th>University PK</th>
                                    <th>Student PK</th>
                                    <th>Data Type</th>
                                </tr>
                            </thead>
                            <tbody>
                                {certificates.map((cert, index) => (
                                    <tr key={index}>
                                        <td title={cert.certHash}>{truncateText(cert.certHash, 10)}</td>
                                        <td title={cert.universitySignature}>{truncateText(cert.universitySignature, 10)}</td>
                                        <td title={cert.studentSignature}>{truncateText(cert.studentSignature, 10)}</td>
                                        <td title={cert.dateOfIssuing}>{truncateText(cert.dateOfIssuing, 10)}</td>
                                        <td title={cert.certUUID}>{truncateText(cert.certUUID, 10)}</td>
                                        <td title={cert.universityPK}>{truncateText(cert.universityPK, 10)}</td>
                                        <td title={cert.studentPK}>{truncateText(cert.studentPK, 10)}</td>
                                        <td title={cert.dataType}>{truncateText(cert.dataType, 10)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>

            <div className="button-group">
                <button className="btn" onClick={handleConnect}>Connect</button>
                <button className="btn" onClick={handleInitLedger}>Init Ledger</button>
                <button className="btn" onClick={handleGetAll}>Get All Certificates</button>
            </div>

            <div className="input-box">
                <h3>Issue Certificate</h3>
                <form onSubmit={handleIssueCertificate}>
                    <input className="input" name="certHash" placeholder="Certificate Hash" value={certInputs.certHash} onChange={updateCertInput} />
                    <input className="input" name="universitySignature" placeholder="University Signature" value={certInputs.universitySignature} onChange={updateCertInput} />
                    <input className="input" name="studentSignature" placeholder="Student Signature" value={certInputs.studentSignature} onChange={updateCertInput} />
                    <input className="input" name="dateOfIssuing" placeholder="Date (e.g., 2025-03-07)" value={certInputs.dateOfIssuing} onChange={updateCertInput} />
                    <input className="input" name="certUUID" placeholder="Certificate UUID" value={certInputs.certUUID} onChange={updateCertInput} />
                    <input className="input" name="universityPK" placeholder="University Public Key" value={certInputs.universityPK} onChange={updateCertInput} />
                    <input className="input" name="studentPK" placeholder="Student Public Key" value={certInputs.studentPK} onChange={updateCertInput} />
                    <button type="submit" className="btn">Issue</button>
                </form>
            </div>

            {/* <div className="input-box">
                <h3>Register University</h3>
                <form onSubmit={handleRegisterUniversity}>
                    <input className="input" name="name" placeholder="University Name" value={univInputs.name} onChange={updateUnivInput} />
                    <input className="input" name="publicKey" placeholder="Public Key" value={univInputs.publicKey} onChange={updateUnivInput} />
                    <input className="input" name="location" placeholder="Location" value={univInputs.location} onChange={updateUnivInput} />
                    <input className="input" name="description" placeholder="Description" value={univInputs.description} onChange={updateUnivInput} />
                    <button type="submit" className="btn">Register</button>
                </form>
            </div> */}
        </div>
    );
}

export default App;