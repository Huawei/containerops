#!/usr/bin/env python3

import subprocess
import os
import sys
import glob
import json
import line_profiler as profiler
import linecache
import inspect
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


def get_python_cmd(version):
    if version == 'py3k' or version == 'python3':
        return 'python3'

    return 'python'


def init_env(version):
    subprocess.run([get_pip_cmd(version), 'install', 'cython', 'line_profiler'])


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


def show_json(stats, unit):
    """ Show text for the given timings.
    """
    retval = {}
    retval['Timer unit'] = '%g s' % unit
    retval['functions'] = []
    for (fn, lineno, name), timings in sorted(stats.items()):
        func = show_func(fn, lineno, name, stats[fn, lineno, name], unit)
        if func:
            retval['functions'].append(func)

    return retval

def show_func(filename, start_lineno, func_name, timings, unit):
    """ Show results for a single function.
    """

    d = {}
    total_time = 0.0
    linenos = []
    for lineno, nhits, time in timings:
        total_time += time
        linenos.append(lineno)

    if total_time == 0:
        return False

    retval = {}

    retval['Total time'] = "%g s" % (total_time * unit)
    if os.path.exists(filename) or filename.startswith("<ipython-input-"):
        retval['File'] = filename
        retval['Function'] = '%s at line %s' % (func_name, start_lineno)
        if os.path.exists(filename):
            # Clear the cache to ensure that we get up-to-date results.
            linecache.clearcache()
        all_lines = linecache.getlines(filename)
        sublines = inspect.getblock(all_lines[start_lineno-1:])
    else:
        # Fake empty lines so we can see the timings, if not the code.
        nlines = max(linenos) - min(min(linenos), start_lineno) + 1
        sublines = [''] * nlines

    for lineno, nhits, time in timings:
        d[lineno] = (nhits, time, '%5.1f' % (float(time) / nhits),
            '%5.1f' % (100*time / total_time))
    linenos = range(start_lineno, start_lineno + len(sublines))
    empty = ('', '', '', '')
    lines = []
    for lineno, line in zip(linenos, sublines):
        nhits, time, per_hit, percent = d.get(lineno, empty)
        line = {
                'Line #': lineno,
                'Hits': nhits,
                'Time': time,
                'Per Hit': per_hit,
                '% Time': percent,
                'Line Contents': line.rstrip('\n').rstrip('\r')
                }
        lines.append(line)

    retval['lines'] = lines
    return retval

def line_profiler(file_name, use_yaml):
    r = subprocess.run(['kernprof', '-l', os.path.join(REPO_PATH, file_name)], stdout=subprocess.PIPE)

    passed = True
    if (r.returncode != 0):
        passed = False

    st = profiler.load_stats('{}/{}.lprof'.format(REPO_PATH, file_name))
    out = show_json(st.timings, st.unit)
    if use_yaml:
        out = bytes(yaml.safe_dump(out), 'utf-8')
        print('[COUT] CO_YAML_CONTENT {}'.format(str(out)[1:]))
    else:
        print('[COUT] CO_JSON_CONTENT {}'.format(json.dumps(out)))

    return passed


def trim_repo_path(n):
    return n[len(REPO_PATH) + 1:]


def parse_argument():
    data = os.environ.get('CO_DATA', None)
    if not data:
        return {}

    validate = ['git-url', 'entry-file', 'version', 'out-put-type']
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

    entry_file = argv.get('entry-file')
    if not entry_file:
        print("[COUT] The entry-file value is null", file=sys.stderr)
        print("[COUT] CO_RESULT = false")
        return

    if not git_clone(git_url):
        return

    for file_name in glob.glob('./*/setup.py'):
        setup(file_name, version)

    for file_name in glob.glob('./*/*/setup.py'):
        setup(file_name, version)

    for file_name in glob.glob('./*/requirements.txt'):
        pip_install(file_name, version)

    for file_name in glob.glob('./*/*/requirements.txt'):
        pip_install(file_name, version)

    use_yaml = argv.get('out-put-type', 'json') == 'yaml'

    out = line_profiler(entry_file, use_yaml)

    if out:
        print("[COUT] CO_RESULT = true")
    else:
        print("[COUT] CO_RESULT = false")


main()
