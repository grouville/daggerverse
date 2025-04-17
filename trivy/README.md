# Dagger Trivy Module

This module provides functions for using the popular [Trivy](https://aquasecurity.github.io/trivy/latest/) open source security scanner maintained by [Aqua Security](https://www.aquasec.com/).

## Status

- Currently supports scanning of container images in a registry or derived from a Dagger `Container` type.
- Future:
  - possibly support more types of scans that Trivy can do: e.g. Filesystem, Git Repository, Virtual Machine Image, Kubernetes, AWS


## Try me

From the `dagger` cli:
```sh
dagger call -m github.com/jpadams/daggerverse/trivy scan-image --image-ref alpine:latest

dagger call -m github.com/jpadams/daggerverse/trivy scan-image --severity MEDIUM --image-ref alpine/git:latest

dagger call -m github.com/jpadams/daggerverse/trivy scan-image --severity HIGH,CRITICAL --exit-code 1 --format json --image-ref alpine/git:latest
```

From a Dagger module:
```go
const (
trivyImageTag = "0.46.1" // semver tag or "latest"
)
_, err = dag.Trivy().ScanContainer(ctx, app, TrivyScanContainerOpts{
	TrivyImageTag: trivyImageTag,
	Severity:      "HIGH,CRITICAL",
	ExitCode:      1,
	})
if err != nil {
	return err
}
```
