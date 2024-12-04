#!/bin/bash

CHARACTERS=("lily" "chad" "mike" "timothy" "squidward")
LOGPATH="logs/"



function start_character {
	character="$1"
	mkdir -p "$LOGPATH"
	pidpath="$character.pid"
	if [ -f "$pidpath" ]; then
		echo "character $character already running!"
		return 1
	fi
	./artifactsMMO "$character" > "$LOGPATH/$character.log" &
}

function start_all_characters {
	for character in "${CHARACTERS[@]}"
	do
		start_character "$character"
	done
}

function kill_character {
	character="$1"
	pidpath="$character.pid"
	kill `cat $pidpath`
	rm "$pidpath"
}

function kill_all_characters {
	for character in "${CHARACTERS[@]}"
	do
		kill_character "$character"
	done
}

function mmo_status {
	for character in "${CHARACTERS[@]}"
	do
		pidpath="$character.pid"
		if [ -f "$pidpath" ]; then
			echo -e "$character:" "\033[32mUP\033[0m"
		else
			echo -e "$character:" "\033[0;31mDOWN\033[0m"
		fi
	done
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then

	function usage {
		echo "Usage: $0 <function_name> [arguments...]"
		echo "Available functions:"
		echo "  mmo_status                            - Print status"
		echo "  start_character <character_name>      - Start a specific character"
		echo "  start_all_characters                  - Start all characters"
		echo "  kill_character <character_name>       - Kill a specific character"
		echo "  kill_all_character <character_name>   - Kill a specific character"
		return 1
	}

	if [ $# -lt 1 ]; then
		usage
	fi

	function_name=$1
	shift

	if declare -f "$function_name" > /dev/null; then
	  $function_name "$@"
	else
		usage
	fi
fi