user=jane
device=phone

payload=$(jo _type=location \
   t=u \
   batt=11 \
   lat=48.856826 \
   lon=2.292713 \
   tid=JJ \
   tst=$(date +%s))

echo ${payload}
curl -v -X POST http://192.168.0.100:8090/location -H "Content-Type: application/json" -d "${payload}" 
