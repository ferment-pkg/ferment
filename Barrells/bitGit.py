import os
from index import Barrells
import subprocess
class bitgit(Barrells):
    def __init__(self):
        Barrells.__init__(self)
        self.git=False
        self.url="https://github.com/chriswalz/bit/archive/v1.1.2.tar.gz"
        self.description="Bit is a modern Git CLI"
        self.homepage="https://github.com/chriswalz/bit"
        self.sha256="563ae6b0fa279cb8ea8f66b4b455c7cb74a9e65a0edbe694505b2c8fc719b2ff"
        self.license="Apache-2.0"
        self.version="1.1.2"
        self.dependencies=["go", "curl", "git"]
        self.binary="bit"
    def install(self):
        subprocess.run(["go", "build"], cwd=self.cwd)
        return True
    def uninstall(self) -> bool:
        os.remove("/usr/local/bin/bit")
        return super().uninstall()
    def test(self):
        subprocess.run(["git", "init", "/tmp/testDir"])
        subprocess.run(["touch", "/tmp/testDir/test.txt"])
        subprocess.run(["/usr/local/bin/bit", "add", "test.txt"], cwd="/tmp/testDir/")
        output=subprocess.check_output(["/usr/local/bin/bit", "status"], cwd="/tmp/testDir/")
        if b"new file:   test.txt" in output:
            print("True")
            return True
        else:
            print("False")
            return False
        