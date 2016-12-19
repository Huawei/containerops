from component_util import *


if __name__ == "__main__":
    # get env 
    init()
    
    # ComponentStart
    ComponentStart(result="component start ...")
    
    # prepare task
    cmd = "git clone --depth=50 https://github.com/pingcap/tikv.git pingcap/tikv"
    execCommand(cmd)
    
    # TaskStart
    TaskStart(result="task start ...")
    
    
    os.chdir("pingcap/tikv")
    cmd = "make test"
    
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