
docker build -t containerops/analysis/java_gradle_checkstyle analysis/checkstyle <br>
docker build -t containerops/analysis/java_gradle_cpd analysis/cpd<br>
docker build -t containerops/analysis/java_gradle_dependencies analysis/dependencies<br>
docker build -t containerops/analysis/java_gradle_findbugs analysis/findbugs<br>
docker build -t containerops/analysis/java_gradle_jdepend analysis/jdepend<br>
docker build -t containerops/analysis/java_gradle_pmd analysis/pmd<br>

docker build -t huawei/compile/java_gradle_war compile/war<br>


docker build -t huawei/document/java_gradle_javadoc document/javadoc<br>

docker build -t containerops/test/java_gradle_jacoco test/jacoco<br>
docker build -t containerops/test/java_gradle_junit test/junit<br>
docker build -t containerops/test/java_gradle_testng test/testng<br>


docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_checkstyle<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=xml" containerops/analysis/java_gradle_cpd<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git" containerops/analysis/java_gradle_dependencies<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_findbugs<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_jdepend<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_pmd<br>

docker run --rm --env CO_DATA="git-url=https://github.com/rominirani/GradleWebAppSample.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/2.2.4/web.war" huawei/compile/java_gradle_war<br>


docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar" huawei/document/java_gradle_javadoc<br>

docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_jacoco<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_junit<br>
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_testng<br>
