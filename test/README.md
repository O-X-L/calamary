# Integration Tests

The functionality of a proxy-server cannot be tested well by just using unit-tests.

Therefore it is essential to also use integration tests to check if the proxy handles actual traffic the way we would expect.

## Setup

We have a tester- (*client*) and a proxy-VM.

1. The tester builds the binary for the version we want to test.
2. The tester deploys the binary on the proxy.
3. The tester starts an instance of a test-service on the proxy (*runs the binary*)
4. Test-scripts are executed

  * They generate network traffic that should pass through the proxy
  * Using the multiple transport-/listener-modes Calamary provides
  * Checking if the responses match our expectations
  * Fail if testing if any fails

5. Clean-up & stop proxy-instance after tests finished or failed

If a timeout is reached - the proxy-instance and tests are terminated.

## Run

```bash
git clone https://github.com/superstes/calamary
bash test/wrapper.sh
```
