import os
import subprocess
from index import Barrells
class libzip(Barrells):
    def __init__(self):
        self.url="https://libzip.org/download/libzip-1.8.0.tar.gz"
        self.git=False
        self.description="A C library for reading, creating, and modifying zip archives"
        self.dependencies=["cmake"]
        self.lib=True
    def install(self) -> bool:
        args=["-DENABLE_GNUTLS=OFF", "-DENABLE_MBEDTLS=OFF", "-DENABLE_OPENSSL=OFF", "-DBUILD_REGRESS=OFF", "-DBUILD_EXAMPLES=OFF"] 
        subprocess.call(["cmake","-DCMAKE_BUILD_TYPE=release"," ".join(args),"."], cwd=self.cwd)
        subprocess.call(["make", "install"], cwd=self.cwd)
        os.symlink(f"{self.cwd}/lib/zip.h","/usr/local/include/zip.h")
        return super().install()
    def uninstall(self) -> bool:
        os.remove("/usr/local/include/zip.h")
        return super().uninstall()