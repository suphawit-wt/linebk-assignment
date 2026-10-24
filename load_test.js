import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 30 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<250'],
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  const payload = JSON.stringify({
    id: '000027b0e9a211ef95a30242ac180002',
    passcode: '123456',
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_URL}/v1/auth/signin/passcode`, payload, params);

  const success = check(res, {
    'setup signin status 200': (r) => r.status === 200,
  });

  if (!success) {
    console.error(`Setup authentication failed with status ${res.status}: ${res.body}`);
    return { token: null };
  }

  const token = res.json('data.accessToken');
  return { token };
}

export default function (data) {
  if (!data.token) {
    errorRate.add(1);
    console.error('No valid token available. Aborting VU execution.');
    return;
  }

  const authParams = {
    headers: {
      'Authorization': `Bearer ${data.token}`,
      'Content-Type': 'application/json',
    },
  };

  group('LINE BK Microservices', function () {
    const userRes = http.get(`${BASE_URL}/v1/user`, authParams);
    check(userRes, { 'get user info status 200': (r) => r.status === 200 }) || errorRate.add(1);

    const accRes = http.get(`${BASE_URL}/v1/accounts`, authParams);
    check(accRes, { 'get accounts status 200': (r) => r.status === 200 }) || errorRate.add(1);

    const txRes = http.get(`${BASE_URL}/v1/transactions/recent`, authParams);
    check(txRes, { 'get recent transactions status 200': (r) => r.status === 200 }) || errorRate.add(1);

    const cardRes = http.get(`${BASE_URL}/v1/debit-cards`, authParams);
    check(cardRes, { 'get debit cards status 200': (r) => r.status === 200 }) || errorRate.add(1);

    const bannerRes = http.get(`${BASE_URL}/v1/banners`, authParams);
    check(bannerRes, { 'get banners status 200': (r) => r.status === 200 }) || errorRate.add(1);
  });

  // Pacing: 1 second pause per VU iteration
  sleep(1);
}