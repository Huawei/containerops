<?php
class Phpunit {
    const reportPath = "/tmp/report.xml";
    const reportFormat = "TEST REPORT";
    
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            if ($input["composer"] == "true") {
                exec("cd " . WORK_DIR . " && " . "composer install", $e, $result);
                stdoutArray($e);
                if ($result != 0) {
                    stderrln("[COUT] Composer install dependence error.");
                    stderrln("[COUT] CO_RESULT = false");
                    return;
                }
            }

            $cmd = "/root/.composer/vendor/bin/phpunit ";

            $params = [
                "bootstrap",
                "include-path",
                "configuration"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "") {
                    $cmd = "$cmd --$value " . $input[$value];
                }
            }

            $cmd = "$cmd --log-junit " . self::reportPath;

            exec("cd " . WORK_DIR . " && " . $cmd, $e, $result);
            stdoutArray($e);
            if ($result != 0) {
                stderrln("[COUT] Test error.");
                stderrln("[COUT] CO_RESULT = false");
                return;
            }
            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    }
}
?>