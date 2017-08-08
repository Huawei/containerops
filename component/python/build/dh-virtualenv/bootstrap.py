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


def upload_file(upload):
    parsed = upload.split('/')
    host = parsed[0]
    namespace = parsed[1]
    repo = parsed[2]
    binary = parsed[3]
    tag = parsed[4]
    url = 'https://{}/binary/v1/{}/{}/binary/{}/{}'.format(host, namespace,
                                                           repo, binary, tag)


    file_name = glob.glob('*.deb')[0]
    r1 = subprocess.run(['curl', '-XPUT', '-d', '@' + file_name, url])
    if r1.returncode != 0:
        print("[COUT] upload error", file=sys.stderr)
        return False
    return True


def build():
    r = subprocess.run('cd {}; yes | mk-build-deps -ri; dpkg-buildpackage -us -uc -b'.format(REPO_PATH), shell=True)

    if r.returncode != 0:
        print("[COUT] build error", file=sys.stderr)
        return False

    return True


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-file', 'upload']
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

    upload = argv.get('upload')
    if not upload:
        print("[COUT] The upload value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    if not git_clone(git_url):
        return

    if not build():
        print("[COUT] CO_RESULT = false")
        return

    if not upload_file(upload):
        print("[COUT] CO_RESULT = false")
        return

    print("[COUT] CO_RESULT = true")


main()
