import React, { useState, useEffect } from "react";
import { ethers } from "ethers";
import ABI from "./abi.json"; // Replace with your actual ABI file path

declare global {
  interface Window {
    ethereum?: any;
  }
}

const FileUploadPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const [ipfsHash, setIpfsHash] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    const getWalletAddress = async () => {
      if (window.ethereum) {
        try {
          const accounts = await window.ethereum.request({
            method: "eth_requestAccounts",
          });
          setWalletAddress(accounts[0]);
        } catch (error) {
          console.error("Error fetching wallet address:", error);
        }
      } else {
        alert("MetaMask is not installed");
      }
    };
    getWalletAddress();
  }, []);

  const provider = new ethers.BrowserProvider(window.ethereum);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      alert("Please select a file to upload");
      return;
    }

    const formData = new FormData();
    formData.append("document", file);
    console.log(formData);

    try {
      setLoading(true);
      const response = await fetch("http://localhost:8080/upload", {
        // Replace with your actual endpoint
        method: "POST",
        body: formData,
      });
      console.log(response);

      const text = await response.text();
      const matches = text.match(/\{"ipfs_hash":"(.*?)"/);
      if (matches && matches[1]) {
        const ipfsHash = matches[1];
        setIpfsHash(ipfsHash);

        const signer = await provider.getSigner();
        const con = import.meta.env.VITE_CONTRACT_ADDRESS;
        const contract = new ethers.Contract(con, ABI, signer);
        await contract.storeCID(ipfsHash);
        alert("IPFS hash stored on the blockchain");
      } else {
        alert("Failed to retrieve IPFS hash from response");
      }
    } catch (error) {
      console.error("Error uploading file:", error);
      alert("Failed to upload file");
    } finally {
      setLoading(false);
    }
  };

  const handleRetrieve = async () => {
    try {
      const signer = await provider.getSigner();
      const con = import.meta.env.VITE_CONTRACT_ADDRESS;
      const contract = new ethers.Contract(con, ABI, signer);
      const cid = await contract.getCID();
      alert(`Retrieved CID: ${cid}`);
    } catch (error) {
      console.error("Error retrieving CID:", error);
    }
  };

  return (
    <div>
      <h1>File Upload Page</h1>
      <p>Connected Wallet: {walletAddress}</p>

      <input type="file" onChange={handleFileChange} />
      <button onClick={handleUpload} disabled={loading}>
        {loading ? "Uploading..." : "Upload File"}
      </button>

      {ipfsHash && <p>IPFS Hash: {ipfsHash}</p>}

      <button onClick={handleRetrieve}>Retrieve CID from Blockchain</button>
    </div>
  );
};

export default FileUploadPage;
