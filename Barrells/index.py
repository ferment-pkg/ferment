from typing import Optional
import os
class Barrells:
    def __init__(self):        
        self.url:str
        self.description:str
        self.homepage:str
        self.version:str
        self.git:bool
        self.license:str
        self.mirror:Optional[list[str]]
        self.sha256:str
        self.supported_OS:list[str]
        self.dependencies:list[str]
        self.binary:str
        self.lib:bool
        #Provided by the cmd
        self.cwd:str
    def install(self)->bool:
        print("True")
        return True
    def uninstall(self)->bool:
        print("True")
        return True
    def build(self)->bool:
        print("True")
        return True
    def update(self)->bool:
        print("True")
        return True
    def test(self)->bool:
        print("True")
        return True
    def download(self)->bool:
        print("True")
        return True
    # Helper Functions DON'T EDIT
    def EditPath(self, path:str):
        home=os.getenv("HOME")
        with open(f"{home}/.zshrc", "a") as f:
            f.write(f"PATH=$PATH:{path}")