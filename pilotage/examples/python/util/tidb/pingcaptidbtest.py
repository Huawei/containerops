# Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from component_util import *


if __name__ == "__main__":
    # get env 
    init()

    # ComponentStart
    ComponentStart(result="component start ...")

    # prepare task
    cmd = "git clone https://github.com/pingcap/tidb.git /root/gopath/src/github.com/pingcap/tidb"
    execCommand(cmd)

    # TaskStart
    TaskStart(result="task start ...")

    os.chdir("/root/gopath/src/github.com/pingcap/tidb")

    execCommand("rm -rf store/tikv/*_slow_test.go")

    # Task exec
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