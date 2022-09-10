# Install
## Description:
The Install command is used to install packages. It can install multiple packages at once, be a prebuild or for you to build yourself. It can also install dependencies for you.
## Usage:
```sh
ferment install <packages>
```
## Example:
```sh
ferment install bitgit
```
## Flags:
| Flag | Description |
| --- | --- |
|   -b, --build-from-source |         Build From Source or use an available prebuild|
|  -h, --help       |               help for install|
|  -v, --verbose    |               Verbose Output|
| --no-cache        |               Don't Use Cached Downloads|
