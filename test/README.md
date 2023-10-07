# Integration Tests

The functionality of a proxy-server cannot be tested well by just using unit-tests.

Therefore it is essential to also use integration tests to check if the proxy handles actual traffic the way we would expect.

## Setup

We need a tester- and a proxy-VM.

* Add 'tester' user on both nodes
* Add service to run proxy-instances on proxy-vm

   ```bash
   # /etc/systemd/system/calamary@.service

   [Unit]
   Description=Service to run an instance of calamary proxy

   [Service]
   Type=simple
   User=proxy
   Group=proxy
   ExecStart=/tmp/calamary_%i -f /tmp/calamary_%i.yml

   StandardOutput=journal
   StandardError=journal
   SyslogIdentifier=cicd_calamary
   ```

   ```bash
   systemctl daemon-reload
   ```

* Add privileges to start/stop proxy on proxy-vm

   ```bash
   # /etc/sudoders.d/tester_calamary

   Cmnd_Alias TESTER_CALAMARY = \
   /bin/systemctl start calamary* ,\
   /bin/systemctl stop calamary*

   tester ALL=(ALL) NOPASSWD: TESTER_CALAMARY
   ```

* Configure tester to be able to connect via SSH from tester-VM to proxy-VM

* Install test-utils on tester-VM:

   ```bash
   sudo apt install git curl python3-pip python3-virtualenv
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

* You may want to add these lines to the `/home/tester/.bashrc` file

   ```bash
   source ~/venv/bin/activate
   export PATH=$PATH:/usr/local/go/bin
   ```
   

## Run

On the tester-VM:

```bash
source ~/venv/bin/activate

TMP_DIR="/tmp/calamary_$(date +%s)"
mkdir "$TMP_DIR"
cd "$TMP_DIR"

git clone https://github.com/superstes/calamary
bash calamary/test/wrapper.sh latest
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
