# weather-monster

### Running the application & tests

---
To run the application:
```bash
$ docker-compose up --build
```

To run the tests
```bash
$ go test -v ./...
```

### How to make API requests

---

Create City request
```bash
curl -XPOST http://localhost:3000/cities \
-d name=Berlin \
-d latitude=52.520008 \
-d longitude=13.404954
```

Update City request
```bash
curl -XPATCH http://localhost:3000/cities/{id} \
-d name=Potsdam \
-d latitude=52.520008 \
-d longitude=13.404954
-d version={version}
```

Delete City request
```bash
curl -XDELETE http://localhost:3000/cities/{id}
```

Create Temperature request
```bash
curl -XPOST http://localhost:3000/temperatures \
-d city_id=1 \
-d max=35 \
-d min=32
```

Get Forecast request 
```bash
curl http://localhost:3000/forecasts/{city_id}
```



### TODO

- [ ]  better test coverage 
- [ ]  refactoring of repetitive code scattered in logic and tests
- [ ]  improve error handling
- [ ]  improve logging
- [ ]  while http handlers are tested, the routes *need* to be tested too
- [ ]  improve code documentation

