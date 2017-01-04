## Quick start

#### Clone codes

```
git clone https://github.com/Huawei/containerops.git  $GOPATH/src/github.com/Huawei/containerops
cd $GOPATH/src/github.com/Huawei/containerops/pilotage
```
	
#### Install Robot framework (GUI,API TEST) depend on Linux (Ubuntu 14.04+)

```
install python2.7
install pip
pip install robotframework
install wxPython
pip install robotframework-ride
pip install robotframework-selenium2library
pip install robotframework-databaselibrary
pip install requests
pip install -U robotframework-requests
install pywin32
install AutoItLibrary
```

#### install multi-mechanize (performance test) depend on Linux (Ubuntu 14.04+)

```
pip install multi-mechanize mechanize numpy matplotlib
```

#### run gui/api test

```
pybot [testproject]
```

#### run performance test

```
multimech-run [testproject]
```