
docker build -t containerops/analysis/java_gradle_checkstyle analysis/checkstyle
docker build -t containerops/analysis/java_gradle_cpd analysis/cpd
docker build -t containerops/analysis/java_gradle_dependencies analysis/dependencies
docker build -t containerops/analysis/java_gradle_findbugs analysis/findbugs
docker build -t containerops/analysis/java_gradle_jdepend analysis/jdepend
docker build -t containerops/analysis/java_gradle_pmd analysis/pmd

docker build -t containerops/test/java_gradle_jacoco test/jacoco
docker build -t containerops/test/java_gradle_junit test/junit
docker build -t containerops/test/java_gradle_testng test/testng


docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_checkstyle
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=xml" containerops/analysis/java_gradle_cpd
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git" containerops/analysis/java_gradle_dependencies
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_findbugs
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_jdepend
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_pmd

docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_jacoco
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_junit
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_testng
