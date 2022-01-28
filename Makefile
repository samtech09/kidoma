APP_VER := "0.1.0"
PREVIOUS_VER := "0.0.0"

GOOS=linux
GOARCH=amd64
BUILDPATH=$(CURDIR)
BINPATH=$(BUILDPATH)/cmd
EXENAME=kidoma


.PHONY: default
default: help ;


.PHONY: clean
## clean: clean the binaries
clean:
	@rm -rf $(BINPATH)

	
.PHONY: build
## build: build application for current OS i.e. Linux
build: clean
	@if [ ! -d $(BINPATH) ] ; then mkdir -p $(BINPATH) ; fi
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 GO111MODULE=off go build -ldflags "-X main.appVer=$(APP_VER) -X main.buildVer=`date -u +b-%Y%m%d.%H%M%S`" -o $(BINPATH)/$(EXENAME) . || (echo "build failed $$?"; exit 1)
	#@go build -ldflags "-X main.appVer=$(APP_VER) -X main.buildVer=`date -u +b-%Y%m%d.%H%M%S`" -o $(BINPATH)/$(EXENAME) . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: buildarm
## buildarm: build application for arm architecture
buildarm: GOARCH=arm GOARM=6
buildarm: clean
	@if [ ! -d $(BINPATH) ] ; then mkdir -p $(BINPATH) ; fi
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=off go build -ldflags "-X main.appVer=$(APP_VER) -X main.buildVer=`date -u +b-%Y%m%d.%H%M%S`" -o $(BINPATH)/$(EXENAME) . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: run
## run: build and run application
run: build --prepfiles-for-package
	@cd $(BINPATH) && ./$(EXENAME)



.PHONY: package
## package: build and package application as zip file
package: build --prepfiles-for-package
	@# create new zip file
	@7z a $(BINPATH)/${EXENAME}.zip $(BINPATH)/$(EXENAME) $(BINPATH)/assets
	@echo ZIP Package created.


.PHONY: package-arm
## package-arm: build for arm and package application as zip file
package-arm: build --prepfiles-for-package
	@# create new zip file
	@7z a $(BINPATH)/${EXENAME}.zip $(BINPATH)/$(EXENAME) $(BINPATH)/config.json $(BINPATH)/assets
	@echo ZIP Package created.


--prepfiles-for-package:
	@-rm $(PUBLISHPATH)/$(EXENAME).zip
	@-mkdir $(BINPATH)/assets
	@cp -r assets $(BINPATH)/
	@cp config.json $(BINPATH)/
	@echo ...files synced


.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

