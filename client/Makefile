CMD_DIR=cmd

SRCS := $(shell find . -name '*.go')
NAME=clit_client

all: $(NAME)

$(NAME): $(SRCS)
	go build -o $(NAME) cmd/main.go

run: $(NAME)
	./$(NAME)
