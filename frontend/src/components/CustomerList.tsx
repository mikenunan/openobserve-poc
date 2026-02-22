import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchCustomers, type Customer } from "../api/client";

export default function CustomerList() {
    const [customers, setCustomers] = useState<Customer[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        fetchCustomers()
            .then(setCustomers)
            .catch((err) => setError(err.message))
            .finally(() => setLoading(false));
    }, []);

    if (loading) return <div className="loading">Loading customers...</div>;
    if (error) return <div className="error">Error: {error}</div>;

    return (
        <div className="customer-list">
            <h1>Customers</h1>
            <table>
                <thead>
                    <tr>
                        <th>Party ID</th>
                        <th>First Name</th>
                        <th>Last Name</th>
                        <th>City</th>
                        <th>Postcode</th>
                    </tr>
                </thead>
                <tbody>
                    {customers.map((c) => (
                        <tr
                            key={c.partyId}
                            onClick={() => navigate(`/customer/${c.partyId}`)}
                        >
                            <td>{c.partyId}</td>
                            <td>{c.firstName}</td>
                            <td>{c.lastName}</td>
                            <td>{c.city}</td>
                            <td>{c.postcode}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
