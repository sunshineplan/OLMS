#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go npm
}

installOLMS() {
    curl -Lo- https://github.com/sunshineplan/olms/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/olms* /var/www/olms
    cd /var/www/olms
    bash build.sh
    ./olms install
}

configOLMS() {
    read -p 'Please enter unix socket (default: /run/olms.sock): ' unix
    [ -z $unix ] && unix=/run/olms.sock
    read -p 'Please enter host (default: 0.0.0.0): ' host
    [ -z $host ] && host=0.0.0.0
    read -p 'Please enter port (default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path (default: /var/log/app/olms.log): ' log
    [ -z $log ] && log=/var/log/app/olms.log
    read -p 'Please enter reCAPTCHA site key (leave blank if not set reCAPTCHA): ' sitekey
    read -p 'Please enter reCAPTCHA secret key (leave blank if not set reCAPTCHA): ' secretkey
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/olms/config.ini.default > /var/www/olms/config.ini
    sed -i "s,\$log,$log," /var/www/olms/config.ini
    sed -i "s/\$host/$host/" /var/www/olms/config.ini
    sed -i "s/\$port/$port/" /var/www/olms/config.ini
    sed -i "s/\$sitekey/$sitekey/" /var/www/olms/config.ini
    sed -i "s/\$secretkey/$secretkey/" /var/www/olms/config.ini
    service olms start
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
    cp -s /var/www/olms/scripts/olms.cron /etc/cron.monthly/olms
    chmod +x /var/www/olms/scripts/olms.cron
}

setupNGINX() {
    cp -s /var/www/olms/scripts/olms.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/olms/scripts/olms.conf
    sed -i "s,\$unix,$unix," /var/www/olms/scripts/olms.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installOLMS
    configOLMS
    writeLogrotateScrip
    createCronTask
    setupNGINX
}

main
