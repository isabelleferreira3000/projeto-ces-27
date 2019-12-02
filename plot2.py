import matplotlib.pyplot as plt
import csv

w = []
z = []

with open('BullyImproved.txt','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=' ')
    for row in plots:
        w.append(float(row[2]))
        z.append(float(row[5]))

plt.plot(w,z,'b')
plt.xlabel('Number of processes')
plt.ylabel('Number of messages sent')
plt.title('Improved Bully')
plt.legend()
plt.show()
