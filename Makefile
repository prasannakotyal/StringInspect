# Makefile for StringInspect

CC = gcc
CFLAGS = -Wall -Wextra -std=c11
TARGET = stringinspect
SRC = src/stringinspect.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

clean:
	rm -f $(TARGET)

install:
	cp $(TARGET) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(TARGET)

.PHONY: all clean install uninstall
