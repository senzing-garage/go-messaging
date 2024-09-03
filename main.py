#! /usr/bin/env python3

from typing import Any, Dict

import iso8601

from python.typedef import SenzingMessage

example_message_dict: Dict[Any, Any] = {
    "time": "2023-04-10T11:00:20.623748617-04:00",
    "level": "TRACE",
    "id": "senzing-99990002",
    "text": "A fake error",
    "duration": 199045,
    "location": "In main() at main.go:36",
    "errors": ["0027E|Unknown DATA_SOURCE value 'DOESNTEXIST'"],
    "details": [
        {"position": 1, "type": "string", "value": "DoesntExist"},
        {"position": 2, "type": "string", "value": "1070", "valueRaw": 1070},
        {"position": 3, "type": "int64", "value": "-1", "valueRaw": -1},
        {
            "position": 4,
            "type": "szengine._Ctype_longlong",
            "value": "-2",
            "valueRaw": -2,
        },
        {
            "position": 5,
            "type": "error",
            "value": "0027E|Unknown DATA_SOURCE value 'DOESNTEXIST'",
        },
    ],
}

print("-- New style --------")
example_message = SenzingMessage.from_json_data(example_message_dict)
print(f"Year: {example_message.time.year}")
print(f"Text: {example_message.text}")
print(example_message.details.value[0].value)
print(example_message.details.value[4].value_raw)

# "Old" style.  Issues:
# - IDEs cannot give meaningful hints.
# - Static analysis cannot be done.
# - Timestamp not returned as datetime object.

print("-- Old style --------")
print(iso8601.parse_date(example_message_dict.get("time", "")).year)
print(example_message_dict.get("details", [])[0].get("value"))
