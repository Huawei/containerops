<?php
class Phpmd {
    const reportPath = "/tmp/phpmd.xml";
    const reportFormat = "REPORT";
    
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "/root/.composer/vendor/bin/phpmd ";

            if ($input["path"] == "" || $input["path"] == null) {
                $input["path"] = ".";
            }

            if ($input["formats"] == "" || $input["formats"] == null) {
                $input["formats"] = "xml";
            }

            if ($input["ruleset"] == "" || $input["ruleset"] == null) {
                $input["ruleset"] = "cleancode,codesize,controversial,design,naming,unusedcode";
            }

            $cmd = "$cmd " . $input['path'] . " " . $input['formats'] . " " . $input['ruleset'];

            $params = [
                "minimumpriority",
                "exclude",
                "suffixes"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "") {
                    $cmd = "$cmd --$value " . $input[$value];
                }
            }

            $params_bool = [
                "minimumpriority",
                "strict",
                "ignore-violations-on-exit"
            ];

            foreach ($params as $value) {
                if ($input[$value] == "true") {
                    $cmd = "$cmd --$value";
                }
            }

            $cmd = "$cmd --reportfile " . self::reportPath;

            exec("cd " . WORK_DIR . " && $cmd");

            stdoutReport(self::reportPath, self::reportFormat);
            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    }
}
?>