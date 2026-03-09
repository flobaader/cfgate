package controller

import (
	"testing"

	"github.com/go-logr/logr"

	cfgatev1alpha1 "cfgate.io/cfgate/api/v1alpha1"
)

func TestConvertAccessRulesRejectsUnsupportedNamedLookups(t *testing.T) {
	_, err := convertAccessRules(logr.Discard(), []cfgatev1alpha1.AccessRule{
		{IPList: &cfgatev1alpha1.AccessIPListRule{Name: "office-ips"}},
	})
	if err == nil {
		t.Fatal("expected ipList.name lookup to fail")
	}

	_, err = convertAccessRules(logr.Discard(), []cfgatev1alpha1.AccessRule{
		{EmailList: &cfgatev1alpha1.AccessEmailListRule{Name: "employees"}},
	})
	if err == nil {
		t.Fatal("expected emailList.name lookup to fail")
	}
}
