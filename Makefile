
build:
	go version
	# go mod tidy
	go build -v -o hx

swag:
	rm -rf docs
	go install github.com/swaggo/swag/cmd/swag@latest
	swag -v
	swag i -g api/landing/v1/api.go --exclude ./controller/userctr,./controller/merchantctr,./go_cache --instanceName lv1
	swag i -g api/merchant/v1/api.go --exclude ./controller/landingctr,./controller/userctr,./go_cache --instanceName mv1
	swag i -g api/user/v1/api.go --exclude ./controller/landingctr,./controller/merchantctr,./go_cache --instanceName uv1
	

deploy:
	systemctl stop usersite
	mv hx usersite
	systemctl restart usersite

restart:
	systemctl restart usersite

stop:
	systemctl stop usersite

status:
	systemctl status usersite --no-pager -l