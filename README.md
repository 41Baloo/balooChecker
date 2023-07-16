# balooChecker

Simple lightweight proxy-checker written in golang

## Usage

This proxy-checker reads directly from the osStdIn. This means usage works as following

**Check from website**

`curl https://your.website/proxies.txt | ./main [timeout] [output name] [threads]`

**Check from file**

`cat your-file.txt | ./main [timeout] [output name] [threads]`

**Check from scanner**

`zmap -p 80 172.217.0.0/24 | ./main [timeout] [output name] [threads]`

**Check directly from terminal**

`echo "1.1.1.1:80" | ./main [timeout] [output name] [threads]`