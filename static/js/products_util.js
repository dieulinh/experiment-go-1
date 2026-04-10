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
    increateCartCount()
    // document.getElementById('cart-item-count').textContent = localStorage.getItem('cart_count') 
  }).catch(error => console.error('Error:', error))

}

function increateCartCount() {
  let count = localStorage.getItem('cart_count')||0
  count++

  localStorage.setItem('cart_count', count)
  document.querySelector('#cart-item-count').textContent = count
}
function setCartCount(cartCount){
  if(!localStorage.getItem('cart_count')) {
    localStorage.setItem('cart_count', 0)    
  }
  localStorage.setItem('cart_count', cartCount)
}