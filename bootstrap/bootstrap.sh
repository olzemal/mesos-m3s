#!/bin/bash

cat /etc/resolv.conf

apt-get update -y 
apt-get install -y containernetworking-plugins containerd tcpdump curl inetutils-ping iptables fuse-overlayfs procps bash
mkdir -p /etc/cni/net.d

export INSTALL_K3S_SKIP_ENABLE=true
export INSTALL_K3S_SKIP_START=true

update-alternatives --set iptables /usr/sbin/iptables-legacy
curl -sfL https://get.k3s.io | sh -
curl https://raw.githubusercontent.com/AVENTER-UG/go-mesos-framework-k3s/master/bootstrap/server > $MESOS_SANDBOX/server
curl https://raw.githubusercontent.com/kubernetes/dashboard/v2.2.0/aio/deploy/recommended.yaml > $MESOS_SANDBOX/dashboard.yaml

if [ "$K3SFRAMEWORK_TYPE" = "server" ]
then
  chmod +x $MESOS_SANDBOX/server
  exec $MESOS_SANDBOX/server &
fi

echo $1
$1 
