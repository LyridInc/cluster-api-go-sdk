package test

import (
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
)

func TestInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("openstack")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}

func TestUnavailableInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("oci")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}
