export const handleScroll = async (
  callback: () => void,
  loading: boolean,
  end: boolean,
  threshold = 200,
) => {
  // const container = document.querySelector('#app')
  let container = document.documentElement

  if (!container) return
  let { scrollTop, clientHeight, scrollHeight } = container
  while (scrollTop + clientHeight > scrollHeight - threshold && !loading && !end) {
    await new Promise((resolve) => setTimeout(resolve, 100))
    callback()
    container = document.documentElement
    scrollTop = container.scrollTop
    clientHeight = container.clientHeight
    scrollHeight = container.scrollHeight
  }
}
