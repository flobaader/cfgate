# Service Mesh Integration

## Overview

cfgate uses the [Gateway API](https://gateway-api.sigs.k8s.io/) standard (see [Gateway API Primer](gateway-api-primer.md) for cfgate-specific concepts)—the same specification that Istio, Envoy Gateway, and Cilium rely on for traffic management. Each mesh implementation registers and manages its own `GatewayClass`, and cfgate registers `cfgate.io/cloudflare-tunnel-controller`. Because Gateway API supports multiple concurrent `GatewayClass` resources by design, cfgate coexists with other implementations without conflict.

## Kiali

[Kiali](https://kiali.io/) provides observability for Istio service meshes. By default, Kiali only recognizes Istio's own `GatewayClass` resources. When cfgate's `GatewayClass` is present in the cluster, Kiali flags it with **KIA1504** validation warnings—indicating an unrecognized Gateway API class. The fix is adding cfgate to Kiali's `gateway_api_classes` configuration.

### Kiali CR

If you manage Kiali through the Kiali Operator, add the `cfgate` class under `spec.external_services.istio.gateway_api_classes`:

```yaml
spec:
  external_services:
    istio:
      gateway_api_classes:
        - class_name: "istio"
          name: "Istio"
        - class_name: "cfgate"
          name: "cfgate"
```

### Kiali ConfigMap

If you deploy Kiali without the Operator, apply the same configuration in the Kiali ConfigMap:

```yaml
external_services:
  istio:
    gateway_api_classes:
      - class_name: "istio"
        name: "Istio"
      - class_name: "cfgate"
        name: "cfgate"
```

> **Note:** Setting `gateway_api_classes` explicitly replaces Kiali's auto-discovery. Include all `GatewayClass` resources you want Kiali to recognize (e.g., `istio`, `istio-remote`, `cfgate`).

See the [Kiali CR Reference](https://kiali.io/docs/configuration/kialis.kiali.io/) for all configuration options.
