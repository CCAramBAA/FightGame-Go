import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/admin',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

apiClient.interceptors.response.use(
  (res) => res.data,
  (err) => {
    console.error('Admin API Error:', err)
    return Promise.reject(err)
  }
)

export default apiClient
