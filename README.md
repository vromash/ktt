# Financing Application Aggregator Service

## Quick Start

The easiest way to run the service is with Docker Compose:

```sh
docker-compose up -d
```

To stop and remove everything:

```sh
docker-compose down -v
```

This will start the backend service and a PostgreSQL database. The backend will be available at `http://localhost:6666` by default.

---

## Application Logic

1. **Submit Application:**
   - The client submits a financing application via the HTTP API.
2. **Background Bank Requests:**
   - The service immediately sends requests to all available banks in the background.
3. **Offer Status Updates (Cron):**
   - Every 30 seconds, a cron job checks for updates on all offers with `DRAFT` status by polling the banks.
   - If an offer status changes, the update is sent to all connected WebSocket clients in real time.
4. **Data Access:**
   - If you are not connected via WebSocket, you can still fetch the latest application data and all available offers using the HTTP API.

---

## API

### Authentication
All endpoints (except `/healthz`) require an `Authorization` header with a Bearer token. For simplicity, it requires any non-empty value.

```
Authorization: Bearer <any-non-empty-value>
```

Requests without this header will receive a 401 Unauthorized response.

### Endpoints

To see the list of endpoints, please refer to `docs/swagger.yaml`. It's an auto-generated file based on annotations.

### Updates via WebSocket
- `GET /ws/applications/{id}`
  - Upgrade to a WebSocket connection to receive real-time updates for offers on a specific application.
  - **Headers:**
    - `Authorization: Bearer <token>`
  - **Usage:** Connect and listen for JSON messages with offer updates as soon as they are available.

## Improvements & Further Development

Here are some thoughts and ideas for how this service could be improved or extended in the future:

1. **WebSocket Connection Model:**
   - For the sake of simplicity, I've set up the WebSocket to connect via application ID. But that's not very convenient for real-world use. It would be much better to introduce a User or Session model, link it with all applications, and use that for WebSocket connections. That way, a user could subscribe to all their applications at once, not just one by one.

2. **User Registration & Authorization:**
   - Continuing on the User model, it would be great to have a registration/authorization system. This would let users save their email and phone once, and then submit applications with just the financial data. It would also make the Authorization middleware meaningful, since we'd be able to generate and validate tokens for each user.

3. **Caching:**
   - If the service needs to scale, adding a cache layer (like Redis) could help reduce the number of read requests to the database and speed up reads for applications, offers, and users.

4. **Offer Delivery:**
   - If WebSockets aren't the preferred way for the frontend to get offer updates, periodic polling is always an option. For API users, it would also be possible to add webhook support, so they can get notified as soon as something changes.

5. **Bank Data Sync:**
   - Right now, I'm using a cron job to fetch the latest data from banks every 30 seconds. This value can be tweaked to better fit the banks' response times and rate limits. The best solution would be to use webhooks from the banks themselves, so we get notified instantly when something changes—if the banks support that, of course.

6. **Request Model Simplification:**
   - It's probably possible to combine or simplify some of the application request fields, but I wasn't sure about the business meaning of each, so I kept all the fields I found in the banks' submit requests. I'd want to clarify this with a PM or PO before releasing the service.

7. **Tracing & Metrics:**
    - The application should be covered with traces for all important functions and business flows, as well as key metrics (e.g., how many cron checks were needed to get an offer from a bank). For this, I would use the OTEL Go package. For fast metrics and traces, Gin-specific tracing middleware and a DB tracing package can also be used.

8. **Test Coverage:**
    - I've covered the essential parts of the application with tests, but there's definitely room to increase coverage (cover HTTP handlers with tests, checking data validations, cover mappers, cover repositories, possibly adding custom mocks for banks).

9. **Health Checks:**
    - The `/healthz` endpoint is just a placeholder for now. It would be good to extend it to check the status of the database, bank integrations, and other dependencies.
