package internal

type MavenInstallInfo struct {
	mavenInstallJsonFile string
	artifactPackagesMap  map[string][]string
}

func NewMavenInstallInfo(mavenInstallJsonFile string) *MavenInstallInfo {
	return &MavenInstallInfo{
		mavenInstallJsonFile: mavenInstallJsonFile,
	}
}

func (mii *MavenInstallInfo) ProcessDeps() {
	// placeholder to process maven_install.json and extract
	// mapping of artifact ID to the set of packages/imports available
}
