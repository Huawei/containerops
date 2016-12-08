from component_util import *


if __name__ == "__main__":
    # get env 
    init()
    
    # ComponentStart
    ComponentStart()
    
    # prepare task
    cmd = "git clone --depth=50 --branch=master https://github.com/pingcap/pd.git /root/gopath/src/github.com/pingcap/pd"
    execCommand(cmd)
    
    # TaskStart
    TaskStart()
    
    
    os.chdir("/root/gopath/src/github.com/pingcap/pd")
    cmd = "make dev"
    
    # Task exec
    execCommand(cmd)
    
    
    # TaskStatus
    TaskStatus({"status":"finish"})
    
    #ComponentStop
    ComponentStop()
