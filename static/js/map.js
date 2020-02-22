function initMap() {
  // init empty map
  var mapEl = document.getElementById('map')
  mapEl.style.height = '500px'
  var map = new google.maps.Map(mapEl, {
    // zoom: 3,
    center: { lng: 0, lat: 0 },
    maxZoom: 5
  })

  var bounds = new google.maps.LatLngBounds()

  // set markers
  var locElements = document.querySelectorAll('.location')

  var locNames = []
  locElements.forEach(el => {
    locNames.push(el.textContent)
  })
  locNames = [...new Set(locNames)]

  var locs = []
  var count = 0
  locNames.forEach(name => {
    fetch(
      `https://maps.googleapis.com/maps/api/geocode/json?address=${name}&key=AIzaSyDju6U6HXWpd9ecaE6y8y_OLEF7v-onlUU`
    )
      .then(res => res.json())
      .then(res => {
        if (res.status === 'OK') {
          locs.push(res.results[0].geometry.location)
        } else {
          console.log(res)
        }
      })
      .catch(err => console.log(err))
      .finally(() => {
        count++
        if (count === locNames.length) {
          locs.forEach((loc, i) => {
            var marker = new google.maps.Marker({
              position: loc,
              map: map,
              title: locNames[i],
              animation: google.maps.Animation.DROP
            })

            bounds.extend(marker.getPosition())
          })

          map.fitBounds(bounds)
        }
      })
  })
}
