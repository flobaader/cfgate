package controller

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gateway "sigs.k8s.io/gateway-api/apis/v1"
)

func TestBuildRulesFromHTTPRouteUsesBackendNamespaceAndAllMatches(t *testing.T) {
	pathType := gateway.PathMatchPathPrefix
	backendNamespace := gateway.Namespace("backend-ns")

	route := &gateway.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: gateway.HTTPRouteSpec{
			Hostnames: []gateway.Hostname{"app.example.com"},
			Rules: []gateway.HTTPRouteRule{
				{
					Matches: []gateway.HTTPRouteMatch{
						{Path: &gateway.HTTPPathMatch{Type: &pathType, Value: ptrTo("/api")}},
						{Path: &gateway.HTTPPathMatch{Type: &pathType, Value: ptrTo("/admin")}},
					},
					BackendRefs: []gateway.HTTPBackendRef{
						{
							BackendRef: gateway.BackendRef{
								BackendObjectReference: gateway.BackendObjectReference{
									Name:      "svc",
									Namespace: &backendNamespace,
									Port:      ptrTo(gateway.PortNumber(8443)),
								},
							},
						},
					},
				},
			},
		},
	}

	reconciler := &CloudflareTunnelReconciler{}
	rules, err := reconciler.buildRulesFromHTTPRoute(route)
	if err != nil {
		t.Fatalf("buildRulesFromHTTPRoute returned error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}

	for _, rule := range rules {
		if rule.Service != "http://svc.backend-ns.svc.cluster.local:8443" {
			t.Fatalf("unexpected service target %q", rule.Service)
		}
	}
	if rules[0].Path != "/api" || rules[1].Path != "/admin" {
		t.Fatalf("unexpected paths: %#v", rules)
	}
}

func TestBuildRulesFromHTTPRouteRejectsMultipleBackends(t *testing.T) {
	route := &gateway.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: gateway.HTTPRouteSpec{
			Hostnames: []gateway.Hostname{"app.example.com"},
			Rules: []gateway.HTTPRouteRule{
				{
					BackendRefs: []gateway.HTTPBackendRef{
						{
							BackendRef: gateway.BackendRef{
								BackendObjectReference: gateway.BackendObjectReference{
									Name: "svc-a",
									Port: ptrTo(gateway.PortNumber(80)),
								},
							},
						},
						{
							BackendRef: gateway.BackendRef{
								BackendObjectReference: gateway.BackendObjectReference{
									Name: "svc-b",
									Port: ptrTo(gateway.PortNumber(80)),
								},
							},
						},
					},
				},
			},
		},
	}

	reconciler := &CloudflareTunnelReconciler{}
	if _, err := reconciler.buildRulesFromHTTPRoute(route); err == nil {
		t.Fatal("expected error for multiple backendRefs")
	}
}
