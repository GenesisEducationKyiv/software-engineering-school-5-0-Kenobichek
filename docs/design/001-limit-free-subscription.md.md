# System Design Doc (Short Form)

## 1. Context
* **Problem:** Users can create multiple free subscriptions, driving costs beyond budget.
* **Goal:** Allow at most **one** free subscription per user and add an upgrade path to a paid tier.
---

## 2. Solution
1. **DB Layer** – Add unique index `(user_id, plan_type='FREE')`.
2. **API** – New endpoint `POST /subscriptions/upgrade` converts free ➜ paid.
3. **Migration** – Background job keeps the oldest free sub, deletes the rest.
---

## 3. Validation

* Unit + integration tests for duplicate prevention.
* Grafana metric `subscription_duplicate_attempt_total` < **0.1 %**.
---


## 4. Ops & Safety
* PagerDuty alert if duplicates > 5/min (10 min window).
* Feature flag `unique_subscription_guard`; rollback by disabling flag and dropping index.
---


## 5. Reviewers
* Backend Team
---


## 6. Deadline
deadline is not defined