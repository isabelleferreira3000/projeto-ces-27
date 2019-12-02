import os
import sys
from time import sleep


def file_len(fname):
    i = -1
    with open(fname, "r") as f:
        for i, l in enumerate(f):
            pass
    return i + 1


def run_algorithm(bullyType, n_proccess, candidate):
    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.close()

    results_filename = "results/" + bullyType + "/" + str(n_proccess) + ".txt"

    f = open(results_filename, "w")
    f.close()

    for i in range(n_proccess):
        command = "x-terminal-emulator -x /usr/local/go/bin/go run " + bullyType + ".go " + str(i+1) + " "
        if i+1 != candidate:
            command += "0"
            print(command)
            os.system(command)

    sleep(5)
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


def buildResults(bullyType):
    final_results_path = "results/" + bullyType + ".txt"

    results_path = "results/" + bullyType + "/"
    
    messages = 0

    with open(final_results_path, "w") as final_file:
        for filename in os.listdir(os.getcwd() + "/" + results_path):
            if filename == ".gitkeep":
                continue
            number_of_proccess = int(filename.split(".txt")[0])

            with open(results_path + filename, "r") as f:
                lines = f.readlines()
                for line in lines:
                    if len(line.split("Total messages = ")) > 1:
                        messages = int(line.split("Total messages = ")[1])
                        final_file.write("n = " + str(number_of_proccess) + ", messages = " + str(messages) + "\n")

    with open(final_results_path, "r") as f:
        lines = f.readlines()
        lines.sort()
    with open(final_results_path, "w") as f:
        for line in lines:
            f.write(line)


if __name__ == "__main__":
    n_proccess = -1
    candidate = -1
    if len(sys.argv) == 1:
        n_proccess = int(input("Qual o numero de processos? "))
        candidate = int(input("Quem iniciara a eleicao? "))

        if candidate > n_proccess or candidate <= 0:
            print("Invalid candidate number = ", candidate, " with n_proccess = ", n_proccess)
            sys.exit()

        run_algorithm("BullyNormal", n_proccess, candidate)
        # run_algorithm("BullyImproved", n_proccess, candidate)
    else:
        n_proccess = int(sys.argv[1])
        # for i in range(1, n_proccess+1):
        #     run_algorithm("BullyNormal", i, 1)
        # buildResults("BullyNormal")
        for i in range(1, n_proccess+1):
            command = "x-terminal-emulator -x /usr/local/go/bin/go run MaxHeap.go"
            print(command)
            os.system(command)
            run_algorithm("BullyImproved", i, 1)
        buildResults("BullyNormal")

        # run_algorithm("BullyImproved", n_proccess, candidate)
        # buildResults("BullyImproved")