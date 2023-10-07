# Integration Tests

The functionality of a proxy-server cannot be tested well by just using unit-tests.

Therefore it is essential to also use integration tests to check if the proxy handles actual traffic the way we would expect.

## Adding

1. Add or extend test in one of the test-category directories:

   general, transparent, transparentTproxy, http, https, proxyproto, socks5

2. If it is a new test - add it to the list inside the main category-script (*per example in: 'testGeneral.sh'*)

General tests are ran against the 'transparent' mode.

### Running

1. You need an instance of the proxy up-and-running
2. Export the `PROXY_HOST` and `PROXY_PORT` variables
3. The `run_test` function should be imported by default (`source ./base.sh`)
4. Execute the test-scripts you want to check

   Per example: `bash testTransparent.sh`

## Setup

We need a tester- and a proxy-VM.

Why use a dedicated proxy-VM? The most-used implementation of transparent-mode is redirecting routed traffic by using the `prerouting` chain. That's also the way it should be tested. In addition to this - the `proxyproto` mode is also used to connect to a remote server.

* Add 'tester' user on both nodes
* Add service to run proxy-instances on proxy-vm

   See: [test/setup_proxy/](https://github.com/superstes/calamary/tree/latest/test/setup_proxy)

   ```bash
   systemctl daemon-reload
   ```

* Add privileges to start/stop proxy on proxy-vm

   ```bash
   # /etc/sudoders.d/tester_calamary

   Cmnd_Alias TESTER_CALAMARY = \
   /bin/systemctl start calamary* ,\
   /bin/systemctl stop calamary*, \
   /usr/bin/chown proxy\:proxy /tmp/calamary*, \
   /usr/bin/chown tester\:tester /tmp/calamary*, \
   /usr/bin/rm -f /tmp/calamary*

   tester ALL=(ALL) NOPASSWD: TESTER_CALAMARY
   ```

* Configure tester to be able to connect via SSH from tester-VM to proxy-VM

* Install test-utils on tester-VM:

   ```bash
   sudo apt install git curl python3-pip python3-virtualenv openssl
   su tester
   python3 -m virtualenv ~/venv/
   source ~/venv/bin/activate
   pip install anybadge
   ```

* [Download & install GO](https://go.dev/doc/install) on the tester-VM

* Create badge-directory, if needed (on the tester-VM)

   ```bash
   mkdir -p /var/www/cicd/calamary
   chown tester /var/www/cicd/calamary
   ```

* You may want to append these lines to the `/home/tester/.bashrc` file

   ```bash
   source ~/venv/bin/activate
   export PATH=$PATH:/usr/local/go/bin
   ```
   

## Run all tests

On the tester-VM:

```bash
source ~/venv/bin/activate

TMP_DIR="/tmp/calamary_$(date +%s)"
mkdir "$TMP_DIR"
cd "$TMP_DIR"

git clone https://github.com/superstes/calamary

# TODO: update the connection-settings in 'target.sh'

bash calamary/test/wrapper.sh latest
```

### As Service

Using a systemd service to run it:

See: [test/setup_client/](https://github.com/superstes/calamary/tree/latest/test/setup_client)

```bash
systemctl daemon-reload
```

To start:

```bash
systemctl start calamary-test@latest.service

# logs
systemctl status calamary-test@latest.service --no-pager --full
journalctl -u calamary-test@latest.service --no-pager --full -n 50
```

Example run:

```bash
systemd[1]: Started calamary-test@latest.service - Service to run integration tests for calamary proxy.
cicd_calamary[20673]: Cloning into 'calamary'...
cicd_calamary[20680]: TESTING VERSION 'latest' WITH TEST-VERSION 'latest-157e6d6b'
cicd_calamary[20680]: BUILDING BINARY (/tmp/calamary_1696688465/calamary)
cicd_calamary[20680]: STARTING TESTS
cicd_calamary[20680]: CLEANUP
cicd_calamary[20680]: STOPPING PROXY
cicd_calamary[20680]: PREPARING FOR TESTS
cicd_calamary[20680]: GENERATING CERTS
cicd_calamary[20680]: COPYING FILES TO PROXY-HOST
cicd_calamary[20680]: STARTING PROXY
cicd_calamary[20680]: STARTING TESTS
cicd_calamary[20680]: ##### RUNNING TESTS: TRANSPARENT #####
cicd_calamary[20680]: RUNNING TEST 'transparent/basic'
cicd_calamary[20680]: RUNNING TEST 'transparent/dummyOk'
cicd_calamary[20822]: Testing....
cicd_calamary[20680]: TEST-RUN FINISHED SUCCESSFULLY!
cicd_calamary[20680]: CLEANUP
cicd_calamary[20680]: STOPPING PROXY
systemd[1]: calamary-test@latest.service: Deactivated successfully.
```

## Workflow

We have a tester- (*client*) and a proxy-VM.

All actions are run by the tester.

* `wrapper.sh`

   Building the binary for the version we want to test.

* `main.sh`

  * Deploying the binary on the proxy.

  * Starting an instance of a test-service on the proxy (*runs the binary*)

  * Test-scripts are executed

    * They generate network traffic that should pass through the proxy
    * Using the multiple transport-/listener-modes Calamary provides
    * Checking if the responses match our expectations
    * Fail if testing if any fails

  * Clean-up & stop proxy-instance after tests finished or failed

If a timeout is reached - the proxy-instance and tests are terminated.
