revision=`git log | head -n 1 | cut  -f 2 -d ' '`
sat_default_registry=registry.staging.example.com
prod_default_registry=registry.example.com
default_version=latest
revision_version=$(revision)

ifndef env
	env=stag_r
endif

ifndef env_alias
	env_alias=_$(env)
endif

ifeq ($(env), stag)
	default_registry=$(sat_default_registry)
	env_alias=
else ifeq ($(env), stag_r)
	default_registry=$(sat_default_registry)
else ifeq ($(env), stag_payout)
    default_registry=$(sat_default_registry)
else ifeq ($(env), stag_payout_r)
    default_registry=$(sat_default_registry)
else ifeq ($(env), prod)
	default_registry=$(prod_default_registry)
else ifeq ($(env), prod_payout)
    default_registry=$(prod_default_registry)
endif
ifndef registry
	registry=$(default_registry)
endif

default_port=$(shell ./port.sh $(service) $(env))
ifndef port
	port=$(default_port)
endif

ifndef version
	version=$(default_version)
endif

image_name=`echo "app/$(service)" | sed "s/\\//_/g"`
ifndef service_name
	service_name=$(image_name)
endif

build:
ifdef service
	cd $(service) && REVISION=$(revision) GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin -ldflags "-w -s -X main.revision="$(revision)
	chmod +x ./port.sh
	docker build -t $(registry)/$(image_name)$(env_alias)_runner:$(version) -t $(registry)/$(image_name)$(env_alias)_runner:$(revision_version) \
    		-f Dockerfile \
    		--build-arg REGISTRY=$(registry) \
            --build-arg ENV=$(env) \
            --build-arg SERVICE=$(service) \
            --build-arg PORT=$(port) \
            .
else
	echo "You have to use: make build service=service/core/metrics"
endif

build-with-chromium:
ifdef service
	cd $(service) && REVISION=$(revision) GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin -ldflags "-w -s -X main.revision="$(revision)
	chmod +x ./port.sh
	docker build -t $(registry)/$(image_name)$(env_alias)_runner:$(version) -t $(registry)/$(image_name)$(env_alias)_runner:$(revision_version) \
    		-f chromium.Dockerfile \
    		--build-arg REGISTRY=$(registry) \
            --build-arg ENV=$(env) \
            --build-arg SERVICE=$(service) \
            --build-arg PORT=$(port) \
            .
else
	echo "You have to use: make build service=service/core/metrics"
endif

deploy-with-chromium: build-with-chromium
ifdef service
	docker push $(registry)/$(image_name)$(env_alias)_runner:$(version)
	docker push $(registry)/$(image_name)$(env_alias)_runner:$(revision_version)
else
	echo "You have to use: make deploy service=service/core/metrics"
endif

deploy: build
ifdef service
	docker push $(registry)/$(image_name)$(env_alias)_runner:$(version)
	docker push $(registry)/$(image_name)$(env_alias)_runner:$(revision_version)
else
	echo "You have to use: make deploy service=service/core/metrics"
endif

deploy-new: deploy
ifdef service
	echo 'docker service create --limit-memory=2GB --update-order start-first --update-parallelism 1 --with-registry-auth -e env=$(env) --replicas 1 --name $(service_name)$(env_alias)_runner  --publish $(port):$(port) --update-delay 10s $(registry)/$(image_name)$(env_alias)_runner:$(version)'
else
	echo "You have to use: make deploy-new service=service/core/metrics"
endif

update:
ifdef service
	echo 'docker service update --force --update-order start-first --update-parallelism 1 --with-registry-auth --image $(registry)/$(image_name)$(env_alias)_runner:$(version) $(image_name)$(env_alias)_runner'
else
	echo "You have to use: make update service=service/core/metrics"
endif
