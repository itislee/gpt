#!/bin/sh
pkill -9 openai
ps aux|grep openai
sleep 1
nohup go build -o openai openai.go && ./openai &
