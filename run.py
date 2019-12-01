import os

if __name__ == "__main__":
    n_proccess = input("Qual o numero de processos? ")
    candidate = input("Quem iniciara a eleicao? ")

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.write("\n")
    f.write(str(candidate))
    f.close()

    string_params = ""
    for i in range(n_proccess):
        tam_i = len(str(i))
        string_params += ":8"
        for j in range(3-tam_i):
            string_params += "0"
        num = int(i)
        for j in range(tam_i):
            string_params += int(str(num)[0])
            num = int(num/10)

    for i in range(n_proccess):
        os.system("x-terminal-emulator -e go run main.go")
