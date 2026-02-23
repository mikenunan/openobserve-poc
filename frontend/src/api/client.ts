export interface Customer {
    partyId: number;
    firstName: string;
    lastName: string;
    addressLine1: string;
    addressLine2?: string;
    city: string;
    postcode: string;
    country: string;
}

export interface Affiliation {
    partyId: number;
    documentName: string;
}

export interface DocumentContent {
    name: string;
    content: string;
}

const API_BASE = "/api";

export async function fetchCustomers(): Promise<Customer[]> {
    const res = await fetch(`${API_BASE}/partyman/customers`);
    if (!res.ok) throw new Error(`Failed to fetch customers: ${res.status}`);
    return res.json();
}

export async function fetchCustomer(partyId: number): Promise<Customer> {
    const res = await fetch(`${API_BASE}/partyman/customer?partyId=${partyId}`);
    if (!res.ok) throw new Error(`Failed to fetch customer: ${res.status}`);
    return res.json();
}

export async function updateCustomer(customer: Customer): Promise<Customer> {
    const res = await fetch(`${API_BASE}/partyman/customer`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(customer),
    });
    if (!res.ok) throw new Error(`Failed to update customer: ${res.status}`);
    return res.json();
}

export async function fetchAffiliations(
    partyId: number,
): Promise<Affiliation[]> {
    const res = await fetch(`${API_BASE}/filing/documents?partyId=${partyId}`);
    if (!res.ok) throw new Error(`Failed to fetch affiliations: ${res.status}`);
    return res.json();
}

export async function fetchDocument(name: string): Promise<DocumentContent> {
    const res = await fetch(
        `${API_BASE}/catalogue/document?name=${encodeURIComponent(name)}`,
    );
    if (!res.ok) throw new Error(`Failed to fetch document: ${res.status}`);
    return res.json();
}
