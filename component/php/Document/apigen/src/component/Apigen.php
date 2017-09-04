<?php
class Apigen {
    public static function exec($input) {
        try {
            Git::clone($input[GIT_REPO]);
        } catch (Exception $ex) {
            stderrln("Git clone error: " . $ex->getMessage());
            stderrln("[COUT] CO_RESULT = false");
        }

        try {
            $cmd = "/root/.composer/vendor/bin/apigen generate -q ";

            if ($input["path"] == "" || $input["path"] == null) {
                $input["path"] = ".";
            }

            $cmd = "$cmd -s " . $input["path"];

            $params = [
                "destination"
            ];

            foreach ($params as $value) {
                if ($input[$value] != "" && $input[$value] != null) {
                    $cmd = "$cmd --$value " . $input[$value];
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