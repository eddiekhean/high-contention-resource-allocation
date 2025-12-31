# High-Contention Resource Allocation System (HCRAS)

## 1. Project Description

HCRAS is a scheduling and allocation system designed for scarce resources in environments where request volume significantly exceeds capacity. Practical applications include promotional voucher distribution, limited event slot booking, software license management, and real-time GPU/CPU time allocation.

The system focuses on resolving four core challenges:

* **Priority:** Prioritizing critical requests based on client identity and status.
* **Fairness:** Ensuring equitable distribution and preventing "starvation" for low-priority requests.
* **Burst Load:** Handling thousands of concurrent requests within millisecond windows.
* **Observability:** Providing real-time tracking of latency metrics and allocation rates.

## 2. Research Foundation

The system architecture is based on principles from the research paper:
**"Priority-based Fair Scheduling in Edge Computing"** (Source: arXiv:2001.09070).

### Problem Statement

In Edge Computing environments or limited resource distribution systems, traditional strategies such as First-Come-First-Serve (FCFS) or Strict Priority often lead to suboptimal performance, high latency for critical tasks, or total starvation of non-priority users.

### Hybrid Scheduling Solution

This system implements a Hybrid strategy that balances Priority and Fairness. Experimental data from the referenced research indicates that this approach yields the highest efficiency under high-load conditions and non-uniform request distributions.

## 3. Architecture and Algorithms

### 3.1. Mapping Model

| Research Component | Practical Mapping | Description |
| --- | --- | --- |
| Edge Resource | Voucher Slots / Finite Resource | The limited resource to be allocated. |
| Job / Task | Resource Request | An individual request from a user or service. |
| Client | User / Service Identity | The unique identifier used to calculate fairness metrics. |
| Hybrid Strategy | Dynamic Scoring Scheduler | The scheduling engine that calculates real-time priority. |

### 3.2. Scoring Algorithm

The system utilizes a Dynamic Scoring mechanism to rank requests within the queue:

Where:

* : Base priority (e.g., VIP = 100, Member = 50, Guest = 10).
* : The duration the request has been pending in the queue.
* : Fairness adjustment factor to balance resource consumption across different clients.
* : Hyper-parameters used to tune the system according to specific business requirements.

## 4. Technical Features

* **Intermediate Queueing:** Utilizes a Message Broker to decouple traffic from the database and absorb burst loads effectively.
* **Starvation Mitigation:** By incorporating the  parameter, low-priority requests gain score over time, ensuring they are eventually processed.
* **Optimized Data Structures:** Implements Priority Queues or Redis Sorted Sets to achieve  complexity when extracting the highest-scoring requests.
* **Resilience and Recovery:** Maintains state persistence to allow for seamless recovery in the event of a system failure.

## 5. Technology Stack

* **Languages:** Go / Java / Node.js.
* **Data Infrastructure:** Redis (Sorted Sets), PostgreSQL.
* **Monitoring:** Prometheus, Grafana (tracking starvation rates, throughput, and latency).

## 6. Expected Outcomes

* **Optimization of conversion rates** by ensuring target customer segments are prioritized.
* **100% processing of valid requests** within committed Service Level Agreements (SLA).
* **System stability** during high-traffic events such as flash sales or major product launches, preventing cascading failures.

---
