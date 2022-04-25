from typing import Optional

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
    
    
    
    