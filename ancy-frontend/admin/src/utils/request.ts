import router from '@/router'
import { useUserStore } from '@/stores/modules/user'
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 5000,
})

const noTokenUrl = ['/user/login', '/user/register']
// 添加请求拦截器
request.interceptors.request.use(
  function (config) {
    const userStore = useUserStore()
    const token = userStore.getToken()
    if (token && !noTokenUrl.includes(config.url as string)) {
      config.headers.token = `${token}`
    }
    return config
  },
  function (error) {
    // 对请求错误做些什么
    return Promise.reject(error)
  },
)

// 添加响应拦截器
request.interceptors.response.use(
  function (response) {
    // 2xx 范围内的状态码都会触发该函数。
    // 对响应数据做点什么
    if (
      response.headers['content-type'] ===
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8'
    ) {
      // download file
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.style.display = 'none'
      link.href = url
      link.setAttribute('download', 'file.xlsx')
      document.body.appendChild(link)
      link.click()
      return response.data
    }
    if (response.data.code !== 200) {
      ElMessage.error(response.data.msg)
      return Promise.reject(response.data)
    }
    if (response.data.data === true) {
      ElMessage.success(response.data.msg)
    }
    return response.data
  },
  function (error) {
    // 超出 2xx 范围的状态码都会触发该函数。
    // 对响应错误做点什么
    if (error.response.data.code === 510) {
      router.push('/login')
      ElMessage.error('登录过期，请重新登录')
      return Promise.reject(error)
    }
    const msg = error.response.data.msg
    ElMessage.error(msg)
    return Promise.reject(error)
  },
)
export default request
