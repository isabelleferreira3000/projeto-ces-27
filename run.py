import os
from time import sleep


def file_len(fname):
    i = -1
    with open(fname, "r") as f:
        for i, l in enumerate(f):
            pass
    return i + 1


if __name__ == "__main__":
    n_proccess = input("Qual o numero de processos? ")
    candidate = input("Quem iniciara a eleicao? ")

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.close()

    f = open("results.txt", "w")
    f.close()

    for i in range(n_proccess):
        command = "x-terminal-emulator -e /usr/local/go/bin/go run BullyNormal.go " + str(i+1) + " "
        if i+1 != candidate:
            command += "0"            
            print(command)
            os.system(command)

    command = "x-terminal-emulator -e /usr/local/go/bin/go run BullyNormal.go " + str(candidate) + " 1"
    print(command)
    os.system(command)

    while file_len("results.txt") != n_proccess:
        pass

    with open("results.txt", "r") as f:
        lines = f.readlines()
        lines.sort()

    total_messages = 0
    with open("results.txt", "w") as f:
        for line in lines:
            f.write(line)
            total_messages += int(line.split("with ")[1].split(" ")[0])
        f.write("Total messages = " + str(total_messages))
