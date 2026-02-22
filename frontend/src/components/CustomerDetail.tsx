import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
    fetchCustomer,
    updateCustomer,
    fetchAffiliations,
    type Customer,
    type Affiliation,
} from "../api/client";
import DocumentViewer from "./DocumentViewer";

export default function CustomerDetail() {
    const { partyId } = useParams<{ partyId: string }>();
    const navigate = useNavigate();

    const [customer, setCustomer] = useState<Customer | null>(null);
    const [editForm, setEditForm] = useState<Customer | null>(null);
    const [affiliations, setAffiliations] = useState<Affiliation[]>([]);
    const [selectedDoc, setSelectedDoc] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [saveMessage, setSaveMessage] = useState<string | null>(null);

    const loadData = useCallback(async () => {
        if (!partyId) return;
        try {
            const id = parseInt(partyId);
            const [cust, affs] = await Promise.all([
                fetchCustomer(id),
                fetchAffiliations(id),
            ]);
            setCustomer(cust);
            setEditForm({ ...cust });
            setAffiliations(affs);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Unknown error");
        } finally {
            setLoading(false);
        }
    }, [partyId]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleChange = (field: keyof Customer, value: string) => {
        if (!editForm) return;
        setEditForm({ ...editForm, [field]: value });
    };

    const handleAmend = async () => {
        if (!editForm) return;
        setSaving(true);
        setSaveMessage(null);
        try {
            const updated = await updateCustomer(editForm);
            setCustomer(updated);
            setEditForm({ ...updated });
            setSaveMessage("Customer updated successfully");
            setTimeout(() => setSaveMessage(null), 3000);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to save");
        } finally {
            setSaving(false);
        }
    };

    const handleCancel = () => {
        if (customer) setEditForm({ ...customer });
    };

    if (loading) return <div className="loading">Loading customer...</div>;
    if (error) return <div className="error">Error: {error}</div>;
    if (!editForm) return <div className="error">Customer not found</div>;

    return (
        <div className="customer-detail">
            <button className="back-btn" onClick={() => navigate("/")}>
                ← Back to List
            </button>

            <h1>Customer {editForm.partyId}</h1>

            {saveMessage && <div className="save-message">{saveMessage}</div>}

            <div className="form-section">
                <h2>Personal Details</h2>
                <div className="form-grid">
                    <label>
                        Party ID
                        <input type="text" value={editForm.partyId} disabled />
                    </label>
                    <label>
                        First Name
                        <input
                            type="text"
                            value={editForm.firstName}
                            onChange={(e) =>
                                handleChange("firstName", e.target.value)
                            }
                        />
                    </label>
                    <label>
                        Last Name
                        <input
                            type="text"
                            value={editForm.lastName}
                            onChange={(e) =>
                                handleChange("lastName", e.target.value)
                            }
                        />
                    </label>
                </div>

                <h2>Address</h2>
                <div className="form-grid">
                    <label>
                        Address Line 1
                        <input
                            type="text"
                            value={editForm.addressLine1}
                            onChange={(e) =>
                                handleChange("addressLine1", e.target.value)
                            }
                        />
                    </label>
                    <label>
                        Address Line 2
                        <input
                            type="text"
                            value={editForm.addressLine2 || ""}
                            onChange={(e) =>
                                handleChange("addressLine2", e.target.value)
                            }
                        />
                    </label>
                    <label>
                        City
                        <input
                            type="text"
                            value={editForm.city}
                            onChange={(e) =>
                                handleChange("city", e.target.value)
                            }
                        />
                    </label>
                    <label>
                        Postcode
                        <input
                            type="text"
                            value={editForm.postcode}
                            onChange={(e) =>
                                handleChange("postcode", e.target.value)
                            }
                        />
                    </label>
                    <label>
                        Country
                        <input
                            type="text"
                            value={editForm.country}
                            onChange={(e) =>
                                handleChange("country", e.target.value)
                            }
                        />
                    </label>
                </div>

                <div className="form-actions">
                    <button
                        className="btn-primary"
                        onClick={handleAmend}
                        disabled={saving}
                    >
                        {saving ? "Saving..." : "Amend"}
                    </button>
                    <button
                        className="btn-secondary"
                        onClick={handleCancel}
                        disabled={saving}
                    >
                        Cancel
                    </button>
                </div>
            </div>

            <div className="affiliations-section">
                <h2>Affiliated Documents</h2>
                {affiliations.length === 0 ? (
                    <p className="no-data">
                        No documents affiliated with this customer.
                    </p>
                ) : (
                    <ul className="doc-list">
                        {affiliations.map((aff) => (
                            <li key={aff.documentName}>
                                <button
                                    className={`doc-link ${selectedDoc === aff.documentName ? "active" : ""}`}
                                    onClick={() =>
                                        setSelectedDoc(
                                            selectedDoc === aff.documentName
                                                ? null
                                                : aff.documentName,
                                        )
                                    }
                                >
                                    📄 {aff.documentName}
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            {selectedDoc && (
                <DocumentViewer
                    documentName={selectedDoc}
                    onClose={() => setSelectedDoc(null)}
                />
            )}
        </div>
    );
}
