import http from 'k6/http';
import { sleep } from 'k6';

export const options3 = {
  vus: 500,           // usuários virtuais simultâneos
  duration: '5s',   // duração total do teste
};

export const options = {
  vus: 500,           // usuários virtuais simultâneos
  stages: [
    { duration: '10s', target: 10 },   // aquecimento
    { duration: '20s', target: 10 },    // manter carga constante
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% das requisições < 500ms
    http_req_failed: ['rate<0.01'],    // < 1% de falhas
  },
};

export default function () {
  // const url = 'http://localhost:8080/transactions?user_id=62d54d96-88f4-4111-8564-c043d710bdcd&limit=10&offset=0';
  const url = 'http://localhost:8080/health';
  http.get(url);
  sleep(1); // aguarda 1 segundo entre as requisições por usuário
}