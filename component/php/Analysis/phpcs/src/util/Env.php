<?php
define("CO_DATA", "CO_DATA");
class Env {
    public static function load() {
        $env = [];
        $co_data = getenv(CO_DATA);
        $kvs = explode(' ', $co_data);

        foreach ($kvs as $kv) {
            $arr = [];
            $arr = explode('=', $kv, 2);
            $env[$arr[0]] = $arr[1];
            echo "$arr[0]=$arr[1]\n";
        }

        return $env;
    }
}
?>