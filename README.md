# easyd
simple webhook deploy

## Getting start
```bash
easyd [token]
```
After run this command,it would serving with port 8082.
```bash
curl your_ip|domain:8082/deploy?script=hello&token=[your token]
```
It will run `hello.sh` on the `scripts` dir


