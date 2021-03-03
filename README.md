# paladinswr

## how to build
```sh
docker build -f config/docker/wrcli/Dockerfile -t yornesek/wrcli:$(git rev-parse --short HEAD) -t yornesek/wrcli:latest .
docker build -f config/docker/wrsvc/Dockerfile -t yornesek/wrsvc:$(git rev-parse --short HEAD) -t yornesek/wrsvc:latest .
docker build -f config/docker/frontend/Dockerfile -t yornesek/wrfrontend:$(git rev-parse --short HEAD) -t yornesek/wrfrontend:latest .
docker push yornesek/wrcli -a
docker push yornesek/wrsvc -a
docker push yornesek/wrfrontend -a
```
## how to deploy
```sh
kubectl -n paladinswr rollout restart deploy wrsvc
kubectl -n paladinswr rollout restart deploy wrfrontend
```