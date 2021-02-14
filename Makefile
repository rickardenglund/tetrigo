FIND ?= find

.PHONY: clean

GO_FILES := $(shell $(FIND) . -name '*.go')

build: build/tetris build/assets/blocks.png build/assets/sound/clock.mp3 build/assets/intuitive.ttf

build/tetris: $(GO_FILES)
	install -d build
	go build -o build/tetris cmd/tetris/main.go

build/assets/blocks.png: assets/blocks.png
	install -d build/assets
	cp $< $@

build/assets/sound/clock.mp3: assets/sound/clock.mp3
	install -d build/assets/sound
	cp $< $@

build/assets/intuitive.ttf: assets/intuitive.ttf
	install -d build/assets
	cp $< $@

dist/mac.zip: build
	install -d dist
	zip -r $@ build/*

clean:
	$(RM) -r build dist
