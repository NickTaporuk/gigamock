package fileProvider

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	createUserYAMLFilePath = "../../examples/rest/create-user.yaml"
)

func toAbsPath(relativePath string) (string, error) {
	return filepath.Abs(relativePath)
}

func TestYAMLProvider_Unmarshal(t *testing.T) {
	t.Run("positive test parse correct yaml file", func(t *testing.T) {
		provider := YAMLProvider{}
		path, err := toAbsPath(createUserYAMLFilePath)

		assert.NoError(t, err)
		baseScenario, err := provider.Unmarshal(path)

		fmt.Printf("%#v", baseScenario)
		assert.NoError(t, err)

		assert.NotNil(t, baseScenario)
	})
}
