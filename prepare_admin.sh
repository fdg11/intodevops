#!/usr/bin/env bash

#Update && Upgrade;
#Base utils install;
#Ansible install;
#IPV6 disable;

if [ $# != 2 ]; then 
echo -e "Enter two parametrs: \n\"email\" for git\n\"name\" for git"
exit 1
fi

DIR='/admin'
echo -e "Enter PUB_KEY(In RSA format):\n"; read PUB_KEY

# Checking for an empty value
if [ -z "$PUB_KEY" ]; then
	echo -e "The public key variable in the script body is not defined!"
	exit 1
fi

# Update & upgrade & base utils install:
apt-get update && apt upgrade -y 
apt-get install software-properties-common apt-transport-https \
	wget curl git htop atop build-essential tree vim -y

# install docker-ce, docker-compose
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt-get update
apt-get install docker-ce -y 
apt-get autoremove -y && apt-get autoclean && apt-get clean
curl -sL https://github.com/docker/compose/releases/download/1.16.1/docker-compose-`uname -s`-`uname -m` \
	-o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose 

# Default dir:
if [ ! -d "$DIR" ]; then
	mkdir $DIR
fi
echo "cd $DIR" >> ~/.bashrc

# Install add tools
git clone https://github.com/fdg11/intodevops.git $DIR
git config --global user.email $1
git config --global user.name "$2"
git config --global push.default simple

cd /$DIR

# Permit root login by the key:
mkdir -p /root/.ssh
touch ~/.ssh/authorized_keys
chmod 644 ~/.ssh/authorized_keys
echo $PUB_KEY >> ~/.ssh/authorized_keys
sed -i '/#AuthorizedKeysFile/s/#//' /etc/ssh/sshd_config
sed -i '/#PasswordAuthentication/s/#//;/PasswordAuthentication/s/yes/no/' /etc/ssh/sshd_config
service ssh restart

# IPv6 disable and enable swap limit support:
sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="ipv6.disable=1 cgroup_enable=memory swapaccount=1"/g' /etc/default/grub
update-grub

# Show ips:
echo -e "\nNetworks:\n"
ip a | grep inet" "

echo -e "\nReboot now? (y|N)\n"; read input_var

if echo $input_var | grep -iq "^y"; then
	reboot
 else
	echo -e "\nDon't forget to do that\n"
fi
