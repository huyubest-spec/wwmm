<template>
  <div id="app">
    <el-container class="app-container" v-if="ready">
      <el-header class="app-header" height="64px">
        <div class="header-inner">
          <div class="brand" @click="$router.push('/')">
            <div class="brand-logo">W</div>
            <div class="brand-text">
              <div class="brand-title">WWMM 摄影存证</div>
              <div class="brand-sub">基于区块链的投票存证平台</div>
            </div>
          </div>
          <div class="nav">
            <router-link to="/" tag="span" class="nav-item">作品广场</router-link>
            <router-link to="/ranking" tag="span" class="nav-item">人气排行</router-link>
            <router-link to="/chain" tag="span" class="nav-item">区块链浏览器</router-link>
            <router-link v-if="user.role===1" to="/upload" tag="span" class="nav-item">上传作品</router-link>
            <router-link v-if="user.role===1" to="/my" tag="span" class="nav-item">我的作品</router-link>
            <router-link v-if="user.role===2" to="/admin" tag="span" class="nav-item">审核中心</router-link>
          </div>
          <div class="user-area">
            <template v-if="!user.userId">
              <el-button type="text" @click="$router.push('/login')">登录</el-button>
              <el-button type="primary" size="small" @click="$router.push('/register')">免费注册</el-button>
            </template>
            <el-dropdown v-else trigger="click" @command="onCommand">
              <span class="user-info">
                <el-avatar :size="32" icon="el-icon-user-solid"></el-avatar>
                <span class="user-name">{{ user.username }}</span>
                <el-tag :type="roleTag(user.role)" size="mini" effect="plain">{{ roleLabel(user.role) }}</el-tag>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item v-if="user.role===1" command="my">我的作品</el-dropdown-item>
                <el-dropdown-item v-if="user.role===2" command="admin">审核中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view :key="$route.fullPath"/>
      </el-main>
      <el-footer height="64px" class="app-footer">
        <div class="footer-line">
          <span>WWMM Photography Chain © 2026 - 基于区块链的摄影作品投票存证系统</span>
        </div>
        <div class="footer-line sub">
          <span>采用 Go Web + Vue 2 + MySQL 8 + 自研 PoW 区块链 · Solidity 智能合约</span>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      user: { userId: 0, username: '', role: 0 },
      ready: false
    }
  },
  created() {
    this.loadUser()
    this.ready = true
  },
  watch: {
    '$route'() {
      this.loadUser()
    }
  },
  methods: {
    loadUser() {
      const raw = localStorage.getItem('wwmm_user')
      if (raw) {
        try { this.user = JSON.parse(raw) } catch (e) {}
      }
    },
    roleLabel(r) {
      return ['普通用户', '摄影师', '管理员'][r] || '游客'
    },
    roleTag(r) {
      return ['', 'success', 'danger'][r] || 'info'
    },
    async onCommand(cmd) {
      if (cmd === 'logout') {
        const t = localStorage.getItem('wwmm_token')
        try { await this.$api.post('/user/logout', {}, { headers: { Token: t } }) } catch (e) {}
        localStorage.removeItem('wwmm_user')
        localStorage.removeItem('wwmm_token')
        this.user = { userId: 0, username: '', role: 0 }
        this.$message.success('已退出登录')
        this.$router.push('/')
      } else if (cmd === 'my') {
        this.$router.push('/my')
      } else if (cmd === 'admin') {
        this.$router.push('/admin')
      }
    }
  }
}
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: #f4f6fb;
}
.app-header {
  background: #ffffff;
  border-bottom: 1px solid #e6ebf5;
  padding: 0 !important;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.header-inner {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 32px;
  max-width: 1280px;
  margin: 0 auto;
}
.brand {
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-right: 48px;
}
.brand-logo {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #2c5fe0, #5b8def);
  color: #fff;
  font-weight: 900;
  font-size: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
}
.brand-title {
  font-weight: 700;
  font-size: 16px;
  color: #1f2a44;
}
.brand-sub {
  font-size: 11px;
  color: #8a96b3;
  margin-top: 2px;
}
.nav {
  flex: 1;
  display: flex;
  gap: 4px;
}
.nav-item {
  padding: 8px 14px;
  font-size: 14px;
  color: #4a5568;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}
.nav-item:hover {
  color: #2c5fe0;
  background: #eef2fb;
}
.nav-item.router-link-active {
  color: #2c5fe0;
  font-weight: 600;
  background: #eef2fb;
}
.user-area {
  display: flex;
  align-items: center;
  gap: 12px;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
}
.user-info:hover {
  background: #eef2fb;
}
.user-name {
  font-size: 14px;
  color: #1f2a44;
}
.app-main {
  max-width: 1280px;
  margin: 0 auto;
  padding: 24px 32px !important;
  width: 100%;
  box-sizing: border-box;
}
.app-footer {
  text-align: center;
  color: #8a96b3;
  font-size: 12px;
  background: #ffffff;
  border-top: 1px solid #e6ebf5;
  padding-top: 12px !important;
}
.footer-line.sub {
  font-size: 11px;
  margin-top: 4px;
}
</style>
