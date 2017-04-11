#!/bin/bash
cmd="echo 'hello world' >> /tmp/lyc-ld" 
passwd="test"
nodeip="192.168.60.150"

expect << EOF
spawn ssh -C root@ip "$CMD"
expect{
    "yes"{send "yes\r";exp_continue}
    "password"{send "$passwd\r";expect eof}
}

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