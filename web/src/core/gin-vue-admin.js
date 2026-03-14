import { register } from './global'
import packageInfo from '../../package.json'
import config from './config'

export default {
  install(app) {
    register(app)
    if (import.meta.env.DEV) {
      console.log(`[${config.appName}] v${packageInfo.version} - 开发模式`)
    }
  }
}
