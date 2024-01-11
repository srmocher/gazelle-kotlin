package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKotlinParseCodeSuccess(t *testing.T) {
	kp, err := NewKotlinParser()
	defer kp.Stop()
	assert.NoError(t, err)
	filePath := filepath.Join(os.Getenv("TEST_SRCDIR"), "gazelle-kotlin/kotlin/test/com/github/srmocher/gazelle_kotlin/kotlinparser/testdata/ExampleFile.kt")
	sinfos, err := kp.ParseKotlinFiles([]string{filePath})
	assert.NoError(t, err)
	assert.Equal(t, len(sinfos), 1)
}
