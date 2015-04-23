gohosts
======

go version of command line tool for hosts manager which only support linux,mac OSX, linux.

#Require

1. go  

#Install 

before install, you should allow go to read and write hosts. 

At mac osX, linux, just execute:

    $sudo chmod 777 /etc/hosts

At windows, make sure that go has permession to read and write hosts file which location at **"C:\Windows\System32\drivers\etc"** directory.

Then:

    $go get github.com/micanzhang/gohosts/hm
    $go install github.com/micanzhang/gohosts/hm

nothing echo means installed succcessfully.

#Useage

	$ go -h //help 
	$ go -l //list all hosts 
	$ go -e [editor] //open emacs to editor hosts file directly
	$ go -s //enable host config 
	$ go -r //disable host config 

Also, you you can combine "-g", "-d", "-p" parameters with "-l", "-s", "-s". -g can be ommited. "-g":group, "-d": domain, "-i":ip address

	$ go -l default

Result:

	#==== default
	127.0.0.1        localhost
	255.255.255.255  broadcasthost
	::1              localhost
	#127.0.0.1 		 test.com
	#====

	$ go -l -d localhost

Result:

	#==== default
	127.0.0.1        localhost
	::1              localhost test.com
	#====

    
#License

The MIT License (MIT)

Copyright (c) 2015 micanzhang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
