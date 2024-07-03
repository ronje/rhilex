# UDP Server
## 数据格式
JSON格式字符串，注意要考虑UDP的数据长度限制。
## 示例
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

int main() {
    int sock = socket(AF_INET, SOCK_DGRAM, 0);
    if (sock == -1) {
        printf("Failed to create socket\n");
        return 1;
    }
    struct sockaddr_in server_addr;
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");
    server_addr.sin_port = htons(12345);
    char message[1024];
    sprintf(message, "{\"key1\":\"value1\", \"key2\":\"value2\"}");
    sendto(sock, message, strlen(message), 0, (struct sockaddr*)&server_addr, sizeof(server_addr));
    char buffer[1024] = {0};
    recvfrom(sock, buffer, sizeof(buffer) - 1, 0, NULL, NULL);
    printf("Received: %s\n", buffer);
    close(sock);
    return 0;
}
```