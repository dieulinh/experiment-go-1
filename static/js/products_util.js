function addToCart(product_id){
  console.log(product_id)
  fetch('/api/cart/items',{
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    product_id: product_id,
    quantity: 1
  })
  
}).then(resp => resp.json()).then(rs=>{
    console.log('add to cart')
  }).catch(error => console.error('Error:', error))

}