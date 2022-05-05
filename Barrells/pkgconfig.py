#Credit: Kareku Sato - https://noknow.info/it/os/install_pkg_config_from_source?lang=en
import subprocess
import requests
import os
import platform
from index import Barrells
class pkgconfig(Barrells):
    def __init__(self):
        self.url="https://pkgconfig.freedesktop.org/releases/pkg-config-0.29.2.tar.gz"
        self.git=False
        self.description="Manage compile and link flags for libraries"
        self.dependencies=["autoconf", "automake", "libtool"]
    def install(self):
            os.chdir(self.cwd)
            args=["--disable-debug", f"--prefix={self.cwd}/built","--disable-host-tool", " --with-internal-glib"]
            #https://github.com/ferment-pkg/pkg-config-prebuilt
            if "x86_64" in platform.machine():
                # Use prebuilt amd64 binary
                content=requests.get("https://github.com/ferment-pkg/pkg-config-prebuilt/archive/refs/tags/v1.tar.gz").content
                with open("pkg-config-prebuilt.tar.gz", "wb") as f:
                    f.write(content)
                os.system("tar -xzf pkg-config-prebuilt.tar.gz")
                os.system("mv pkg-config-prebuilt-1 built")
                os.symlink(f"{self.cwd}/built/amd64/bin/pkg-config", "/usr/local/bin/pkg-config")
                os.symlink(f"{self.cwd}/built/amd64/share/aclocal/pkg.m4", "/usr/local/share/aclocal/pkg.m4")
            else: 
                subprocess.call(["sh","configure", *args])
                subprocess.call(["make"])
                subprocess.call(["make","install"])
                os.symlink(f"{self.cwd}/built/bin/pkg-config", "/usr/local/bin/pkg-config")
                os.symlink(f"{self.cwd}/built/share/aclocal/pkg.m4", "/usr/local/share/aclocal/pkg.m4")
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/pkg-config")
            os.remove("/usr/local/share/aclocal/pkg.m4")
        finally:
            return super().uninstall()