import http from 'k6/http';
import { check } from 'k6';

// Konfigurasi test scenarios
export const options = {
    stages: [
        { duration: '10s', target: 50 },   // Ramp up to 50 users over 10 seconds
        { duration: '50s', target: 100 },  // Ramp up to 100 users over 50 seconds
        { duration: '10s', target: 10 },   // Ramp down to 10 users over 10 seconds
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% of requests must be below 200ms
        http_req_failed: ['rate<0.05'],   // Error rate must be below 5%
    },
};

// Test function yang akan dijalankan oleh setiap virtual user
export default function () {
    const url = 'http://localhost:1312/api/auth/login';

    const payload = JSON.stringify({
        email: 'admin@gmail.com',
        password: 'admin123',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Kirim POST request
    const response = http.post(url, payload, params);

    // Validasi response
    check(response, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'response has token': (r) => {
            try {
                const json = JSON.parse(r.body);
                return json.token !== undefined || json.access_token !== undefined;
            } catch (e) {
                return false;
            }
        },
    });

    // Optional: sleep untuk simulasi user behavior yang lebih realistis
    // sleep(1);
}