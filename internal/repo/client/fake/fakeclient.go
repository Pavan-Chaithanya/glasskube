package fake

import (
	"errors"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/client"
	"github.com/glasskube/glasskube/internal/repo/types"
)

// FakeClient is a mock implementation of RepoClient for use in tests
type FakeClient struct {
	Packages map[string]map[string]v1alpha1.PackageManifest
}

func (f *FakeClient) AddPackage(name, version string, manifest v1alpha1.PackageManifest) {
	if f.Packages == nil {
		f.Clear()
	}
	if _, ok := f.Packages[name]; !ok {
		f.Packages[name] = map[string]v1alpha1.PackageManifest{}
	}
	f.Packages[name][version] = manifest
}

func (f *FakeClient) Clear() {
	f.Packages = map[string]map[string]v1alpha1.PackageManifest{}
}

var _ client.RepoClient = &FakeClient{}

// FetchLatestPackageManifest implements client.RepoClient.
func (f *FakeClient) FetchLatestPackageManifest(repoURL string, name string, target *v1alpha1.PackageManifest) (version string, err error) {
	if versions, ok := f.Packages[name]; ok {
		for v, m := range versions {
			*target = m
			return v, nil
		}
	}
	return "", errors.New("not found")
}

// FetchPackageIndex implements client.RepoClient.
func (f *FakeClient) FetchPackageIndex(repoURL string, name string, target *types.PackageIndex) error {
	if versions, ok := f.Packages[name]; ok {
		var result types.PackageIndex
		for v := range versions {
			result.LatestVersion = v
			result.Versions = append(result.Versions, types.PackageIndexItem{Version: v})
		}
		*target = result
		return nil
	}
	return errors.New("not found")
}

// FetchPackageManifest implements client.RepoClient.
func (f *FakeClient) FetchPackageManifest(repoURL string, name string, version string, target *v1alpha1.PackageManifest) error {
	if versions, ok := f.Packages[name]; ok {
		if manifest, ok := versions[version]; ok {
			*target = manifest
			return nil
		}
	}
	return errors.New("not found")
}

// FetchPackageRepoIndex implements client.RepoClient.
func (f *FakeClient) FetchPackageRepoIndex(repoURL string, target *types.PackageRepoIndex) error {
	var result types.PackageRepoIndex
	for pkg, versions := range f.Packages {
		item := types.PackageRepoIndexItem{Name: pkg}
		for v := range versions {
			item.LatestVersion = v
		}
		result.Packages = append(result.Packages, item)
	}
	*target = result
	return nil
}

// GetLatestVersion implements client.RepoClient.
func (f *FakeClient) GetLatestVersion(repoURL string, pkgName string) (string, error) {
	if versions, ok := f.Packages[pkgName]; ok {
		for v := range versions {
			return v, nil
		}
	}
	return "", errors.New("not found")
}

// GetPackageManifestURL implements client.RepoClient.
func (f *FakeClient) GetPackageManifestURL(repoURL string, name string, version string) (string, error) {
	return "fake url", nil
}
