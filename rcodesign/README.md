# How to use this module:

Integrated in your pipeline, you can either rely on the shell command.

#### How to repro and use this package

1. Export locally the certificate: 
```shell
export CRT="$(cat demo_signature)" 
```

1. Check if our example Go app is signed:
```shell
codesign -dv -r- ./mygoapp/test 
./mygoapp/test: code object is not signed at all
```

1. Sign it using the module
```shell
dagger download --output ./out with-source --src ./mygoapp  with-pem-signature --pem-key-cert "$CRT" sign-and-export --path test
```

with:
- `--src`: path to project including the binary to sign
- `--pem-key-cert`: signature to use, stored as a secret and passed as an env variable, so safe to use
- `--path`: path to binary to sign relative to the `--src` directory


4. Check the result: 
```shell
codesign -dv -r- ./out/test
Executable=/Users/home/Documents/daggerverse/rcodesign/out/test
Identifier=test
Format=Mach-O thin (arm64)
CodeDirectory v=20400 size=16189 flags=0x2(adhoc) hashes=501+2 location=embedded
Signature=adhoc
Info.plist=not bound
TeamIdentifier=not set
Sealed Resources=none
# designated => cdhash H"6dc875054abd876d330dfb4f864d3a2b50e8d17d"
```