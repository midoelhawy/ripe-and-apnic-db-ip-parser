#!/bin/bash

set -o verbose

input_file="db/ripe.db.inetnum"
output_root_folder="db/chunks"
max_blocks_per_file=10000
max_files_per_dir=200
block_counter=0
file_counter=0
dir_counter=1
found_inetnum=false
last_dir_num=-1
last_file_num=-1
mkdir -p "$output_root_folder"

create_new_dir() {
    mkdir -p "$output_root_folder/$dir_counter"
}
skipped_lines=0
current_block_data=""
while IFS= read -r line; do
    start_new_block=false
    
    if [[ $line == "inetnum:"* ]]; then
        # echo "Found inetnum"
        found_inetnum=true
    fi
    
    if ! $found_inetnum; then
        ((skipped_lines++))
        
        # echo "Skipped $skipped_lines"
        continue;
    fi
    
    
    if [[ $line == "remarks:"* ]]; then
        #echo "skipping remarks"
        continue
    fi
    
    if [[ $line == "inetnum:"* ]]; then
        # echo "Found inetnum"
        
        
        ((block_counter++))
        start_new_block=true
    else
        start_new_block=false
    fi
    
    
    file_counter=$(( $block_counter / $max_blocks_per_file ))
    dir_counter=$(( $file_counter / $max_files_per_dir ))
    
    if ((last_dir_num != dir_counter)); then
        #echo "create new dir $dir_counter; $file_counter; $start_new_block"
        
        
        
        create_new_dir
    fi
    
    if ((last_file_num != file_counter)); then
        if [ -n "$current_block_data" ]; then
            echo "output: $output;block_counter:$block_counter"
            echo -en "$current_block_data" >> "$output"
        fi
        current_block_data=""
    fi
    
    last_dir_num=$dir_counter
    last_file_num=$file_counter
    output="$output_root_folder/$dir_counter/$file_counter.db"

    
    current_block_data+="$line\n"
done < "$input_file"

echo "Split complete."
