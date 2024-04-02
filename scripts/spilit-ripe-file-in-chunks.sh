#!/bin/bash

set -o verbose

input_file="db/ripe.db.inetnum"
output_root_folder="db/chunks"
max_blocks_per_file=400
max_files_per_dir=200
block_counter=0
file_counter=0
dir_counter=1
found_inetnum=false

mkdir -p "$output_root_folder"

create_new_dir() {
    mkdir -p "$output_root_folder/dir_$dir_counter"
}
skipped_lines=0
while IFS= read -r line; do
    #cho "$line"
    if [[ $line == "inetnum:"* ]]; then
        # echo "Found inetnum"
        found_inetnum=true
    fi
    
    if ! $found_inetnum; then
        ((skipped_lines++))
        
        echo "Skipped $skipped_lines"
        continue;
    fi
    

    if [[ $line == "remarks:"* ]]; then
        echo "skipping remarks"
        continue
    fi

    if [[ $line == "inetnum:"* ]]; then
        # echo "Found inetnum"
        found_inetnum=true
    fi


    #echo "Start $skipped_lines"
    
    if ((block_counter % max_blocks_per_file == 0)); then
        if ((block_counter > 0 && block_counter % (max_blocks_per_file * max_files_per_dir) == 0)); then
            ((file_counter % max_files_per_dir == 0)) && ((dir_counter++))
            create_new_dir
        fi
        ((file_counter++))
        output="$output_root_folder/dir_$dir_counter/block_$file_counter.db"
        echo "Write output in $output"
        create_new_dir
    fi
    
    echo "$line" >> "$output"
    ((block_counter++))
done < "$input_file"

echo "Split complete."
