<template>
  <div style="max-width: 720px; margin: 0 auto;">
    <h2 class="section-title">⛓ 哈希存证校验</h2>
    <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px;">
      <div slot="title">输入作品的 SHA-256 哈希，验证其是否被记录在区块链上</div>
      任何人都可以独立验证作品存证的真实性。
    </el-alert>
    <el-input v-model="hash" placeholder="请输入 64 字符的 SHA-256 哈希值" size="medium">
      <el-button slot="append" type="primary" :loading="loading" @click="verify">校验</el-button>
    </el-input>

    <div v-if="result" class="block-card" style="margin-top: 20px;">
      <el-result :icon="result.verified?'success':'error'" :title="result.verified?'校验通过 · 存证有效':'未找到该哈希'">
        <div slot="extra" v-if="result.verified">
          <p>该图片哈希已在区块链上记录，对应作品：</p>
          <el-descriptions :column="1" border size="small" style="text-align: left;">
            <el-descriptions-item label="作品标题">{{ result.photo.title }}</el-descriptions-item>
            <el-descriptions-item label="作者ID">{{ result.photo.photographerId }}</el-descriptions-item>
            <el-descriptions-item label="上链交易">
              <router-link :to="'/chain/tx/' + result.photo.chainTxHash" class="hash-cell" style="color:#2c5fe0;">{{ result.photo.chainTxHash }}</router-link>
            </el-descriptions-item>
            <el-descriptions-item label="上传时间">{{ result.photo.createdAt }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-result>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Verify',
  data() {
    return { hash: '', loading: false, result: null }
  },
  methods: {
    async verify() {
      if (this.hash.length !== 64) {
        this.$message.warning('请输入 64 字符的 SHA-256 哈希')
        return
      }
      this.loading = true
      try {
        const r = await this.$api.get('/chain/verify/' + this.hash)
        this.result = r.data
      } catch (e) {
        this.result = { verified: false, message: e.message }
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
