#!/bin/bash


# 假设网卡设备名为 `eth0`，IP地址为 `172.17.0.2`
# 要限制带宽的设备
DEV=eth0

# 清除已有规则
tc qdisc del dev $DEV root

# 添加一个根qdisc，并指定一个handle
tc qdisc add dev $DEV root handle 1: htb default 12

# 添加一个class，指定带宽
tc class add dev $DEV parent 1: classid 1:1 htb rate 1mbit ceil 1mbit
tc class add dev $DEV parent 1: classid 1:12 htb rate 1mbit ceil 1mbit

# 基于IP过滤网络包并指定对应的class
tc filter add dev $DEV protocol ip parent 1:0 prio 1 u32 match ip dst 172.17.0.2/32 flowid 1:1
tc filter add dev $DEV protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.2/32 flowid 1:1
