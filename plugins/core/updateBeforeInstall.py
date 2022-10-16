from dataclasses import dataclass

import ferment


def beforeInstall():
    version=ferment.getManagerVersion()
    print(f"Current Version Of Ferment Is: {version.currentVersion}")
    if version.currentVersion != version.latestVersion:
        print("Updating Ferment")
        ferment.upgrade()
        print("Updated Ferment")
    else:
        print("Ferment Is Up To Date")
