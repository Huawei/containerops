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
