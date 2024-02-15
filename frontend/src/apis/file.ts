import {get, Instance, post} from "./http";

export interface Config {
  configID: string
  configType: string
  version: string
  configData: any
}

export const startUploadFile = (hash: string, fileName: string, chunkNums: number, segmentSize: number) => {
  return post<Config>(`/files/start`, {
    hash,
    fileName,
    chunkNums,
    segmentSize
  })
}

export const uploadChunk = (chunk: Blob, hash: string, chunkIndex: number) => {
  const formData = new FormData()
  formData.set("chunk", chunk)
  formData.set("hash", hash)
  formData.set("chunkIndex", chunkIndex + "")
  return Instance.postForm("/files/chunk", formData)
}

interface FileInfo {
  hash: string;
  status: 0 | 1 | 2
  current_index: number
}

export const queryFile = (hash: string) => {
  return get<FileInfo>(`/files/${hash}`)
}

export const queryFileURL = (hash: string) => {
  return get<{ url: string }>(`/files/${hash}/public_url`)
}