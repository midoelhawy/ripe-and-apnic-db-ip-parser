

from pathlib import Path
from lib.db import SQLiteHandler
from lib.ripe_parser import RIPE_PARSER

blocks = []
total_blocks_processed= 0
if __name__ == "__main__":
    db_name = 'ripe_data.db'


    default_ripe_data = str(Path.joinpath(Path(__file__).parents[0],'db/ripe.db.inetnum'))
    db_handler = SQLiteHandler(db_name)
    db_handler.create_table()

    def on_single_block_process(block):
        global blocks,total_blocks_processed
        blocks.append(block)
        if len(blocks) >= 1000:
            
            db_handler.insert_data(blocks)
            total_blocks_processed += len(blocks)
            blocks = []
            print(f"Total blocks processed: {total_blocks_processed}")
    
    RIPE_PARSER.parse_file(default_ripe_data,on_single_block_process)
    

    print("Done")