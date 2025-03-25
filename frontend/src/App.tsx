import { Route, Routes, BrowserRouter as Router } from "react-router-dom";
import SignupPage from "./pages/auth";
import LoginPage from "./pages/login";
import FileUploadPage from "./pages/cid";
import FileVerification from "./pages/verify";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/cid" element={<FileUploadPage />} />
        <Route path="/verify" element={<FileVerification />} />
      </Routes>
    </Router>
  );
}

export default App;
