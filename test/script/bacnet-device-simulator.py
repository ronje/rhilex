import socket
import struct

# BACnet Whois 服务的请求和响应数据结构
BACNET_WHOIS_REQUEST = 0x00
BACNET_WHOIS_RESPONSE = 0x01

# 设备ID
DEVICE_ID = 2654428

# 设备信息
DEVICE_NAME = "Example Device"
DEVICE_TYPE = "DEVICE"
DEVICE_ADDRESS = "192.168.1.100"


# BACnet Whois 请求的数据结构
def bacnet_whois_request(device_id):
    return struct.pack(">BHH", BACNET_WHOIS_REQUEST, 0, device_id)


# BACnet Whois 响应的数据结构
def bacnet_whois_response(device_id, device_name, device_type, device_address):
    # 设备信息长度
    device_info_length = len(device_name) + len(device_type) + len(device_address)
    # 设备信息
    device_info = device_name + device_type + device_address
    # 构建响应报文
    return (
        struct.pack(">BHH", BACNET_WHOIS_RESPONSE, device_info_length, device_id)
        + device_info
    )


# 创建一个UDP套接字
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

# 绑定到广播地址和端口
sock.bind(("", 47808))

# 监听广播
print("Listening for Whois requests on port 47808...")

while True:
    # 接收数据
    data, addr = sock.recvfrom(1024)
    print(data, addr)
    # 检查是否为Whois请求
    if data[0] == BACNET_WHOIS_REQUEST:
        # 解包Whois请求
        request_id = struct.unpack(">H", data[1:3])[0]
        if request_id == DEVICE_ID:
            # 构建Whois响应
            response = bacnet_whois_response(
                DEVICE_ID, DEVICE_NAME, DEVICE_TYPE, DEVICE_ADDRESS
            )
            # 发送响应
            sock.sendto(response, addr)
            print(f"Sent Whois response to {addr}")

    # 检查是否需要退出
    if data == b"exit\n":
        print("Exiting...")
        break

# 关闭套接字
sock.close()
