go build
./slave start -name slave1 -port 11001 -master http://localhost:11000 -override &
./slave start -name slave2 -port 11002 -master http://localhost:11000 -override &
./slave start -name slave3 -port 11003 -master http://localhost:11000 -override &
./slave start -name slave4 -port 11004 -master http://localhost:11000 -override & 
./slave start -name slave5 -port 11005 -master http://localhost:11000 -override &
