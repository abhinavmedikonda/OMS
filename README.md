http://localhost:8000/graphql
http://localhost:8000/playground

ACCOUNT:

query q {
  q1: accounts {
    id
    name
  }

  q2: accounts(id: "35a4KUgUPgx1wC1zX2pPsLHXGyt"){
    name
    orders{
      totalPrice
    }
  }

	q3: accounts(id: "35a4KUgUPgx1wC1zX2pPsLHXGyt"){
    id
    name
    orders{
      id
      createdAt
      totalPrice
      products{
        id
        name
        description
        price
        quantity
      }
    }
  }
}

mutation m {
  createAccount(account: {name: "abhiacc"}){
    id
    name
  }
}



CATALOG:

query q {
  q1: products {
    id
    name
    description
    price
  }

  q2: products(pagination: {skip: 0, take: 5}, query: "ab"){
    id
    name
    description
    price
  }
}


mutation m {
  createProduct(product: {name: "abcp rod", description: "abc desc", price: 35}){
    id
    name
    price
    description
  }
}



ORDER:

mutation{
  createOrder(order: {accountId: "35a4KUgUPgx1wC1zX2pPsLHXGyt",
    products: [{id: "35a51vlRTojJnNPLQQi54b5tX0B", quantity: 2}]}){
    id
    totalPrice
    products{
      name
      quantity
    }
  }
}