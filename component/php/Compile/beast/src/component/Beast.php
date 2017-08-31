<?php
$nfiles = 0;
$finish = 0;

class Beast {
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

            
            $conf = [
                "src_path" => WORK_DIR,
                "dst_path" => "/root/dist",
                "expire" => "",
                "encrypt_type" => $input["encrypt_type"],
            ];
            $src_path     = trim($conf['src_path']);
            $dst_path     = trim($conf['dst_path']);
            $expire       = trim($conf['expire']);
            $encrypt_type = strtoupper(trim($conf['encrypt_type']));
            if (empty($src_path) || !is_dir($src_path)) {
                exit("Fatal: source path `{$src_path}' not exists\n\n");
            }
            if (empty($dst_path)
                || (!is_dir($dst_path)
                && !mkdir($dst_path, 0777)))
            {
                exit("Fatal: can not create directory `{$dst_path}'\n\n");
            }
            switch ($encrypt_type)
            {
            case 'AES':
                $entype = BEAST_ENCRYPT_TYPE_AES;
                break;
            case 'BASE64':
                $entype = BEAST_ENCRYPT_TYPE_BASE64;
                break;
            case 'DES':
            default:
                $entype = BEAST_ENCRYPT_TYPE_DES;
                break;
            }
            stdoutln(sprintf("Source code path: %s", $src_path));
            stdoutln(sprintf("Destination code path: %s", $dst_path));
            stdoutln(sprintf("Expire time: %s", $expire));
            stdoutln(sprintf("------------- start process -------------"));
            $expire_time = 0;
            if ($expire) {
                $expire_time = strtotime($expire);
            }
            $time = microtime(TRUE);
            calculate_directory_schedule($src_path);
            encrypt_directory($src_path, $dst_path, $expire_time, $entype);
            $used = microtime(TRUE) - $time;
            stdoutln(printf("\nFinish processed encrypt files, used %f seconds", $used));

            stdoutln("[COUT] CO_RESULT = true");
        } catch (Exception $ex) {
            stderrln("[COUT] CO_RESULT = false");
        }
    }
}

function calculate_directory_schedule($dir)
{
    global $nfiles;
    $dir = rtrim($dir, '/');
    $handle = opendir($dir);
    if (!$handle) {
        return false;
    }
    while (($file = readdir($handle))) {
        if ($file == '.' || $file == '..') {
            continue;
        }
        $path = $dir . '/' . $file;
        if (is_dir($path)) {
            calculate_directory_schedule($path);
        } else {
            $infos = explode('.', $file);
            if (strtolower($infos[count($infos)-1]) == 'php') {
                $nfiles++;
            }
        }
    }
    closedir($handle);
}

function encrypt_directory($dir, $new_dir, $expire, $type)
{
    global $nfiles, $finish;
    $dir = rtrim($dir, '/');
    $new_dir = rtrim($new_dir, '/');
    $handle = opendir($dir);
    if (!$handle) {
        return false;
    }
    while (($file = readdir($handle))) {
        if ($file == '.' || $file == '..') {
            continue;
        }
        $path = $dir . '/' . $file;
        $new_path =  $new_dir . '/' . $file;
        if (is_dir($path)) {
            if (!is_dir($new_path)) {
                mkdir($new_path, 0777);
            }
            encrypt_directory($path, $new_path, $expire, $type);
        } else {
            $infos = explode('.', $file);
            if (strtolower($infos[count($infos)-1]) == 'php'
                && filesize($path) > 0)
            {
                if ($expire > 0) {
                    $result = beast_encode_file($path, $new_path,
                                                $expire, $type);
                } else {
                    $result = beast_encode_file($path, $new_path, 0, $type);
                }
                if (!$result) {
                    echo "Failed to encode file `{$path}'\n";
                }
                $finish++;
                $percent = intval($finish / $nfiles * 100);
                printf("\rProcessed encrypt files [%d%%] - 100%%", $percent);
            } else {
                copy($path, $new_path);
            }
        }
    }
    closedir($handle);
}
?>