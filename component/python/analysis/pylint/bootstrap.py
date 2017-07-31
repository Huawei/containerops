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


def pylint(file_name):
    r = subprocess.run(['pylint', '-f', 'json', file_name], stdout=subprocess.PIPE)

    passed = True
    if (r.returncode != 0):
        passed = False

    out = str(r.stdout, 'utf-8').strip()
    for o in json.loads(out):
        o['path'] = trim_repo_path(o['path'])
        print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(o)))

    return passed


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
                o = pylint(os.path.join(root, file_name))
                all_true = all_true and o

    if all_true:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")


main()
