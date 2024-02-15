import axios from "axios";

export const baseURL = 'http://localhost:80/v1'

export const Instance = axios.create({
  baseURL: baseURL
})

const Http = <D>(method: string, url: string, params?: any, headers?: any) => {
  return new Promise((resolve: (arg0: D) => void, reject: (arg0: any) => void) => {
    Instance.request<D>({
      url: url,
      method: method,
      headers: headers,
      data: params,
    }).then((res) => {
      if (res.status < 300) {
        resolve(res.data)
      } else {
        reject(res.data)
      }
    }).catch((err) => {
      reject(err)
    })
  })
}

export const get = <D>(url: string, headers?: any) => {
  return Http<D>('get', url, {}, headers)
}

export const post = <D>(url: string, params?: any, headers?: any) => {
  return Http<D>('post', url, params, headers)
}

export const put = <D>(url: string, params?: any, headers?: any) => {
  return Http<D>('put', url, params, headers)
}

export const del = <D>(url: string, headers?: any) => {
  return Http<D>('delete', url, {}, headers)
}


