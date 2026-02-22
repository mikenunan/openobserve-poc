import { useEffect, useState } from "react";
import { fetchDocument, type DocumentContent } from "../api/client";

interface Props {
    documentName: string;
    onClose: () => void;
}

export default function DocumentViewer({ documentName, onClose }: Props) {
    const [doc, setDoc] = useState<DocumentContent | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        setLoading(true);
        setError(null);
        fetchDocument(documentName)
            .then(setDoc)
            .catch((err) => setError(err.message))
            .finally(() => setLoading(false));
    }, [documentName]);

    return (
        <div className="document-viewer">
            <div className="doc-header">
                <h3>{documentName}</h3>
                <button className="close-btn" onClick={onClose}>
                    ✕
                </button>
            </div>
            <div className="doc-content">
                {loading && <p>Loading document...</p>}
                {error && <p className="error">Error: {error}</p>}
                {doc && <pre>{doc.content}</pre>}
            </div>
        </div>
    );
}
