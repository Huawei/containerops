#!/usr/bin/env python3

import subprocess
import os
import sys
import json

REPO_PATH = 'git-repo'


def git_clone(url):
    r = subprocess.run(['git', 'clone', url, REPO_PATH])

    if (r.returncode == 0):
        return True
    else:
        print("[COUT] Git clone error: Invalid argument to exit", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return False


def flake8(file_name):
    r = subprocess.run(['flake8', file_name], stderr=subprocess.PIPE,
                       stdout=subprocess.PIPE)

    passed = True
    if (r.returncode != 0):
        passed = False

    out = str(r.stdout, 'utf-8').strip().split('\n')
    for o in out:
        o = parse_flake8_result(o)
        if o:
            print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(o)))

    return passed


def parse_flake8_result(line):
    line = line.split(':')
    if (len(line) < 4):
        return False
    return {'file': trim_repo_path(line[0]), 'line': line[1], 'col': line[2], 'msg': line[3]}

def trim_repo_path(n):
    return n[len(REPO_PATH) + 1:]


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url']
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


# git_clone('https://github.com/Lupino/python-aio-periodic.git')
# flake8('git-repo/aio_periodic/utils.py')
# print(parse_argument())

def main():
    argv = parse_argument()
    git_url = argv.get('git-url')
    if not git_url:
        print("[COUT] The git-url value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    if not git_clone(git_url):
        return

    all_true = True

    for root, dirs, files in os.walk(REPO_PATH):
        for file_name in files:
            if file_name.endswith('.py'):
                o = flake8(os.path.join(root, file_name))
                all_true = all_true and o

    print("[COUT] CO_RESULT = false")

main()
