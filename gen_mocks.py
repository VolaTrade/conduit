from os import system, chdir, getcwd 


if __name__ == '__main__':
    
    chdir("internal")
    for direct in ["cache", "client", "connections", "service", "stats"]:
        chdir(direct)
        call: str = f"mockgen -destination=../mocks/mock_{direct}.go -package=mocks . {direct.title()}"
        print(getcwd())
        print(call)
        system(call)
        chdir("../")


    print("Finished generating mocks")
