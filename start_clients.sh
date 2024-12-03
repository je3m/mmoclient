#!/bin/bash

CHARACTERS=("lily" "chad" "mike" "timothy" "squidward")
LOGPATH="logs/"

mkdir "$LOGPATH"

for character in "${CHARACTERS[@]}"
do
	./artifactsMMO "$character" > "$LOGPATH/$character.log" &
done


