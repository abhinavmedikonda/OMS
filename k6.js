import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "20s", target: 30 },
    { duration: "40s", target: 100 },
    { duration: "20s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<1000"],
    http_req_failed: ["rate<0.05"],
  },
};

const URL = "http://oms-sbx.local/graphql";
const headers = { "Content-Type": "application/json" };

// --------------------
// GLOBAL ID STORAGE (per VU copy)
// --------------------
let accountIds = [];
let productIds = [];

// --------------------
// HELPER
// --------------------
function gql(query) {
  return http.post(URL, JSON.stringify({ query }), { headers });
}

// --------------------
// SETUP (runs once)
// --------------------
export function setup() {
  let accRes = gql(`query { accounts(pagination: { take: 50 }) { id } }`);
  let prodRes = gql(`query { products(pagination: { take: 50 }) { id } }`);

  let accounts =
    JSON.parse(accRes.body).data?.accounts?.map(a => a.id) || [];

  let products =
    JSON.parse(prodRes.body).data?.products?.map(p => p.id) || [];

  // Seed if empty
  if (accounts.length === 0) {
    for (let i = 0; i < 200; i++) {
      let r = gql(`
        mutation {
          createAccount(account: { name: "seed-${Math.random()}" }) {
            id
          }
        }
      `);
      let body = JSON.parse(r.body);
      if (body.errors) {
        console.error("createAccount error:", body.errors);
        continue;
      }
      let id = body.data?.createAccount?.id;
      if (id) accounts.push(id);
    }
  }

  if (products.length === 0) {
    for (let i = 0; i < 500; i++) {
      let r = gql(`
        mutation {
          createProduct(product: {
            name: "seed-prod-${Math.random()}"
            description: "seed"
            price: ${(Math.random() * 100).toFixed(2)}
          }) {
            id
          }
        }
      `);
      let body = JSON.parse(r.body);
      if (body.errors) {
        console.error("createProduct error:", body.errors);
        continue;
      }
      let id = body.data?.createProduct?.id;
      if (id) products.push(id);
    }
  }

  return { accounts, products };
}

// --------------------
// MAIN
// --------------------
export default function (data) {

  // Initialize per-VU ID lists once
  if (accountIds.length === 0) {
    accountIds = [...data.accounts];
  }

  if (productIds.length === 0) {
    productIds = [...data.products];
  }

  const rand = Math.random();
  let res;

  // 30% Accounts read
  if (rand < 0.3) {
    res = gql(`
      query {
        accounts(pagination: { take: 5 }) {
          id
          name
          orders {
            id
            totalPrice
          }
        }
      }
    `);
  }

  // 30% Product read
  else if (rand < 0.6) {
    res = gql(`
      query {
        products(pagination: { take: 5 }) {
          id
          name
          price
        }
      }
    `);
  }

  // 15% Create Account
  else if (rand < 0.75) {
    res = gql(`
      mutation {
        createAccount(account: { name: "user-${Math.random()}" }) {
          id
        }
      }
    `);

    const id = JSON.parse(res.body).data?.createAccount?.id;
    if (id) {
      accountIds.push(id);
    }
  }

  // 15% Create Product
  else if (rand < 0.9) {
    res = gql(`
      mutation {
        createProduct(product: {
          name: "prod-${Math.random()}"
          description: "load"
          price: ${(Math.random() * 100).toFixed(2)}
        }) {
          id
        }
      }
    `);

    const id = JSON.parse(res.body).data?.createProduct?.id;
    if (id) {
      productIds.push(id);
    }
  }

  // 10% Create Order (USES REAL IDS)
  else {
    const accountId =
      accountIds[Math.floor(Math.random() * accountIds.length)];

    const productId =
      productIds[Math.floor(Math.random() * productIds.length)];

    res = gql(`
      mutation {
        createOrder(input: {
          accountId: "${accountId}",
          products: [{ id: "${productId}", quantity: 1 }]
        }) {
          id
          totalPrice
        }
      }
    `);
  }

  check(res, {
    "status 200": (r) => r.status === 200,
  });

  sleep(0.3);
}