# Ferment  ![Build](https://github.com/NotTimIsReal/Ferment/actions/workflows/build.yml/badge.svg)
## Fast and efficent package manager written in GO
## Uses Python For Installation And Uninstallation
# Installation:

## Requirements:
- python3 (Should Be Pre-Installed On Most Macs)

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
ferment install <package>
```
## Uninstall
```sh
ferment uninstall <package>
```
## List
```sh
ferment list
```
## Search
```sh
ferment search <package>
```
# FAQ
## Why Is Ferment Faster Than Brew?
Ferment is written in GO which is compiled to native code which is faster than the interpreted language ruby.
## My Package has dependencies, How Do I Install Them?
simply adding `self.dependencies=[<dependency>]` to your Barrell's '\__init\__ function.
## How Do i add my own package to Ferment?
Create a new file in the Barrells folder and name it the same as the package you want to add. Create A Class with the same name as the file, you can look at index.py in barrells to see what variables are read. 
