## ServiceComb Chassis Component

### What's the Component?
This image is java runtime image, used for analysis your ServiceComb Chassis java coding style. Checkstyle is a development tool to help programmers write Java code that adheres to a coding standard. It automates the process of checking Java code to spare humans of this boring (but important) task. This makes it ideal for projects that want to enforce a coding standard. 

gradle checkstyleMain 
gradle checkstyleTest

### Learn how to build it?

Use the `docker build` command build the image, and ServiceComb Chassis must build with gradle

```
docker build -t containerops/serviceComb/serviceComb_Chassis_gradle_checkstyle:v1.0 ./
```
### Component Usage
```
docker run  \
--rm \--env CO_DATA="\
version=gradle4 \
git-url=https://github.com/ServiceComb/ServiceComb-Java-Chassis.git \
out-put-type=json \
report-path=./webapp/build/reports/checkstyle” \
containerops/analysis/serviceComb_Chassis_gradle_checkstyle
  
```

### Parameters 

### Versions 1.0.0
