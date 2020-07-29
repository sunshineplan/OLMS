#! /bin/bash

installSoftware() {
    apt -qq -y install nginx default-mysql-client
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installOLMS() {
    curl -Lo- https://github.com/sunshineplan/olms-go/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/olms-go* /var/www/olms-go
    cd /var/www/olms-go
    go build
}

configOLMS() {
    read -p 'Please enter unix socket(default: /run/olms-go.sock): ' unix
    [ -z $unix ] && unix=/var/www/olms-go/olms-go.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/olms-go.log): ' log
    [ -z $log ] && log=/var/log/app/olms-go.log
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/olms-go/config.ini.default > /var/www/olms-go/config.ini
    sed -i "s,\$log,$log," /var/www/olms-go/config.ini
    sed -i "s/\$host/$host/" /var/www/olms-go/config.ini
    sed -i "s/\$port/$port/" /var/www/olms-go/config.ini
}

setupsystemd() {
    cp -s /var/www/olms-go/scripts/olms-go.service /etc/systemd/system
    systemctl enable olms-go
    service olms-go start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

createCronTask() {
    cp -s /var/www/olms-go/scripts/olms-go.cron /etc/cron.monthly/olms-go
    chmod +x /var/www/olms-go/scripts/olms-go.cron
}

setupNGINX() {
    cp -s /var/www/olms-go/scripts/olms-go.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/olms-go/scripts/olms-go.conf
    sed -i "s,\$unix,$unix," /var/www/olms-go/scripts/olms-go.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installOLMS
    configOLMS
    setupsystemd
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main
