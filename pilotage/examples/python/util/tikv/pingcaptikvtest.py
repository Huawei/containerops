from component_util import *


if __name__ == "__main__":
    # get env 
    init()
    
    # ComponentStart
    ComponentStart()
    
    # prepare task
    cmd = "git clone --depth=50 https://github.com/pingcap/tikv.git pingcap/tikv"
    execCommand(cmd)
    
    # TaskStart
    TaskStart()
    
    
    os.chdir("pingcap/tikv")
    cmd = "make test"
    
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
