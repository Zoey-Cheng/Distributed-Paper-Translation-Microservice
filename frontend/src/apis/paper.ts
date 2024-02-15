import {baseURL, del, get, post} from "./http";

interface Paper {
  paperID: string;
  status: 0 | 1 | 2 | 3;
  createAt: number;
  resultText: string
  fileHash: string
}

export const queryPapers = () => {
  return get<Paper[]>("/papers/")
}

export const queryPaper = (id: string) => {
  return get<Paper>(`/papers/${id}`)
}

export const createPapers = (fileHash: string, targetLanguage: string, emailTo: string) => {
  return post<{ paperID: string }>("/papers/", {
    fileHash,
    emailTo,
    targetLanguage
  })
}

export const deletePaper = (id: string) => {
  return del<any>(`/papers/${id}`)
}

export const getPaperDownloadURL = (id: string) => {
  return baseURL + `/papers/${id}/download_txt`
}