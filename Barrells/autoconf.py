import os
import subprocess
from index import Barrells
class autoconf(Barrells):
    def __init__(self):
        self.url="http://ftp.gnu.org/gnu/autoconf/autoconf-2.69.tar.gz"
        self.git=False
        self.description="Autoconf -- system configuration part of autotools"
    def install(self) -> bool:
        os.chdir(self.cwd)
        subprocess.call(["sh","configure", f"--prefix={self.cwd}/built"])
        subprocess.call(["make"])
        subprocess.call(["make","install"])
        os.symlink(f"{self.cwd}/bin/autoconf", "/usr/local/bin/autoconf")
        os.symlink(f"{self.cwd}/bin/autoheader", "/usr/local/bin/autoheader")
        os.symlink(f"{self.cwd}/bin/autom4te", "/usr/local/bin/autom4te")
        os.symlink(f"{self.cwd}/bin/autoreconf", "/usr/local/bin/autoreconf")
        os.symlink(f"{self.cwd}/bin/autoscan", "/usr/local/bin/autoscan")
        os.symlink(f"{self.cwd}/bin/autoupdate", "/usr/local/bin/autoupdate")
        os.symlink(f"{self.cwd}/bin/ifnames", "/usr/local/bin/ifnames")
        return super().install()
    def uninstall(self) -> bool:
        os.remove("/usr/local/bin/autoconf")
        os.remove("/usr/local/bin/autoheader")
        os.remove("/usr/local/bin/autom4te")
        os.remove("/usr/local/bin/autoreconf")
        os.remove("/usr/local/bin/autoscan")
        os.remove("/usr/local/bin/autoupdate")
        os.remove("/usr/local/bin/ifnames")
        return super().uninstall()