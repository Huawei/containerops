
docker build -t containerops/analysis/java_gradle_checkstyle -f analysis/checkstyle/Dockerfile
docker build -t containerops/analysis/java_gradle_cpd -f analysis/cpd/Dockerfile
docker build -t containerops/analysis/java_gradle_dependencies -f analysis/dependencies/Dockerfile
docker build -t containerops/analysis/java_gradle_findbugs -f analysis/findbugs/Dockerfile
docker build -t containerops/analysis/java_gradle_jdepend -f analysis/jdepend/Dockerfile
docker build -t containerops/analysis/java_gradle_pmd -f analysis/pmd/Dockerfile

docker build -t containerops/test/java_gradle_jacoco -f test/jacoco/Dockerfile
docker build -t containerops/test/java_gradle_junit -f test/junit/Dockerfile
docker build -t containerops/test/java_gradle_testng -f test/testng/Dockerfile


docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_checkstyle
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=xml" containerops/analysis/java_gradle_cpd
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git" containerops/analysis/java_gradle_dependencies
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_findbugs
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_jdepend
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_pmd

docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_jacoco
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_junit
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_testng
