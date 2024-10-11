# Copyright (C) 2024 wwhai
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
import time

# List of schema IDs
schema_ids = [
    "SCHEMAETNVHNTE",
    "SCHEMAZI8HSTFO",
    "SCHEMAKV3LGTYE",
    "SCHEMAXVJRXTZE",
    "SCHEMAQP28GTEV",
    "SCHEMA6YTDTFQS",
    "SCHEMAHNHXRGWY",
    "SCHEMA4AKSLNEY",
    "SCHEMA4QFSCZYC",
    "SCHEMAOQP6A6LP",
    "SCHEMAZJMARNBX",
    "SCHEMA6DC2G7UP",
]

# List of template IDs
template_ids = [
    "TEMP_HUMIDITY",
    "SWITCH_STATUS",
    "WATER_QUALITY",
    "AIR_QUALITY",
    "MOTION_SENSOR",
    "SMART_METER",
    "SOIL_MOISTURE",
    "GPS_TRACKER",
    "SMOKE_DETECTOR",
    "SMART_LOCK",
    "SIX_AXIS_ACCELEROMETER",
]

# Headers for the request
headers = {"User-Agent": "Apifox/1.0.0 (https://apifox.com)"}

# Iterate over both lists
for schema_id in schema_ids:
    for template_id in template_ids:
        # Construct the URL using schema_id and template_id
        url = f"http://127.0.0.1:2580/api/v1/schema/genTemplate?schemaId={schema_id}&templateId={template_id}"

        # Sending the POST request
        response = requests.post(url, headers=headers, data={})

        # Printing the response for each request
        print(f"Response for schemaId {schema_id} and templateId {template_id}:")
        print(response.text)
        print("-" * 50)
        time.sleep(1)
