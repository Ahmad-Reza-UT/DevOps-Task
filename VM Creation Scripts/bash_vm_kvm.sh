#!/bin/bash

### System Update ###
sudo apt-get update && sudo apt-get -y upgrade


sudo dnf install epel-release wget curl net-tools vim
sudo dnf install cloud-init cloud-utils-growpart acpid

### Installing KVM in Linux ###
sudo systemctl enable libvirtd
sudo systemctl start libvirtd

# Verifying if the libvirtd daemon is running.
sudo systemctl status libvirtd


# In case you are using Ubuntu/Debian, make sure that the vhost-net image is loaded
sudo modprobe vhost_net

### Create a KVM Virtual Image ###
sudo qemu-img create -o preallocation=metadata -f qcow2 /var/lib/libvirt/images/centos8.qcow2 25G


sudo virt-install --virt-type kvm --name centos8 --ram 2048 \
--disk /var/lib/libvirt/images/centos8.qcow2,format=qcow2 \
--network network=default \
--graphics vnc,listen=0.0.0.0 --noautoconsole \
--os-type=linux --os-variant=rhel7.0 \
--location=/home/tecmint/Downloads/CentOS-8-x86_64-1905-dvd1.iso

### Creating KVM Virtual Machine Template Image ###
sudo dnf update

# Disabling the zeroconf route.
echo "NOZEROCONF=yes" >> /etc/sysconfig/network

sudo virt-sysprep -d centos8

sudo virsh undefine centos8