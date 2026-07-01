<template>
  <div v-loading="loading">
    <el-page-header @back="$router.push('/chain')" content="交易详情" style="margin-bottom: 16px;"></el-page-header>
    <div v-if="tx.txId" class="block-card">
      <h2 style="margin: 0 0 16px 0;">
        <el-tag :type="tx.txType===1?'primary':'warning'" size="medium">{{ tx.txType===1?'作品存证':'投票' }}</el-tag>
        <span style="margin-left: 10px;">交易详情</span>
      </h2>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="交易哈希"><span class="hash-cell">{{ tx.txHash }}</span></el-descriptions-item>
        <el-descriptions-item label="所在区块">
          <router-link v-if="tx.blockId" :to="'/chain/block/' + tx.blockId" style="color:#2c5fe0;">#{{ tx.blockId }}</router-link>
          <span v-else>待打包</span>
        </el-descriptions-item>
        <el-descriptions-item label="发送方">{{ tx.sender }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="tx.status===1?'success':'info'" size="mini">
            {{ tx.status===1?'已打包':'待打包' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ tx.createdAt }}</el-descriptions-item>
      </el-descriptions>
      <h3 class="section-title" style="margin-top: 18px;">交易载荷 (Payload)</h3>
      <pre class="tx-payload">{{ tx.payload }}</pre>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TxDetail',
  props: ['hash'],
  data() { return { loading: true, tx: {} } },
  async mounted() { await this.load() },
  methods: {
    async load() {
      this.loading = true
      try {
        const r = await this.$api.get('/chain/tx/' + this.hash)
        this.tx = r.data
      } catch (e) {
        this.$message.error(e.message)
      } finally { this.loading = false }
    }
  }
}
</script>
