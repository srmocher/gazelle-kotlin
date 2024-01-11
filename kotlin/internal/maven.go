package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// MavenArtifact represents a maven artifact
// It's maven coordinate, the set of packages it provides and the resolved Bazel label
type MavenArtifact struct {
	Coord      string
	Packages   []string
	BazelLabel string
}

// MavenInstallInfo represents the information extracted from maven_install.json
type MavenInstallInfo struct {
	mavenRepoName        string
	mavenInstallJsonFile string
	artifacts            []MavenArtifact
}

// mavenInstallJson represents the structure of maven_install.json and the fields we need
// for extracting dependency information
type mavenInstallJson struct {
	depTree struct {
		deps []struct {
			coord    string   `json:"coord"`
			packages []string `json:"packages"`
		} `json:"dependencies"`
	} `json:"dependency_tree"`
}

func NewMavenInstallInfo(mavenInstallJsonFile string, mavenRepoName string) *MavenInstallInfo {
	return &MavenInstallInfo{
		mavenInstallJsonFile: mavenInstallJsonFile,
		mavenRepoName:        mavenRepoName,
	}
}

func (mii *MavenInstallInfo) coordToBazelLabel(coord string) string {
	parts := strings.Split(coord, ":")
	groupId := nonAlphanumericRegex.ReplaceAllString(parts[0], "_")
	artifactId := nonAlphanumericRegex.ReplaceAllString(parts[1], "_")
	return fmt.Sprintf("@%s://%s_%s", mii.mavenRepoName, groupId, artifactId)
}

func (mii *MavenInstallInfo) ProcessDeps() error {
	if _, err := os.Stat(mii.mavenInstallJsonFile); err != nil {
		log.Fatalf("maven file %s does not exist!", mii.mavenInstallJsonFile)
	}

	f, err := os.Open(mii.mavenInstallJsonFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("unable to open file %s", mii.mavenInstallJsonFile)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading file %s", mii.mavenInstallJsonFile)
	}

	var mvn mavenInstallJson
	if err = json.Unmarshal(b, &mvn); err != nil {
		return fmt.Errorf("Error extracting maven deps from json: %s", err)
	}

	var artifacts []MavenArtifact
	for _, dep := range mvn.depTree.deps {
		artifacts = append(artifacts, MavenArtifact{
			Coord:      dep.coord,
			Packages:   dep.packages,
			BazelLabel: mii.coordToBazelLabel(dep.coord),
		})
	}

	mii.artifacts = artifacts
	return nil
}

// GetMavenArtifacts returns the set of MavenArtifacts extracted from the maven_install.json file
func (mii *MavenInstallInfo) GetMavenArtifacts() []MavenArtifact {
	return mii.artifacts
}
