<template>
  <div v-loading="loading">
    <el-page-header @back="$router.push('/chain')" :content="'区块 #' + index" style="margin-bottom: 16px;"></el-page-header>
    <div v-if="block.index" class="block-card">
      <h2 style="margin: 0 0 16px 0;">⛓ 区块 #{{ block.index }}</h2>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="区块高度">{{ block.index }}</el-descriptions-item>
        <el-descriptions-item label="时间戳">{{ formatTime(block.timestamp) }}</el-descriptions-item>
        <el-descriptions-item label="交易数">{{ block.txCount }}</el-descriptions-item>
        <el-descriptions-item label="挖矿难度">{{ block.difficulty }}</el-descriptions-item>
        <el-descriptions-item label="Nonce" :span="2">
          <span class="hash-cell">{{ block.nonce }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="区块哈希" :span="2">
          <span class="hash-cell">{{ block.hash }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="前一区块哈希" :span="2">
          <span class="hash-cell">{{ block.prevHash }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="Merkle 根" :span="2">
          <span class="hash-cell">{{ block.merkleRoot }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="打包者" :span="2">{{ block.miner }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <h2 class="section-title">📦 区块内交易 ({{ txs.length }} 笔)</h2>
    <div v-if="!txs.length" class="empty-hint">该区块是创世区块，没有交易</div>
    <div v-for="t in txs" :key="t.txId" class="block-card">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
        <div>
          <el-tag :type="t.txType===1?'primary':'warning'" size="small">{{ t.txType===1?'作品存证':'投票' }}</el-tag>
          <span class="hash-cell" style="margin-left: 10px;">{{ t.txHash }}</span>
        </div>
        <span style="color: #8a96b3; font-size: 12px;">{{ t.createdAt }}</span>
      </div>
      <div style="color: #4a5568; font-size: 13px; margin-bottom: 6px;">发送方：<b>{{ t.sender }}</b></div>
      <pre class="tx-payload">{{ t.payload }}</pre>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BlockDetail',
  props: ['index'],
  data() { return { loading: true, block: {}, txs: [] } },
  async mounted() { await this.load() },
  methods: {
    formatTime(ts) {
      if (!ts) return '-'
      return new Date(ts * 1000).toLocaleString('zh-CN')
    },
    async load() {
      this.loading = true
      try {
        const r = await this.$api.get('/chain/block/' + this.index)
        this.block = r.data.block
        this.txs = r.data.txs || []
      } catch (e) {
        this.$message.error(e.message)
      } finally { this.loading = false }
    }
  }
}
</script>
