package porter

import (
	"os"
	"testing"

	"get.porter.sh/porter/pkg/config"
	"get.porter.sh/porter/pkg/experimental"
	"github.com/stretchr/testify/assert"

	"get.porter.sh/porter/tests"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	p := NewTestPorter(t)
	defer p.Teardown()

	err := p.Create()
	require.NoError(t, err)

	configFileStats, err := p.FileSystem.Stat("porter.yaml")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, "porter.yaml", os.FileMode(0600), configFileStats.Mode())

	// Verify that helpers is present and executable
	helperFileStats, err := p.FileSystem.Stat("helpers.sh")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, "helpers.sh", os.FileMode(0700), helperFileStats.Mode())

	dockerfileStats, err := p.FileSystem.Stat("Dockerfile.tmpl")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, "Dockerfile.tmpl", os.FileMode(0600), dockerfileStats.Mode())

	readmeStats, err := p.FileSystem.Stat("README.md")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, "README.md", os.FileMode(0600), readmeStats.Mode())

	gitignoreStats, err := p.FileSystem.Stat(".gitignore")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, ".gitignore", os.FileMode(0600), gitignoreStats.Mode())

	dockerignoreStats, err := p.FileSystem.Stat(".dockerignore")
	require.NoError(t, err)
	tests.AssertFilePermissionsEqual(t, ".dockerignore", os.FileMode(0600), dockerignoreStats.Mode())

}

func TestCreateWithBuildkit(t *testing.T) {
	p := NewTestPorter(t)
	defer p.Teardown()

	p.SetExperimentalFlags(experimental.FlagBuildDrivers)
	p.Data.BuildDriver = config.BuildDriverBuildkit

	err := p.Create()
	require.NoError(t, err)

	dockerfile, err := p.FileSystem.ReadFile("Dockerfile.tmpl")
	require.NoError(t, err, "could not read template dockerfile")

	assert.Contains(t, string(dockerfile), "# syntax=docker/dockerfile:1.2")
}
