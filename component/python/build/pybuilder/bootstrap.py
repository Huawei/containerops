#!/usr/bin/env python3

import subprocess
import os
import sys
import glob
import yaml
import json

REPO_PATH = 'git-repo'


def git_clone(url):
    r = subprocess.run(['git', 'clone', url, REPO_PATH])

    if r.returncode == 0:
        return True
    else:
        print("[COUT] Git clone error: Invalid argument to exit",
              file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return False


def setup(path):
    file_name = os.path.basename(path)
    dir_name = os.path.dirname(path)
    r = subprocess.run('cd {}; python3 {} install'.format(dir_name, file_name),
                       shell=True)

    if r.returncode != 0:
        print("[COUT] install dependences failed: {}".format(path), file=sys.stderr)
        return False

    return True


def pip_install(file_name):
    r = subprocess.run(['pip3', 'install', '-r', file_name])

    if r.returncode != 0:
        print("[COUT] install dependences failed: {}".format(path), file=sys.stderr)
        return False

    return True


def pybuilder(dir_name, task):
    if dir_name and dir_name != '.':
        dir_name = '{}/{}'.format(REPO_PATH, dir_name)
    else:
        dir_name = REPO_PATH

    r = subprocess.run('cd {}; pyb {}'.format(dir_name, task), shell=True)

    if r.returncode != 0:
        print("[COUT] pybuilder error", file=sys.stderr)
        return False

    return True


def echo_yaml(dir_name):
    if dir_name and dir_name != '.':
        dir_name = '{}/{}'.format(REPO_PATH, dir_name)
    else:
        dir_name = REPO_PATH
    for root, dirs, files in os.walk('{}/target'.format(dir_name)):
        for file_name in files:
            if file_name.endswith('.json'):
                data = json.load(open(os.path.join(root, file_name)))
                lines = yaml.safe_dump(data)
                for line in lines.split('\n'):
                    print('[COUT] CO_YAML_CONTENT {}'.format(line))

    return True

def echo_xml(dir_name):
    if dir_name and dir_name != '.':
        dir_name = '{}/{}'.format(REPO_PATH, dir_name)
    else:
        dir_name = REPO_PATH
    for root, dirs, files in os.walk('{}/target'.format(dir_name)):
        for file_name in files:
            if file_name.endswith('.xml'):
                with open(os.path.join(root, file_name), 'r') as f:
                    for line in f.readlines():
                        line = line.rstrip()
                        print('[COUT] CO_XML_CONTENT {}'.format(line))

    return True


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-path', 'task']
    ret = {}
    for s in data.split(' '):
        s = s.strip()
        if not s:
            continue
        arg = s.split('=')
        if len(arg) < 2:
            print('[COUT] Unknown Parameter: [{}]'.format(s))
            continue

        if arg[0] not in validate:
            print('[COUT] Unknown Parameter: [{}]'.format(s))
            continue

        ret[arg[0]] = arg[1]

    return ret


def main():
    argv = parse_argument()
    git_url = argv.get('git-url')
    if not git_url:
        print("[COUT] The git-url value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    entry_path = argv.get('entry-path', '.')
    task = argv.get('task')

    if not git_clone(git_url):
        return

    for file_name in glob.glob('{}/setup.py'.format(REPO_PATH)):
        setup(file_name)

    for file_name in glob.glob('{}/*/setup.py'.format(REPO_PATH)):
        setup(file_name)

    for file_name in glob.glob('{}/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name)

    for file_name in glob.glob('{}/*/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name)

    out = pybuilder(entry_path, task)
    echo_yaml(entry_path)
    echo_xml(entry_path)

    if not out:
        print("[COUT] CO_RESULT = false")
    else:
        print("[COUT] CO_RESULT = true")


main()
