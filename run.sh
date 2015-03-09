./bin/authserver -log=auth-log.xml &
./bin/timeserver -log=seelog.xml -port=8081 -maxinflight=80 -response=500 -deviation=300 &
./bin/loadgen -url='http://localhost:8081/time' -runtime=10 -rate=200 --burst=20 -timeout=1000