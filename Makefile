DISTDIR	:= ./_dist
APPDIR	:= ./cmd
APPS	:= $(notdir $(shell find $(APPDIR) -mindepth 1 -maxdepth 1 -type d))

.PHONY: $(APPS)
$(APPS):
	go build -o $(DISTDIR)/$@ $(APPDIR)/$@

.PHONY: build
build: clean $(APPS)

.PHONY: clean
clean:
	rm -rf $(DISTDIR)/*
