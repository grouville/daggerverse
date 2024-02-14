package main

type Zip struct{}

// zips the given directory and returns the path to the zip file
// usage: dagger call compress --file <path> export --path <path.zip>
func (m *Zip) Compress(file *File) *File {
	filePath := "/tmp/file"

	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "zip"}).
		WithMountedFile(filePath, file).
		WithExec([]string{"zip", "-r", filePath + ".zip", filePath}).
		File(filePath + ".zip")
}
