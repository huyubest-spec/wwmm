<template>
  <div style="max-width: 560px; margin: 40px auto;">
    <el-card>
      <div slot="header" class="clearfix">
        <span style="font-size: 18px; font-weight: 600;">注册新账号</span>
        <span style="float: right; color: #8a96b3; font-size: 12px;">链上不可篡改 · 公开透明</span>
      </div>
      <el-form :model="form" :rules="rules" ref="form" label-width="80px" size="medium">
        <el-form-item label="账号" prop="username">
          <el-input v-model="form.username" placeholder="3-16位字母数字下划线"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="至少6位" show-password></el-input>
        </el-form-item>
        <el-form-item label="真实姓名">
          <el-input v-model="form.realName" placeholder="选填"></el-input>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="选填"></el-input>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="选填"></el-input>
        </el-form-item>
        <el-form-item label="性别">
          <el-radio-group v-model="form.sex">
            <el-radio :label="0">未设置</el-radio>
            <el-radio :label="1">男</el-radio>
            <el-radio :label="2">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="注册为">
          <el-radio-group v-model="form.role">
            <el-radio :label="0">普通用户（投票者）</el-radio>
            <el-radio :label="1">摄影师（可上传作品）</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%;" @click="onSubmit">注册并登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    return {
      loading: false,
      form: {
        username: '',
        password: '',
        realName: '',
        phone: '',
        email: '',
        sex: 0,
        role: 0
      },
      rules: {
        username: [
          { required: true, message: '请输入账号', trigger: 'blur' },
          { min: 3, max: 16, message: '长度 3-16', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: '只允许字母数字下划线', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '至少6位', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    async onSubmit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          await this.$api.post('/user/register', this.form)
          this.$message.success('注册成功，正在登录...')
          const r = await this.$api.post('/user/login', { username: this.form.username, password: this.form.password })
          const d = r.data
          localStorage.setItem('wwmm_token', d.token)
          localStorage.setItem('wwmm_user', JSON.stringify({ userId: d.userId, username: d.username, role: d.role, realName: d.realName }))
          this.$router.replace('/')
          window.location.reload()
        } catch (e) {
          this.$message.error(e.message || '注册失败')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>
