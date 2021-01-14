go build client.go
for i in {1..100}
do
    echo "Running client ${i}"
    ./client &
done
