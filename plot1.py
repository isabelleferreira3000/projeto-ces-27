import matplotlib.pyplot as plt
import csv

x = []
y = []

with open('BullyNormal.txt','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=' ')
    for row in plots:
        x.append(float(row[2]))
        y.append(float(row[5]))

plt.plot(x,y,'r')
plt.xlabel('Number of processes')
plt.ylabel('Number of messages sent')
plt.title('Normal Bully')
plt.legend()
plt.show()
