import os
import sys
from time import sleep


def file_len(fname):
    i = -1
    with open(fname, "r") as f:
        for i, l in enumerate(f):
            pass
    return i + 1


def run_algorithm(bullyType):
    n_proccess = -1
    candidate = -1
    if len(sys.argv) == 1:
        n_proccess = input("Qual o numero de processos? ")
        candidate = input("Quem iniciara a eleicao? ")
    else:
        n_proccess = int(sys.argv[1])
        candidate = int(sys.argv[2])

    if candidate > n_proccess or candidate <= 0:
        print("Invalid candidate number = ", candidate, " with n_proccess = ", n_proccess)
        sys.exit()

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.close()

    results_filename = "results/results_" + str(n_proccess) + "_" + str(candidate) + ".txt"

    f = open(results_filename, "w")
    f.close()

    for i in range(n_proccess):
        command = "x-terminal-emulator -x /usr/local/go/bin/go run " + bullyType + ".go " + str(i+1) + " "
        if i+1 != candidate:
            command += "0"
            print(command)
            os.system(command)

    command = "x-terminal-emulator -x /usr/local/go/bin/go run " + bullyType + ".go " + str(candidate) + " 1"
    print(command)
    os.system(command)

    while file_len(results_filename) != n_proccess:
        pass

    with open(results_filename, "r") as f:
        lines = f.readlines()
        lines.sort()

    total_messages = 0
    with open(results_filename, "w") as f:
        for line in lines:
            f.write(line)
            total_messages += int(line.split("with ")[1].split(" ")[0])
        f.write("Total messages = " + str(total_messages))


if __name__ == "__main__":
    run_algorithm("BullyNormal")
