EXEC=regulatorTree

$(EXEC):main.go regulator.go
	GOARCH=arm GOOS=linux go build -o $(EXEC)

clean:
	rm -rf $(EXEC)

scp:
	expect ./scp.sh
