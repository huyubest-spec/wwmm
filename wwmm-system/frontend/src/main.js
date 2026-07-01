import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import router from './router'
import api from './api'
import './styles/global.css'

Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.prototype.$api = api
Vue.prototype.$ELEMENT = { size: 'medium' }

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
