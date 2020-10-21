from datetime import datetime
from subprocess import call 


def spinup():
    call()


def spindown():


def run():
    
    start = datetime.now()
    while True:
        now = datetime.now()
        if now.hour == start - 1 or (start == 1 and now.hour == 24):
            spinup()
            spindown()
    


if __name__ == '__main__':
    run()
