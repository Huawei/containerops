## Build, Release pylint component

```bash
docker build -t containerops/pylint .
(workdir:contianerops/component/images/python/pylint)
```

## Set pylint parameter
```bash
The default pylint prameter in contianerops/component/images/python/pylint/src/pylint.conf.   
You can modify it befor build the component.
```

## Run and test pylint component
```bash
docker run --env CO_CODERPO="coderepo=https://github.com/haijunTan/pyhello.git"  containerops/pylint:latest

Output result:
[COUT] CO_TEST = coderepo=https://github.com/haijunTan/pyhello.git
Cloning into '/var/opt/gopath/src/tmp'...
------------------------------
[COUT] Start lint file：pylinttest/pylintnowarning.py
[COUT] pylinttest/pylintnowarning.py isn't any warning
[COUT] End lint file：pylinttest/pylintnowarning.py end
------------------------------
[COUT] Start lint file：pylinttest/pylinttest.py
************* Module pylinttest
W:  3, 0: Found indentation with tabs instead of spaces (mixed-indentation)
C:  3, 0: Unnecessary parens after 'print' keyword (superfluous-parens)
W:  6, 0: Found indentation with tabs instead of spaces (mixed-indentation)
W:  7, 0: Found indentation with tabs instead of spaces (mixed-indentation)
W: 10, 0: Found indentation with tabs instead of spaces (mixed-indentation)
C:  1, 0: Missing module docstring (missing-docstring)
C:  2, 0: Missing function docstring (missing-docstring)
C:  5, 0: Missing function docstring (missing-docstring)
C:  7, 1: Invalid variable name "x" (invalid-name)
W:  7, 1: Unused variable 'x' (unused-variable)
[COUT] End lint file：pylinttest/pylinttest.py end
------------------------------
[COUT] Start lint file：pylinttest/pylinttest1.py
************* Module pylinttest1
W:  3, 0: Found indentation with tabs instead of spaces (mixed-indentation)
C:  3, 0: Unnecessary parens after 'print' keyword (superfluous-parens)
W:  4, 0: Found indentation with tabs instead of spaces (mixed-indentation)
W:  7, 0: Found indentation with tabs instead of spaces (mixed-indentation)
C:  7, 0: Unnecessary parens after 'print' keyword (superfluous-parens)
W: 10, 0: Found indentation with tabs instead of spaces (mixed-indentation)
W: 11, 0: Found indentation with tabs instead of spaces (mixed-indentation)
W: 14, 0: Found indentation with tabs instead of spaces (mixed-indentation)
C:  1, 0: Missing module docstring (missing-docstring)
C:  2, 0: Missing function docstring (missing-docstring)
W:  4, 1: Unused variable 'unusevary' (unused-variable)
C:  6, 0: Missing function docstring (missing-docstring)
C:  9, 0: Missing function docstring (missing-docstring)
W: 10, 1: Unused variable 'unusevarx' (unused-variable)
[COUT] End lint file：pylinttest/pylinttest1.py end
[COUT] CO_RESULT = true
```
