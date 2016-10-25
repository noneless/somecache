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
	go get github.com/756445638/somecache
	cd ${GOPATH}/src/github.com/756445638/somecache
	go run  examples/pub.go 
	go run slave.go -cachesize 500 -worker=3   //worker is like a donwload thread
	go run slave.go -cachesize 500 -worker 2 -master-tcp-address 192.168.1.8:4000
	
	there are some other messages dump into stdout,dont`t mind it to much
	pub.go is tool that read from stdio,and get or put message 
	support cmds are:
		get remote aaa
		get aaa
		put remote aaa bbb 
		put aaa bbb
	put&get with "remote" will  work wite slave,it is a debug option,in production should look localcache first


somecache is written in pure golang,somecache is not stable currently,welcome to test,debug,pull, fork,and advise!!
approiate anything