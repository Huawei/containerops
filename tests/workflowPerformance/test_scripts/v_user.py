import mechanize
import random
import time
import uuid
import json

class Transaction(object):
    def __init__(self):
        self.custom_timers = {}

    def run(self):
	br = mechanize.Browser()
	br.set_handle_robots(False)
	start_time = time.time()
	url = "http://192.168.10.131:10000/pipeline/v1/demo/demo/pythoncheck"
	data = {"gitUrl":"https://github.com/xiechuanj/python-sonar-runner.git"} 
	data = json.dumps(data)
	resp = br.open(url, data)
	resp.read()

	latency = time.time() - start_time

	self.custom_timers['workflow'] = latency
	assert (resp.code == 200), 'Bad HTTP Response'
	assert ('pipeline start' in resp.get_data()), 'Failed Content Verification'


if __name__ == '__main__':
    trans = Transaction()
    trans.run()
    print trans.custom_timers
