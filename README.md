# somecache
high performance and scalable memory cache system

#master 
	it is a lib that manage you cluster  
	it a lib that reuse connection reuse connection make master looks like gateway
	it is a monitor your slaves and decide which slave you will use
	it also a localcache that store data in local machine`s memory

#slave 
	it is a node save you data in memory
	slave is a standalone program
	
#use case 
	example.go is program that after it has been started ,it will push a lot of data in server,you can see a lot of
	"print" dont`t be nervous
	go run slave.go -cachesize 500 -worker=3   //worker is like a donwload thread
	go run slave.go -cachesize 500 -worker 2 -master-tcp-address 192.168.1.8:4000
	go run example.go 
	

somecache is written in pure golang