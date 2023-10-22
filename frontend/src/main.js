const calcSize = () => {
  const input = htmx.findAll('input')[0]
  const inputHeight = input.getBoundingClientRect().height
  const option = htmx.findAll('option')[0]
  const optionHeight = option.getBoundingClientRect().height
  const listSize = (window.innerHeight - inputHeight) / optionHeight
  htmx.findAll('select')[0].size = listSize
}
window.addEventListener('resize', calcSize)
document.body.addEventListener('htmx:afterOnLoad', calcSize)
