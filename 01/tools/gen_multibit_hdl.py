import sys

line = sys.stdin.readline()
for i in range(1,16):
  print(line.replace("0",str(i)).strip())
