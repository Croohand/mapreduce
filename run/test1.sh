rm /Users/therion/.mapreduce/logger/requests.log
simple_logger start &

coordinator start -name coordinator1 &
m1=$!
sleep 5
coordinator start -name coordinator2 &
m2=$!
coordinator start -name coordinator3 &
m3=$!

worker start -name scheduler1 &
h1=$!
worker start -name scheduler2 &
h2=$!

worker start -name worker1 &
s1=$!
worker start -name worker2 &
s2=$!
worker start -name worker3 &
s3=$!
worker start -name worker4 &
s4=$!
worker start -name worker5 &
s5=$!

sleep 2


MR_CLIENT="client1"; client read -path long_file | client write -path long_file1 &
MR_CLIENT="client1"; client read -path long_file | client write -path long_file2 &
MR_CLIENT="client2"; client mapreduce -in long_file -out long_file_count1 -mappers 5 -reducers 3 -srcs ~/go/src/github.com/Croohand/mapreduce/operations/word_count/mruserlib
sleep 2
kill $s3
kill $s4
sleep 1
MR_CLIENT="client3"; client mapreduce -in long_file -out long_file_count2 -mappers 3 -reducers 3 -srcs ~/go/src/github.com/Croohand/mapreduce/operations/word_count/mruserlib
sleep 2
kill $m1
sleep 5
MR_CLIENT="client1"; client read -path long_file | client write -path nlong_file


sleep 2

./kill.sh

echo "done"
