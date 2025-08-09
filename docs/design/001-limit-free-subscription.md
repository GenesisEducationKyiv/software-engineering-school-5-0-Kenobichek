# Design-001: Enforcing Free Subscription Limits

## 1. Context

Problem: Users can currently create multiple FREE subscriptions, increasing our operational costs.

Goal: Allow **only one free subscription per user**, while enabling upgrade to a paid plan.

---

## 2. Solution

### a. Database Enforcement

```sql
CREATE UNIQUE INDEX idx_user_free_sub 
  ON subscriptions(user_id) 
  WHERE plan_type = 'FREE';
```

This ensures at most one free subscription per user in the DB layer.

---

### b. API Endpoint: Upgrade to Paid

**POST** `/subscriptions/upgrade`

#### Request:
```json
{
  "user_id": 123,
  "from_plan": "FREE",
  "to_plan": "PAID"
}
```

#### Response (Success):
```json
{
  "status": "upgraded",
  "subscription_id": 456
}
```

#### Response (Errors):
- `409 Conflict`: Already has an active PAID subscription.
- `400 Bad Request`: Invalid plan type.
- `404 Not Found`: No FREE subscription found.

---

### c. Migration Plan

Batch cleanup job to retain only the oldest FREE subscription per user:

```sql
DELETE FROM subscriptions
WHERE user_id = ?
  AND plan_type = 'FREE'
  AND id != (
    SELECT id FROM subscriptions
    WHERE user_id = ? AND plan_type = 'FREE'
    ORDER BY created_at ASC LIMIT 1
  );
```

- Use batching (e.g., `LIMIT 1000`) and throttling to avoid DB contention.

---

## 3. Validation & Monitoring

- **Testing**: Unit + integration tests including concurrent signup/upgrade simulation.
- **Metric**: `subscription_duplicate_attempt_total` to monitor rejected duplicate attempts.
---

## 4. Operations & Safety

PagerDuty alert: If >5 duplicate attempts per minute over a 10-minute window.

---

## 5. Reviewers

Backend team

---

## 6. Deadline

No deadline specified