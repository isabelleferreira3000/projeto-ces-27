if __name__ == "__main__":
    n_proccess = input("Qual o numero de processos? ")
    candidate = input("Quem iniciara a eleicao? ")

    f = open("params.txt", "w")
    f.write(str(n_proccess))
    f.write(" ")
    f.write(str(candidate))
    f.close()
