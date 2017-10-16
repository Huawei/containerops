#!/usr/bin/env python3

import subprocess
import os
import sys
import json
import anymarkup

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
        print("[COUT] install dependences failed: {}".format(file_name), file=sys.stderr)
        return False

    return True


def tox(file_name):
    r = subprocess.run('cd {}/{}; tox --result-json /tmp/output.json'.format(REPO_PATH, file_name), shell=True)

    if r.returncode != 0:
        return False

    return True


def echo_json(use_yaml):
    data = json.load(open('/tmp/output.json'))
    if use_yaml:
        data = anymarkup.serialize(data, 'yaml')
        print('[COUT] CO_YAML_CONTENT {}'.format(str(data)[1:]))
    else:
        print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(data)))

    return True


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-path', 'out-put-type']
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

    entry_path = argv.get('entry-path')
    if not entry_path:
        print("[COUT] The entry-path value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    if not git_clone(git_url):
        return

    out = tox(entry_path)

    use_yaml = argv.get('out-put-type', 'json') == 'yaml'

    echo_json(use_yaml)

    if not out:
        print("[COUT] CO_RESULT = false")
        return

    print("[COUT] CO_RESULT = true")


main()
