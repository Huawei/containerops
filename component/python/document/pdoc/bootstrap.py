#!/usr/bin/env python3

import subprocess
import os
import sys
import glob

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


def pdoc(mod):
    r = subprocess.run('pdoc --html-dir /tmp/output --html {} --all-submodules'.format(mod), shell=True)

    if r.returncode != 0:
        print("[COUT] pdoc error", file=sys.stderr)
        return False

    return True


def echo_xml():
    for root, dirs, files in os.walk('/tmp/output'):
        for file_name in files:
            if file_name.endswith('.html'):
                with open(os.path.join(root, file_name), 'r') as f:
                    for line in f.readlines():
                        line = line.rstrip()
                        print('[COUT] CO_XML_CONTENT {}'.format(line))

    return True


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-mod']
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

    entry_mod = argv.get('entry-mod')
    if not entry_mod:
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


    out = pdoc(entry_mod)
    echo_xml()

    if out:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")


main()
