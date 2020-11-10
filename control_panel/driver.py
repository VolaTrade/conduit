from datetime import datetime
from subprocess import call 
import sys
import os 
import time


def spinup():
    os.system("docker run --network=\"host\" --log-opt max-size=10m --log-opt max-file=5 --restart always -d tickers >> id.txt")
    
    with open("id.txt", "r") as fr:
        container: str = fr.readlines()[-1].replace("\n", "")
    os.system("rm id.txt")
    return container

def start_db(container: str):
    time.sleep(20)
    os.system(f"bash switch_state.sh {container}")

def destroy(container: str):
    os.system(f"docker kill {container}")

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
            update_time = False

        if datetime.now().hour % 4 == 0 and datetime.now().minute == 0 and datetime.now().second <=1:

            update_time = True

if __name__ == '__main__':
    run(None, None)
