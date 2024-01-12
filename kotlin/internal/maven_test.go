package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMavenArtifactsSuccess(t *testing.T) {
	mii := NewMavenInstallInfo("testdata/maven_install.json", "maven")
	err := mii.ProcessDeps()
	assert.NoError(t, err)
	assert.Equal(t, len(mii.artifacts), 135)
	assert.Equal(t, mii.artifacts[0].Coord, "biz.aQute.bnd:biz.aQute.bnd.util:6.4.0")
	assert.Equal(t, mii.artifacts[0].Packages, []string{"aQute.bnd.classfile",
		"aQute.bnd.classfile.builder",
		"aQute.bnd.classfile.preview",
		"aQute.bnd.exceptions",
		"aQute.bnd.memoize",
		"aQute.bnd.result",
		"aQute.bnd.signatures",
		"aQute.bnd.stream",
		"aQute.bnd.unmodifiable",
		"aQute.lib.io",
		"aQute.lib.stringrover",
		"aQute.libg.glob"})
	assert.Equal(t, mii.artifacts[2].Coord, "com.benasher44:uuid-jvm:0.2.2")
	assert.Equal(t, mii.artifacts[2].Packages, []string{"com.benasher44.uuid"})
}
