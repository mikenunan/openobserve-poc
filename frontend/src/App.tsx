import { BrowserRouter, Routes, Route } from "react-router-dom";
import CustomerList from "./components/CustomerList";
import CustomerDetail from "./components/CustomerDetail";

export default function App() {
    return (
        <BrowserRouter>
            <div className="app-shell">
                <header className="app-header">
                    <div className="header-content">
                        <span className="logo">A</span>
                        <h1>A-Bank</h1>
                        <span className="subtitle">Customer Management</span>
                    </div>
                </header>
                <main className="app-main">
                    <Routes>
                        <Route path="/" element={<CustomerList />} />
                        <Route
                            path="/customer/:partyId"
                            element={<CustomerDetail />}
                        />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    );
}
