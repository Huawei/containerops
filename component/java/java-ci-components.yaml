# Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

uri: containerops/component/java-component-build
title: Component for java gradle project with composer
version: 1
tag: latest
timeout: 0
stages:
  -
    type: start
    name: start
    title: Start
  -
    type: normal
    name: analysis_java_gradle_checkstyle
    title: Component, java code analysis with checkstyle
    sequencing: sequence
    actions:
      -
        name: coala-test
        title: analysis your java coding style
        jobs:
          -
            type: component
            kubectl: analysis_java_gradle_checkstyle.yaml
            endpoint: hub.opshub.sh/containerops/analysis_java_gradle_checkstyle:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json"
      -
        name: flake8-test
        title: analysis your java coding style
        jobs:
          -
            type: component
            kubectl: analysis_java_gradle_cpd.yaml
            endpoint: hub.opshub.sh/containerops/analysis_java_gradle_cpd:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com//gradle-demo.git out-put-type=xml"
      -
        name: line-profiler-test
        title: analysis your java function line-by-line
        jobs:
          -
            type: component
            kubectl: analysis_java_gradle_dependencies.yaml
            endpoint: hub.opshub.sh/containerops/analysis_java_gradle_dependencies:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json"
      -
        name: analysis_java_gradle_jdepend
        title: analysis your java memory usage
        jobs:
          -
            type: component
            kubectl: analysis_java_gradle_jdepend.yaml
            endpoint: hub.opshub.sh/containerops/analysis_java_gradle_findbugs:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json"
      -
        name: analysis_java_gradle_pmd
        title: analysis your java coding style
        jobs:
          -
            type: component
            kubectl: analysis_java_gradle_pmd.yaml
            endpoint: hub.opshub.sh/containerops/analysis_java_gradle_jdepend:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json"
      -
        name: compile_java_gradle_war
        title: java call graph visualizations
        jobs:
          -
            type: component
            kubectl: compile_java_gradle_war.yaml
            endpoint: hub.opshub.sh/containerops/compile_java_gradle_war:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/rominirani/GradleWebAppSample.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/2.2.4/web.war"
      -
        name: compile_java_gradle_jar
        title: analysis your java coding style
        jobs:
          -
            type: component
            kubectl: analysis/pylama/pylama-test.yaml
            endpoint: hub.opshub.sh/containerops/compile_java_gradle_jar:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.jar" 
      -
        name: compile_java_gradle_ear
        title: analysis your java coding style
        jobs:
          -
            type: component
            kubectl: compile_java_gradle_ear.yaml
            endpoint: hub.opshub.sh/containerops/compile_java_gradle_ear:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.ear" 
      -
        name: document_java_gradle_javadoc
        title: build your Debian package with dh-virtualenv
        jobs:
          -
            type: component
            kubectl: document_java_gradle_javadoc.yaml
            endpoint: hub.opshub.sh/containerops/document_java_gradle_javadoc:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar" 
      -
        name: test_java_gradle_jacoco
        title: compile your java code with jacoco
        jobs:
          -
            type: component
            kubectl: test_java_gradle_jacoco.yaml
            endpoint: hub.opshub.sh/containerops/test_java_gradle_jacoco:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" 
      -
        name: test_java_gradle_junit
        title: build your java code with pybuilder
        jobs:
          -
            type: component
            kubectl: test_java_gradle_junit.yaml
            endpoint: hub.opshub.sh/containerops/test_java_gradle_junit:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" 
      -
        name: test_java_gradle_testng
        title: build your java code with gradle
        jobs:
          -
            type: component
            kubectl: test_java_gradle_testng.yaml
            endpoint: hub.opshub.sh/containerops/test_java_gradle_testng:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "version=gradle3 git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json"
      -
  -
    type: end
    name: end
    title: End