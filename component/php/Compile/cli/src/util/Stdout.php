<?php
function stdoutln ($string) {
    fwrite(STDOUT, $string . "\n");
}

function stderrln ($string) {
    fwrite(STDERR, $string . "\n");
}

function stdoutReport ($name, $type) {
    $report = fopen($name, "r");
    fwrite(STDOUT, "[COUT] $type = \n");
    fwrite(STDOUT, fread($report, filesize($name)));
}

function stdoutArray ($strings) {
    foreach ($strings as $string) {
        stdoutln($string);
    }
}
?>