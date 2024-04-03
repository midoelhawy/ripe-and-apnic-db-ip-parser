import os
from typing import Callable
from lib.ripe_parser import RIPE_PARSER

def explore_folder(folder_path, cb:Callable[[dict],None]):
    
    chunk_groups = sorted(os.listdir(folder_path))
    for dir_name in chunk_groups:
        sub_folder_path = os.path.join(folder_path, dir_name)
        if os.path.isdir(sub_folder_path):
            for root, dirs, files in os.walk(sub_folder_path):
                for file in sorted(files):
                    if file.endswith(".db"):
                        file_path = os.path.join(root, file)
                        print(f"Processing file {file_path}")
                        

                        
                        RIPE_PARSER.parse_file(file_path,cb)
                        


        
