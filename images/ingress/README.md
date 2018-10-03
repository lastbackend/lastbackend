Build and run lastbackend/ingress image

```bash 
docker build -t lastbackend/ingress .
docker run -i -d --restart=always --name=ingress -p 80:80 -p 443:443 -v /etc/haproxy:/etc/haproxy lastbackend/ingress
```