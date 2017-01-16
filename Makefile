.PHONY: all binary 

all: binary

binary:
	go build  -o docker-auth-plugin .

clean:
	rm -f docker-auth-plugin
