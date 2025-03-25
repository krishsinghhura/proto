import React, { useState, useEffect } from "react";
import { ethers } from "ethers";
import ABI from "./abi.json"; // Replace with your actual ABI file path

declare global {
  interface Window {
    ethereum?: any;
  }
}

const LoginPage: React.FC = () => {
  const [walletAddress, setWalletAddress] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);

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

  const handleLogin = async () => {
    if (!walletAddress) {
      alert("Wallet address not found");
      return;
    }
    try {
      setLoading(true);
      const signer = await provider.getSigner();
      const contract = new ethers.Contract(
        import.meta.env.VITE_CONTRACT_ADDRESS,
        ABI,
        signer
      );
      const isUserLoggedIn = await contract.login();
      if (isUserLoggedIn) {
        alert("Login successful!");
        setIsLoggedIn(true);
      } else {
        alert("Login failed. Please try again.");
      }
    } catch (error) {
      console.error(error);
      alert("Login failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>Login Page</h1>
      <p>Connected Wallet: {walletAddress}</p>
      <button onClick={handleLogin} disabled={loading}>
        {loading ? "Logging In..." : "Login"}
      </button>
      {isLoggedIn && <p>You are now logged in.</p>}
    </div>
  );
};

export default LoginPage;
