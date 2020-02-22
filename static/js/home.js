var targets = ['members-range', 'first-album-range', 'creation-range']

var artistListEl = document.getElementById('artist-list')

let [mRange, faRange, cRange] = createRanges()

var filterForm = document.getElementById('filter-form')

addListeners()

function addListeners() {
  targets.forEach(target => {
    document.getElementById(target).noUiSlider.on('update', function(values) {
      var min = document.querySelector('.min[data-target=' + target + ']')
      var max = document.querySelector('.max[data-target=' + target + ']')
      min.innerText = values[0]
      max.innerText = values[1]
    })
  })
}

filterForm.addEventListener('submit', function(e) {
  e.preventDefault()
  var data = {
    countries: [],
    creationRange: { min: 0, max: 0 },
    membersRange: { min: 0, max: 0 },
    firstAlbumsRange: { min: 0, max: 0 }
  }

  var inputs = document.querySelectorAll('#filter-form input:checked')
  inputs.forEach(input =>
    data.countries.push(input.closest('label').textContent.trim())
  )

  let min, max
  ;[min, max] = mRange.get()
  data.membersRange.min = Number(min)
  data.membersRange.max = Number(max)
  ;[min, max] = cRange.get()
  data.creationRange.min = Number(min)
  data.creationRange.max = Number(max)
  ;[min, max] = faRange.get()
  data.firstAlbumsRange.min = Number(min)
  data.firstAlbumsRange.max = Number(max)

  fetch('/filter', {
    headers: {
      'Content-Type': 'application/json'
    },
    method: 'POST',
    body: JSON.stringify(data)
  })
    .then(res => res.text())
    .then(res => {
      artistListEl.innerHTML = res
      //   console.log(res)
    })
})

function createRanges() {
  var res = []
  targets.forEach(target => {
    var membersRangeEl = document.getElementById(target)
    var min = Number(membersRangeEl.getAttribute('data-min'))
    var max = Number(membersRangeEl.getAttribute('data-max'))

    res.push(
      noUiSlider.create(membersRangeEl, {
        start: [min, max],
        connect: true,
        step: 1,
        orientation: 'horizontal', // 'horizontal' or 'vertical'
        range: {
          min: min,
          max: max
        },
        format: wNumb({
          decimals: 0
        })
      })
    )
  })

  return res
}

document
  .querySelector('.reset.locations.remove')
  .addEventListener('click', function(e) {
    document.querySelectorAll('.loc-input').forEach(el => (el.checked = false))
  })

document
  .querySelector('.reset.locations.add')
  .addEventListener('click', function(e) {
    document.querySelectorAll('.loc-input').forEach(el => (el.checked = true))
  })

document
  .querySelector('.reset.creation')
  .addEventListener('click', function(e) {
    cRange.reset()
  })

document.querySelector('.reset.members').addEventListener('click', function(e) {
  mRange.reset()
})

document.querySelector('.reset.fa').addEventListener('click', function(e) {
  faRange.reset()
})
