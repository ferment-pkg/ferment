import subprocess
from dataclasses import dataclass


@dataclass
class ManagerVersion():
    latestVersion: str
    currentVersion:str


def install(package)->int:
    if package is None:
        return 1
    return subprocess.check_call(["ferment", "install", package])
def upgrade()->int:
    return subprocess.check_call(["ferment", "upgrade"])
def remove(package)->int:
    if package is None:
        return 1
    return subprocess.check_call(["ferment", "uninstall", package])

def update()->int:
    return subprocess.check_call(["ferment", "update"])
def getManagerVersion()->ManagerVersion:
    version=subprocess.check_output(["ferment", "--version"])
    versions=version.decode("utf-8").split("\n")
    versions[0]=versions[0].replace("Ferment:  ", "")
    versions[1]=versions[1].replace("Latest Version:  ", "")
    return ManagerVersion(versions[0], versions[1])

