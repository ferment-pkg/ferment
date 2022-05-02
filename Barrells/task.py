import os
import subprocess
import sys
from index import Barrells
class task(Barrells):
    def __init__(self):
        self.url="https://github.com/GothenburgBitFactory/taskwarrior/releases/download/v2.6.2/task-2.6.2.tar.gz"
        self.git=False
        self.license="MIT"
        self.description="Feature-rich console based todo list manager"
        self.homepage="https://taskwarrior.org/"
        self.dependencies=["cmake","gcc", "gnutils"]
    def install(self):
        subprocess.call(["cmake","-DCMAKE_BUILD_TYPE=release","."], cwd=self.cwd, stdout=sys.stdout)
        subprocess.call(["make"], cwd=self.cwd, stdout=sys.stdout)
        subprocess.call(["make","install"], cwd=self.cwd, stdout=sys.stdout)
        return super().install()
    def uninstall(self) -> bool:
        try:
            os.remove("/usr/local/bin/task")
        finally:
            return super().uninstall()
        
        