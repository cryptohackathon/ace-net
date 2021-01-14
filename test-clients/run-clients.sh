go build client.go
for i in {1..500}
do
    echo "Running client ${i}"
    sleep 0.1
    ./client &
done
