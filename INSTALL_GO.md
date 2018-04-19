# How to go

For full experience with go, you have to do both
- install go toolchain 
- set up your go workspace

## Installation 
TL;DR:
 - Download go version 1.10 or newer for your platform: https://golang.org/dl/
 - Then follow these steps https://golang.org/doc/install

### Linux / Mac OSX

Location `/usr/local` used in examples for installation.

**download tarball**

`wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz`

**extract**

`tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz`

**set system path** in your `~/.profile` or `~/.bashrc`

`export PATH=$PATH:/usr/local/go/bin`

**check your installation**

`which go` - should be /usr/local/go/bin/go

`go version` - should be whatever version you downloaded


### Windows

Either download windows `.msi` installer from here https://golang.org/dl/

Or install manually following following steps:

  - Download newest `.zip` file for windows from https://golang.org/dl/
  - Extract zip file to your prefered location (`C:\go` is recommended)
  - Add `C:\go\bin` to your system PATH variable, here is how in case you don't know:
     - press `win` key, type in `system variables`
     - click `Edit the system enviromental variables`
     - click `Enviromental Variables…`
     - double click line with `Path`
     - add `;C:\go\bin` to end of your path variable
     - `ok`, `ok`, `ok`
     steps for windows 8 and newer might by slightly different

## Set up workspace system variables
TL;DR:
 - Just create `go` folder in your home directory.

**GOPATH**

Third party package sources and binaries, and your own packages, have to take place somewhere.

Default location is `~/go`, if this does not suits you, change it by setting `GOPATH` variable.

**GOBIN**

If you want to be able to run binaries installed with `go get`/`go install` you need to add `GOBIN` to your `PATH` - value of GOBIN is `$GOPATH/bin` by default.

### Linux / Max OSX

`mkdir ~/go`

Add following to your `~/.profile` or `./bashrc` file
```
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

### Windows

`mkdir %USERPROFILE%/go`

add followint system variables
```
GOPATH=%USERPROFILE%/go
GOBIN=%GOPATH%/bin
PATH=%PATH%;%GOBIN%
```