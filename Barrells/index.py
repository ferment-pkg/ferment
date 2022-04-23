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
        pass
    def uninstall(self)->bool:
        pass
    def build(self)->bool:
        pass
    def update(self)->bool:
        pass
    def test(self)->bool:
        pass
    
    
    
    