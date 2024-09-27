document.addEventListener('keyup', (e) => {
  if (e.key === 'Escape') {
    const input = document.getElementsByName('filter').item(0)
    if (document.activeElement.nodeName === 'INPUT') {
      input.value = ''
    }
    input.focus()
  }
})
