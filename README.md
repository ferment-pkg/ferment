# Ferment  ![Build](https://github.com/NotTimIsReal/Ferment/actions/workflows/build.yml/badge.svg)
<image src="images/logo.svg" width="150px">

## Fast and efficent package manager written in GO
## Uses Python For Installation And Uninstallation
# Installation:

## Requirements:
- [python3](https://www.python.org/) (Should Be Pre-Installed On Most Macs)
- [xcode-cli](https://www.freecodecamp.org/news/install-xcode-command-line-tools/)

## Install:
```sh
git clone https://github.com/NotTimIsReal/Ferment.git /usr/local/Ferment/
cd /usr/local/Ferment/
./install.sh
```

# How Does This Work
Macos doesn't have a package manager on default, and the one that most would install is brew. Brew itself is a very good package manager but the one down side is speed, it's written in ruby which is an interpreted language which results in slow speeds and wouldn't be able to effectively use the performance of the system. Ferment on the other hand is written in GO which is a compiled language which has amazing speeds and can use multiple cores. The only interpreted language used is python which is used for the installation and uninstallation of Packages.

# Supported Systems
Operating System: MacOS

Architecture: amd64, arm64

# Usage
## Install
```sh
ferment install <packages>
```
<image src="images/output.gif" >

## Uninstall
```sh
ferment uninstall <packages>
```
## List
```sh
ferment list
```
## Search
```sh
ferment search <package>
```
## Reinstall
```sh
ferment reinstall <packages>
```
# FAQ
## Why Is Ferment Faster Than Brew?
Ferment is written in GO which is compiled to native code which is faster than the interpreted language ruby.
## My Package has dependencies, How Do I Install Them?
simply adding `self.dependencies=[<dependency>]` to your Barrell's '\__init\__ function.
## How Do i add my own package to Ferment?
Create a new file in the Barrells folder and name it the same as the package you want to add. Create A Class with the same name as the file, you can look at index.py in barrells to see what variables are read. 

## How to update to a newer version of Ferment?
```sh
./update.sh
```

# Trouble-Shooting
## Automake Doesn't Install
**Fix:** Uninstall the package and then close your terminal and clear your cache at `/tmp/ and ~/Library/Caches/`, open a new terminal window and run `ferment install <package>`
## Pkg-Config Doesn't Build On M1
**Fix:** No current fix is available except running the build command manually. On Intel, a pre-compiled will be used.

**PS:** It would be incredibly helpful if someone after manually building pkg-config would submit a pull request on https://github.com/ferment-pkg/pkg-config-prebuilt with the prefix's output (bin and share).




