import React, { useState } from "react";

const FileVerification: React.FC = () => {
  const [verificationSuccess, setVerificationSuccess] = useState(false);

  const handleFileUpload = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("document", file);

    try {
      const response = await fetch("http://localhost:8080/verify", {
        method: "POST",
        body: formData,
      });
      const data = await response.json();

      if (data.verification === "SUCCESS") {
        setVerificationSuccess(true);
      } else {
        alert("Verification failed.");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("Something went wrong.");
    }
  };

  const openFile = () => {
    const fileURL = URL.createObjectURL(new Blob()); // Placeholder for file URL
    window.open(fileURL, "_blank");
  };

  return (
    <div className="p-4 flex flex-col items-center justify-center min-h-screen bg-gray-900 text-white">
      <input
        type="file"
        accept=".pdf,.doc,.docx"
        onChange={handleFileUpload}
        className="mb-4 cursor-pointer"
      />
      {verificationSuccess && (
        <button
          onClick={openFile}
          className="px-4 py-2 bg-green-500 rounded-xl text-white hover:bg-green-600"
        >
          Open File
        </button>
      )}
    </div>
  );
};

export default FileVerification;
