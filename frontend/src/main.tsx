import { initTelemetry } from "./telemetry/init";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./styles/index.css";

// Initialise browser tracing before React renders
try {
    initTelemetry();
} catch (e) {
    console.warn("[OTel] Failed to initialise telemetry:", e);
}

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <App />
    </StrictMode>,
);
