<?php
class Phpcs {
    const reportPath = "/tmp/phpcs.xml";
    const reportFormat = "REPORT";
    
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "/root/.composer/vendor/bin/phpcs ";

            if ($input["file"] == "" || $input["file"] == null) {
                $input["file"] = ".";
            }

            $cmd = "$cmd " . $input["file"];

            $params = [
                "report",
                "basepath",
                "bootstrap",
                "severity",
                "error-severity",
                "warning-severity",
                "standard",
                "sniffs",
                "encoding",
                "parallel",
                "generator",
                "extensions",
                "ignore",
                "file-list"
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

            $cmd = "$cmd --report-file=" . self::reportPath;

            $lastLine = exec("cd " . WORK_DIR . " && $cmd", $output, $result);
            if ($result != 0) {
                stderrln($lastLine);
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