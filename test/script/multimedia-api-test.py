# Copyright (C) 2025 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

import requests

# 基础 URL，根据实际情况修改
base_url = "http://localhost:2580/api/v1"


def test_create_multimedia():
    url = f"{base_url}/multimedia/create"
    data = {
        "uuid": "",
        "type": "RTSP",
        "name": "Test Multimedia",
        "description": "This is a test multimedia",
        "config": {
            "streamUrl": "rtsp://example.com/stream",
            "enablePush": True,
            "pushUrl": "rtsp://example.com/push",
            "enableAi": True,
            "aiModel": "YOLOV8",
        },
    }
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, json=data, headers=headers)
    print(f"Create MultiMedia Response: {response.status_code}, {response.json()}")


def test_update_multimedia():
    url = f"{base_url}/multimedia/update"
    data = {
        "uuid": "12345678-1234-5678-1234-567812345678",
        "type": "RTMP",
        "name": "Updated Multimedia",
        "description": "This is an updated test multimedia",
        "config": {
            "streamUrl": "rtmp://example.com/stream",
            "enablePush": False,
            "pushUrl": "rtmp://example.com/push",
            "enableAi": False,
            "aiModel": "FACENET",
        },
    }
    headers = {"Content-Type": "application/json"}
    response = requests.put(url, json=data, headers=headers)
    print(f"Update MultiMedia Response: {response.status_code}, {response.json()}")


def test_multimedia_detail():
    uuid = "12345678-1234-5678-1234-567812345678"
    url = f"{base_url}/multimedia/detail?uuid={uuid}"
    response = requests.get(url)
    print(f"MultiMedia Detail Response: {response.status_code}, {response.json()}")


def test_list_multimedia():
    url = f"{base_url}/multimedia/list"
    response = requests.get(url)
    print(f"List MultiMedia Response: {response.status_code}, {response.json()}")


def test_delete_multimedia():
    uuid = "12345678-1234-5678-1234-567812345678"
    url = f"{base_url}/multimedia/del?uuid={uuid}"
    response = requests.delete(url)
    print(f"Delete MultiMedia Response: {response.status_code}, {response.json()}")


if __name__ == "__main__":
    test_create_multimedia()
    test_update_multimedia()
    test_multimedia_detail()
    test_list_multimedia()
    test_delete_multimedia()
