import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate } from 'k6/metrics';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.1.0/index.js";

const errorRate = new Rate('errors');

export const options = {
  // 4 minute 30 second optimized profile to clear 1M requests safely under strict latency
  stages: [
    { duration: '30s', target: 500 },    // Warm-up phase
    { duration: '30s', target: 1200 },   // Rapid ramp-up
    { duration: '3m', target: 1350 },    // Stable sustained load for 3 minutes (sweet spot for ~3,800 req/sec)
    { duration: '30s', target: 0 },      // Graceful ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // Strict p95 latency threshold under 500ms
    errors: ['rate<0.01'],           // Error rate strictly under 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  const payload = JSON.stringify({
    id: '000027b0e9a211ef95a30242ac180002',
    passcode: '123456',
  });

  const res = http.post(`${BASE_URL}/v1/auth/signin/passcode`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (res.status !== 200) {
    console.error(`Setup signin failed: ${res.status}`);
    return { token: null };
  }

  return { token: res.json('data.accessToken') };
}

export default function (data) {
  if (!data.token) {
    errorRate.add(1);
    return;
  }

  const authParams = {
    headers: {
      'Authorization': `Bearer ${data.token}`,
      'Content-Type': 'application/json',
    },
  };

  group('Stress Test - 1M Optimized Workload', function () {
    const accountsRes = http.get(`${BASE_URL}/v1/accounts`, authParams);
    check(accountsRes, { 'status is 200': (r) => r.status === 200 }) || errorRate.add(1);

    const txRes = http.get(`${BASE_URL}/v1/transactions/recent`, authParams);
    check(txRes, { 'status is 200': (r) => r.status === 200 }) || errorRate.add(1);

    const cardRes = http.get(`${BASE_URL}/v1/debit-cards`, authParams);
    check(cardRes, { 'status is 200': (r) => r.status === 200 }) || errorRate.add(1);
  });

  // Balanced sleep pacing (15ms) to maintain high RPS without socket exhaustion
  sleep(0.015);
}

export function handleSummary(data) {
  return {
    "stdout": textSummary(data, { indent: " ", enableColors: true }),
    "stress-test-report.html": htmlReport(data),
  };
}