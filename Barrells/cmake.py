import os
from index import Barrells
class cmake(Barrells):
    def __init__(self):
        self.url="https://github.com/Kitware/CMake/releases/download/v3.23.1/cmake-3.23.1-macos10.10-universal.tar.gz"
        self.git=False
        self.description="Cross-platform build system"
    def install(self) -> bool:
        os.symlink(f"{self.cwd}/CMake.app/Contents/bin/cmake","/usr/local/bin/cmake")
        os.symlink(f"{self.cwd}/CMake.app/Contents/bin/ccmake","/usr/local/bin/ccmake")
        os.symlink(f"{self.cwd}/CMake.app/Contents/bin/cpack","/usr/local/bin/cpack")
        os.symlink(f"{self.cwd}/CMake.app/Contents/bin/ctest","/usr/local/bin/ctest")
        os.symlink(f"{self.cwd}/CMake.app","/Applications/CMake.app")
        return super().install()
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/cmake")
            os.remove("/usr/local/bin/ccmake")
            os.remove("/usr/local/bin/cpack")
            os.remove("/usr/local/bin/ctest")
            os.remove("/Applications/CMake.app")
        finally:
            return super().uninstall()