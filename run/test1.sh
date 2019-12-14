rm /Users/therion/.mapreduce/logger/requests.log
simple_logger start &

master start -name master1 &
m1=$!
sleep 5
master start -name master2 &
m2=$!
master start -name master3 &
m3=$!

slave start -name scheduler1 &
h1=$!
slave start -name scheduler2 &
h2=$!

slave start -name slave1 &
s1=$!
slave start -name slave2 &
s2=$!
slave start -name slave3 &
s3=$!
slave start -name slave4 &
s4=$!
slave start -name slave5 &
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
