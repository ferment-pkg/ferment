import os
from index import Barrells

class htop(Barrells):
    def __init__(self):
        self.url="https://github.com/htop-dev/htop/archive/3.2.0.tar.gz"
        self.description="Improved top (interactive process viewer)"
        self.dependencies=["autoconf", "automake", "libtool"]
        self.homepage="https://htop.dev/"
        self.git=False
        self.version="3.2.0"
    def install(self) -> bool:
        os.chdir(self.cwd)
        self.runcmdincwd(["sh", "autogen.sh"])
        args=[f"--prefix={self.cwd}/built"]
        self.runcmdincwd(["./configure", *args])
        self.runcmdincwd(["make"])
        self.runcmdincwd(["make", "install"])
        os.symlink(f"{self.cwd}/built/bin/htop", "/usr/local/bin/htop")
        return super().install()
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/htop")
        finally: 
            return super().uninstall()
        
