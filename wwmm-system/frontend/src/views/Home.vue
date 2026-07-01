<template>
  <div>
    <div class="banner">
      <h1>基于区块链的摄影作品投票存证平台</h1>
      <p>每一张作品上链、每一次投票存证。我们用 Go 语言自研 PoW 区块链引擎，配合 Solidity 智能合约设计稿，
        让你的作品在去中心化的世界中拥有不可篡改的存在证明。</p>
      <span class="chain-pill">⛓ 当前区块高度：{{ chainState.latestIndex || '-' }} · 累计交易：{{ chainState.totalTxs || 0 }}</span>
    </div>

    <div class="stat-grid">
      <div class="stat-card">
        <div class="label">⛓ 区块高度</div>
        <div class="value">{{ chainState.latestIndex || 0 }}</div>
        <div class="chain-hash" v-if="chainState.latestHash">{{ chainState.latestHash.substr(0, 32) }}…</div>
      </div>
      <div class="stat-card">
        <div class="label">📦 累计交易</div>
        <div class="value">{{ chainState.totalTxs || 0 }}</div>
        <div class="text-mute" style="font-size: 12px; margin-top: 4px;">不可篡改 · 公开可查</div>
      </div>
      <div class="stat-card">
        <div class="label">🖼 作品存证</div>
        <div class="value">{{ chainState.txCertifyCount || 0 }}</div>
        <div class="text-mute" style="font-size: 12px; margin-top: 4px;">SHA-256 上链</div>
      </div>
      <div class="stat-card">
        <div class="label">🗳 投票记录</div>
        <div class="value">{{ chainState.txVoteCount || 0 }}</div>
        <div class="text-mute" style="font-size: 12px; margin-top: 4px;">真实可信</div>
      </div>
    </div>

    <h2 class="section-title" style="margin-top: 28px;">热门作品</h2>
    <div v-loading="loading" element-loading-text="正在加载作品...">
      <div v-if="!photos.length" class="empty-hint">
        <i class="el-icon-picture-outline" style="font-size: 40px; color: #cdd5e6;"></i>
        <p>暂无作品，快上传第一张吧！</p>
      </div>
      <div v-else class="photo-grid">
        <div v-for="p in photos" :key="p.photoId" class="photo-card" @click="$router.push('/photo/' + p.photoId)">
          <div class="image" :style="{ backgroundImage: 'url(' + (p.imageUrl.startsWith('http') ? p.imageUrl : ($root.$api.defaults.baseURL.replace('/api','') + p.imageUrl)) + ')' }">
            <span class="chain-tag" v-if="p.isOnChain">⛓ 已上链</span>
          </div>
          <div class="info">
            <div class="title">{{ p.title }}</div>
            <div class="author">📷 {{ p.photographerName || '匿名' }}</div>
            <div class="meta">
              <span class="vote">❤ {{ p.voteCount || 0 }}</span>
              <span>👁 {{ p.viewCount || 0 }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home',
  data() {
    return {
      photos: [],
      chainState: {},
      loading: true
    }
  },
  async mounted() {
    await this.load()
  },
  methods: {
    async load() {
      this.loading = true
      try {
        const [p, s] = await Promise.all([
          this.$api.get('/photo/list'),
          this.$api.get('/chain/state')
        ])
        this.photos = p.data.list || []
        this.chainState = s.data || {}
      } catch (e) {
        this.$message.error('数据加载失败：' + e.message)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
