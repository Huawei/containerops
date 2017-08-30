<?php
class Phploc {
    const reportPath = "/tmp/report.xml";
    const reportFormat = "XML_REPORT";
    
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "/root/.composer/vendor/bin/phploc ";

            if ($input["path"] == "" || $input["path"] == null) {
                $input["path"] = ".";
            }

            $cmd = "$cmd " . $input["file"];

            $params = [
                "exclude",
                "names",
                "names-exclude"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "") {
                    $cmd = "$cmd --$value=" . $input[$value];
                }
            }

            $params_bool = [
                "count_tests",
            ];

            foreach ($params as $value) {
                if ($input[$value] == "true") {
                    $cmd = "$cmd --$value";
                }
            }

            $cmd = "$cmd --log-xml=" . self::reportPath;

            exec("cd " . WORK_DIR . " && $cmd");

            stdoutReport(self::reportPath, self::reportFormat);
            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    }
}
?>