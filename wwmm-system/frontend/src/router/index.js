import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import PhotoDetail from '../views/PhotoDetail.vue'
import UploadPhoto from '../views/UploadPhoto.vue'
import MyPhotos from '../views/MyPhotos.vue'
import Admin from '../views/Admin.vue'
import ChainExplorer from '../views/ChainExplorer.vue'
import BlockDetail from '../views/BlockDetail.vue'
import TxDetail from '../views/TxDetail.vue'
import Ranking from '../views/Ranking.vue'
import Verify from '../views/Verify.vue'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'hash',
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/login', name: 'login', component: Login },
    { path: '/register', name: 'register', component: Register },
    { path: '/photo/:id', name: 'photo-detail', component: PhotoDetail, props: true },
    { path: '/upload', name: 'upload', component: UploadPhoto, meta: { auth: true, role: 1 } },
    { path: '/my', name: 'my', component: MyPhotos, meta: { auth: true, role: 1 } },
    { path: '/admin', name: 'admin', component: Admin, meta: { auth: true, role: 2 } },
    { path: '/chain', name: 'chain', component: ChainExplorer },
    { path: '/chain/block/:index', name: 'block-detail', component: BlockDetail, props: true },
    { path: '/chain/tx/:hash', name: 'tx-detail', component: TxDetail, props: true },
    { path: '/chain/verify', name: 'verify', component: Verify },
    { path: '/ranking', name: 'ranking', component: Ranking }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta && to.meta.auth) {
    const raw = localStorage.getItem('wwmm_user')
    if (!raw) return next({ path: '/login', query: { redirect: to.fullPath } })
    try {
      const u = JSON.parse(raw)
      if (to.meta.role !== undefined && u.role !== to.meta.role && u.role !== 2) {
        return next({ path: '/' })
      }
    } catch (e) {
      return next({ path: '/login' })
    }
  }
  next()
})

export default router
