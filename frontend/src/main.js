const calcSize = () => {
  const input = htmx.findAll('input')[0]
  const inputHeight = input.getBoundingClientRect().height
  const option = htmx.findAll('option')[0]
  const optionHeight = option.getBoundingClientRect().height
  const listSize = (window.innerHeight - inputHeight) / optionHeight
  htmx.findAll('select')[0].size = listSize
}
window.addEventListener('pageshow', calcSize)
window.addEventListener('resize', calcSize)

document.addEventListener('keyup', (e) => {
  if (e.key === 'Escape') {
    const input = document.getElementsByName('filter').item(0)
    input.focus()
    input.value = ''
  }
})
