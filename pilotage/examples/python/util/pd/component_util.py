#encoding=utf-8
             
import logging, logging.handlers, json, urllib2, urllib, os, io, commands, shutil
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
    CO_INFO['CO_status'] = True   
    eventlsts = CO_INFO[CO_EVENT_LIST].split(";")
    events = {}
    for eventlst in eventlsts:
        events[eventlst.split(",")[0]] = eventlst.split(",")[1]
        CO_INFO['CO_EVENT_LIST'] = events
    logger.debug("init done...")
    logger.debug(CO_INFO)

def NotifyEvent(url, paras):
    logger.debug(url,paras)
    jdata = json.dumps(paras) #{'':''}
    request = urllib2.Request(url, jdata)
    request.add_header('Content-Type', 'application/json')
    request.get_method = lambda:'POST'
    request = urllib2.urlopen(request)
    result = request.read() 
    logger.debug(result)
    return result  


def  ComponentStart():    
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_COMPONENT_START,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_COMPONENT_START])}   
    result = NotifyEvent(CO_INFO[CO_COMPONENT_START], paras) 
    return result


def  ComponentStop():    
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_COMPONENT_STOP,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_COMPONENT_STOP])}    
    result = NotifyEvent(CO_INFO[CO_COMPONENT_STOP], paras)
    return result      

def  TaskStart():    
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_TASK_START,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_TASK_START])}    
    result = NotifyEvent(CO_INFO[CO_TASK_START], paras)    
    return result  

def  TaskResult(info={}):    
    paras = {"RUN_ID":CO_INFO[CO_RUN_ID],"EVENT":CO_TASK_RESULT,"EVENT_ID":int(CO_INFO[CO_EVENT_LIST][CO_TASK_RESULT]),"INFO":info}    
    result = NotifyEvent(CO_INFO[CO_TASK_RESULT], paras)
    logger.debug(result)    
    return result  

def  TaskStatus(info={}):    
    paras = {"RUN_ID":CO_INFO["CO_RUN_ID"],"EVENT":"CO_COMPONENT_START","EVENT_ID":int(CO_INFO['CO_EVENT_LIST']['CO_COMPONENT_START']), "INFO":info}    
    result = NotifyEvent(CO_INFO["CO_COMPONENT_START"], paras)
    logger.debug(result)    
    return result  


def execCommand(cmd):
    (status, output) = commands.getstatusoutput(cmd)      
    if status == 0:
        CO_INFO['CO_status'] = True
    else:
        CO_INFO['CO_status'] = False
      
    CO_INFO['co_out'] = output
    logger.debug(CO_INFO['CO_status'])
    logger.debug(CO_INFO['co_out'])
    TaskResult({"result":CO_INFO['co_out'], "status":CO_INFO['CO_status']})
    return  CO_INFO['CO_status'], CO_INFO['co_out'] 


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
