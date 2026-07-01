<template>
  <div style="max-width: 480px; margin: 40px auto;">
    <el-card>
      <div slot="header" class="clearfix">
        <span style="font-size: 18px; font-weight: 600;">欢迎登录 WWMM</span>
        <span style="float: right; color: #8a96b3; font-size: 12px;">基于区块链的摄影存证平台</span>
      </div>
      <el-form :model="form" :rules="rules" ref="form" label-width="0" size="medium">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="账号" prefix-icon="el-icon-user" clearable></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" prefix-icon="el-icon-lock" show-password></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%;" @click="onSubmit">登录</el-button>
        </el-form-item>
      </el-form>
      <div style="text-align: center; color: #8a96b3; font-size: 13px; margin-top: 10px;">
        还没有账号？<router-link to="/register" style="color: #2c5fe0;">立即注册</router-link>
      </div>

      <el-divider>演示账号</el-divider>
      <div style="font-size: 12px; color: #8a96b3; line-height: 1.9;">
        <div>👤 管理员：admin / admin123</div>
        <div>📷 摄影师：photographer / photo123</div>
        <div>🗳 投票用户：voter / vote123</div>
        <el-button size="mini" type="text" @click="quickFill('admin', 'admin123')">填入管理员</el-button>
        <el-button size="mini" type="text" @click="quickFill('photographer', 'photo123')">填入摄影师</el-button>
        <el-button size="mini" type="text" @click="quickFill('voter', 'vote123')">填入投票用户</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      loading: false,
      form: { username: '', password: '' },
      rules: {
        username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      }
    }
  },
  methods: {
    quickFill(u, p) {
      this.form.username = u
      this.form.password = p
    },
    async onSubmit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          const r = await this.$api.post('/user/login', this.form)
          const d = r.data
          localStorage.setItem('wwmm_token', d.token)
          localStorage.setItem('wwmm_user', JSON.stringify({ userId: d.userId, username: d.username, role: d.role, realName: d.realName }))
          this.$message.success('登录成功')
          const redirect = this.$route.query.redirect || '/'
          this.$router.replace(redirect)
          window.location.reload()
        } catch (e) {
          this.$message.error(e.message || '登录失败')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>
