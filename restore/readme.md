sequenceDiagram
    participant A as Client
    participant B as WayBackServer
    A->>B: GET: All Urls
    B->>A: Return: All Urls
    loop Return Url
        B-->A: Return Url!
    end