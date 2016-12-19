#encoding=utf-8
'''
http server 
'''
from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import io,shutil  
import urllib,time
import getopt,string

class MyRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.process(2)

    def do_POST(self):
        self.process(1)

    def do_PUT(self):
	self.process(1)
        
    def process(self,type):
        content =""
        if type==1: #put, post method
            datas = self.rfile.read(int(self.headers['content-length']))
            datas = urllib.unquote(datas).decode("utf-8", 'ignore')
	    print datas
#            datas = eval(data)
            content = str(datas)+"\r\n"
                
                    
            
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

def transDicts(params):
    dicts={}
    if len(params)==0:
        return
    params = params.split('&')
    for param in params:
        dicts[param.split('=')[0]]=param.split('=')[1]
    return dicts
       
if __name__=='__main__':
    
    try:
        server = HTTPServer(('', 8001), MyRequestHandler)
        print 'started httpserver...'
        server.serve_forever()

    except KeyboardInterrupt:
        server.socket.close()
    
    pass

