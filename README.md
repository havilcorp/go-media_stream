services:
mysql:
image: mysql:latest
container_name: mysql
restart: always
ports: - "3306:3306"
volumes: - ./data/mysql:/var/lib/mysql
environment: # MYSQL_USER: "user" # MYSQL_PASSWORD: "password"
MYSQL_ROOT_PASSWORD: "rock4how"
networks:
br0:
ipv4_address: 10.0.0.6

phpmyadmin:
image: phpmyadmin
container_name: "phpmyadmin"
restart: always
ports: - 3003:80
environment:
PMA_HOST: mysql
PMA_PORT: 3306
UPLOAD_LIMIT: 1024M
MEMORY_LIMIT: 1024M
MAX_EXECUTION_TIME: 300
networks:
br0:
ipv4_address: 10.0.0.7
networks:
br0:
external: true
