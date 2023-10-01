.. _modes:

.. include:: ../_inc/head.rst

.. include:: ../_inc/in_progress.rst

#####
Modes
#####

Transparent
###########

Behaviour
=========

DNAT - TCP (plaintext)
----------------------

**Server**

.. code-block:: bash

    2023-10-01 23:43:01 | INFO | 192.168.11.104 => 135.181.170.219:80 | Accept


**Client**

.. code-block:: bash

    curl http://superstes.eu -v
    *   Trying 135.181.170.219:80...
    * Connected to superstes.eu (135.181.170.219) port 80 (#0)
    > GET / HTTP/1.1
    > Host: superstes.eu
    ...
    < 
    <html>
    <head><title>301 Moved Permanently</title></head>
    <body>
    <center><h1>301 Moved Permanently</h1></center>
    <hr><center>nginx</center>
    </body>
    </html>
    * Connection #0 to host superstes.eu left intact


DNAT - TLS
----------

**Server**

.. code-block:: bash

    2023-10-01 23:43:09 | INFO | 192.168.11.104 => 135.181.170.219:443 | Accept


**Client**

.. code-block:: bash

    host@calamary$ curl https://superstes.eu -v

    *   Trying 135.181.170.219:443...
    * Connected to superstes.eu (135.181.170.219) port 443 (#0)
    ...
    < HTTP/2 302 
    < server: nginx
    ...
    < 
    <html>
    <head><title>302 Found</title></head>
    <body>
    <center><h1>302 Found</h1></center>
    <hr><center>nginx</center>
    </body>
    </html>
    * Connection #0 to host superstes.eu left intact


HTTP Proxy
##########

Behaviour
=========

HTTP
----

**Server**

.. code-block:: bash

    2023-10-01 23:40:34 | INFO | 127.0.0.1 => 135.181.170.219:80 | Accept


**Client**

.. code-block:: bash

    host@calamary$ http_proxy=http://localhost:4130 curl http://superstes.eu -v

    * Uses proxy env variable http_proxy == 'http://localhost:4130'
    *   Trying 127.0.0.1:4130...
    * Connected to (nil) (127.0.0.1) port 4130 (#0)
    > GET http://superstes.eu/ HTTP/1.1
    > Host: superstes.eu
    > User-Agent: curl/7.81.0
    > Accept: */*
    > Proxy-Connection: Keep-Alive
    > 
    ...
    < 
    <html>
    <head><title>301 Moved Permanently</title></head>
    <body>
    <center><h1>301 Moved Permanently</h1></center>
    <hr><center>nginx</center>
    </body>
    </html>
    * Connection #0 to host (nil) left intact

HTTPS
-----

**Server**

.. code-block:: bash

    2023-10-01 23:40:43 | INFO | 127.0.0.1 => 135.181.170.219:443 | Accept


**Client**

.. code-block:: bash

    host@calamary$ https_proxy=http://localhost:4130 curl https://superstes.eu -v

    * Uses proxy env variable https_proxy == 'http://localhost:4130'
    *   Trying 127.0.0.1:4130...
    * Connected to (nil) (127.0.0.1) port 4130 (#0)
    * allocate connect buffer!
    * Establish HTTP proxy tunnel to superstes.eu:443
    > CONNECT superstes.eu:443 HTTP/1.1
    > Host: superstes.eu:443
    > User-Agent: curl/7.81.0
    > Proxy-Connection: Keep-Alive
    > 
    < HTTP/1.1 200 OK
    < Content-Length: 0
    * Ignoring Content-Length in CONNECT 200 response
    < 
    * Proxy replied 200 to CONNECT request
    * CONNECT phase completed!
    ...
    > GET / HTTP/2
    > Host: superstes.eu
    > user-agent: curl/7.81.0
    > accept: */*
    > 
    ...
    < HTTP/2 302 
    < server: nginx
    ...
    < 
    * TLSv1.2 (IN), TLS header, Supplemental data (23):
    <html>
    <head><title>302 Found</title></head>
    <body>
    <center><h1>302 Found</h1></center>
    <hr><center>nginx</center>
    </body>
    </html>
    * Connection #0 to host (nil) left intact


HTTPS Proxy
###########

Behaviour
=========

tbd

Proxy Protocol
##############

Behaviour
=========

tbd

SOCKS5
######

Behaviour
=========

tbd
