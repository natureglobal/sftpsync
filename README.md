sftpsync
=======

[![Test Status](https://github.com/natureglobal/sftpsync/workflows/test/badge.svg?branch=main)][actions]
[![Coverage Status](https://codecov.io/gh/natureglobal/sftpsync/branch/main/graph/badge.svg)][codecov]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/natureglobal/sftpsync)][PkgGoDev]

[actions]: https://github.com/natureglobal/sftpsync/actions?workflow=test
[codecov]: https://codecov.io/gh/natureglobal/sftpsync
[license]: https://github.com/natureglobal/sftpsync/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/natureglobal/sftpsync

sftpsync short description

## Synopsis

```console
% SFTP_PASSWORD=xxxx sftpsync -P 2222 -src ./public -dst htdocs user@example.com
```

## Description

## Installation

```console
# Install the latest version. (Install it into ./bin/ by default).
% curl -sfL https://raw.githubusercontent.com/natureglobal/sftpsync/main/install.sh | sh -s

# Specify installation directory ($(go env GOPATH)/bin/) and version.
% curl -sfL https://raw.githubusercontent.com/natureglobal/sftpsync/main/install.sh | sh -s -- -b $(go env GOPATH)/bin [vX.Y.Z]

# In alpine linux (as it does not come with curl by default)
% wget -O - -q https://raw.githubusercontent.com/natureglobal/sftpsync/main/install.sh | sh -s [vX.Y.Z]

# go get
% go install github.com/natureglobal/sftpsync/cmd/sftpsync@latest
```

## Author

[Songmu](https://github.com/Songmu)
