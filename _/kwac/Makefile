.PHONY: index.html mem o
index.html: s h o
	sh mkhtml

deploy: index.html
	cp index.html ../../ktye.github.io

o: s
	#k s>0
	i s>o

mem:
	i -kwac > mem

# https://github.com/WebAssembly/wabt/releases
WABT=~/wabt
valid: o
	$(WABT)/wasm-validate o

dump:	o
	$(WABT)/wasm-objdump -sd o

wat:	o
	$(WABT)/wasm2wat o
	
