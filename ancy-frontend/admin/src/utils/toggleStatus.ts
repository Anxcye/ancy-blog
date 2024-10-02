export const toggleStatus = (row: any) => {
  if (row.status === '1') {
    row.status = '0'
  } else {
    row.status = '1'
  }
}

export const toggleTop = (row: any) => {
  if (row.isTop === '1') {
    row.isTop = '0'
  } else {
    row.isTop = '1'
  }
}
