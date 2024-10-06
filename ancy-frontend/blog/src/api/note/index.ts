import request from '@/utils/request'
import type { NotePageParams, NotePageRes } from './type'

export function reqNotePage(params: NotePageParams): Promise<NotePageRes> {
  return request.get('/note/page', { params })
}
