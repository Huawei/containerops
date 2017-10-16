#!/usr/bin/env python3

import subprocess
import os
import sys
import json
import pycco.main as pyccoLib
import anymarkup

REPO_PATH = 'git-repo'


def git_clone(url):
    r = subprocess.run(['git', 'clone', url, REPO_PATH])

    if (r.returncode == 0):
        return True
    else:
        print("[COUT] Git clone error: Invalid argument to exit", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return False


def pycco(file_name, use_yaml):
    code = open(file_name, 'rb').read().decode('utf-8')
    language = pyccoLib.get_language(file_name, code)
    sections = pyccoLib.parse(code, language)
    pyccoLib.highlight(sections, language, outdir="docs")

    data = { 'file_path': trim_repo_path(file_name), 'sections': sections }
    if use_yaml:
        data = anymarkup.serialize(data, 'yaml')
        print('[COUT] CO_YAML_CONTENT {}'.format(str(data)[1:]))
    else:
        print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(data)))


    return True


def trim_repo_path(n):
    return n[len(REPO_PATH) + 1:]


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'out-put-type']
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

    use_yaml = argv.get('out-put-type', 'json') == 'yaml'

    for root, dirs, files in os.walk(REPO_PATH):
        for file_name in files:
            if file_name.endswith('.py'):
                o = pycco(os.path.join(root, file_name), use_yaml)
                all_true = all_true and o

    if all_true:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")

main()
