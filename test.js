const http = require('http')

const key = 'AIzaSyBUUmTPz_PJpUW_hfJggbIwzbthAs1Sjk0'
let address = 'Sao+Paulo,+Brazil'

http
  .get(
    `maps.googleapis.com/maps/api/geocode/json?address=${address}&key=${key}`
  )
  .then(res => console.log(res))
  .catch(err => console.log(err))
