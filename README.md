# iot-backend

## API documentation

### POST /data
An endpoint to insert data into InfluxDB.

This endpoint expects the following **JSON** request body:
```
{
    "measurement": "",
    "location": "",
    "temp": 0.0,
    "rh": 0.0,
    "radiance": 0,
    "timestamp": ""
}
```

### GET /data
An endpoint to retrieve data from InfluxDB.
