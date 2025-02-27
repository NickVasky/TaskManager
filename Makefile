all: TaskManager

clean: 
	rm -f TaskManager.out

rebuild: clean all

TaskManager:
	go build -o $@.out