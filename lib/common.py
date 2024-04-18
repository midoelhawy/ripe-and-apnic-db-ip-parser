import ipaddress


def netmask_from_first_last_ip(first_ip, last_ip):
    first_ip = ipaddress.IPv4Address(first_ip)
    last_ip = ipaddress.IPv4Address(last_ip)
    
    num_hosts = int(last_ip) - int(first_ip)
    num_bits = num_hosts.bit_length()

    netmask_bits = 32 - num_bits
    return netmask_bits