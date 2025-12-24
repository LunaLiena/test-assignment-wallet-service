// loadtest/wallet-stress.js
import http from "k6/http";
import { check, sleep, fail } from "k6";

export const options = {
  vus: 30, // 30 виртуальных пользователей
  duration: "45s", // 45 секунд нагрузки
  thresholds: {
    http_req_duration: ["p(95) < 1000"], // 95% < 1 сек
    http_req_failed: ["rate < 0.02"], // < 2% ошибок
  },
};

const BASE_URL = "http://app:8080/api/v1";
const WALLET_ID = "550e8400-e29b-41d4-a716-446655440000";

export function setup() {
  console.log("📝 Initializing wallet...");
  const res = http.post(
    `${BASE_URL}/wallets/init`,
    JSON.stringify({ walletId: WALLET_ID }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );
  if (res.status !== 201) {
    fail(`Failed to init wallet: ${res.status} ${res.body}`);
  }
  console.log("✅ Wallet ready");
}

export default function () {
  // Чередуем операции
  const op = Math.random() > 0.5 ? "DEPOSIT" : "WITHDRAW";
  const amount = op === "DEPOSIT" ? 10 : 5;

  const payload = JSON.stringify({
    walletId: WALLET_ID,
    operationType: op,
    amount: amount,
  });

  const res = http.post(`${BASE_URL}/wallet`, payload, {
    headers: { "Content-Type": "application/json" },
  });

  check(res, {
    "status is 200": (r) => r.status === 200,
  });

  // Иногда читаем баланс
  if (Math.random() < 0.15) {
    const balanceRes = http.get(`${BASE_URL}/wallets/${WALLET_ID}`);
    check(balanceRes, { "balance ok": (r) => r.status === 200 });
  }

  sleep(0.3);
}
