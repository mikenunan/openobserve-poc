import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { ZoneContextManager } from "@opentelemetry/context-zone";
import { Resource } from "@opentelemetry/resources";
import { ATTR_SERVICE_NAME } from "@opentelemetry/semantic-conventions";
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";
import { registerInstrumentations } from "@opentelemetry/instrumentation";

/**
 * Initialises OpenTelemetry browser tracing.
 *
 * Traces are sent via OTLP/HTTP to the OTel Collector, which forwards
 * them to OpenObserve alongside the Go service traces.
 */
export function initTelemetry(): void {
    const collectorUrl = "/api/otel/v1/traces";

    const exporter = new OTLPTraceExporter({
        url: collectorUrl,
    });

    const provider = new WebTracerProvider({
        resource: new Resource({
            [ATTR_SERVICE_NAME]: "frontend",
        }),
        spanProcessors: [new BatchSpanProcessor(exporter)],
    });

    provider.register({
        contextManager: new ZoneContextManager(),
    });

    // Auto-instrument fetch() calls to add trace context headers
    registerInstrumentations({
        instrumentations: [
            new FetchInstrumentation({
                // Only trace our own API calls, not external resources
                ignoreUrls: [/fonts\.googleapis/, /fonts\.gstatic/],
                propagateTraceHeaderCorsUrls: [/.*/],
            }),
        ],
    });

    console.log("[OTel] Browser tracing initialised → " + collectorUrl);
}
