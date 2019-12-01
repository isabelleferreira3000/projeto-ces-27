import os

if __name__ == "__main__":
    n_proccess = input("Qual o numero de processos? ")
    candidate = input("Quem iniciara a eleicao? ")

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.write("\n")
    f.write(str(candidate))
    f.close()

    for i in range(n_proccess):
        command = "x-terminal-emulator -e go run main.go " + str(i+1) + " "
        if i+1 == candidate:
            command += "1"
        else:
            command += "0"

        print(command)
        os.system(command)
