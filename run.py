import os

if __name__ == "__main__":
    n_proccess = input("Qual o numero de processos? ")
    candidate = input("Quem iniciara a eleicao? ")

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.close()

    f = open("results.txt", "w")
    f.close()

    for i in range(n_proccess):
        command = "x-terminal-emulator -e go run main.go " + str(i+1) + " "
        if i+1 != candidate:
            command += "0"            
            print(command)
            os.system(command)

    command = "x-terminal-emulator -e go run main.go " + str(candidate) + " 1"
    print(command)
    os.system(command)
