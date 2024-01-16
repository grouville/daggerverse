package main

import "errors"

const (
	DIR_INPUT  = "/src"
	DIR_OUTPUT = "/out"
)

type Rcodesign struct {
	Source     *Directory
	PemKeyCert *Secret
	P12File    *Secret
	P12Pass    *Secret
}

func (m *Rcodesign) WithSource(src *Directory) *Rcodesign {
	m.Source = src
	return m
}

func (m *Rcodesign) WithPemSignature(pemKeyCert *Secret) *Rcodesign {
	m.PemKeyCert = pemKeyCert
	return m
}

func (m *Rcodesign) WithP12Signature(p12File *Secret, P12Pass *Secret) *Rcodesign {
	m.P12File = p12File
	m.P12Pass = P12Pass
	return m
}

func (m *Rcodesign) Sign(path string) (*Container, error) {
	baseContainer := dag.
		Cargo().
		InstallFromGit("https://github.com/indygreg/apple-platform-rs", "main", "rcodesign", "apple-codesign").
		Container()

	if m.Source != nil {
		baseContainer = baseContainer.WithMountedDirectory(DIR_INPUT, m.Source)
	} else {
		return nil, errors.New("no source provided")
	}

	if m.PemKeyCert != nil {
		baseContainer = baseContainer.WithMountedSecret("/keycert.pem", m.PemKeyCert)

		return baseContainer.WithExec([]string{"rcodesign", "sign", DIR_INPUT + "/" + path, DIR_OUTPUT + "/" + path, "--pem-file", "/keycert.pem"}), nil

	} else if m.P12File != nil && m.P12Pass != nil {
		baseContainer = baseContainer.WithMountedSecret("/cert.p12", m.P12File).
			WithMountedSecret("/pass", m.P12Pass)
		// TODO: handle this case
	}

	// encodes with no signature
	return baseContainer.WithExec([]string{"rcodesign", "sign", DIR_INPUT + "/" + path, DIR_OUTPUT + "/" + path}), nil
}

func (m *Rcodesign) SignAndExport(path string) (*Directory, error) {
	ctr, err := m.Sign(path)
	if err != nil {
		return nil, err
	}
	return ctr.Directory(DIR_OUTPUT), nil
}
