#!/bin/bash
#cmd="echo 'hello world' >> /tmp/lyc-ld" 
cmd1="touch /tmp/ldld" 

passwd="root"
nodeip="192.168.60.150"
expect "yes"{send "yes\r";exp_continue}
expect "password"{send "$passwd\r";expect eof}
ssh -C root@192.168.60.150 "$cmd1"

expect <<EOF
expect "yes"{send "yes\r";exp_continue}
expect "password"{send "$passwd\r";expect eof}
EOF


expect <<EOF
expect "yes"{send "yes\r";exp_continue}
expect "password"{send "$passwd\r";expect eof}
spawn ssh -C root@192.168.60.150 "touch /tmp/ldld"
EOF

#local 
ssh-keygen -b 2048 -t rsa                    
cat /roor/.ssh/id_rsa.pub

#remote
#ssh root@nodeip "mkdir tes1 ; ls"

ssh root@nodeip "mkdir /root/.ssh""

ssh root@nodeip "touch /root/.ssh/authorized_keys""

vi /etc/ssh/sshd_config
PasswordAuthentication no


ssh root@nodeip "service sshd restart"