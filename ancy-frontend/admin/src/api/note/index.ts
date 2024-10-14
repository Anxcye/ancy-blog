import request from '@/utils/request'
import type {
  NoteAddParams,
  NoteAddRes,
  NoteDeleteRes,
  NoteGetByIdRes,
  NotePageParams,
  NotePageRes,
  NoteUpdateRes,
} from './type'

export function reqNotePage(params: NotePageParams): Promise<NotePageRes> {
  return request.get(`/note/page`, { params })
}

export function reqNoteAdd(params: NoteAddParams): Promise<NoteAddRes> {
  return request.post(`/note`, params)
}

export function reqNoteGetById(id: number): Promise<NoteGetByIdRes> {
  return request.get(`/note/${id}`)
}

export function reqNoteUpdate(id: number, params: NoteAddParams): Promise<NoteUpdateRes> {
  return request.put(`/note/${id}`, params)
}

export function reqNoteDelete(id: number): Promise<NoteDeleteRes> {
  return request.delete(`/note/${id}`)
}
