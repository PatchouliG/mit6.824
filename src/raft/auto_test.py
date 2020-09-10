# This is a sample Python script.

# Press ⇧F10 to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.

import sys

# Press the green button in the gutter to run the script.
import os

if __name__ == '__main__':

    # test_name = ["2A", "2B"]
    test_name = ["2A"]
    for test in test_name:
        print("start test " + test)
        for i in range(1, 20):
            print(str(i) + "start")
            res = os.system("go test -run " + test + " > tmp")
            if res != 0:
                print("fail int trun " + str(i))
                break
