
build:
	go version
	# go mod tidy
	go build -v -o hx

swag:
    ifeq ($(findstring No,$(swag -v)),No)
	go install github.com/swaggo/swag/cmd/swag@latest
    endif

	rm -rf docs
	swag fmt
	swag i -g api/merchant/v1/api.go --exclude ./controller/userctr,./go_cache --instanceName mv1
	swag i -g api/user/v1/api.go --exclude ./controller/merchantctr,./go_cache --instanceName uv1
	

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