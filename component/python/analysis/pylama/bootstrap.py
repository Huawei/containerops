#!/usr/bin/env python3

import subprocess
import os
import sys
import json
import yaml

REPO_PATH = 'git-repo'


def git_clone(url):
    r = subprocess.run(['git', 'clone', url, REPO_PATH])

    if (r.returncode == 0):
        return True
    else:
        print("[COUT] Git clone error: Invalid argument to exit", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return False


def get_pip_cmd(version):
    if version == 'py3k' or version == 'python3':
        return 'pip3'

    return 'pip'


def init_env(version):
    subprocess.run([get_pip_cmd(version), 'install', 'pylama'])


def validate_version(version):
    valid_version = ['python', 'python2', 'python3', 'py3k']
    if version not in valid_version:
        print("[COUT] Check version failed: the valid version is {}".format(valid_version), file=sys.stderr)
        return False

    return True


def pylama(file_name, use_yaml):
    r = subprocess.run(['pylama', file_name], stderr=subprocess.PIPE,
                       stdout=subprocess.PIPE)

    passed = True
    if (r.returncode != 0):
        passed = False

    out = str(r.stdout, 'utf-8').strip().split('\n')
    retval = []
    for o in out:
        o = parse_pylama_result(o)
        if o:
            retval.append(o)

    if len(retval) > 0:
        out = {"results": { "cli": retval }}
        if use_yaml:
            out = bytes(yaml.safe_dump(out), 'utf-8')
            print('[COUT] CO_YAML_CONTENT {}'.format(str(out)[1:]))
        else:
            print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(out)))

    return passed


def parse_pylama_result(line):
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

    validate = ['git-url', 'version', 'out-put-type']
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

    version = argv.get('version', 'py3k')

    if not validate_version(version):
        print("[COUT] CO_RESULT = false")
        return

    init_env(version)

    if not git_clone(git_url):
        return

    all_true = True

    use_yaml = argv.get('out-put-type', 'json') == 'yaml'
    for root, dirs, files in os.walk(REPO_PATH):
        for file_name in files:
            if file_name.endswith('.py'):
                o = pylama(os.path.join(root, file_name), use_yaml)
                all_true = all_true and o

    if all_true:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")

main()
