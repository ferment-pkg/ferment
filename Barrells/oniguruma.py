import subprocess
from index import Barrells
import os
class oniguruma(Barrells):
    def __init__(self):
        self.url="https://github.com/kkos/oniguruma/releases/download/v6.9.8/onig-6.9.8.tar.gz"
        self.git=False
        self.description="Oniguruma is a modern and flexible regular expressions library."
        self.dependencies=["autoconf", "automake", "libtool"]
        self.lib=True
    def install(self) -> bool:
        os.chdir(self.cwd)
        subprocess.run(["autoreconf", "-vfi"], timeout=1200)
        subprocess.run(["sh", "./configure", f"--prefix={self.cwd}/built"], timeout=1200)
        subprocess.call(["make"], timeout=1200)
        subprocess.call(["make", "install"], timeout=1200)
        os.symlink(f"{self.cwd}/built/bin/onig-config", "/usr/local/bin/onig-config")
        os.symlink(f"{self.cwd}/built/include/oniggnu.h", "/usr/local/include/oniggnu.h")
        os.symlink(f"{self.cwd}/built/include/onigguruma.h", "/usr/local/include/onigguruma.h")
        dirs = filter(os.path.isdir, os.listdir(f"{self.cwd}/built/lib"))
        for f in dirs:
            os.symlink(f"{self.cwd}/built/lib/{f}", f"/usr/local/lib/{f}")
        os.symlink(f"{self.cwd}/built/lib/pkgconfig/oniguruma.pc", "/usr/local/lib/pkgconfig/oniguruma.pc")
        
        return super().install()
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/onig-config")
            os.remove("/usr/local/include/oniggnu.h")
            os.remove("/usr/local/include/onigguruma.h")
            dirs = filter(os.path.isdir, os.listdir(f"{self.cwd}/built/lib"))
            for f in dirs:
                os.remove(f"/usr/local/lib/{f}")
            os.remove("/usr/local/lib/pkgconfig/oniguruma.pc")
        finally:
            return super().uninstall()