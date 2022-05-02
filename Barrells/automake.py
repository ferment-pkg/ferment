import os
import subprocess
from index import Barrells
class automake(Barrells):
    def __init__(self):
        self.url="https://ftp.gnu.org/gnu/automake/automake-1.14.tar.gz"
        self.git=False
        self.description="Automake -- Makefile generator"
        self.dependencies=["autoconf"]
    def install(self) -> bool:
        os.chdir(self.cwd)
        os.environ["PERL"]="/usr/bin/perl"
        subprocess.call(["sh","configure", f"--prefix={self.cwd}/built"])
        subprocess.call(["make"])
        subprocess.call(["make","install"])
        os.symlink(f"{self.cwd}/built/bin/aclocal", "/usr/local/bin/aclocal")
        os.symlink(f"{self.cwd}/built/bin/automake", "/usr/local/bin/automake")
        os.symlink(f"{self.cwd}/built/bin/automake-1.14", "/usr/local/bin/automake-1.14")
        os.symlink(f"{self.cwd}/built/bin/aclocal-1.14", "/usr/local/bin/aclocal-1.14")
        super().install()
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/automake")
            os.remove("/usr/local/bin/aclocal")
            os.remove("/usr/local/bin/automake-1.16")
            os.remove("/usr/local/bin/aclocal-1.16")
        finally:
            return super().uninstall()