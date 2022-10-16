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
curl -SsL https://fermentpkg.tech/install.sh | sh
```
  or
```sh
git clone --recurse-submodules -j8 https://github.com/NotTimIsReal/Ferment.git /usr/local/Ferment/
cd /usr/local/Ferment/
./install.sh
```

# How Does This Work
Macos doesn't have a package manager on default, and the one that most would install is brew. Brew itself is a very good package manager but the one down side is speed, it's written in ruby which is an interpreted language which results in slow speeds and wouldn't be able to effectively use the performance of the system. Ferment on the other hand is written in GO which is a compiled language which has amazing speeds and can use multiple cores. The only interpreted language used is python which is used for the installation and uninstallation of Packages.

# Supported Systems
Operating System: MacOS

Architecture: amd64, arm64

# Usage
<image src="images/output.gif" >

## Docs
- [Install](Docs/install.md)
- [Uninstall](Docs/uninstall.md)
- [List](Docs/list.md)
- [Reinstall](Docs/reinstall.md)
- [Update](Docs/update.md)
- [Upgrade](Docs/upgrade.md)
- [Clean](Docs/clean.md)
- [Search](Docs/search.md)
- [Own](Docs/own.md)

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
## "Unknown Download" Loop
This can be fixed by entering the `/tmp/ferment` directory and locating the directory named after the package you are trying to install. In there look for a tar.gz file, if there are multiple you might need to look for the newest file in there. With the file's name known you need to rename it to something else and quit the command. Then restore it's original name and install the package again and it should run normally.




