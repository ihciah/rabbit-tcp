NAME=rabbit
BINDIR=bin
VERSIONPARAM=
ifdef RABBITVERSION
	VERSIONPARAM=-X 'main.Version=$(RABBITVERSION)'
endif
GOBUILD=CGO_ENABLED=0 go build -ldflags "-w -s $(VERSIONPARAM)"
BUILDFILE=cmd/rabbit.go

current:
	$(GOBUILD) -o $(BINDIR)/$(NAME) $(BUILDFILE)

all: linux-amd64 linux-386 linux-arm64 linux-arm darwin-amd64 darwin-386 windows-amd64 windows-386

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

linux-386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

linux-arm:
	GOARCH=arm GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

darwin-386:
	GOARCH=386 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@ $(BUILDFILE)

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(BUILDFILE)

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe $(BUILDFILE)

releases: linux-amd64 linux-386 linux-arm64 linux-arm darwin-amd64 darwin-386 windows-amd64 windows-386
	chmod +x $(BINDIR)/$(NAME)-*
	gzip $(BINDIR)/$(NAME)-linux-amd64
	gzip $(BINDIR)/$(NAME)-linux-386
	gzip $(BINDIR)/$(NAME)-linux-arm64
	gzip $(BINDIR)/$(NAME)-linux-arm
	gzip $(BINDIR)/$(NAME)-darwin-amd64
	gzip $(BINDIR)/$(NAME)-darwin-386
	zip -m -j $(BINDIR)/$(NAME)-windows-amd64.zip $(BINDIR)/$(NAME)-windows-amd64.exe
	zip -m -j $(BINDIR)/$(NAME)-windows-386.zip $(BINDIR)/$(NAME)-windows-386.exe

clean:
	rm $(BINDIR)/*