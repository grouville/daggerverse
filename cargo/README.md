A module for working with Rust projects, with utilities to lint, test, build 

```go
dag.
	Cargo().
	WithProject(dir).
	Build([]string{})
```

and even install packages:

```go
dag.
	Cargo().
	InstallFromGit(ctx, "https://github.com/your/repo.git", "master", "binary_name", "package_name")
```