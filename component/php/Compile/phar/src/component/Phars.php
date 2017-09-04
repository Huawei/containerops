<?php
class Phars {
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "php -f ";

            if ($input["entry-file"] == "" || $input["entry-file"] == null) {
                stderrln("[COUT] Entry file could not be null.");
                stderrln("[COUT] CO_RESULT = false");
                return;
            }

            $cmd = "$cmd " . $input["entry-file"];
            if ($input["composer"] == "true") {
                exec("cd " . WORK_DIR . " && " . "composer install", $e, $result);
                stdoutArray($e);
                if ($result != 0) {
                    stderrln("[COUT] Composer install dependence error.");
                    stderrln("[COUT] CO_RESULT = false");
                    return;
                }
            }

            exec("cd " . WORK_DIR . " && " . $cmd, $e, $result);
            stdoutArray($e);
            if ($result != 0) {
                stderrln("[COUT] Compile error.");
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