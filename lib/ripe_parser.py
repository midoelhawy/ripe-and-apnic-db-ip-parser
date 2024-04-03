import ipaddress
from typing import Callable


class RIPE_PARSER:
    def __init__(self):
        pass

        

    def format_block(block):
        new_block = {}
        inetnum_splited = block["inetnum"].split(" - ")
        new_block["first_ip"] = inetnum_splited[0]
        new_block["last_ip"] = inetnum_splited[1] if len(inetnum_splited) > 1 else inetnum_splited[0] 
        firstIp = ipaddress.ip_address(new_block["first_ip"])
        lastIp = ipaddress.ip_address(new_block["last_ip"])
        new_block["first_ip_int"] = int(firstIp)
        new_block["last_ip_int"] = int(lastIp)
        new_block["ip_count"] = new_block["last_ip_int"] - new_block["first_ip_int"] + 1
        new_block["netname"] = block.get("netname", "Unknown")
        new_block["country"] = block.get("country", "Unknown")
        new_block["descr"] = block.get("descr", "Unknown")
        new_block["mnt-by"] = block.get("mnt-by", "Unknown")
        
        return new_block
    
    def parse_file(file_path,cb:Callable[[dict],None]):
        data = []
        with open(file_path, 'r',-1,"latin-1") as file:
            block = {}
            for line in file:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if line.startswith("inetnum:"):
                    if block:
                        # data.append(RIPE_PARSER.format_block(block))
                        cb(RIPE_PARSER.format_block(block))
                        block = {}
                if line and line.find(":") >= 0:
                    key, value = line.split(":", 1)
                    #Note : This block is to avoid overwrite the information like mnt-by 
                    if key in block:
                        if key == "descr":
                            block[key.strip()] += "\n" + value.strip()
                    else:
                        block[key.strip()] = value.strip()
            if block:  
                cb(RIPE_PARSER.format_block(block))
                # data.append(RIPE_PARSER.format_block(block))
        return data
    