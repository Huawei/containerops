#!/bin/bash
cmd="echo 'hello world' >> /tmp/lyc-ld" 
passwd="root"
nodeip="192.168.60.150"

expect << EOF
spawn ssh -C root@nodeip "$CMD"
expect{
    "yes"{send "yes\r";exp_continue}
    "password"{send "$passwd\r";expect eof}
}




expect << EOF
spawn ssh -C root@192.168.60.150 "touch /tmp/lyc-ld"
expect{
    "yes"{send "yes\r";exp_continue}
    "password"{send "$passwd\r";expect eof}
}
EOF


