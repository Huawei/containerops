from component_util import *


if __name__ == "__main__":
    # get env 
    init()

    # ComponentStart
    ComponentStart(result="component start ...")

    # prepare task
    cmd = "git clone https://github.com/pingcap/pd.git /root/gopath/src/github.com/pingcap/pd"
    execCommand(cmd)

    # TaskStart
    TaskStart(result="task start ...")


    os.chdir("/root/gopath/src/github.com/pingcap/pd")
    cmd = "make dev"

    # Task exec
    status = execCommand(cmd)

    # reesult status
    TaskResult(result="task result ...", status=status)

    # TaskStatus
    TaskStatus(result="task status ...", status=status)

    #ComponentStop
    ComponentStop(result="component stop ...")

    # wait
    holdWait()

