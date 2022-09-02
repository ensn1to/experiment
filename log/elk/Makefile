ps:
	docker-compose -f docker-compose.yml ps -a

pull:
	docker-compose -f docker-compose.yml pull

up:
	docker-compose -f docker-compose.yml up 

upd:
	docker-compose -f docker-compose.yml up -d

down:
	docker-compose -f docker-compose.yml down

downv:
	docker-compose -f docker-compose.yml down -v


# usage example: make upsrv service=srv1
service=
upsrv:
	docker-compose -f my-docker-compose.yml up -d ${service}

config:
	docker-compose -f my-docker-compose.yml config


# 把标准错误输出重定向到标准输出，都输出到日志up.log
up2log:
	docker-compose -f my-docker-compose.yml up > up.log 2>&1

log:
	docker-compose -f my-docker-compose.yml logs -f
