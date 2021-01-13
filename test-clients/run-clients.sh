for i in {1..100}
do
    echo "Running client ${i}"
    go run client.go &
done
