# paladinswr

## how to build
```sh
docker build -f config/docker/wrcli/Dockerfile -t yornesek/wrcli:$(git rev-parse --short HEAD) -t yornesek/wrcli:latest .
docker push yornesek/wrcli -a
```
## how to deploy
```sh
kubectl -n hello rollout restart deploy hello
```