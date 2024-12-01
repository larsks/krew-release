# krew-release

This is a tool for producing [krew] plugin manifests from GitHub releases.

## Requirements

You must have the [`gh`][ghcli] cli installed and configured.

[ghcli]: https://cli.github.com/

## Usage

[krew]: https://krew.sigs.k8s.io/

When run from inside a git repository, krew uses the [`gh`][ghcli] command line tool to get a list of assets from the given release. It iterates over these assets, downloading each one to calculate the required SHA256 checksum. `krew-release` requires that your assets are named using the format `<name>-<os>-<arch>.<anything>`; e.g. for my [`kubectl-saconfig`][saconfig] plugin, the assets in a given release are:

- `kubectl-saconfig-darwin-amd64.tar.gz`
- `kubectl-saconfig-darwin-arm64.tar.gz`
- `kubectl-saconfig-linux-amd64.tar.gz`
- `kubectl-saconfig-linux-arm.tar.gz`
- `kubectl-saconfig-linux-arm64.tar.gz`

[saconfig]: https://github.com/larsks/kubectl-saconfig

This information is used to render a template written in Go's [template language]. The template is provided with the following top-level attributes:

[template language]: https://pkg.go.dev/text/template

- `.Release` -- this is the release version specified on the command line.
- `.Assets` -- this is a list of assets associated with given release. Each entry has the following attributes:
  - `Name` -- name
  - `Url` -- download url
  - `Sha256` -- sha256 checksum
  - `Arch` -- architecture
  - `Os` -- operating system

## Example

The template for my [`kubectl-saconfig`][saconfig] project looks like this:

```
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: saconfig
spec:
  homepage: https://github.com/larsks/kubectl-saconfig
  shortDescription: Generate a kubeconfig file for authenticating as a service account
  version: {{ .Release }}
  description: |
    Request a token using the [TokenRequest] API and generate a kubeconfig file
    for authenticating as a service account. Outputs the generated
    configuration to stdout (default) or to a file of your choice (--output).

    [tokenrequest]: https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/
  platforms:
{{ range .Assets }}
    - selector:
        matchLabels:
          os: {{ .Os }}
          arch: {{ .Arch }}
      uri: {{ .Url }}
      sha256: "{{ .Sha256 }}"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .
{{ end }}
```

Running the following command from a working copy of this repository:

```
krew-release v0.1.3 krew/krew.yaml.tmpl
```

Produces the following output:

```
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: saconfig
spec:
  homepage: https://github.com/larsks/kubectl-saconfig
  shortDescription: Generate a kubeconfig file for authenticating as a service account
  version: v0.1.3
  description: |
    Request a token using the [TokenRequest] API and generate a kubeconfig file
    for authenticating as a service account. Outputs the generated
    configuration to stdout (default) or to a file of your choice (--output).

    [tokenrequest]: https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-request-v1/
  platforms:

    - selector:
        matchLabels:
          os: darwin
          arch: amd64
      uri: https://github.com/larsks/kubectl-saconfig/releases/download/v0.1.3/kubectl-saconfig-darwin-amd64.tar.gz
      sha256: "d3c1b261a31cf4812c9b41d88891b0cf9b0428915ae1b7e8d165886be9af2b5a"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .

    - selector:
        matchLabels:
          os: darwin
          arch: arm64
      uri: https://github.com/larsks/kubectl-saconfig/releases/download/v0.1.3/kubectl-saconfig-darwin-arm64.tar.gz
      sha256: "b193071172fbfdbf760e542c4b20cdf90ddadbd0c9713c9b48c06a8399c77442"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .

    - selector:
        matchLabels:
          os: linux
          arch: amd64
      uri: https://github.com/larsks/kubectl-saconfig/releases/download/v0.1.3/kubectl-saconfig-linux-amd64.tar.gz
      sha256: "3d339dc5fc9a830c4d463b040ea19d975f08ee464640008eda01192bb2c73d2d"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .

    - selector:
        matchLabels:
          os: linux
          arch: arm
      uri: https://github.com/larsks/kubectl-saconfig/releases/download/v0.1.3/kubectl-saconfig-linux-arm.tar.gz
      sha256: "45443c024cb6aed5a2c2b62cccde4be616796c9c194993bbf910184db1e8fa83"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .

    - selector:
        matchLabels:
          os: linux
          arch: arm64
      uri: https://github.com/larsks/kubectl-saconfig/releases/download/v0.1.3/kubectl-saconfig-linux-arm64.tar.gz
      sha256: "750f5b78dd4da1aabc6362154afa9cf42920e8dc57f43fbe4f1e48a1b399c1ca"
      bin: "./kubectl-saconfig"
      files:
        - from: kubectl-saconfig
          to: .
        - from: LICENSE
          to: .

```
