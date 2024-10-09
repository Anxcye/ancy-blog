import { message } from 'ant-design-vue'
import { useAdminStore } from '@/stores/admin'
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 5000,
})

const adminRequest = ['/comment']

request.interceptors.request.use(
  function (config) {
    if (useAdminStore().token && adminRequest.some((url) => config.url?.startsWith(url))) {
      config.headers.token = useAdminStore().token
    }
    return config
  },
  function (error) {
    return Promise.reject(error)
  },
)

request.interceptors.response.use(
  function (response) {
    if (response.data.code !== 200) {
      message.error(response.data.msg)
    }
    return response.data
  },
  function (error) {
    if (error.response.status === 510) {
      message.error('管理员登录过期，请重新登录')
      useAdminStore().logout()
    }
    message.error(error.response.data.msg)
    return Promise.reject(error)
  },
)
export default request
