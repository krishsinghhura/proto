import React, { useState, useEffect } from "react";
import { ethers } from "ethers";
import ABI from "./abi.json"; // Replace with your actual ABI file path

declare global {
  interface Window {
    ethereum?: any;
  }
}

const SignupPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState<string>("");
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

  const provider = new ethers.BrowserProvider(window.ethereum); // Uses browser provider

  const handleSignup = async () => {
    if (!walletAddress) {
      alert("Wallet address not found");
      return;
    }
    try {
      setLoading(true);
      const signer = await provider.getSigner();
      const con = import.meta.env.VITE_CONTRACT_ADDRESS;
      const contract = new ethers.Contract(con, ABI, signer); // Replace with your contract address
      const tx = await contract.signup();
      await tx.wait();
      alert("Signup successful!");
    } catch (error) {
      console.error(error);
      alert("Signup failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>Signup Page</h1>
      <p>Connected Wallet: {walletAddress}</p>
      <button onClick={handleSignup} disabled={loading}>
        {loading ? "Signing Up..." : "Signup"}
      </button>
    </div>
  );
};

export default SignupPage;
