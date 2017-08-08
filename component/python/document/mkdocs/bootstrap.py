#!/usr/bin/env python3

import subprocess
import os
import sys

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


def mkdocs(dir_name):
    r = subprocess.run('cd {}/{}; mkdocs json'.format(REPO_PATH, dir_name), shell=True)

    if r.returncode != 0:
        print("[COUT] mkdocs error", file=sys.stderr)
        return False

    return True


def echo_json(dir_name):
    r = subprocess.run('cd {}/{}; find -name "*.json" | while read F;do echo -n \'CO_JSON_CONTENT \'; cat $F | tr -d \'\n\'; echo; done'.format(REPO_PATH, dir_name), shell=True)

    if r.returncode != 0:
        print("[COUT] echo json failed", file=sys.stderr)
        return False

    return True


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

    if not mkdocs(entry_path):
        print("[COUT] CO_RESULT = false")
        return

    if not echo_json(entry_path):
        print("[COUT] CO_RESULT = false")
        return

    print("[COUT] CO_RESULT = true")


main()
