from datetime import datetime
from subprocess import call 
import sys
import os 
import time


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

def run(blue_container, green_container):
    update_time: bool = False 
    while True: 
        if blue_container == None:
            blue_container = spinup()
            start_db(blue_container)

        if update_time:
            green_container = spinup()
            destroy(blue_container)
            start_db(green_container)
            blue_container = green_container
            green_container = None

        if datetime.now().minute() == 42:
            update_time = True 
    


if __name__ == '__main__':
 
    run(None, None)

