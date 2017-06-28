#!/usr/bin/expect -f 

spawn scp ./regulatorTree ubuntu@test:~/setupKernel
expect "*password:" 
send "ubuntu\r"
expect eof 
