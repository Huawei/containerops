<?php
class Phpcpd {
    const reportPath = "/tmp/PMD-CPD.xml";
    const reportFormat = "XML_REPORT";

    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            if ($input["path"] == "") {
                $input["path"] = ".";
            }
            $cmd = "/root/.composer/vendor/bin/phpcpd " . $input['path'];

            $params = [
                "names",
                "names-exclude",
                "regexps-exclude",
                "exclude",
                "min-lines",
                "min-tokens"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "") {
                    $cmd = "$cmd --$value=" . $input[$value];
                }
            }

            $cmd = "$cmd --log-pmd=" . self::reportPath;

            exec("cd " . WORK_DIR . " && $cmd");

            stdoutReport(self::reportPath, self::reportFormat);
            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    } 
}
?>