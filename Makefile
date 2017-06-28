EXEC=regulatorView

$(EXEC):main.go regulator.go ui.go
	GOARCH=arm GOOS=linux go build -o $(EXEC)

clean:
	rm -rf $(EXEC)

scp:
	expect ./scp.sh
