/**
 * Axios 实例：统一 baseURL、请求/响应拦截、全局异常处理
 * - 默认使用相对路径 /api/v1，与页面同源，便于前后端一体部署、避免跨域
 * - 本地开发若前后端分离（前端 Vite 与后端不同端口），可设置 VITE_API_BASE_URL 指向后端
 */
import axios from 'axios'
import { ElMessage } from 'element-plus'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

const request = axios.create({
  baseURL: API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理业务码与异常（200/201 视为成功，避免 201 Created 被当成错误）
request.interceptors.response.use(
  (response) => {
    const { data } = response
    const code = data && data.code
    if (typeof code !== 'undefined' && code !== 200 && code !== 201) {
      const msg = data.message || '请求失败'
      ElMessage.error(msg)
      return Promise.reject(new Error(msg))
    }
    return response
  },
  (error) => {
    const message =
      error.response?.data?.message ||
      error.message ||
      '网络异常，请稍后重试'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default request
export { API_BASE_URL }
