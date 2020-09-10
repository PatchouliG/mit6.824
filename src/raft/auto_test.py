# This is a sample Python script.

# Press ⇧F10 to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.

import sys

# Press the green button in the gutter to run the script.
import os

if __name__ == '__main__':
    import subprocess

    # test_name = sys.argv[1]
    for i in range(1, 20):
        res = os.system("go test -run TestUnreliableChurn2C > tmp")
        if res != 0:
            break

# See PyCharm help at https://www.jetbrains.com/help/pycharm/
