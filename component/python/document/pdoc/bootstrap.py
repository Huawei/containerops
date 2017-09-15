#!/usr/bin/env python3

import subprocess
import os
import sys
import glob
from bs4 import BeautifulSoup
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


def get_pip_cmd(version):
    if version == 'py3k' or version == 'python3':
        return 'pip3'

    return 'pip'


def get_python_cmd(version):
    if version == 'py3k' or version == 'python3':
        return 'python3'

    return 'python'


def init_env(version):
    subprocess.run([get_pip_cmd(version), 'install', 'pdoc'])


def validate_version(version):
    valid_version = ['python', 'python2', 'python3', 'py3k']
    if version not in valid_version:
        print("[COUT] Check version failed: the valid version is {}".format(valid_version), file=sys.stderr)
        return False

    return True


def setup(path, version='py3k'):
    file_name = os.path.basename(path)
    dir_name = os.path.dirname(path)
    r = subprocess.run('cd {}; {} {} install'.format(dir_name, get_python_cmd(version), file_name),
                       shell=True)

    if r.returncode != 0:
        print("[COUT] install dependences failed", file=sys.stderr)
        return False

    return True


def pip_install(file_name, version='py3k'):
    r = subprocess.run([get_pip_cmd(version), 'install', '-r', file_name])

    if r.returncode != 0:
        print("[COUT] install dependences failed", file=sys.stderr)
        return False

    return True


def pdoc(mod):
    r = subprocess.run('pdoc --html-dir /tmp/output --html {} --all-submodules'.format(mod), shell=True)

    if r.returncode != 0:
        print("[COUT] pdoc error", file=sys.stderr)
        return False

    return True


def echo_json(use_yaml):
    for root, dirs, files in os.walk('/tmp/output'):
        for file_name in files:
            if file_name.endswith('.html'):
                with open(os.path.join(root, file_name), 'r') as f:
                    data = f.read()
                    soup = BeautifulSoup(data, 'html.parser')
                    title = soup.find('title').text
                    body = soup.find('body').renderContents()
                    data = {
                        "title": title,
                        "body": str(body, 'utf-8', errors='ignore'),
                        "file": file_name
                    }
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

    validate = ['git-url', 'entry-mod', 'version', 'out-put-type']
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

    entry_mod = argv.get('entry-mod')
    if not entry_mod:
        print("[COUT] The entry-path value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    if not git_clone(git_url):
        return

    for file_name in glob.glob('{}/setup.py'.format(REPO_PATH)):
        setup(file_name, version)

    for file_name in glob.glob('{}/*/setup.py'.format(REPO_PATH)):
        setup(file_name, version)

    for file_name in glob.glob('{}/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name, version)

    for file_name in glob.glob('{}/*/requirements.txt'.format(REPO_PATH)):
        pip_install(file_name, version)


    out = pdoc(entry_mod)

    use_yaml = argv.get('out-put-type', 'json') == 'yaml'

    echo_json(use_yaml)

    if out:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")


main()
