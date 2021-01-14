go build client.go
for i in {1..10}
do
    echo "Running client ${i}"
    ./client &
done
