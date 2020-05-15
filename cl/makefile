# makefile
# Makefile for building clTrace system
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved


TARGET=cltrace

CC=g++
CFLAGS=-Iinclude -std=c++11 -Wall

OBJ_FILES= $(patsubst obj/%.cpp,obj/%.o,$(patsubst src/%,obj/%,$(wildcard src/*.cpp) ))
#$(info $$OBJ_FILES is [${OBJ_FILES}])

all: $(TARGET)

$(TARGET): $(OBJ_FILES)
	$(CC) $(CFLAGS) $^ -o $@

obj/%.o : src/%.cpp | obj
	$(CC) $(CFLAGS)  $<  -c  -o $@

.PHONY: clean
clean:
	rm -rf obj
	rm $(TARGET)

obj:
	mkdir -p obj
