GIT_REV:=$(shell git rev-parse --short HEAD)

LOCAL_MODS:=$(shell git status --short)
ifeq ($(LOCAL_MODS),)
VERSION:=${GIT_REV}
else
VERSION:=${GIT_REV}+dev
endif

# directories and remove vendor directory
# used for running unit tests
NOVENDOR := $(shell go list -e ./... | grep -v /vendor/ )

TEST_COVERAGE := $(shell go list -e ./... | grep -v /vendor/ )

# files and remove vendor directory, auto generated files, and mock files
# used for static analysis, code linting, and code formatting
NOVENDOR_FILES := $(shell find . -name "*.go" | grep -v /vendor/ | grep -v "_gen\.go" )

all: deploy

deploy: build
	rm -f howto.zip
	zip howto.zip howto

build:
	go build --ldflags '-X github.com/thoughtdealership/howto/app/version.Version=$(VERSION)' -o howto
	go test -run ^$$ $(NOVENDOR) > /dev/null

test: vet
	@# use redirect instead of tee to preserve exit code
	go test -short -p=1 -cover $(NOVENDOR) > report.out; \
	code=$$?; \
	cat report.out; \
	grep -e 'FAIL' report.out; \
	exit $${code}

cover:
	@echo "NOTE: make cover does not exit 1 on failure, don't use it to check for tests success!"
	mkdir -p .cover
	rm -f .cover/*.out .cover/all.merged
	@for MOD in $(TEST_COVERAGE); do \
		go test -coverpkg=`echo $(TEST_COVERAGE)|tr " " ","` \
			-coverprofile=.cover/unit-`echo $$MOD|tr "/" "_"`.out \
			$$MOD 2>&1 | grep -v "no packages being tested depend on"; \
	done
	gocovmerge .cover/*.out > .cover/all.merged
	go tool cover -html .cover/all.merged
	@echo ""
	@echo "=====> Total test coverage: <====="
	@echo ""
	$Q go tool cover -func .cover/all.merged

cloc:
	find . -name "*.go" | grep -v vendor | grep -v "_test.go" | xargs cloc

cloc-test:
	find . -name "*.go" | grep -v vendor | grep "_test.go" | xargs cloc

cloc-vendor:
	@for DIR in $(shell find vendor -mindepth 3 -maxdepth 3); do \
		echo $$DIR `find $$DIR -name "*.go" | xargs wc -l | tail -n1 | sed -e 's/^[ \t]*//' | cut -d ' ' -f1`; \
	done;

fmt:
	for FILE in $(NOVENDOR_FILES); do go fmt $$FILE; done;

vet:
	@echo 'Running go tool vet -shadow'
	@for FILE in $(NOVENDOR_FILES); do go tool vet -shadow $$FILE || exit 1; done;

lint:
	for FILE in $(NOVENDOR_FILES); do golint $$FILE; done;

clean:
	rm -rf howto howto.zip report.out .cover

version:
	@echo '$(VERSION)'