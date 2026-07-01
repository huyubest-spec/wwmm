<template>
  <div v-loading="loading">
    <h2 class="section-title">⛓ 区块链浏览器</h2>
    <el-alert type="info" :closable="false" show-icon style="margin-bottom: 20px;">
      <div slot="title">底层：自研 PoW 区块链 · Solidity 智能合约</div>
      本区块链由 Go 语言完整实现：采用 SHA-256 哈希算法 + 工作量证明（PoW）机制 + Merkle 树形结构。
      所有数据公开透明、不可篡改、可追溯。
    </el-alert>

    <div class="stat-grid" style="grid-template-columns: repeat(4, 1fr);">
      <div class="stat-card">
        <div class="label">⛓ 最新区块高度</div>
        <div class="value">{{ state.latestIndex || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="label">📦 累计区块</div>
        <div class="value">{{ state.totalBlocks || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="label">💱 累计交易</div>
        <div class="value">{{ state.totalTxs || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="label">⚙️ 挖矿难度</div>
        <div class="value">4 <span style="font-size: 12px; color: #8a96b3;">前缀零</span></div>
      </div>
    </div>

    <h2 class="section-title" style="margin-top: 28px;">最新区块</h2>
    <el-table :data="blocks" stripe>
      <el-table-column label="高度" prop="index" width="80"></el-table-column>
      <el-table-column label="区块哈希" min-width="200">
        <template slot-scope="s">
          <router-link :to="'/chain/block/' + s.row.index" class="hash-cell" style="color:#2c5fe0;">
            {{ (s.row.hash || '').substr(0, 40) }}…
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="前一区块" min-width="200">
        <template slot-scope="s">
          <span class="hash-cell">{{ (s.row.prevHash || '0'.repeat(64)).substr(0, 40) }}…</span>
        </template>
      </el-table-column>
      <el-table-column label="Merkle 根" min-width="200">
        <template slot-scope="s">
          <span class="hash-cell">{{ (s.row.merkleRoot || '').substr(0, 40) }}…</span>
        </template>
      </el-table-column>
      <el-table-column label="Nonce" prop="nonce" width="100"></el-table-column>
      <el-table-column label="交易数" prop="txCount" width="80"></el-table-column>
      <el-table-column label="时间" width="170">
        <template slot-scope="s">{{ formatTime(s.row.timestamp) }}</template>
      </el-table-column>
    </el-table>

    <h2 class="section-title" style="margin-top: 28px;">最新交易</h2>
    <el-table :data="txs" stripe>
      <el-table-column label="类型" width="120">
        <template slot-scope="s">
          <el-tag :type="s.row.txType===1?'primary':'warning'" size="mini">{{ s.row.txType===1?'作品存证':'投票' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="交易哈希" min-width="200">
        <template slot-scope="s">
          <router-link :to="'/chain/tx/' + s.row.txHash" class="hash-cell" style="color:#2c5fe0;">
            {{ (s.row.txHash || '').substr(0, 40) }}…
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="发送方" prop="sender" width="160"></el-table-column>
      <el-table-column label="状态" width="100">
        <template slot-scope="s">
          <el-tag :type="s.row.status===1?'success':'info'" size="mini" effect="plain">
            {{ s.row.status===1?'已打包':'待打包' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="区块" width="100">
        <template slot-scope="s">
          <router-link v-if="s.row.blockId" :to="'/chain/block/' + s.row.blockId" style="color:#2c5fe0;">#{{ s.row.blockId }}</router-link>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="时间" width="170">
        <template slot-scope="s">{{ s.row.createdAt }}</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'ChainExplorer',
  data() {
    return { loading: true, state: {}, blocks: [], txs: [] }
  },
  async mounted() { await this.load() },
  methods: {
    formatTime(ts) {
      if (!ts) return '-'
      const d = new Date(ts * 1000)
      return d.toLocaleString('zh-CN')
    },
    async load() {
      this.loading = true
      try {
        const [s, b, t] = await Promise.all([
          this.$api.get('/chain/state'),
          this.$api.get('/chain/blocks?size=10'),
          this.$api.get('/chain/txs?size=15')
        ])
        this.state = s.data || {}
        this.blocks = b.data.list || []
        this.txs = t.data.list || []
      } catch (e) {
        this.$message.error(e.message)
      } finally { this.loading = false }
    }
  }
}
</script>
