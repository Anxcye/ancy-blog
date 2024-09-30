import request from '@/utils/request'
import type { UploadRes } from './type'

export function reqUpload(img: File): Promise<UploadRes> {
  const formData = new FormData()
  formData.append('file', img)
  return request.post('/upload', formData)
}
