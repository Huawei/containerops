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

# encoding=utf-8

import logging, logging.handlers, json, urllib2, urllib, os, io, commands, shutil, time
from BaseHTTPServer import BaseHTTPRequestHandler

# define const

LOG_FILE = 'tst.log'
CO_POD_NAME     = "CO_POD_NAME"
CO_RUN_ID       = "CO_RUN_ID"
CO_EVENT_LIST   = "CO_EVENT_LIST"
CO_SERVICE_ADDR = "CO_SERVICE_ADDR"
CO_COMPONENT_START = "CO_COMPONENT_START"
CO_COMPONENT_STOP  = "CO_COMPONENT_STOP"
CO_TASK_START  = "CO_TASK_START"
CO_TASK_RESULT = "CO_TASK_RESULT"
CO_TASK_STATUS = "CO_TASK_STATUS"
CO_REGISTER_URL = "CO_register"
CO_DATA = "CO_DATA"
CO_SET_GLOBAL_VAR_URL = "CO_SET_GLOBAL_VAR_URL"
CO_LINKSTART_TOKEN = "CO_LINKSTART_TOKEN"
CO_LINKSTART_URL   = "CO_LINKSTART_URL"
INFO = "INFO"
#define info
CO_INFO = {}


# define log format
handler = logging.handlers.RotatingFileHandler(LOG_FILE, maxBytes = 1024*1024, backupCount = 5)
fmt = '%(asctime)s - %(filename)s:%(lineno)s - %(name)s - %(message)s'

formatter = logging.Formatter(fmt)
handler.setFormatter(formatter)

logger = logging.getLogger('tst')
logger.addHandler(handler)
logger.setLevel(logging.DEBUG)


def init():
    logger.debug("start init...")
    CO_INFO[CO_DATA] = os.getenv("CO_DATA")
    CO_INFO[CO_COMPONENT_START] = os.getenv("CO_COMPONENT_START")
    CO_INFO[CO_COMPONENT_STOP] = os.getenv("CO_COMPONENT_STOP")
    CO_INFO[CO_TASK_START] = os.getenv("CO_TASK_START")
    CO_INFO[CO_TASK_RESULT] = os.getenv("CO_TASK_RESULT")
    CO_INFO[CO_TASK_STATUS] = os.getenv("CO_TASK_STATUS")
    CO_INFO[CO_REGISTER_URL] = os.getenv("CO_REGISTER_URL")
    CO_INFO[CO_EVENT_LIST] = os.getenv("CO_EVENT_LIST")
    CO_INFO[CO_RUN_ID] = os.getenv("CO_RUN_ID")
    CO_INFO[CO_POD_NAME] = os.getenv("CO_POD_NAME")
    eventlsts = CO_INFO[CO_EVENT_LIST].split(";")
    events = {}
    for eventlst in eventlsts:
        events[eventlst.split(",")[0]] = eventlst.split(",")[1]
        CO_INFO['CO_EVENT_LIST'] = events
    logger.debug("init done...")
    logger.debug(CO_INFO)

def NotifyEvent(url, paras):
    jdata = json.dumps(paras) #{'':''}
    request = urllib2.Request(url, jdata)
    request.add_header('Content-Type', 'application/json')
    request.get_method = lambda:'POST'
    request = urllib2.urlopen(request)
    result = request.read()
    logger.debug("%s,%s,%s"%(url,paras,result))
    return result


def  ComponentStart(output="",result="component start ...",status=True):
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_COMPONENT_START,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_COMPONENT_START]), INFO:{"output":output,"result":result,"status":status}}
    result = NotifyEvent(CO_INFO[CO_COMPONENT_START], paras)
    return result


def  ComponentStop(output="",result="component stop ...",status=True):
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_COMPONENT_STOP,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_COMPONENT_STOP]), INFO:{"output":output,"result":result,"status":status}}
    result = NotifyEvent(CO_INFO[CO_COMPONENT_STOP], paras)
    return result

def  TaskStart(output="",result="task start ...",status=True):
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_TASK_START,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_TASK_START]), INFO:{"output":output,"result":result,"status":status}}
    result = NotifyEvent(CO_INFO[CO_TASK_START], paras)
    return result

def  TaskResult(output="",result="task result ...",status=True):
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_TASK_RESULT,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_TASK_RESULT]), INFO:{"output":output,"result":result,"status":status}}
    result = NotifyEvent(CO_INFO[CO_TASK_RESULT], paras)
    return result

def  TaskStatus(output="",result="task status ...",status=True):
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_TASK_STATUS,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_TASK_STATUS]), INFO:{"output":output,"result":result,"status":status}}
    result = NotifyEvent(CO_INFO[CO_TASK_STATUS], paras)
    return result


def execCommand(cmd):
    TaskResult(output=cmd, result=os.getcwd())
    import subprocess
    subp=subprocess.Popen(cmd,shell=True,stdout=subprocess.PIPE)
    while subp.poll()==None:
        message = subp.stdout.readline()
        logger.debug(message)
        TaskResult(result=message)

    status =  subp.returncode
    if status == 0:
        return True
    else:
        return False


def holdWait():
    while True:
        time.sleep(100)

class Listener(BaseHTTPRequestHandler):

    def getData(self):

        datas = self.rfile.read(int(self.headers['content-length']))
        datas = urllib.unquote(datas).decode("utf-8", 'ignore')
        content = str(datas) +"\r\n"


        enc="UTF-8"
        content = content.encode(enc)
        f = io.BytesIO()
        f.write(content)
        f.seek(0)
        self.send_response(200)
        self.send_header("Content-type", "text/html; charset=%s" % enc)
        self.send_header("Content-Length", str(len(content)))
        self.end_headers()
        shutil.copyfileobj(f,self.wfile)

