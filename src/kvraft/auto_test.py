# This is a sample Python script.

# Press ⇧F10 to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.

import sys

# Press the green button in the gutter to run the script.
import os

if __name__ == '__main__':

    all_test = ["2A", "2B", "2C"]
    test_name = ["3A"]
    for test in test_name:
        print("start test " + test)
        times = 1
        for i in range(0, times):
            print(str(i) + "start")
            res = os.system("go test -run " + test + " > test_output")
            if res != 0:
                print("fail int turn " + str(i))
                break
