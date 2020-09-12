# This is a sample Python script.

# Press ⇧F10 to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.

import sys

# Press the green button in the gutter to run the script.
import os

if __name__ == '__main__':

    # test_name = ["2A", "2B","",]
    # test_name = ["2A", "2B"]
    all_test = ["2A", "2B", "2C"]
    # test_name = ["TestUnreliableAgree2C"]
    # test_name = ["TestFigure8Unreliable2C"]
    test_name = ["2C"]
    for test in all_test:
        print("start test " + test)
        times = 7
        for i in range(0, times):
            print(str(i) + "start")
            res = os.system("go test -run " + test + " > test_output")
            if res != 0:
                print("fail int turn " + str(i))
                break
