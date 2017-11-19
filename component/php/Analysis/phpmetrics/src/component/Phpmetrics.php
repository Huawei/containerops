<?php
class Phpmetrics {
    const reportPath = "/tmp/phpmetrics.xml";
    const reportFormat = "REPORT";
    
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "/root/.composer/vendor/bin/phpmetrics ";

            if ($input["path"] == "" || $input["path"] == null) {
                $input["path"] = ".";
            }

            $cmd = "$cmd " . $input["path"];

            $params = [
                "exclude"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "") {
                    $cmd = "$cmd --$value=" . $input[$value];
                }
            }

            $params_bool = [
                "ignore-annotations",
            ];

            foreach ($params as $value) {
                if ($input[$value] == "true") {
                    $cmd = "$cmd --$value";
                }
            }

            $cmd = "$cmd --report-violations=" . self::reportPath;

            $lastLine = system("cd " . WORK_DIR . " && $cmd", $result);
            
            if ($result != 0) {
                throw new Exception();
            }

            stdoutReport(self::reportPath, self::reportFormat);
            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    }
}
?>