
# Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
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
             
from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import io,shutil
import urllib,time
import getopt,string

import threading
import logging
import Queue
import urllib2
import json
import os,time
import shutil
import subprocess
import time,signal
import threading

import logging  
import logging.handlers  
  
LOG_FILE = 'tst.log'

handler = logging.handlers.RotatingFileHandler(LOG_FILE, maxBytes = 1024*1024, backupCount = 5)
fmt = '%(asctime)s - %(filename)s:%(lineno)s - %(name)s - %(message)s'

formatter = logging.Formatter(fmt)
handler.setFormatter(formatter)

logger = logging.getLogger('tst')
logger.addHandler(handler)
logger.setLevel(logging.DEBUG)

import Queue
q = Queue.Queue()




class Listener(BaseHTTPRequestHandler):

    def do_POST(self):
        content =""
      
        datas = self.rfile.read(int(self.headers['content-length']))
        datas = urllib.unquote(datas).decode("utf-8", 'ignore')
        content = str(datas) +"\r\n" #datas is dict string, if no data, input {"":""}
	q.put(content) # store in queue
	
        
        enc="UTF-8"  
        content = content.encode(enc)          
        f = io.BytesIO()  
        f.write(content)  
        f.seek(0)  
        self.send_response(200)  
        self.send_header("Content-type", "text/html; charset=%s" % enc)  
        self.send_header("Content-Length", str(len(content)))  
        self.end_headers()  
        shutil.copyfileobj(f,self.wfile)  # return
	       



class ExecTask(threading.Thread):
    
    def __init__(self, info):
        threading.Thread.__init__(self)
        self.info = info
	
    """
	get require	
    """
    def _http_get(self, url):
	
        response = urllib2.urlopen(url)
        return response.read()

    """
	put require
    """
    def _http_put(self, url, paras):
	try:
            jdata = json.dumps(paras) #{'':''}
            request = urllib2.Request(url, jdata)
            request.add_header('Content-Type', 'application/json')
            request.get_method = lambda:'PUT'
            request = urllib2.urlopen(request)
            return request.read()
        except:
	    self.info['status'] = False
            return ""

	
    """
	exec shell script
    """    
    def _execCommand(self, cmd):
        try:
            p = subprocess.Popen(cmd, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE, shell = True)
	    out, err = p.communicate()
            self.info["out"] = out

        except:
            self.info['status'] = False

    """
	generate sonar project file
    """    
    def _writeFile(self):
        with open(self.info['CO_DATA']["filename"], 'w') as f:            
            f.write(self.info['CO_DATA']["contents"])
 
    """
	parse git url, get repo
    """
    def _parserGitUrl(self):
        if self.info["gitUrl"][-4:] == ".git":
            self.info["project"]=self.info["gitUrl"][:-4].split("/")[-1]    
	else:
	    self.info['status'] = False
            

    """
	parser EventList, get eventid, event
    """
    def _parserEventList(self):
	eventlsts = self.info["EVENT_LIST"].split(";")
        events = {}
	for eventlst in eventlsts:
	    events[eventlst.split(",")[0]] = eventlst.split(",")[1]
        self.info['EVENT_LIST'] = events

    """
	init env
    """
    def _initEvent(self):
	try:
            self.info['CO_DATA'] = eval(os.getenv("CO_DATA"))
            self.info["COMPONENT_START"] = os.getenv("COMPONENT_START")
            self.info["COMPONENT_STOP"] = os.getenv("COMPONENT_STOP")
            self.info["TASK_START"] = os.getenv("TASK_START")
            self.info["TASK_RESULT"] = os.getenv("TASK_RESULT")
            self.info["TASK_STATUS"] = os.getenv("TASK_STATUS")
            self.info["REGISTER_URL"] = os.getenv("REGISTER_URL")
	    self.info['EVENT_LIST'] = os.getenv("EVENT_LIST")
	    self.info['RUN_ID'] = os.getenv("RUN_ID")
	    self.info['SERVICE_PORT'] = ":" + (os.getenv('SERVICE_ADDR')).split(":")[1] + "/"
	    logger.debug(self.info['SERVICE_PORT']) 
            self.info['POD_NAME'] = os.getenv("POD_NAME") 
        except:

            self.info = {"COMPONENT_START": "http://172.17.0.1:8001","COMPONENT_STOP": "http://172.17.0.1:8001", "TASK_START": "http://172.17.0.1:8001", "TASK_RESULT": "http://172.17.0.1:8001", "TASK_STATUS": "http://172.17.0.1:8001", "REGISTER_URL": "http://172.17.0.1:8001", "EVENT_LIST": "COMPONENT_START,1;COMPONENT_STOP,2;TASK_START,8;TASK_RESULT,9;TASK_STATUS,10;REGISTER_URL,11","RUN_ID": "1","POD_NAME":"go-doc1-2-1-178-1-pod", "CO_DATA":{"serverities":"MAJOR","SERVICE_ADDR": "127.0.0.1:30001:9999", "RUN_ID": "1","POD_NAME":"go-doc1-2-1-178-1-pod","gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git", "contents": "sonar.projectKey=python-sonar-runner\nsonar.projectName=python-sonar-runner\nsonar.projectVersion=1.0\nsonar.sources=src\nsonar.language=py\nsonar.sourceEncoding=UTF-8","filename": "sonar-project.properties"}}
	self.info['status'] = True   
        self.info['CO_DATA']['sonarurl'] = "http://localhost:9000/api/issues/search?projectKeys=python-sonar-runner&severities=%s&statuses=OPEN" 

	    

    
    """
	run check 
    """	
    def run(self):
       	logger.debug("start .....") 
	self._initEvent() 
	logger.debug("init end")
	self._parserEventList()
	logger.debug( "send component start, %s, %s" % (self.info["COMPONENT_START"],{"RUN_ID":self.info["RUN_ID"],"EVENT":"COMPONENT_START","EVENTID":self.info['EVENT_LIST']['COMPONENT_START']}))
	result = self._http_put(self.info["COMPONENT_START"], {"RUN_ID":self.info["RUN_ID"],"EVENT":"COMPONENT_START","EVENTID":int(self.info['EVENT_LIST']['COMPONENT_START'])})
	logger.debug("start sonar ...%r"%result)
	out = os.popen("./bin/run.sh > sonar.log &")
	time.sleep(10)
	logger.debug("register ....")
	result = self._http_put(self.info["REGISTER_URL"], {"RUN_ID":self.info["RUN_ID"],"POD_NAME":self.info['POD_NAME'],"RECEIVE_URL":self.info['SERVICE_PORT']})
	logger.debug("wait data...%r"%result)	
	     
        try:
            server = HTTPServer(('', 8000), Listener)
            server.handle_request()
	    self.info['gitUrl'] =  eval(q.get())['gitUrl']
	    
	    logger.debug("TASK_START...")
	    result = self._http_put(url=self.info["TASK_START"], paras={"RUN_ID":self.info["RUN_ID"],"EVENT":"TASK_START","EVENTID":int(self.info['EVENT_LIST']['TASK_START'])})
	    logger.debug("TASK_START...%r"%result)
    	except KeyboardInterrupt:
	    server.socket.close()
	    self.info['status'] = False
	logger.debug("status...%s"%self.info['status'])
        self._parserGitUrl()
	logger.debug("download code from git...%s"%self.info["gitUrl"])
        self._execCommand("git clone " + self.info["gitUrl"])
	logger.debug("download finish")
	path = os.getcwd() + "/" +  self.info["project"]
        os.chdir(path)
        self._writeFile()
        self._http_put(url=self.info["TASK_STATUS"], paras={"RUN_ID":self.info["RUN_ID"],"EVENT":"TASK_STATUS","EVENTID":int(self.info['EVENT_LIST']['TASK_STATUS']), "INFO":{"TASK_STATUS":"RUNNING"}})
	self._execCommand("/opt/sonar-scanner/bin/sonar-scanner")
        self._http_put(url=self.info["TASK_STATUS"], paras={"RUN_ID":self.info["RUN_ID"],"EVENT":"TASK_STATUS","EVENTID":int(self.info['EVENT_LIST']['TASK_STATUS']), "INFO":{"TASK_STATUS":"GET RESULT"}})
	time.sleep(10)
        searchRes = {}
        if self.info['status'] == True:
	    searchRes = self._http_get((self.info['CO_DATA']['sonarurl'])%self.info['CO_DATA']['serverities'])
            if '"total":0' in searchRes:
		pass
	    else:
		self.info['status'] = False

        self._http_put(url=self.info["TASK_RESULT"], paras={"RUN_ID":self.info["RUN_ID"],"EVENT":"TASK_RESULT","EVENTID":int(self.info['EVENT_LIST']['TASK_RESULT']), "INFO":{'status':self.info['status'],"result":searchRes}})
      

        self._http_put(url=self.info["COMPONENT_STOP"], paras={"RUN_ID":self.info["RUN_ID"],"EVENT":"COMPONENT_STOP","EVENTID":int(self.info['EVENT_LIST']['COMPONENT_STOP'])})



if __name__ == "__main__":

    ExecTask({}).start()
    while True:
	time.sleep(5)
