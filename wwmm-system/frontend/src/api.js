import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 15000
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('wwmm_token')
  if (token) {
    config.headers['Token'] = token
    config.headers['Authorization'] = 'Bearer ' + token
  }
  return config
})

api.interceptors.response.use(
  resp => {
    const data = resp.data
    if (data && data.code !== 0 && data.code !== undefined) {
      return Promise.reject(new Error(data.message || '请求失败'))
    }
    return data
  },
  err => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('wwmm_user')
      localStorage.removeItem('wwmm_token')
    }
    return Promise.reject(err)
  }
)

export default api
