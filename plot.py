import matplotlib.pyplot as plt
import csv

x = []
y = []
w = []
z = []

with open('BullyNormal.txt','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=' ')
    for row in plots:
        x.append(float(row[2]))
        y.append(float(row[5]))

with open('BullyImproved.txt','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=' ')
    for row in plots:
        w.append(float(row[2]))
        z.append(float(row[5]))

plt.plot(x,y,'r',label='Normal')
plt.plot(w,z,'b',label='Improved')
plt.xlabel('Number of processes')
plt.ylabel('Number of messages sent')
plt.title('Comparing methods')
plt.legend()
plt.show()
