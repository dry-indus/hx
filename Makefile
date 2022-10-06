
build:
	pwd
	go version
	go install github.com/swaggo/swag/cmd/swag@latest
	rm -rf docs
	swag i -g api/merchant/v1/api.go --exclude ./controller/userctr --instanceName mv1
	swag i -g api/user/v1/api.go --exclude ./controller/merchantctr --instanceName uv1
	go build -v -a -o hx


deploy:
	systemctl stop usersite
	mv hx usersite
	systemctl restart usersite

restart:
	systemctl restart usersite

stop:
	systemctl stop usersite

status:
	systemctl status usersite --no-pager