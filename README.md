# rct

# Summary
The project has two services
1. Parking
2. Booking

Both services have in memory datastore. *store.go from each of these services can be easily extended to create DB backed store
Parking handles the single responsibility of managing the parking slots.
Booking handles the single responsibility of booking/reserving slots. It makes use of the APIs of the parking service through a dependency injection
Each service has unit tests.
The project has a Makefile to build/run/test

# Build
````
make local-build
````

Below is a list of APIs that are implemented along with the response

# Get all parking slots
````
curl -X GET http://localhost:8080/parking/v1/getAll/
{"spots":[{"id":1,"lat":"44.968046","lon":"-94.420307","cost":"100","isReserved":false,"address":"address 1"},{"id":2,"lat":"44.33328","lon":"-89.132008","cost":"10","isReserved":false,"address":"address 2"},{"id":3,"lat":"33.755787","lon":"-116.359998","cost":"80","isReserved":false,"address":"address 3"},{"id":4,"lat":"33.844843","lon":"-116.54911","cost":"70","isReserved":false,"address":"address 4"},{"id":5,"lat":"44.92057","lon":"-93.44786","cost":"90","isReserved":false,"address":"address 5"}]}
````

# Get free/vacant parking slots
````
curl -X GET http://localhost:8080/parking/v1/getFree/
{"spots":[{"id":2,"lat":"44.33328","lon":"-89.132008","cost":"10","isReserved":false,"address":"address 2"},{"id":3,"lat":"33.755787","lon":"-116.359998","cost":"80","isReserved":false,"address":"address 3"},{"id":4,"lat":"33.844843","lon":"-116.54911","cost":"70","isReserved":false,"address":"address 4"},{"id":5,"lat":"44.92057","lon":"-93.44786","cost":"90","isReserved":false,"address":"address 5"},{"id":1,"lat":"44.968046","lon":"-94.420307","cost":"100","isReserved":false,"address":"address 1"}]}
````

# Get reserved parking slots when no slots are reserved
````
curl -X GET http://localhost:8080/parking/v1/getReserved/
{"spots":[]}
````

# Get reserved parking slots when some slots are reserved
````
curl -X GET http://localhost:8080/parking/v1/getReserved/
{"spots":[{"id":1,"lat":"44.968046","lon":"-94.420307","cost":"100","isReserved":true,"address":"address 1"},{"id":2,"lat":"44.33328","lon":"-89.132008","cost":"10","isReserved":true,"address":"address 2"}]}
````

# Find parking slot by incorrectid gives an error
````
curl -X GET http://localhost:8080/parking/v1/find/10
{"error":"not found"}
````

# Find parking slot by id
````
curl -X GET http://localhost:8080/parking/v1/find/1
{"spots":[{"id":1,"lat":"44.968046","lon":"-94.420307","cost":"100","isReserved":false,"address":"address 1"}]}
````

# Search for a vacant spot by cost metric
````
curl -d '{"lat":"33.755787", "lon":"-116.359998", "rad":"10000", "metric":"cost"}' -X GET http://localhost:8080/parking/v1/search/
{"spots":[{"id":3,"lat":"33.755787","lon":"-116.359998","cost":"80","isReserved":false,"address":"address 3","distance":0}]}
````

# Search for a vacant spot by distance metric
````
curl -d '{"lat":"33.755787", "lon":"-116.359998", "rad":"10000", "metric":"dist"}' -X GET http://localhost:8080/parking/v1/search/
````


# Book spotId 1
````
curl -d '{"id":"1"}' -X POST http://localhost:8080/booking/v1/
{"booking":{"id":1,"spotId":1,"startTime":"2018-07-27T10:52:07.596575833+05:30","duration":1800000000000}}
````

# Book an already booked spot results in error
````
curl -d '{"id":"1"}' -X POST http://localhost:8080/booking/v1/
{"error":"spot already reserved"}
````

# Delete booking id 1
````
curl -X DELETE http://localhost:8080/booking/v1/1
{}
````

# View bookings
````
curl -X GET http://localhost:8080/booking/v1/
{"bookings":[{"id":1,"spotId":1,"startTime":"2018-07-27T11:28:27.413230484+05:30","duration":1800000000000}]}
````


# Additional features
## Automated tests
### Run test
````
make test
````

## Dockerized images
### Create a docker build
````
make docker-build
````

### Run in a docker container - prerequisite that docker is installed locally
````
make docker-run
````

## API metrics
### Get api metrics - results in a ton of metrics including those like GC etc
````
curl -X GET http://localhost:8080/metrics
````

### metrics for a service
````
curl -X GET http://localhost:8080/metrics | grep booking

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  6220  100  6220    0     0  1099k      0 --:--:-- --:--:-- --:--:-- 1214k
# HELP api_booking_service_request_count Number of requests received.
# TYPE api_booking_service_request_count counter
api_booking_service_request_count{method="Book"} 2
# HELP api_booking_service_request_latency_microseconds Total duration of requests in microseconds.
# TYPE api_booking_service_request_latency_microseconds summary
api_booking_service_request_latency_microseconds{method="Book",quantile="0.5"} 0.000105819
api_booking_service_request_latency_microseconds{method="Book",quantile="0.9"} 0.000406599
api_booking_service_request_latency_microseconds{method="Book",quantile="0.99"} 0.000406599
api_booking_service_request_latency_microseconds_sum{method="Book"} 0.0005124179999999999
api_booking_service_request_latency_microseconds_count{method="Book"} 2
````

## Prometheus logging indicating the time each api execution took and the api params
````
ts=2018-07-27T05:58:22.727276098Z caller=main.go:94 transport=http address=:8080 msg=listening
ts=2018-07-27T05:58:23.877170071Z caller=middleware.go:28 method=GetAll took=2.094µs err=null
ts=2018-07-27T05:58:27.413240784Z caller=middleware.go:57 method=Find id=1 took=2.093µs err=null
ts=2018-07-27T05:58:27.413358681Z caller=middleware.go:64 method=Update id=1 took=1.531µs err=null
ts=2018-07-27T05:58:27.413402461Z caller=middleware.go:35 method=Book spotId=1 startTime=2018-07-27T11:28:27.413230484+05:30 duration=30m0s took=169.667µs err=null
ts=2018-07-27T05:58:28.796801982Z caller=middleware.go:28 method=GetAll took=6.792µs err=null
ts=2018-07-27T06:01:49.309385489Z caller=middleware.go:50 method=Search lat=33.755787 lon=-116.359998 radius=10000 metric=cost took=14.971µs err=null
ts=2018-07-27T06:01:56.348313267Z caller=middleware.go:50 method=Search lat=33.755787 lon=-116.359998 radius=10000 metric=dist took=13.445µs err=null
ts=2018-07-27T06:04:37.668735885Z caller=middleware.go:36 method=GetFreeParking took=18µs err=null
ts=2018-07-27T06:04:47.990176162Z caller=middleware.go:57 method=Find id=1 took=2.497µs err=null
ts=2018-07-27T06:04:47.990262103Z caller=middleware.go:35 method=Book spotId=1 startTime=2018-07-27T11:34:47.990165975+05:30 duration=30m0s took=93.305µs err="spot already reserved"
ts=2018-07-27T06:04:47.990344384Z caller=server.go:112 component=HTTP err="spot already reserved"
ts=2018-07-27T06:04:53.859233491Z caller=middleware.go:57 method=Find id=2 took=1.945µs err=null
ts=2018-07-27T06:04:53.859274782Z caller=middleware.go:64 method=Update id=2 took=1.182µs err=null
ts=2018-07-27T06:04:53.859312417Z caller=middleware.go:35 method=Book spotId=2 startTime=2018-07-27T11:34:53.859227032+05:30 duration=30m0s took=82.515µs err=null
ts=2018-07-27T06:05:05.570680725Z caller=middleware.go:43 method=GetReservedParking took=3.001µs err=null
````

# TODO
User service to enable logging and booking by a user and validation etc
