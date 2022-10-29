# Ferment  ![Build](https://github.com/NotTimIsReal/Ferment/actions/workflows/build.yml/badge.svg)
<image src="images/logo.svg" width="150px">

## Fast and efficent package manager written in GO
## Uses Python For Extentions
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
Macos doesn't have a package manager on default, and the one that most would install is brew. Brew itself is a very good package manager but the one down side is speed, it's written in ruby which is an interpreted language which results in slow speeds and wouldn't be able to effectively use the performance of the system. Ferment on the other hand is written in GO which is a compiled language which has amazing speeds and can use multiple cores. The only interpreted language used is python which is used for the future extention system for ferment.

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

## How to update to a newer version of Ferment?
```sh
ferment upgrade
```

# Trouble-Shooting
## "Unknown Download" Loop
This can be fixed by entering the `/tmp/ferment` directory and locating the directory named after the package you are trying to install. In there look for a tar.gz file, if there are multiple you might need to look for the newest file in there. With the file's name known you need to rename it to something else and quit the command. Then restore it's original name and install the package again and it should run normally.


# Behind The Scenes
## How Does Ferment Work?
Behind the scenes ferment mainly does network requests to the api to figure out the versions and the architechtures. Once that is done ferment downloads a .ferment file which is effecrively a .tar.xz file that consists of a `.fpkg, built, .ferment_watcher`. .fpkg is a shell like script that allows for simple package creations, the built directory contains prebuilds of the package and the .ferment_watcher is a text file that has every installed directory saved in here. Depending if build-from-source is specified or not, ferment will use the fpkg file to build the package or use the prebuilds.

## Languages Used
Ferment is exclusive to just python and go, with go being used for the most part. The only python used is for the extention system. In the ferment-uploader repo, go is also used to maintain multiple uploads at once. While the main api is made in nestJS. The Ferment-Store app is made with Tauri and VueJS. 

## Linux Support
While linux is not supported officially, one can definately make it possible and fairly easily as well. This is mainly due to linux and mac having a similar file system, the only main thing that might need to be changed across these platforms is likely the build commands as right now it's using macos specific commands.

## Extentions
Extentions don't exist yet but is something that will definately happen in the time coming soon, extentions allows for extending the main packagee manager code but it doesn't require a full recompilation of the entire codebase and it also allows for an opt-in opt-out system.

## RoadMap
- Late 2022-2023: Extentions
- 2023: Workspace packages

