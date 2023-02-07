package test

import (
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// go test ./test -v -run ^TestReadYamlFromUrl$
func TestReadYaml(t *testing.T) {
	t.Run("read flannel manifest from url", func(t *testing.T) {
		yaml, err := model.ReadYamlFromUrl(option.FLANNEL_MANIFEST_URL)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(yaml)
	})
}
