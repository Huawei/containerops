## Build, Release golint component

```bash
docker build -t containerops/golint .
```

## Run and test golint component
```
docker run --env CO_DATA="coderepo=https://github.com/haijunTan/gohello.git" containerops/golint:latest

output:
Cloning into '/var/opt/gopath/src/tmp'...
[COUT] Lint codeRepo(s%!)(string=https://github.com/haijunTan/gohello.git) and report：
/var/opt/gopath/src/tmp/calcproj/src/calc/calc.go:8:5: exported var Usage should have comment or be unexported
/var/opt/gopath/src/tmp/calcproj/src/simplemath/add.go:1:1: package comment should be of the form "Package simplemath ..."
/var/opt/gopath/src/tmp/calcproj/src/simplemath/add.go:4:1: exported function Add should have comment or be unexported
/var/opt/gopath/src/tmp/calcproj/src/simplemath/aqrt.go:1:1: package comment should be of the form "Package simplemath ..."
/var/opt/gopath/src/tmp/calcproj/src/simplemath/aqrt.go:5:1: exported function Sqrt should have comment or be unexported
/var/opt/gopath/src/tmp/copsgolinttest/golintEtest.go:6:9: should omit type int from declaration of var x; it will be inferred from the right-hand side
/var/opt/gopath/src/tmp/copsgolinttest/golintEtest.go:7:3: should replace x += 1 with x++
/var/opt/gopath/src/tmp/copsgolinttest/golintEtest.go:8:11: should omit type string from declaration of var str; it will be inferred from the right-hand side
[COUT] CO_RESULT = true

```
