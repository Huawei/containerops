<?php
define("GIT_REPO", "git-url");

class Git {
    public static function clone($repo) {
        if (!is_string($repo) || $repo == "") {
            stderrln("Git repo url error.");
        }

        exec("git clone $repo " . WORK_DIR);
    }
}
?>