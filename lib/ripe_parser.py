import ipaddress
from typing import Callable


class RIPE_PARSER:
    def __init__(self):
        pass

    
    def get_ip_v6_first_and_last_ip(sub_net):
        ip,subnet = sub_net.split("/",1)
        first_ip = ipaddress.IPv6Address(ip)
        last_ipv6 = (ipaddress.IPv6Address(first_ip) + (2 ** (128 - int(subnet)) - 1))
        return first_ip.exploded,last_ipv6.exploded,ip,int(first_ip),int(last_ipv6),subnet
        

    def format_block(block):
        new_block = {}
        if block.get("ipVersion",4) == 6:
            first_ip,last_ip,prefex,first_ip_int,last_ip_int,subnet  = RIPE_PARSER.get_ip_v6_first_and_last_ip(block["inetnum"])
            new_block["first_ip"] = first_ip
            new_block["last_ip"] = last_ip
            new_block["network_prefix"] = prefex
            new_block["first_ip_int"] = first_ip_int
            new_block["last_ip_int"] = last_ip_int
            new_block["subnet"] = subnet
        else:
            inetnum_splited = block["inetnum"].split(" - ")
            new_block["first_ip"] = inetnum_splited[0].strip()
            new_block["last_ip"] = inetnum_splited[1].strip() if len(inetnum_splited) > 1 else inetnum_splited[0].strip() 
            firstIp = ipaddress.ip_address(new_block["first_ip"])
            lastIp = ipaddress.ip_address(new_block["last_ip"])
            new_block["first_ip_int"] = int(firstIp)
            new_block["last_ip_int"] = int(lastIp)
        new_block["netname"] = block.get("netname", "Unknown")
        new_block["country"] = block.get("country", "Unknown")
        new_block["descr"] = block.get("descr", "Unknown")
        new_block["mnt-by"] = block.get("mnt-by", "Unknown")
        new_block["ip_version"] = block.get("ipVersion", 4)
        
        return new_block
    
    def parse_file(file_path,cb:Callable[[dict],None]):
        data = []
        with open(file_path, 'r',-1,"latin-1") as file:
            block = {}
            for line in file:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if line.startswith("inetnum:") or line.startswith("inet6num:"):
                    if block:
                        # data.append(RIPE_PARSER.format_block(block))
                        cb(RIPE_PARSER.format_block(block))
                        block = {}
                    if line.startswith("inet6num:"):
                        block["inetnum"] = line[8:]
                if line and line.find(":") >= 0:
                    key, value = line.split(":", 1)
                    if key == "inet6num":
                        block["inetnum"] = value.strip()
                        block["ipVersion"] = 6
                        continue
                    elif key == "inetnum":
                        block["inetnum"] = value.strip()
                        block["ipVersion"] = 4
                        continue
                    
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
    