# AX

### Automatically create password-protected archives with an additional layer of AES Encryption on top of it! :lock:

<img align="right" width="180px" src="https://raw.githubusercontent.com/kaynetik/dotfiles/master/svg-resources/ashleymcnamara@GOPHER_STAR_WARS.png">

![GolangCI](https://github.com/kaynetik/ax/workflows/golangci/badge.svg?branch=main)
![Build](https://github.com/kaynetik/ax/workflows/Build/badge.svg?branch=main)
[![Version](https://img.shields.io/badge/version-v0.0.4-blue.svg)](https://github.com/kaynetik/ax/releases)

AX provides the ability to easily create password-protected archives from chosen directory, and apply an additional
layer of AES encryption once first iteration has been completed. Passwords which are used for Archive protection and AES
encryption can (and SHOULD) be different. Recommended length is `> 14` for better entropy, and as complex as you can get
it.

----

## Table of Contents

- [Getting Started](#getting-started)

## Getting Started

Currently, in active development - should become stable by second week of April '21. Do not use until version badge
becomes :GREEN: and until adequate testing coverage has been applied.

If you wish to use it as a module for your own product, you can

```sh
$ go get -u github.com/kaynetik/ax
```

If you want to use it as an CLI, you can download & put to path `ax` in the following way

```sh
$ curl #TODO Once build is completed
$ ln -s #TODO Link with bin to the $PATH
```

Refer to makefile for examples of CLI usage. Proper examples will be provided before first stable release.

## This is still WIP

#### Current focus are the following:

1. Add more flexibility to the `PushToGIT` functionality
2. Add support to `PullFromGIT` & automatically decrypt and extract, given proper credentials were given
3. Add automated release build and generate portable executables
4. Cover with unit-tests
5. Start the GUI wrapper [will be in a separate repo]