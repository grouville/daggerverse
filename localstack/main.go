package main

type Localstack struct{}

// LocalStack returns a new Localstack service
// exposed on port 4566
// usage:  dagger call local-stack up --ports 4566:4566
func (m *Localstack) LocalStack() *Service {
	return dag.Container().
		From("localstack/localstack:latest").
		WithExposedPort(4566).
		AsService()
}
