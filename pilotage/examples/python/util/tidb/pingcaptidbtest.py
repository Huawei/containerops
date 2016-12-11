from component_util import *


if __name__ == "__main__":
    # get env 
    init()

    # ComponentStart
    ComponentStart()

    # prepare task
    cmd = "git clone https://github.com/pingcap/tidb.git"
    execCommand(cmd)

    # TaskStart
    TaskStart()

    os.chdir("/root/gopath/src/github.com/pingcap/tidb")

    execCommand("rm -rf store/tikv/*_slow_test.go")

    # Task exec
    cmd = "make dev"


    # Task exec
    status = execCommand(cmd)

    # reesult status
    TaskResult({"status": status,})


    # TaskStatus
    TaskStatus({"status":status,})

    #ComponentStop
    ComponentStop()

    # wait
    holdWait()

