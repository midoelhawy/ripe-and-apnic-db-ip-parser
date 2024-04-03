from pathlib import Path
import sys
path_root = Path(__file__).parents[1]
sys.path.append(str(path_root))
print(sys.path)

from lib.ripe_parser import RIPE_PARSER

# file_path = str(Path.joinpath(Path(__file__).parents[1],'test/chunk.db'))
# file_path = "/home/ahmedhekal/work/personal/ripe-ip-parser/db/chunks/0/0.db"
file_path = "./ripe-ip-parser/db/ripe.db.inetnum"

def on_single_block_process(block):
    print(block)

parse_data = RIPE_PARSER.parse_file(file_path,on_single_block_process)
print(parse_data)