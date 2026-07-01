<template>
  <div>
    <h2 class="section-title">🏆 人气排行榜</h2>
    <el-alert type="success" :closable="false" show-icon style="margin-bottom: 20px;">
      <div slot="title">排行榜基于区块链上链的投票数据实时生成</div>
      每一次投票都会写入区块链，因此排名是公开透明、不可伪造的。
    </el-alert>
    <div v-if="!list.length" class="empty-hint">暂无上榜作品</div>
    <div v-for="(p, idx) in list" :key="p.photoId" class="ranking-row" @click="$router.push('/photo/' + p.photoId)" style="cursor: pointer;">
      <div :class="['rank-no', idx===0?'rank-1':idx===1?'rank-2':idx===2?'rank-3':'rank-other']">
        {{ idx + 1 }}
      </div>
      <div :style="{ width:'80px', height:'60px', backgroundImage:'url(' + getImg(p.imageUrl) + ')', backgroundSize:'cover', backgroundPosition:'center', borderRadius:'4px', marginRight: '16px' }"></div>
      <div style="flex: 1;">
        <div style="font-weight: 600; font-size: 15px;">{{ p.title }}</div>
        <div style="color: #8a96b3; font-size: 12px; margin-top: 4px;">📷 {{ p.photographerName }} · {{ p.category || '未分类' }}</div>
      </div>
      <div style="text-align: right;">
        <div style="font-size: 18px; font-weight: 700; color: #d44030;">❤ {{ p.voteCount }}</div>
        <div style="font-size: 12px; color: #8a96b3;">累计票数</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Ranking',
  data() { return { list: [] } },
  async mounted() { await this.load() },
  methods: {
    getImg(url) {
      if (!url) return ''
      if (url.startsWith('http')) return url
      const base = (this.$api.defaults.baseURL || '').replace('/api', '')
      return base + url
    },
    async load() {
      try {
        const r = await this.$api.get('/photo/list?size=50')
        this.list = (r.data.list || []).slice().sort((a, b) => b.voteCount - a.voteCount)
      } catch (e) {
        this.$message.error(e.message)
      }
    }
  }
}
</script>
