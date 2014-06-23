xhprof-import-proxy
===================

Import request queue for xhprof.io 

XHProf is a great util to profile php requests and xhprof.io is great to visualize then. If you want to profile 
every request to your application, the SQL request to store the data will most likely slow down your application. 

This tool creates a http service that proxies and queues the request to sore data in xhprof.io, so that the client 
application just have to send the data and will get a 202 Accepted status immediately.

Installation
------------

```
go build -o bin/xhprof-import-proxy -a -i main
```

Copy import.php to xhprof.io install path.

Start the server
----------------

```
bin/xhprof-import-proxy --listen=":8080" --path="/xhprof.io/import" --xhProfIoUrl="http://localhost/xhprof.io/import.php"
```

