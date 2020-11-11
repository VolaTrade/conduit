from datetime import datetime
from subprocess import call 
import sys
import os 
import time


def get_version() -> str:
   
    with open("version", "r") as f:
        version = f.readlines()[0]
    
    return version

def is_new_version(curr_version) -> bool:
    if curr_version != get_version():
        return True
    
    return False
        
def spinup():
    os.system("make docker-run >> id.txt")
    
    with open("id.txt", "r") as fr:
        container: str = fr.readlines()[-1].replace("\n", "")
    os.system("rm id.txt")
    return container

def start_db(container: str):
    time.sleep(20)
    os.system(f"bash switch_state.sh {container}")

def destroy(container: str):
    os.system(f"kill {container}")

def run(blue_container: str, green_container: str):
    update_time: bool = False
    while True: 
        if blue_container is None:
            blue_container = spinup()
            start_db(blue_container)

        if update_time is True:
            green_container = spinup()
            destroy(blue_container)
            time.sleep(10)
            start_db(green_container)
            blue_container = green_container
            green_container = None
            update_time = False

        if datetime.now().hour % 4 == 0 and datetime.now().minute == 0 and datetime.now().second == 0:
            update_time = True


if __name__ == '__main__':
    run(None, None)

