package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	cfgatev1alpha1 "cfgate.io/cfgate/api/v1alpha1"
)

func TestResolveSelectedNamespaces(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme: %v", err)
	}

	reader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "apps", Labels: map[string]string{"team": "edge"}}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ops", Labels: map[string]string{"team": "ops"}}},
	).Build()

	reconciler := &CloudflareDNSReconciler{APIReader: reader}
	namespaces, err := reconciler.resolveSelectedNamespaces(context.Background(), &cfgatev1alpha1.DNSNamespaceSelector{
		MatchLabels: map[string]string{"team": "edge"},
		MatchNames:  []string{"manual"},
	})
	if err != nil {
		t.Fatalf("resolveSelectedNamespaces returned error: %v", err)
	}

	if !namespaces["apps"] || !namespaces["manual"] {
		t.Fatalf("expected selector to include apps and manual, got %#v", namespaces)
	}
	if namespaces["ops"] {
		t.Fatalf("did not expect ops namespace in %#v", namespaces)
	}
}

func TestRenderExplicitTargetRequiresReadyTunnel(t *testing.T) {
	if _, err := renderExplicitTarget("{{ .TunnelDomain }}", nil); err == nil {
		t.Fatal("expected error when template is used without tunnelRef")
	}

	tunnel := &cfgatev1alpha1.CloudflareTunnel{
		Status: cfgatev1alpha1.CloudflareTunnelStatus{
			TunnelDomain: "abc.cfargotunnel.com",
		},
	}
	target, err := renderExplicitTarget("https://{{ .TunnelDomain }}", tunnel)
	if err != nil {
		t.Fatalf("renderExplicitTarget returned error: %v", err)
	}
	if target != "https://abc.cfargotunnel.com" {
		t.Fatalf("unexpected rendered target %q", target)
	}
}

func TestBuildDesiredDNSRecordPreservesType(t *testing.T) {
	record := buildDesiredDNSRecord("app.example.com", cfgatev1alpha1.RecordTypeAAAA, "2001:db8::1", false, 120, "managed by cfgate")
	if record.Type != "AAAA" {
		t.Fatalf("expected AAAA record type, got %q", record.Type)
	}
	if record.Content != "2001:db8::1" {
		t.Fatalf("unexpected record content %q", record.Content)
	}
}
