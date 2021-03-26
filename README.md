# AX

### Automatically create password-protected archives with an additional layer of AES Encryption on top of it! :lock:

<img align="right" width="180px" src="https://raw.githubusercontent.com/kaynetik/dotfiles/master/svg-resources/ashleymcnamara@GOPHER_STAR_WARS.png">

![GolangCI](https://github.com/kaynetik/ax/workflows/golangci/badge.svg?branch=main)
![Build](https://github.com/kaynetik/ax/workflows/Build/badge.svg?branch=main)
[![release](https://github.com/kaynetik/ax/actions/workflows/release.yml/badge.svg)](https://github.com/kaynetik/ax/actions/workflows/release.yml)
[![Version](https://img.shields.io/badge/version-v0.0.14-purple.svg)](https://github.com/kaynetik/ax/releases)
[![Coverage Status](https://coveralls.io/repos/github/kaynetik/ax/badge.svg?branch=main)](https://coveralls.io/github/kaynetik/ax?branch=main)

#### _NOTE: This is a Work In Progress_

##### Check out the [GH Project Dashboard](https://github.com/kaynetik/ax/projects/1) for more info on all the tasks that have to be completed before AX will reach a stable version.

AX provides the ability to easily create password-protected archives from chosen directory, and apply an additional
layer of AES encryption once first iteration has been completed. Passwords which are used for Archive protection and AES
encryption can (and SHOULD) be different. Recommended length is `> 14` for better entropy, and as complex as you can get
it.

----

## Table of Contents

- [Getting Started](#getting-started)

## Getting Started

Currently, in active development - should become stable by second week of April '21. Do not use until version badge
becomes :green_heart: until adequate testing coverage has been applied.

Feel free to download preferred version from the [releases page](https://github.com/kaynetik/ax/releases). Most
operating systems & architectures are supported.

If you wish to use it as a module for your own product, you can

```sh
$ go get -u github.com/kaynetik/ax
```

If you want to use it as an CLI, you can download & run `ax` in the following way:

```sh
$ curl -s https://api.github.com/repos/kaynetik/ax/releases/latest | grep "browser_download_url.*linux_x86_64.*\"" | cut -d : -f 2,3 | tr -d \" | wget -qi - 
$ tax -xvzf ax_linux_x86_64.tar.gz
$ ./ax -help
$ ./ax # To enter interactive mode don't provide any flags
$ ./ax -arc-in ../tmp_to_archive -arc-pass on -arc-out ../tmp_archive_out # Example usage with flags
```

Example above was to download latest release for `linux x86_64` architecture. Same oneliner will work for any OS, you
just need to alter `grep` counterpart, i.e. instead of ` linux_x86_64` place `windows_x86_64`.

Refer to makefile for examples of CLI usage. Proper examples will be provided before first stable release.

## This is still WIP

#### Current focus are the following:

Precise info on priorities can be reached on the [GH Issues Page](https://github.com/kaynetik/ax/issues).

1. Add more flexibility to the `PushToGIT` functionality
2. Add support to `PullFromGIT` & automatically decrypt and extract, given proper credentials were given
3. ~Add automated release build and generate portable executables~
   ~+ Fix issue with the build for win executable - can't use `term.ReadPassword`~ -> issue was that Win doesn't
   support `syscall`, fixed now.
4. ~Cover with unit-tests~
5. Start the GUI wrapper [will be in a separate repo]