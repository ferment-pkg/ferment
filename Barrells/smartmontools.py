import os
import subprocess
from index import Barrells
class smartmontools(Barrells):
    def __init__(self):
        self.url='https://downloads.sourceforge.net/project/smartmontools/smartmontools/7.3/smartmontools-7.3.tar.gz'
        self.git=False
        self.description="smartmontools is a set of utilities to monitor hard drives and other storage devices."
        self.dependencies=["autoconf", "automake"]
    def install(self) -> bool:
        os.chdir(self.cwd)
        args=["--disable-dependency-tracking", "--with-savestates", "--with-attributelog"]
        subprocess.call(["sh","configure", f"--prefix={self.cwd}/built", *args])
        subprocess.call(["make"])
        subprocess.call(["make","install"])
        os.symlink(os.path.join(self.cwd, "built", "sbin", "smartctl"), '/usr/local/bin/smartctl')
        os.symlink(os.path.join(self.cwd, "built", "sbin", "smartd"),  '/usr/local/bin/smartd')
        return super().install()
    def uninstall(self) -> bool:
        try:
            os.remove(os.path.join("/usr/local/", "bin", "smartctl"))
            os.remove(os.path.join("/usr/local/" "bin", "smartd"))
        finally:
            return super().uninstall()
    def test(self) -> bool:
        try:
            subprocess.call(["smartctl", "--version"])
            subprocess.call(["smartd", "--version"])
            return super().test()
        except:
            return False
