#!/usr/bin/env python3

import subprocess
import os
import sys
import glob
import json
import yaml

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


def get_dir_name(dir_name):
    if dir_name and dir_name != '.':
        dir_name = '{}/{}'.format(REPO_PATH, dir_name)
    else:
        dir_name = REPO_PATH

    return dir_name


def green(dir_name):
    r = subprocess.run('cd {}; green -r'.format(get_dir_name(dir_name)), shell=True)

    if r.returncode != 0:
        return False

    return True


def echo_yaml(dir_name):
    file_name = '{}/.coverage'.format(get_dir_name(dir_name))
    data = open(file_name).read()
    idx = data.find('{')
    data = data[idx:]
    data = json.loads(data)
    lines = yaml.safe_dump(data)
    for line in lines.split('\n'):
        print('[COUT] CO_YAML_CONTENT {}'.format(line))


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-path']
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

    for file_name in glob.glob('{}/setup.py'.format(REPO_PATH)):
        setup(file_name)

    for file_name in glob.glob('{}/*/setup.py'.format(REPO_PATH)):
        setup(file_name)

    for file_name in glob.glob('{}/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name)

    for file_name in glob.glob('{}/*/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name)

    for file_name in glob.glob('{}/requirements_dev.txt'.format(REPO_PATH)):
        pip_install(file_name)

    for file_name in glob.glob('{}/*/requirements_dev.txt'.format(REPO_PATH)):
        pip_install(file_name)


    out = green(entry_path)

    echo_yaml(entry_path)

    if not out:
        print("[COUT] CO_RESULT = false")
        return

    print("[COUT] CO_RESULT = true")


main()
